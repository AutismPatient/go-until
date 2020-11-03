package driver

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type Option struct {
	Name string
	Pass string
	Port int64
	Pex  string
}

func NewMySQL(opt Option) {
	source := fmt.Sprintf("")
	conn, err := sql.Open("mysql", source)
	if err != nil {
		panic(err)
	}
	fmt.Print(conn.Ping())
}
func NewRedis() {

}
func NewMongoDB() {

}
func NewMSSQL() {

}
func NewPostgreSQL() {

}
func NewDocker() {

}
