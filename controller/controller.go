package controller

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/okoshiyoshinori/twigolf-server/db"
	"github.com/okoshiyoshinori/twigolf-server/model"
	"github.com/okoshiyoshinori/twigolf-server/util"
)

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
  d.Debug().Order("id asc").Where("class = ?",class).Where("MATCH (name,address) AGAINST (? IN NATURAL LANGUAGE MODE)",word).Find(&clubs)
  if len(clubs) <= 0 {
    c.JSON(http.StatusNotFound,gin.H{"cause":"データがありません"})
    return
  }
  c.JSON(http.StatusOK,&clubs)
}

//都道府県データ送信
func GetPref(c *gin.Context) {
  d := db.GetConn()
  prefs := make([]*model.Prefecture,0)
  d.Find(&prefs)
  if len(prefs) <= 0 {
    c.JSON(http.StatusNotFound,gin.H{"cause":"都道府県データがありません"})
    return
  }
  c.JSON(http.StatusOK,&prefs)
}

//ユーザー情報
func GetUser(c *gin.Context) {
  uid := c.Param("uid")
  user := &model.User{}
  d := db.GetConn()
  if d.Where("id =?",uid).Find(user).RecordNotFound() {
    c.JSON(http.StatusNotFound,gin.H{"cause":"該当ユーザーは存在しません"})
    return
  }
  c.JSON(http.StatusOK,user)
}

//コンペ一覧
func GetCompetition(c *gin.Context) {
  today := time.Now().Local().Format("2006-01-02 15:04:05")
  fmt.Println(today)
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
  state := d.Where(wherestr,today).Order(sortstr).Preload("User",func(db *gorm.DB) *gorm.DB {
    return db.Select([]string{"id","screen_name","sns_id","avatar","description"})
  })
  var count int
  state.Find(&model.Competition{}).Count(&count)
  state.Limit(limit).Offset(offset).Find(&compes)
  if len(compes) <= 0 {
    c.JSON(http.StatusNotFound,gin.H{"cause":"データがありません"})
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
  }).Preload("Keyword").Preload("Club").Where("id = ?",id).Find(&compe).RecordNotFound() {
    c.JSON(http.StatusNotFound,gin.H{"cause":"データがありません"})
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
    c.JSON(http.StatusNotFound,gin.H{"cause":"データがありません"})
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

  var datestr string
  var status string
  if mode == "1" {
    status = "event_day >= ?"
    datestr = day 
  }else {
    status = "event_day IS ?"
    datestr = "NULL"
  }

  var count int
  var state *gorm.DB
  compe := make([]model.Competition,0)
  d := db.GetConn()
  joinq := d.Debug().Joins("left join keywords on competitions.id = keywords.competition_id left join clubs on competitions.club_id = clubs.id")
  if len(keywords) <= 0 {
    state = joinq.Where(status,datestr).Preload("User",func(db *gorm.DB)*gorm.DB {
            return db.Select([]string{"id","sns_id","screen_name","avatar"}) 
            })
  } else {
    state = joinq.Where("MATCH (competitions.title,competitions.contents,competitions.place_text) AGAINST (? IN NATURAL LANGUAGE MODE)",
            strings.Join(keywords," ")).Or("MATCH (clubs.name,clubs.address) AGAINST (? IN NATURAL LANGUAGE MODE)",strings.Join(keywords," ")).
            Or("MATCH(keywords.word) AGAINST (? IN NATURAL LANGUAGE MODE)",strings.Join(keywords," ")).Where(status,datestr).
            Preload("User",func(db *gorm.DB)*gorm.DB {
            return db.Select([]string{"id","sns_id","screen_name","avatar"}) 
           })
  }
  state.Find(&model.Competition{}).Count(&count)
  state.Limit(limit).Offset(offset).Find(&compe)
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
    UpdateAt: time.Now(),
  }
  d := db.GetConn()
  if d.Model(&compe).Where("id = ?",compe.ID).Updates(&compe).RowsAffected == 0 {
     d.Create(&compe)
  }
  c.JSON(http.StatusOK,gin.H{"cause":"正常に終了しました"})
}

