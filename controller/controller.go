package controller

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/okoshiyoshinori/twigolf-server/db"
	"github.com/okoshiyoshinori/twigolf-server/logger"
	"github.com/okoshiyoshinori/twigolf-server/model"
	"github.com/okoshiyoshinori/twigolf-server/util"
)

//ログアウト
func Logout(c *gin.Context) {
  //session := sessions.Default(c) 
  c.JSON(http.StatusOK,gin.H{"cause":"ログアウトしました"})
}


//クラブ情報送信
func GetClubs(c *gin.Context) {
  d := db.GetConn()
  word := c.Query("word") 
  class:= c.Query("class")
  clubs := make([]*model.Club,0)
  if class == "" {
    class = "1"
  }
  if word == "" {
    c.JSON(http.StatusInternalServerError,gin.H{"cause":"不正なデータです"})
    return
  }
  d.Order("id asc").Where("class = ?",class).Where("MATCH (name,address) AGAINST (? IN NATURAL LANGUAGE MODE)",word).Find(&clubs)
  if len(clubs) <= 0 {
    c.JSON(http.StatusNotFound,gin.H{"cause":"該当データがありません"})
    return
  }
  c.JSON(http.StatusOK,&clubs)
}

func GetSession(c *gin.Context) {
  id := 1 //sessionから取得
  user := &model.User{}
  d := db.GetConn()
  if d.Select("id,sns_id,screen_name,avatar,description").Where("id = ?",id).Find(user).RecordNotFound() {
    c.JSON(http.StatusNotFound,gin.H{"cause":"該当ユーザーは存在しません"})
    return
  }
  c.JSON(http.StatusOK,user)
}


//ユーザー情報
func GetUser(c *gin.Context) {
  id := c.Param("snsid")
  user := &model.User{}
  d := db.GetConn()
  if d.Select("id,sns_id,screen_name,avatar,description").Where("sns_id = ?",id).Find(user).RecordNotFound() {
    c.JSON(http.StatusNotFound,gin.H{"cause":"該当ユーザーは存在しません"})
    return
  }
  c.JSON(http.StatusOK,user)
}

func GetUserCompetitions(c *gin.Context) {
  today := time.Now().Local().Format("2006-01-02 15:04:05")
  p := c.Query("p")
  sort := c.Query("sort")
  id := c.Param("snsid")
  offset,limit,err := util.LimitNum(p)
  if err != nil {
    c.JSON(http.StatusInternalServerError,gin.H{"cause":"不正なデータです"})
    return
  }
  if (sort != "1") && (sort !="2") && (sort != "3") {
    c.JSON(http.StatusInternalServerError,gin.H{"cause":"不正なデータです"})
    return
  }
  var sortstr string
  var wherestr string
  if sort == "1" {
    sortstr = "event_day asc"
    wherestr = "event_day >= ?"
  } else if sort == "2" {
    sortstr = "event_day desc"
    wherestr = "event_day < ?"
  } else if sort == "3" {
    sortstr = "update_at asc"
    wherestr = "event_day IS ?"
  }
  d := db.GetConn()
  var state *gorm.DB
  compe := make([]model.Competition,0,limit)
  var num uint

  state = d.Joins("left join participants on competitions.id = participants.competition_id left join users on users.id = participants.user_id")

  if sort == "1" || sort == "2" {
    state = state.Where(wherestr,today)
  } else {
    state = state.Where(wherestr,gorm.Expr("NULL"))
  }
  state.Where("users.sns_id = ?",id).Model(&model.Competition{}).Count(&num)
state.Where("users.sns_id = ?",id).Preload("User",func(db *gorm.DB)*gorm.DB{
  return db.Select([]string{"id","screen_name","sns_id","avatar","description"})
}).Preload("Club").Order(sortstr).Limit(limit).Offset(offset).Find(&compe)
  if len(compe) <= 0 {
    c.JSON(http.StatusNotFound,gin.H{"cause":"該当のデータはありません"})
    return
  }
  c.JSON(http.StatusOK,gin.H{"allNumber":num,"payload":compe})
}

//コンペ一覧
func GetCompetition(c *gin.Context) {
  today := time.Now().Local().Format("2006-01-02 15:04:05")
  //ページ
  p := c.Query("p")
  offset,limit,err := util.LimitNum(p)
  if err != nil {
    c.JSON(http.StatusInternalServerError,gin.H{"cause":"不正なデータです"})
    return
  }
  //ソート
  var sortstr string
  var wherestr string
  mode := c.Query("mode")
  if mode == "1" { //最新
    sortstr = "update_at desc"
    wherestr = "update_at <= ?" 
  } else if mode == "2" { //開催日
    sortstr = "event_day asc"
    wherestr = "event_day >= ?" 
  } else if mode == "3"{ //締め切り
    sortstr = "event_deadline asc"
    wherestr = "event_deadline >= ?"
  } else {
    sortstr = "update_at desc"
    wherestr = "update_at <= ?" 
  }
  d := db.GetConn()
  compes := make([]*model.Competition,0,limit)
  state := d.Where(wherestr,today).Order(sortstr).Preload("Club").Preload("User",func(db *gorm.DB) *gorm.DB {
    return db.Select([]string{"id","screen_name","sns_id","avatar","description"})
  })
  var count int
  state.Model(&model.Competition{}).Count(&count)
  state.Limit(limit).Offset(offset).Find(&compes)
  if len(compes) <= 0 {
    c.JSON(http.StatusNotFound,gin.H{"cause":"該当データがありません"})
    return
  }
  c.JSON(http.StatusOK,gin.H{"allNumber":count,"payload":compes})
}

//コンペ詳細
func GetCompetitonDetail(c *gin.Context) {
  id := c.Param("id")
  d := db.GetConn()
  compe := model.Competition{}
  if d.Preload("User",func(db *gorm.DB)*gorm.DB{
    return db.Select([]string{"id","screen_name","sns_id","avatar","description"})
  }).Preload("Club").Where("id = ?",id).Find(&compe).RecordNotFound() {
    c.JSON(http.StatusNotFound,gin.H{"cause":"該当データがありません"})
    return
  }
  c.JSON(http.StatusOK,&compe)
}

//コンペ >> コメント
func GetComment(c *gin.Context) {
  cid := c.Param("cid")
  d := db.GetConn()
  comments := make([]*model.Comment,0) 
  d.Where("competition_id = ?",cid).Preload("User",func(db *gorm.DB)*gorm.DB {
    return db.Select([]string{"id","screen_name","sns_id","avatar","description"})
  }).Find(&comments)
  if len(comments) <= 0 {
    c.JSON(http.StatusNotFound,gin.H{"cause":"該当データがありません"})
    return
  }
  c.JSON(http.StatusOK,&comments)
}

//コンペ >> 参加者一覧
func GetPaticipants(c *gin.Context) {
  cid := c.Param("cid")
  d := db.GetConn()
  participants := make([]*model.Participant,0)
  d.Where("competition_id = ?",cid).Preload("User",func(db *gorm.DB)*gorm.DB {
    return db.Select([]string{"id","screen_name","sns_id","avatar"})
  }).Find(&participants)
 if len(participants) <= 0 {
   c.JSON(http.StatusNotFound,gin.H{"cause":"参加者がいません"})
   return
 }
 c.JSON(http.StatusOK,&participants)
}

func SearchCompetition(c *gin.Context) {
  p := c.Query("p") //ページ
  offset,limit,err := util.LimitNum(p)
  if (err != nil) {
    c.JSON(http.StatusInternalServerError,gin.H{"cause":"不正なデータです"})
    return
  }
  q := c.Query("q") //キーワード
  keywords := strings.Fields(q) //半角で区切り
  day := c.Query("date") //日付
  mode := c.Query("mode") // 1: 日付指定 2:未定

  if (mode != "1") && (mode !="2") {
    c.JSON(http.StatusInternalServerError,gin.H{"cause":"不正なデータです"})
    return
  }

  var datestr interface{}
  var status string
  if mode == "1" {
    status = "event_day >= ?"
    datestr = day 
  }else {
    status = "event_day IS ?"
    datestr = gorm.Expr("NULL") 
  }

  var count int
  var state *gorm.DB
  compe := make([]model.Competition,0,limit)
  d := db.GetConn()
  joinq := d.Joins("left join clubs on competitions.club_id = clubs.id")
  if len(keywords) <= 0 {
    state = joinq.Where(status,datestr)
  } else {
    state =joinq.
           Where("MATCH (competitions.title,competitions.contents,competitions.place_text,competitions.keyword) AGAINST (? IN NATURAL LANGUAGE MODE)",
           strings.Join(keywords," ")).Or("MATCH (clubs.name,clubs.address) AGAINST (? IN NATURAL LANGUAGE MODE)",strings.Join(keywords," ")).
           Where(status,datestr)
  }
  state.Model(&model.Competition{}).Count(&count)
  state.Preload("User",func(db *gorm.DB)*gorm.DB {
              return db.Select([]string{"id","sns_id","screen_name","avatar"})
            }).Order("competitions.event_day asc").Limit(limit).Offset(offset).Find(&compe)
  if count <= 0 {
    c.JSON(http.StatusNotFound,gin.H{"cause":"該当データがありません"})
    return
  }
  c.JSON(http.StatusOK,gin.H{"payload":&compe,"allNumber":count})
}

//create competitions
func PostCompetiton(c *gin.Context) {
  form := model.CompetitionForm{}
  if err := c.Bind(&form); err != nil {
    c.JSON(http.StatusBadRequest,gin.H{"cause":err.Error()})
    return
  }
  compe := model.Competition {
    ID:form.ID,
    UserID: form.UserID,
    Status: model.OPEN,
    Title: form.Title,
    Capacity: form.Capacity,
    Contents: form.Contents,
    ClubID: form.ClubID,
    PlaceText: form.PlaceText,
    EventDay: form.EventDay,
    EventDeadline: form.EventDeadline,
    Keyword: form.Keyword,
    UpdateAt: time.Now(),
  }
  d := db.GetConn()
  if d.Where("id = ?",compe.ID).Find(&model.Competition{}).RecordNotFound() {
    result := d.Create(&compe)
    if result.Error != nil {
      logger.Error.Println(result.Error.Error())
      c.JSON(http.StatusInternalServerError,gin.H{"cause":result.Error.Error()})
      return
    }
  } else {
    result := d.Save(&compe)
    if result.Error != nil {
      logger.Error.Println(result.Error.Error())
      c.JSON(http.StatusInternalServerError,gin.H{"cause":result.Error.Error()})
      return
    }
  }
  c.JSON(http.StatusOK,gin.H{"cause":"正常に終了しました"})
}

//post participants 
func PostParticipant(c *gin.Context) {
  var form model.ParticipantForm
  if err := c.Bind(&form); err != nil {
    c.JSON(http.StatusBadRequest,gin.H{"cause":"データに不備があります"})
    return
  }
  participant := model.Participant{
    ID: form.ID,
    CompetitionID: form.CompetitionID,
    UserID: form.UserID,
    Status: form.Status,
    UpdateAt: time.Now(),
  }
  d := db.GetConn()
  if d.Where("id = ?",participant.ID).Find(&model.Participant{}).RecordNotFound() {
    result := d.Create(&participant)
    if result.Error != nil {
      logger.Error.Println(result.Error.Error())
      c.JSON(http.StatusInternalServerError,gin.H{"cause":"データ作成に失敗しました"})
      return
    } 
  } else {
    result := d.Save(&participant)
    if result.Error != nil {
      logger.Error.Println(result.Error.Error())
      c.JSON(http.StatusInternalServerError,gin.H{"cause":"データ更新に失敗しました"})
      return
    }
  }
  c.JSON(http.StatusOK,gin.H{"cause":"正常に処理されました"})
}

//post comment
func PostComment(c *gin.Context) {
  var form model.CommentForm
  if err := c.Bind(&form); err != nil {
    logger.Error.Println(err.Error())
    c.JSON(http.StatusBadRequest,gin.H{"cause":"データに不備があります"})
    return
  }
  comment := model.Comment {
    ID: form.ID,
    UserID:form.UserID,
    CompetitionID:form.CompetitionID,
    Message:form.Message,
    UpdateAt: time.Now(),
  }
  d := db.GetConn()
  result := d.Create(&comment) 
  if result.Error != nil {
    c.JSON(http.StatusInternalServerError,gin.H{"cause":"データの更新に失敗しました"})
    return
  }
  c.JSON(http.StatusOK,gin.H{"cause":"正常に処理が完了しました"})
}

