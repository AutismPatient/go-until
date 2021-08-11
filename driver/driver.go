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
	"go-until/config"
	_ "go-until/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"sync"
	"time"
)

var (
	RedisClient *RedisStrut
)

type IDriverHelper interface {
}

func init() {
	client, err := NewRedis(RedisOptions{
		BaseConnStrut: BaseConnStrut{
			UserName: config.Config.Redis.User,
			Pass:     config.Config.Redis.Password,
			DB:       "",
			Addr:     config.Config.Redis.Host,
			Port:     config.Config.Redis.Port,
			SSLMode:  "",
		},
		DialOptions: DialOptions{
			DBNum:          config.Config.Redis.Db,
			ConnectTimeout: config.Config.Redis.ConnectTimeOut,
			WriteTimeout:   config.Config.Redis.WriteTimeOut,
			ReadTimeout:    config.Config.Redis.ReadTimeOut,
			KeepAlive:      0,
			UseTLS:         false,
		},
		PoolOptions: PoolOptions{
			MaxIdle:         config.Config.Redis.MaxIdle,
			MaxActive:       0,
			IdleTimeout:     0,
			Wait:            false,
			MaxConnLifetime: 0,
		},
	})
	if err != nil {
		panic(err)
	}
	RedisClient = client
}

type BaseConnStrut struct {
	UserName string
	Pass     string
	DB       string
	Addr     string
	Port     string
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

type FieldInfo struct {
	Name     string
	TypeInfo interface{}
}

/*
	获取单条记录 TODO 2020年11月3日20:52:53 仅支持常用内置类型
	指明要返回的字段，在参数rely里填入默认值以明确要返回的列
*/
func (itf *RelationSqlStrut) GetSingleCtx(ctx context.Context, query string, rely interface{}, args ...interface{}) (single map[string]interface{}, err error) {
	var (
		lds []FieldInfo
		tf  = reflect.TypeOf(rely)
		tv  = reflect.ValueOf(rely)
	)
	switch tf.Kind() {
	case reflect.Struct:
		for i := 0; i < tv.NumField(); i++ {
			field := tf.Field(i).Name
			if !tv.Field(i).IsNil() {
				lds = append(lds, FieldInfo{
					Name:     field,
					TypeInfo: tf.Field(i).Type,
				})
			}
		}
	default:
		return nil, errors.New("类型错误，并非 reflect.Struct 类型")
	}
	var scans, pointer = make([]interface{}, len(lds)), make([]interface{}, len(lds))
	for i, _ := range pointer {
		scans[i] = &pointer[i]
	}
	err = itf.DB.QueryRowContext(ctx, query, args).Scan(scans...)
	if err != nil {
		return nil, err
	}
	for k, v := range pointer {
		fields := lds[k]
		switch v.(type) {
		case int:
			single[fields.Name] = v.(int)
		case int64:
			single[fields.Name] = v.(int64)
		case string:
			single[fields.Name] = v.(string)
		case float64:
			single[fields.Name] = v.(float64)
		case float32:
			single[fields.Name] = v.(float32)
		case bool:
			single[fields.Name] = v.(bool)
		case time.Time:
			single[fields.Name] = v.(time.Time)
		}
	}
	return
}

/*
	获取多条记录 TODO 2020年11月15日20:21:13
	指明要返回的字段，在参数rely里填入默认值以明确要返回的列
*/
func (itf *RelationSqlStrut) FindOfCtx(ctx context.Context, query string, rely interface{}, args ...interface{}) (list []map[string]interface{}, err error) {
	var (
		lds []FieldInfo
		tf  = reflect.TypeOf(rely)
		tv  = reflect.ValueOf(rely)
	)
	rows, err := itf.DB.QueryContext(ctx, query, args)
	if err != nil {
		return nil, err
	}
	switch tf.Kind() {
	case reflect.Struct:
		for i := 0; i < tv.NumField(); i++ {
			field := tf.Field(i).Name
			if !tv.Field(i).IsNil() {
				lds = append(lds, FieldInfo{
					Name:     field,
					TypeInfo: tf.Field(i).Type,
				})
			}
		}
	default:
		return nil, errors.New("类型错误，并非 reflect.Struct 类型")
	}
	var scans, pointer = make([]interface{}, len(lds)), make([]interface{}, len(lds))
	for i, _ := range pointer {
		scans[i] = &pointer[i]
	}
	for rows.Next() {
		err = rows.Scan(scans...)
		if err != nil {
			log.Fatal(err)
			continue
		}
		l := make(map[string]interface{})
		for k, v := range pointer {
			fields := lds[k]
			switch v.(type) {
			case int:
				l[fields.Name] = v.(int)
			case int64:
				l[fields.Name] = v.(int64)
			case string:
				l[fields.Name] = v.(string)
			case float64:
				l[fields.Name] = v.(float64)
			case float32:
				l[fields.Name] = v.(float32)
			case bool:
				l[fields.Name] = v.(bool)
			case time.Time:
				l[fields.Name] = v.(time.Time)
			}
		}
		list = append(list, l)
	}
	defer rows.Close()
	return
}

/*
	执行操作 TODO 2020年11月16日13:47:21
*/
type ResultRowInfo struct {
	RowsAffected int64 // 影响行数
	LastInsertId int64 // 最高行ID
}

func (itf *RelationSqlStrut) ExecCtx(ctx context.Context, query string, args ...interface{}) (resultRow ResultRowInfo, err error) {
	result, err := itf.ExecContext(ctx, query, args)
	if err != nil {
		return ResultRowInfo{}, nil
	}
	if rowAf, ok := result.LastInsertId(); ok != nil {
		return ResultRowInfo{}, ok
	} else {
		resultRow.RowsAffected = rowAf
	}
	if rowLd, ok := result.LastInsertId(); ok != nil {
		return ResultRowInfo{}, ok
	} else {
		resultRow.LastInsertId = rowLd
	}
	return
}

/*
	事务相关 TODO 2020年11月18日20:39:24
*/
type SqlFunc func(tx *sql.Tx, wg sync.WaitGroup)

const (
	Query = iota
	QueryRow
	Exec
)

func (itf *RelationSqlStrut) Affair(ctx context.Context, opt *sql.TxOptions, f []SqlFunc) (err error) {
	var (
		wg sync.WaitGroup
	)
	tx, err := itf.DB.BeginTx(ctx, opt)

	for i := 0; i < len(f); i++ {
		wg.Add(1)
		go f[i](tx, wg)
	}

	wg.Wait()

	err = tx.Commit()

	if err != nil {
		err = tx.Rollback()
		return
	}
	return
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
		if err != nil {
			log.Fatal(err.Error())
			return nil, err
		}
		err = conn.Err()
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
		if opt.Port == "" {
			opt.Port = "5432"
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
