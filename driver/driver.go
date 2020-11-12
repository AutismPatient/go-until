package driver

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/docker/docker/client"
	"github.com/garyburd/redigo/redis"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"
)

type IDriverHelper interface {
}

type BaseConnStrut struct {
	UserName string
	Pass     string
	DB       string
	Addr     string
	Port     int
	SSLMode  string
}

type RelationSqlStrut struct {
	*sql.DB
}
type RelationSqlOption struct {
	BaseConnStrut
	ConnectString string
	MysqlSetting
}
type MysqlSetting struct {
	MaxOpenConn int
	MaxLifetime time.Duration
	MaxIdleConn int
}

/*
	MYSQL 连接使用TCP协议
	打开由数据库驱动程序名称和参数指定的数据库
	驱动程序特定的数据源名称，通常至少由一个
	数据库名称和连接信息。

	大多数用户会通过驱动程序特定的连接来打开数据库
	帮助函数，返回一个*DB。不包括数据库驱动程序
	在Go标准库中。参见https:golang.org/s/sqldrivers
	第三方驱动程序列表。

	Open可能只验证它的参数，而不创建连接
	到数据库。若要验证数据源名称是否有效，请调用
	/ /平。

	返回的DB对于多个goroutines并发使用是安全的
	并维护自己的空闲连接池。因此,开放
	函数只调用一次。很少有必要这样做
	关闭一个DB。

	TODO 待实现更方便的命令执行方法
*/

func NewMySQL(opt RelationSqlOption) (itf *RelationSqlStrut, err error) {
	var (
		source = ""
		conn   *sql.DB
	)
	if opt.ConnectString != "" {
		source = opt.ConnectString
	} else {
		source = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", opt.UserName, opt.Pass, opt.Addr, opt.Port, opt.DB)
	}
	conn, err = sql.Open("mysql", source)
	if err != nil || conn.Ping() != nil {
		log.Fatal(err.Error())
		return nil, err
	}
	conn.SetMaxOpenConns(opt.MaxOpenConn)
	conn.SetConnMaxLifetime(opt.MaxLifetime)
	conn.SetMaxIdleConns(opt.MaxIdleConn)
	itf.DB = conn
	return
}

/*
	获取单条记录 TODO 2020年11月3日20:52:53
*/
func (itf *RelationSqlStrut) GetSingleByID(db string, rely *interface{}, id string) (err error) {
	var (
		fields string
	)
	tf := reflect.TypeOf(rely)
	tv := reflect.ValueOf(rely)
	switch tf.Kind() {
	case reflect.Struct:
		var (
			lds []string
		)
		for i := 0; i < tv.NumField(); i++ {
			field := tv.Field(i).String()
			lds = append(lds, field)
		}
		fields = strings.Join(lds, ",")
	default:
		return errors.New("类型错误，并非 reflect.Struct 类型")
	}

	return itf.DB.QueryRow("SELECT ? FROM ? WHERE id=?", fields, db, id).Scan()
}

/*
	MYSQL 运行状态
*/
func (itf *RelationSqlStrut) Status() (status sql.DBStats) {
	return itf.DB.Stats()
}

type RedisStrut struct {
	*redis.Pool
}
type RedisOptions struct {
	BaseConnStrut
	DialOptions
	PoolOptions
}
type DialOptions struct {
	DBNum          int
	ConnectTimeout time.Duration
	WriteTimeout   time.Duration
	ReadTimeout    time.Duration
	KeepAlive      time.Duration
	UseTLS         bool
}
type PoolOptions struct {
	MaxIdle         int
	MaxActive       int
	IdleTimeout     time.Duration
	Wait            bool
	MaxConnLifetime time.Duration
}

func NewRedis(opt RedisOptions) (pool *RedisStrut, err error) {
	connFunc := func() (conn redis.Conn, err error) {
		conn, err = redis.Dial("tcp", opt.Addr,
			redis.DialConnectTimeout(opt.ConnectTimeout),
			redis.DialWriteTimeout(opt.WriteTimeout),
			redis.DialPassword(opt.Pass),
			redis.DialKeepAlive(opt.KeepAlive), // 默认的5分钟用于确保检测到半关闭的TCP会话
			redis.DialReadTimeout(opt.ReadTimeout),
			redis.DialDatabase(opt.DBNum),
			redis.DialUseTLS(opt.UseTLS), // 指定当连接到的时候是否应该使用TLS
		)
		if err != nil || conn.Err() != nil {
			log.Fatal(err.Error())
			return nil, err
		}
		return
	}
	pool.Pool = &redis.Pool{
		Dial:            connFunc,
		MaxIdle:         opt.MaxIdle,         // 池中的最大空闲连接数
		MaxActive:       opt.MaxActive,       // 在给定时间池分配的最大连接数。当为0时，池中的连接数没有限制。
		IdleTimeout:     opt.IdleTimeout,     // 在此期间保持空闲状态后关闭连接。如果该值为零，则空闲连接未关闭。应用程序应该设置将超时设置为小于服务器超时的值。
		Wait:            opt.Wait,            // 如果Wait为真，并且池处于MaxActive限制，则Get()等待 ;用于在返回之前将连接返回到池。
		MaxConnLifetime: opt.MaxConnLifetime, // 比这段时间更久的紧密联系。如果值为零，则这个连接池不会根据时间来拉近联系。
	}
	return
}
func (p *RedisStrut) Do(commandName string, args ...interface{}) (reply interface{}, err error) {
	conn := p.Get()
	defer conn.Close()
	return conn.Do(commandName, args)
}

/*
	发送至缓冲区
*/
func (p *RedisStrut) Send(commandName string, args ...interface{}) error {
	conn := p.Get()
	defer conn.Close()
	return conn.Send(commandName, args)
}

/*
	清空缓冲区
*/
func (p *RedisStrut) Flush() error {
	conn := p.Get()
	defer conn.Close()
	return conn.Flush()
}

/*
	读取队列,未响应时堵塞
*/
func (p *RedisStrut) Receive() (reply interface{}, err error) {
	conn := p.Get()
	defer conn.Close()
	return conn.Receive()
}

type MongoDBStrut struct {
	*mongo.Client
	*mongo.Database
	*mongo.Collection
}

/*
	mongodb+srv://<username>:<password>@<cluster-address>/test?w=majority

	详情DOC：https://docs.mongodb.com/drivers/go

	TODO 待实现更方便的命令执行方法
*/

func NewMongoDB(database string, opt *options.ClientOptions, collection string) (client *MongoDBStrut, err error) {
	var (
		ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	)
	defer cancel()
	client = &MongoDBStrut{}
	client.Client, err = mongo.Connect(ctx, opt)
	if err != nil {
		return nil, err
	}
	if err = client.Client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}
	if database != "" {
		client.Database = client.Client.Database(database)
		if collection != "" {
			client.Collection = client.Database.Collection(collection)
		}
	}
	return
}

/*
	驱动地址： https://github.com/denisenkom/go-mssqldb
	建议的连接字符串使用 URL 格式：下面列出了其他受支持的格式。sqlserver://username:password@host/instance?param1=value&param2=value
*/
func NewMSSQL(opt RelationSqlOption, query url.Values) (rst *RelationSqlStrut, err error) {
	var (
		scheme     = "sqlserver"
		connecting = ""
	)
	if opt.ConnectString != "" {
		scheme = "mssql"
		connecting = opt.ConnectString
	} else {
		u := &url.URL{
			Scheme:   scheme,
			User:     url.UserPassword(opt.UserName, opt.Pass),
			Host:     fmt.Sprintf("%s:%d", opt.Addr, opt.Port),
			RawQuery: query.Encode(),
		}
		connecting = u.String()
	}

	conn, err := sql.Open(scheme, connecting)
	if err = conn.Ping(); err != nil {
		log.Fatal(err.Error())
		return nil, err
	}
	conn.SetMaxOpenConns(opt.MaxOpenConn)
	conn.SetConnMaxLifetime(opt.MaxLifetime)
	conn.SetMaxIdleConns(opt.MaxIdleConn)
	rst.DB = conn
	return
}

/*
	https://github.com/lib/pq
	postgres://pqgotest:password@localhost/pqgotest?sslmode=verify-full
	doc : https://godoc.org/github.com/lib/pq

	相关连接参数：

	* dbname - The name of the database to connect to
	* user - The user to sign in as
	* password - The user's password
	* host - The host to connect to. Values that start with / are for unix
	  domain sockets. (default is localhost)
	* port - The port to bind to. (default is 5432)
	* sslmode - Whether or not to use SSL (default is require, this is not
	  the default for libpq)
	* fallback_application_name - An application_name to fall back to if one isn't provided.
	* connect_timeout - Maximum wait for connection, in seconds. Zero or
	  not specified means wait indefinitely.
	* sslcert - Cert file location. The file must contain PEM encoded data.
	* sslkey - Key file location. The file must contain PEM encoded data.
	* sslrootcert - The location of the root certificate file. The file
	  must contain PEM encoded data.

	Valid values for sslmode are:

	* disable - No SSL
	* require - Always SSL (skip verification)
	* verify-ca - Always SSL (verify that the certificate presented by the
	  server was signed by a trusted CA)
	* verify-full - Always SSL (verify that the certification presented by
	  the server was signed by a trusted CA and the server host name
	  matches the one in the certificate)

	TODO 待实现更方便的命令执行方法
*/
func NewPostgreSQL(opt RelationSqlOption) (rst *RelationSqlStrut, err error) {
	var (
		source = ""
		conn   *sql.DB
	)
	rst = &RelationSqlStrut{}
	if opt.ConnectString != "" {
		source = opt.ConnectString
	} else {
		if opt.Port == 0 {
			opt.Port = 5432
		}
		if opt.SSLMode == "" {
			opt.SSLMode = "disable"
		}
		source = fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s", opt.UserName, opt.Pass, opt.Addr, opt.Port, opt.DB, opt.SSLMode)
	}
	conn, err = sql.Open("postgres", source)
	if err != nil {
		return nil, err
	}
	if err = conn.Ping(); err != nil {
		return nil, err
	}
	conn.SetMaxOpenConns(opt.MaxOpenConn)
	conn.SetConnMaxLifetime(opt.MaxLifetime)
	conn.SetMaxIdleConns(opt.MaxIdleConn)
	rst.DB = conn
	return
}

/*
	DOC ：https://docs.docker.com/engine/api/sdk/examples/
	docker命令详情见官网
*/
type DockerOptions struct {
	Host    string
	Version string
	*http.Client
	HttpHeader map[string]string
}
type DockerClient struct {
	*client.Client
}

func NewDocker(opt DockerOptions) (cli *DockerClient, err error) {
	cli.Client, err = client.NewClient(opt.Host, opt.Version, opt.Client, opt.HttpHeader)
	if err != nil {
		return nil, err
	}
	if pg, err := cli.Ping(context.Background()); err != nil {
		log.Fatal(pg)
		return nil, err
	}
	return
}
