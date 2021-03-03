package model

import "time"

//DB User Table
type User struct {
  UId int `json:"u_id" gorm:"primary_key"`
  Password string `json:"password"`
  Email string `json:"email"`
  SnsId string `json:"sns_id"`
  ScreenName string `json:"screen_name"`
  Avatar string `json:"avatar"`
  Description string `json:"description"`
  LoginType string `json:"login_type"`
  Active bool `json:"active"`
  Token string `json:"token"`
  Secret string `json:"secret"`
  UpdateAt time.Time `json:"update_at"`
}

//DB Competition Table
type Competition struct {
  CompeId int `json:"compe_id" gorm:"primary_key"`
  Register int `json:"register"`
  Status int `json:"status"`
  Title string `json:"title"`
  Capacity int `json:"capacity"`
  Contents string `json:"contents"`
  Place int `json:"place"`
  PlaceText string `json:"place_text"`
  EventDay time.Time `json:"event_day"`
  EventDeadline time.Time `json:"event_deadline"`
  UpdateAt time.Time `json:"update_at"`
}

//DB participants table
type Participant struct {
  PId int `json:"p_id" gorm:"primary_key"`
  CompeId int `json:"compe_id"`
  UId int `json:"u_id"`
  Status int `json:"status"`
  UpdateAt time.Time `json:"update_at"`
}

type KeyWord struct {
  KId int `json:"k_id" gorm:"primary_key"`
  CompeId int `json:"compe_id"` 
  Word string `json:"word"`
  UpdateAt time.Time `json:"update_at"`
}

type Comment struct {
  MId int `json:"m_id" gorm:"primary_key"`
  CompeId int `json:"compe_id"`
  UId int `json:"u_id"`
  Message string `json:"message"`
  UpdateAt time.Time `json:"update_at"`
}

type Prefecture struct {
  PrefId int `json:"pref_id"`
  Name string `json:"name"`
}

type Clubs struct {
  ClubId int `json:"club_id" gorm:"primary_key"`
  Class int `json:"class"`
  Name string `json:"name"`
  Address string `json:"address"`
  Other string `json:"other"`
}
