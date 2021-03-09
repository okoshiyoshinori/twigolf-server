package model

import (
	"time"
)

//login type
const (
  Twitter = "TWITTER"
  Other = "OTHER"
)

//competitons accept status type
const (
  OPEN = 1
  CLOSE = 2
)


//DB User Table
type User struct {
  ID uint `json:"id" gorm:"primaryKey"`
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
  Update_at time.Time `json:"update_at"`
}

//DB Competition Table
type Competition struct {
  ID uint `json:"id" gorm:"primaryKey"`
  UserID uint `json:"user_id"`
  Status int `json:"status"`
  Title string `json:"title"`
  Capacity int `json:"capacity"`
  Contents string `json:"contents"`
  ClubID uint `json:"club_id"`
  PlaceText *string `json:"place_text"`
  EventDay *time.Time `json:"event_day"`
  EventDeadline *time.Time `json:"event_deadline"`
  User *User `json:"user"`
  Keyword []*Keyword `json:"keywords"`
  Club *Club `json:"club"`
  UpdateAt time.Time `json:"update_at"`
}

//DB participants table
type Participant struct {
  ID uint `json:"id" gorm:"primaryKey"`
  CompetitionID uint `json:"competition_id"`
  UserID uint `json:"user_id"`
  Status int `json:"status"`
  User *User `json:"user"`
  UpdateAt time.Time `json:"update_at"`
}

//DB kayword table
type Keyword struct {
  ID uint `json:"id" gorm:"primaryKey"`
  CompetitionID uint `json:"competition_id"` 
  Word string `json:"word"`
  UpdateAt time.Time `json:"update_at"`
}
//DB Comment table
type Comment struct {
  ID uint `json:"id" gorm:"primaryKey"`
  CompetitionID uint `json:"competition_id"`
  UserID int `json:"user_id"`
  Message string `json:"message"`
  User *User `json:"user"`
  UpdateAt time.Time `json:"update_at"`
}

//DB Prefecture table
type Prefecture struct {
  ID uint `json:"id" gorm:"primaryKey"`
  Name string `json:"name"`
}

//DB Club table
type Club struct {
  ID uint `json:"club_id" gorm:"primaryKey"`
  Class int `json:"class"`
  Name string `json:"name"`
  Address string `json:"address"`
  Other string `json:"other"`
}

//form userdata
type UserForm struct {
  Sns_id string `json:"sns_id" binding:"required"`
  ScreenName string `json:"screen_name" binding:"required"`
  Avatar string `json:"avatar"`
  Description string `json:"description"`
  Logintype string `json:"login_type" binding:"required"`
  Mail string `json:"mail" binding:"require,email"`
}

//form competitons 
type CompetitionForm struct {
  ID uint `json:"id"` 
  UserID uint `json:"user_id" binding:"required"`
  Status int `json:"status" binding:"required"`
  Title string `json:"title" binding:"required"`
  Capacity int `json:"capacity" binding:"required"`
  Contents string `json:"contents" binding:"required"`
  ClubID uint `json:"club_id"`
  PlaceText *string `json:"place_text"`
  EventDay *time.Time `json:"event_day"`
  EventDeadline *time.Time `json:"event_deadline"`
  KeyWord  KeywordForm `json:"keywords"`
}

//form comment
type CommentForm struct {
  CompetitionID string `json:"competition_id" binding:"required"`
  UserID string `json:"user_id" binding:"required"`
  Message string `json:"message" binding:"required"`
}

//form keyword
type KeywordForm struct {
  Words []string `json:"words"`
}

