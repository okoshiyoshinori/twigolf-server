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
  Host string `toml:"host"`
  TwitterAccount string `toml:"twitter_account"`
  TwitterToken string `toml:"twitter_token"`
  TwitterSecret string `toml:"twitter_secret"`
  SiteName string `toml:"site_name"`
  Image string `toml:"image"`
  NumPerPage int `toml:"numperpage"`
  Origin []string `toml:"origin"`
  Port string `toml:"port"`
  CookieKey string `toml:"cookie_key"`
  Salt string `toml:"salt"`
  CsKey string `toml:"cs_key"`
  CsSecretKey string `toml:"cs_secret_key"`
  TwitterCallBackUrl string `toml:"twitter_callback_url"`
  TwitterTempCredentialRequestURI string `toml:"twitter_temporary_credential_requestURI"`
  TwitterResourceOwnerAuthorizationURI string `toml:"twitter_resource_owner_authorizationURI"`
  TwitterTokenRequestURI string `toml:"twitter_token_requestURI"`
  TwitterVerifyCredentialsURI string `toml:"twitter_verify_credentialsURI"`
  RedirectUrl string `toml:"redirect_url"`
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

