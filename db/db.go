package db

import (
  _ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/okoshiyoshinori/twigolf-server/config"
)

var conn *gorm.DB = newDbConn()

func newDbConn() *gorm.DB {
  conf := config.GetConfig().Dbserver 
  connection := conf.User + ":" + conf.Pass + "@" + conf.Protocol + "/" + conf.Dbname + "?" + conf.Parsetime + "&" + conf.Charset
  conn,err := gorm.Open(conf.Dbms,connection)
  if err != nil {
    panic(err.Error())
  }
  return conn
}

//singleton
func GetConn() *gorm.DB {
  return conn
}
