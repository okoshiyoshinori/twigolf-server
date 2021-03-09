package config

import (
  "github.com/BurntSushi/toml"
)

type conf struct {
  Dbserver dbConfig `toml:"dbserver"`
  Apserver apConfig `toml:"apserver"`
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
  Port string `toml:"port"`
  CookieKey string `toml:"cookie_key"`
  Salt string `toml:"salt"`
}

var systemConfig *conf = readConfig() 

func readConfig() *conf {
  config := conf{}
  _,err := toml.DecodeFile("Config.toml",&config)
  if err != nil {
    panic(err.Error())
  }
  return &config
}

func GetConfig() *conf {
  return systemConfig
}

