package config

import (
  "github.com/BurntSushi/toml"
)

type conf struct {
  Dbserver *dbConfig
  Apserver *apConfig
}

type dbConfig struct {
  Dbms string `toml:"dbms"`
  User string `toml:"user"`
  Pass string `toml:"pass"`
  Protocol string `toml:"protocol"`
  Dbname string `toml:"dbname"`
  Charset string `toml:"charset"`
  Parsetime string `toml:"parsetime"`
}

type apConfig struct {
  NumPerPage int `toml:"numperpage"`
  Origin string `toml:"origin"`
  CookieKey string `toml:"cookie_key"`
}

var systemConfig *conf = readConfig() 

func readConfig() *conf {
  config := conf{}
  _,err := toml.Decode("Config.toml",&config)
  if err != nil {
    panic(err.Error())
  }
  return &config
}

func GetConfig() *conf {
  return systemConfig
}

