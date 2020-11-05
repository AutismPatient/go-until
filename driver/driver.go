package driver

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/garyburd/redigo/redis"
	_ "github.com/go-sql-driver/mysql"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"reflect"
	"strings"
	"time"
)

type BaseConnStrut struct {
	UserName string
	Pass     string
	DB       string
	Addr     string
}

type MysqlStrut struct {
	*sql.DB
}
type MysqlOption struct {
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

*/

func NewMySQL(opt MysqlOption) (itf *MysqlStrut, err error) {
	var (
		source = ""
		conn   *sql.DB
	)
	if opt.ConnectString != "" {
		source = opt.ConnectString
	} else {
		source = fmt.Sprintf("%s:%s@tcp(%s)/%s", opt.UserName, opt.Pass, opt.Addr, opt.DB)
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
func (itf *MysqlStrut) GetSingleByID(db string, rely *interface{}, id string) (err error) {
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
func (itf *MysqlStrut) Status() (status sql.DBStats) {
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
type MongoDBOptions struct {
	*options.ClientOptions
}

/*
	mongodb+srv://<username>:<password>@<cluster-address>/test?w=majority
*/

func NewMongoDB(database, url string, opt MongoDBOptions, collection ...string) (client *MongoDBStrut, err error) {
	var (
		context = context.TODO()
	)
	client.Client, err = mongo.Connect(context,
		opt.SetRetryReads(*opt.RetryReads),
		opt.ApplyURI(url),
		opt.SetAppName(*opt.AppName),
		opt.SetAuth(*opt.Auth),
		opt.SetConnectTimeout(*opt.ConnectTimeout),
		opt.SetSocketTimeout(*opt.SocketTimeout),
		opt.SetDirect(*opt.Direct),
		opt.SetDisableOCSPEndpointCheck(*opt.DisableOCSPEndpointCheck),
		opt.SetCompressors(opt.Compressors),
		opt.SetHosts(opt.Hosts),
		opt.SetMaxPoolSize(*opt.MaxPoolSize),
		opt.SetMaxConnIdleTime(*opt.MaxConnIdleTime),
		opt.SetReplicaSet(*opt.ReplicaSet),
		opt.SetRetryWrites(*opt.RetryWrites),
		opt.SetZlibLevel(*opt.ZlibLevel),
		opt.SetZstdLevel(*opt.ZstdLevel),
		opt.SetTLSConfig(opt.TLSConfig),
		opt.SetWriteConcern(opt.WriteConcern),
		opt.SetReadConcern(opt.ReadConcern),
	)
	if err = client.Client.Ping(context, readpref.Primary()); err != nil {
		return nil, err
	}
	client.Database = client.Client.Database(database)
	if len(collection) > 0 {
		client.Collection = client.Database.Collection(collection[0])
	}
	return
}

func NewMSSQL() {

}
func NewPostgreSQL() {

}
func NewDocker() {

}
