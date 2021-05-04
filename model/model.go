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
//participant
const (
  PARTICIPANT = 1
  CONCERN = 2
  DECLINE = 3
)


//DB User Table
type User struct {
  ID uint `json:"id" gorm:"primaryKey"`
  Password string `json:"password"`
  Email string `json:"email"`
  SnsId string `json:"sns_id"`
  RealName *string `json:"real_name"`
  RealNameKana *string `json:"real_name_kana"`
  Sex *uint `json:"sex"`
  Birthday *time.Time `json:"birthday"`
  ScreenName string `json:"screen_name"`
  Avatar string `json:"avatar"`
  Description string `json:"description"`
  LoginType string `json:"login_type"`
  Token string `json:"token"`
  Secret string `json:"secret"`
  UpdateAt time.Time `json:"update_at"`
  DeletedAt *time.Time `json:"deleted_at"`
}

//DB Competition Table
type Competition struct {
  ID uint `json:"id" gorm:"primaryKey"`
  UserID uint `json:"user_id"`
  Status int `json:"status"`
  Title string `json:"title"`
  Capacity *int `json:"capacity"`
  Contents string `json:"contents"`
  ClubID *uint `json:"club_id"`
  PlaceText *string `json:"place_text"`
  EventDay *time.Time `json:"event_day"`
  EventDeadline *time.Time `json:"event_deadline"`
  Keyword *string `json:"keyword"`
  CombinationOpen bool `json:"combination_open"`
  User *User `json:"user"`
  Club *Club `json:"club"`
  UpdateAt time.Time `json:"update_at"`
  DeletedAt *time.Time `json:"deleted_at"`
}

//DB participants table
type Participant struct {
  ID uint `json:"id" gorm:"primaryKey"`
  CompetitionID uint `json:"competition_id"`
  UserID uint `json:"user_id"`
  Status int `json:"status"`
  User *User `json:"user"`
  UpdateAt time.Time `json:"update_at"`
  DeletedAt *time.Time `json:"deleted_at"`
}

//DB Comment table
type Comment struct {
  ID uint `json:"id" gorm:"primaryKey"`
  CompetitionID uint `json:"competition_id"`
  UserID uint `json:"user_id"`
  Message string `json:"message"`
  User *User `json:"user"`
  UpdateAt time.Time `json:"update_at"`
  DeletedAt *time.Time `json:"deleted_at"`
}

//DB Club table
type Club struct {
  ID uint `json:"id" gorm:"primaryKey"`
  Class int `json:"class"`
  Name string `json:"name"`
  Address string `json:"address"`
  Other string `json:"other"`
}

//DB Combinations table
type Combination struct {
  ID uint `json:"id" gorm:"primaryKey"`
  CompetitionID uint `json:"competition_id"`
  StartTime time.Time `json:"start_time"`
  StartInOut uint `json:"start_in_out"`
  Member1 *uint `json:"member1"`
  Member2 *uint `json:"member2"`
  Member3 *uint `json:"member3"`
  Member4 *uint `json:"member4"`
  UpdateAt time.Time `json:"update_at"`
  DeletedAt *time.Time `json:"deleted_at"`
}

type BundleCombinationWithOpen struct {
  CombinationOpen bool `json:"combination_open"`
  Combinations []Combination `json:"combinations"`
}

//form combinarion 
type CombinationForm struct {
  ID uint `json:"id" binding:"required"` 
  CompetitionID uint `json:"competition_id" binding:"required"`
  StartTime time.Time `json:"start_time" binding:"required"`
  StartInOut uint `json:"start_in_out" binding:"required"`
  Member1 *uint `json:"member1"`
  Member2 *uint `json:"member2"`
  Member3 *uint `json:"member3"`
  Member4 *uint `json:"member4"`
}

//bundleCombinationForm
type BundleCombination struct {
  Open bool `json:"open"`
  Transaction []CombinationForm `json:"transaction" binding:"required"`
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
  Capacity *int `json:"capacity"`
  Contents string `json:"contents" binding:"required"`
  ClubID *uint `json:"club_id"`
  PlaceText *string `json:"place_text"`
  EventDay *time.Time `json:"event_day"`
  EventDeadline *time.Time `json:"event_deadline"`
  Keyword *string `json:"keyword"`
  twitter bool `json:"twitter"`
}

//form participants
type ParticipantForm struct {
  ID uint `json:"id"`
  CompetitionID uint `json:"competition_id" binding:"required"`
  UserID uint `json:"user_id" binding:"required"`
  Status int `json:"status" binding:"required"`
}

type BundleParticipants struct {
  Transaction []ParticipantForm `json:"transaction" binding:"required"`
}

type CommentForm struct {
 // ID uint `json:"id"`
  CompetitionID uint `json:"competition_id" binding:"required"`
  UserID uint `json:"user_id" binding:"required"`
  Message string `json:"message" binding:"required"`
}

type UserBasicInfoForm struct {
  ID uint `json:"id"`
  RealName *string `json:"real_name" binding:"required"`
  RealNameKana *string `json:"real_name_kana" binding:"required"`
  Sex *uint `json:"sex" binding:"required"`
  BirthDay *time.Time `json:"birthday" binding:"required"`
}
