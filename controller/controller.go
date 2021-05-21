package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	twicli "github.com/dghubble/go-twitter/twitter"
	"github.com/garyburd/go-oauth/oauth"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/okoshiyoshinori/twigolf-server/config"
	"github.com/okoshiyoshinori/twigolf-server/db"
	"github.com/okoshiyoshinori/twigolf-server/excel"
	"github.com/okoshiyoshinori/twigolf-server/logger"
	"github.com/okoshiyoshinori/twigolf-server/model"
	"github.com/okoshiyoshinori/twigolf-server/twitter"
	"github.com/okoshiyoshinori/twigolf-server/util"
)

//ログアウト
func Logout(c *gin.Context) {
  session := sessions.Default(c)
  fmt.Println("eeee")
  session.Delete("user_id")
  if err := session.Save(); err != nil {
    c.JSON(http.StatusInternalServerError,gin.H{"cause":"予期せぬエラーが発生しました"})
    return
  }
  c.JSON(http.StatusOK,gin.H{"cause":"ログアウトしました"})
}

//組み合わせファイル転送
func GetCombinationExcel(c *gin.Context) {
  cid := c.Param("cid")
  d := db.GetConn()
  competition := model.Competition{}
  participants := make([]*model.Participant,0)
  combination := make([]*model.Combination,0)
  if d.Where("id = ?",cid).Preload("User").Find(&competition).RecordNotFound() {
    c.JSON(http.StatusNotFound,gin.H{"cause":"データがありません"})
    return
  }
  if d.Where("competition_id = ?",cid).Where("status = ?",1).Preload("User").Find(&participants).RecordNotFound() {
    c.JSON(http.StatusNotFound,gin.H{"cause":"データがありません"})
    return
  }
  if d.Where("competition_id = ?",cid).Find(&combination).RecordNotFound() {
    c.JSON(http.StatusNotFound,gin.H{"cause":"データがありません"})
    return
  }
  //テンプレート読み込み
  const activeSheet string = "Sheet1"
  f,err := excelize.OpenFile("./templates/template.xlsx")
  if err != nil {
    c.JSON(http.StatusInternalServerError,gin.H{"cause":"予期せぬエラーが発生しました"})
    return
  }
 
  ex := excel.NewExcel(f,participants,&competition,combination)
  ex.Make(activeSheet)
  c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+"Workbook.xlsx")
	c.Header("Content-Transfer-Encoding", "binary")
  _ = ex.Fd.Write(c.Writer) 
}

//クラブ情報送信
func GetClubs(c *gin.Context) {
  d := db.GetConn()
  word := c.Query("word") 
  clubs := make([]*model.Club,0)
  if word == "" {
    c.JSON(http.StatusInternalServerError,gin.H{"cause":"不正なデータです"})
    return
  }
  d.Order("id asc").Where("MATCH (name,address) AGAINST (? IN BOOLEAN MODE)",word).Find(&clubs)
  if len(clubs) <= 0 {
    c.JSON(http.StatusNotFound,gin.H{"cause":"該当データがありません"})
    return
  }
  c.JSON(http.StatusOK,&clubs)
}

func GetSession(c *gin.Context) {
  session := sessions.Default(c)
  id :=  session.Get("user_id") //sessionから取得
  user := &model.User{}
  d := db.GetConn()
  if d.Select("id,screen_name,name,real_name,real_name_kana,sex,birthday,avatar,description").Where("id = ?",id).Find(user).RecordNotFound() {
    c.JSON(http.StatusNotFound,gin.H{"cause":"該当ユーザーは存在しません"})
    return
  }
  c.JSON(http.StatusOK,user)
}


//ユーザー情報
func GetUser(c *gin.Context) {
  id := c.Param("screen_name")
  user := &model.User{}
  d := db.GetConn()
  if d.Select("id,name,screen_name,avatar,description").Where("screen_name = ?",id).Find(user).RecordNotFound() {
    c.JSON(http.StatusNotFound,gin.H{"cause":"該当ユーザーは存在しません"})
    return
  }
  c.JSON(http.StatusOK,user)
}

func GetUserCompetitions(c *gin.Context) {
  today := time.Now().Local().Format("2006-01-02 15:04:05")
  p := c.Query("p")
  sort := c.Query("sort")
  id := c.Param("screen_name")
  offset,limit,err := util.LimitNum(p)
  if err != nil {
    c.JSON(http.StatusInternalServerError,gin.H{"cause":"不正なデータです"})
    return
  }
  if (sort != "1") && (sort !="2") && (sort != "3") && (sort != "4") {
    c.JSON(http.StatusInternalServerError,gin.H{"cause":"不正なデータです"})
    return
  }
  var sortstr string
  var wherestr string
  if sort == "1" {
    sortstr = "competitions.event_day asc"
    wherestr = "competitions.event_day >= ?"
  } else if sort == "2" {
    sortstr = "competitions.event_day desc"
    wherestr = "competitions.event_day < ?"
  } else if sort == "3" {
    sortstr = "competitions.update_at asc"
    wherestr = "competitions.event_day IS ?"
  } else if sort == "4" {
    sortstr = "competitions.update_at desc"
  }
  d := db.GetConn()
  var state *gorm.DB
  compe := make([]model.Competition,0,limit)
  var num uint

  if sort != "4" {
    state = d.Joins("left join participants on competitions.id = participants.competition_id left join users on users.id = participants.user_id")
  } else {
    state = d.Joins("left join users on competitions.user_id = users.id")
  }

  if sort == "1" || sort == "2" {
    state = state.Where(wherestr,today)
  } else if sort == "3" {
    state = state.Where(wherestr,gorm.Expr("NULL"))
  } 
  state.Where("users.screen_name = ?",id).Model(&model.Competition{}).Count(&num)
  state.Where("users.screen_name = ?",id).Preload("User",func(db *gorm.DB)*gorm.DB{
   return db.Select([]string{"id","screen_name","name","avatar","description"})
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
    return db.Select([]string{"id","screen_name","name","avatar","description"})
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
    return db.Select([]string{"id","screen_name","name","avatar","description"})
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
  joinq := d.Joins("left join competitions on competitions.id = comments.competition_id")
  joinq.Where("comments.competition_id = ?",cid).Preload("User",func(db *gorm.DB)*gorm.DB {
    return db.Select([]string{"id","screen_name","name","avatar","description"})
  }).Order("update_at desc").Find(&comments)
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
  joinq := d.Joins("left join competitions on competitions.id = participants.competition_id")
  joinq.Where("participants.competition_id = ?",cid).Preload("User",func(db *gorm.DB)*gorm.DB {
    return db.Select([]string{"id","screen_name","name","avatar"})
  }).Find(&participants)
 if len(participants) <= 0 {
   c.JSON(http.StatusNotFound,gin.H{"cause":"参加者がいません"})
   return
 }
 c.JSON(http.StatusOK,&participants)
}

func GetPaticipantsWithRealName(c *gin.Context) {
  //test
  session := sessions.Default(c)
  uid := session.Get("user_id") 
  cid := c.Param("cid")
  d := db.GetConn()
  participants := make([]*model.Participant,0)
  joinq := d.Joins("left join competitions on competitions.id = participants.competition_id")
  joinq.Where("competitions.id = ?",cid).Where("competitions.user_id = ?",uid).Preload("User",func(db *gorm.DB)*gorm.DB {
    return db.Select([]string{"id","screen_name","name","avatar","real_name","real_name_kana,sex,birthday"})
  }).Find(&participants)
 if len(participants) <= 0 {
   c.JSON(http.StatusNotFound,gin.H{"cause":"参加者がいません"})
   return
 }
 c.JSON(http.StatusOK,&participants)
}

func GetCombination(c *gin.Context) {
  cid := c.Param("cid")
  d := db.GetConn()
  //compe
  compe := model.Competition{}
  if d.Where("id = ?",cid).Select("combination_open").Find(&compe).RecordNotFound() {
    c.JSON(http.StatusNotFound,gin.H{"cause":"データがありません"})
    return
  }
  combinations := make([]model.Combination,0)
  d.Where("competition_id = ?",cid).Order("start_time asc").Find(&combinations)
  if len(combinations) <= 0 {
    c.JSON(http.StatusNotFound,gin.H{"cause":"まだデータがありません"})
    return
  }
  c.JSON(http.StatusOK,gin.H{"combination_open":compe.CombinationOpen,"payload":combinations})
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
           Where("MATCH (competitions.title,competitions.contents,competitions.place_text,competitions.keyword) AGAINST (? IN BOOLEAN MODE)",
           strings.Join(keywords," ")).Or("MATCH (clubs.name,clubs.address) AGAINST (? IN BOOLEAN MODE)",strings.Join(keywords," ")).
           Where(status,datestr)
  }
  state.Model(&model.Competition{}).Count(&count)
  state.Preload("User",func(db *gorm.DB)*gorm.DB {
              return db.Select([]string{"id","name","screen_name","avatar"})
            }).Order("competitions.event_day asc").Limit(limit).Offset(offset).Find(&compe)
  if count <= 0 {
    c.JSON(http.StatusNotFound,gin.H{"cause":"該当データがありません"})
    return
  }
  c.JSON(http.StatusOK,gin.H{"payload":&compe,"allNumber":count})
}

func PostUserBasicInfo(c *gin.Context) {
  form := model.UserBasicInfoForm{}
  if err:= c.Bind(&form); err != nil {
    c.JSON(http.StatusBadRequest,gin.H{"cause":"データに不備があります"})
    return
  }
  d := db.GetConn()
  user := model.User{}
  if d.Where("id = ?",form.ID).Find(&user).RecordNotFound() {
    c.JSON(http.StatusNotFound,gin.H{"cause":"該当のユーザーは存在しません"})
    return
  }
  user.RealName = form.RealName
  user.RealNameKana = form.RealNameKana
  user.Sex = form.Sex
  user.Birthday = form.BirthDay
  user.UpdateAt = time.Now()

  result := d.Save(&user)
  if result.Error != nil {
    logger.Error.Println(result.Error.Error())
    c.JSON(http.StatusInternalServerError,gin.H{"cause":"予期せぬエラーが発生しました"})
    return
  }
  c.JSON(http.StatusOK,gin.H{"cause":"正常に完了しました"})
}

//create competitions
func PostCompetiton(c *gin.Context) {
  //test
  session := sessions.Default(c)
  uid := session.Get("user_id") 
  form := model.CompetitionForm{}
  if err := c.Bind(&form); err != nil {
    c.JSON(http.StatusBadRequest,gin.H{"cause":"データに不備があります"})
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
    DeletedAt: nil, 
  }
  d := db.GetConn()
  if d.Where("id = ?",compe.ID).Where("user_id = ?",uid).Find(&model.Competition{}).RecordNotFound() {
    if compe.ID != 0 {
      c.JSON(http.StatusBadRequest,gin.H{"cause":"不正なアクセスがありました"})
      return
    } 
    result := d.Create(&compe)
    if result.Error != nil {
      logger.Error.Println(result.Error.Error())
      c.JSON(http.StatusInternalServerError,gin.H{"cause":"予期せぬエラーが発生しました"})
      return
    }
    participants := model.Participant {
        CompetitionID:compe.ID,
        UserID: compe.UserID, 
        Status:1,
        UpdateAt:time.Now(),
      }
    result = d.Create(&participants)
    if result.Error != nil {
      logger.Error.Println(result.Error.Error())
      c.JSON(http.StatusInternalServerError,gin.H{"cause":"予期せぬエラーが発生しました"})
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
  if form.Twitter {
    url := fmt.Sprintf("%s/events/%d",config.GetConfig().Apserver.Host,compe.ID)
    title := compe.Title
    eventDay := "未定"
    if compe.EventDay != nil {
      eventDay = util.DateToString(compe.EventDay)
    }
    d := db.GetConn()
    user := model.User{}
   //ユーザーツイート
   if !d.Where("id = ?",uid).Find(&user).RecordNotFound() {
     token := twitter.Token{
       AccessToken: user.Token,
       AccessSecretToken: user.Secret,
     }
     client := twitter.NewClient(token)
     tweet := fmt.Sprintf("[twigolf]\n%s\n%s\n%s",eventDay,title,url) 
     client.Statuses.Update(tweet,nil)
   }
   //公式ツイート
   token := twitter.Token{
     AccessToken: config.GetConfig().Apserver.TwitterToken,
     AccessSecretToken:config.GetConfig().Apserver.TwitterSecret,
   }
   client := twitter.NewClient(token)
   tweet := fmt.Sprintf("%s\n%s\n%s",eventDay,title,url) 
   client.Statuses.Update(tweet,nil)
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
  d := db.GetConn()
  //チェック
  compe := model.Competition{}
  d.Where("id = ?",form.CompetitionID).Select("capacity").Find(&compe)
  limit := compe.Capacity
  var nowParticipant int
  d.Model(model.Participant{}).Where("competition_id = ?",form.CompetitionID).Where("status = ?",1).Count(&nowParticipant)
  if limit != nil {
    if form.Status == 1 {
      if nowParticipant >= *limit {
        c.JSON(http.StatusBadRequest,gin.H{"cause":"定員を超過しています。"})
        return
      }
    }
  } 

  participant := model.Participant{
    ID: form.ID,
    CompetitionID: form.CompetitionID,
    UserID: form.UserID,
    Status: form.Status,
    UpdateAt: time.Now(),
  }
  first := model.Participant{}
  if d.Where("id = ?",participant.ID).Find(&first).RecordNotFound() {
    result := d.Create(&participant)
    if result.Error != nil {
      logger.Error.Println(result.Error.Error())
      c.JSON(http.StatusInternalServerError,gin.H{"cause":"データ作成に失敗しました"})
      return
    } 
  } else {
    first.Status = participant.Status
    result := d.Save(&first)
    if result.Error != nil {
      logger.Error.Println(result.Error.Error())
      c.JSON(http.StatusInternalServerError,gin.H{"cause":"データ更新に失敗しました"})
      return
    }
  }
  var send = make([]*model.Participant,0)
  d.Where("competition_id = ?",participant.CompetitionID).Preload("User",func(db *gorm.DB)*gorm.DB {
              return db.Select([]string{"id","name","screen_name","avatar"})
            }).Find(&send)
  c.JSON(http.StatusOK,gin.H{"payload":send,"cause":"正常に完了しました"})
}

//delete competition
func DeleteCompetiton(c *gin.Context) {
  //test
  session := sessions.Default(c)
  uid := session.Get("user_id") 
  idstr := c.Param("id")
  u_id,err := strconv.ParseUint(idstr,10,32 << (^uint(0) >> 63)) 
  if err != nil {
    c.JSON(http.StatusBadRequest,gin.H{"cause":"データに不備があります"})
    return
  }
  comp := model.Competition {
    ID:uint(u_id),
    UserID:uid.(uint),
  }
  d := db.GetConn()
  if d.Where("id = ?",u_id).Where("user_id = ?",uid).Find(&model.Competition{}).RecordNotFound() {
    c.JSON(http.StatusNotFound,gin.H{"cause":"データがありません"})
    return
  }
  result := d.Delete(&comp)
  if result.Error != nil {
    logger.Error.Println(result.Error.Error())
    c.JSON(http.StatusInternalServerError,gin.H{"cause":"削除に失敗しました"})
    return
  }
  c.JSON(http.StatusOK,gin.H{"cause":"正常に処理が完了しました"})
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
    //  ID: form.ID,
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
  com := make([]*model.Comment,0)
  d.Where("competition_id = ?",comment.CompetitionID).Preload("User",func(db *gorm.DB)*gorm.DB {
              return db.Select([]string{"id","name","screen_name","avatar"})
            }).Order("update_at desc").Find(&com)
  c.JSON(http.StatusOK,gin.H{"payload":com,"cause":"正常に完了しました"})
}

//post bundleparticipant
func BundleParticipant(c *gin.Context) {
  form := model.BundleParticipants{}
  if err := c.Bind(&form); err != nil {
    c.JSON(http.StatusBadRequest,gin.H{"cause":"データに不備があります"})
    return
  }
  //参加人数取得
  d := db.GetConn()
  compe := model.Competition{}
  if d.Where("id = ?",form.Transaction[0].CompetitionID).Select("capacity").Find(&compe).RecordNotFound() {
    c.JSON(http.StatusBadRequest,gin.H{"cause":"データに不備があります"})
    return
  }
  //参加人数
  limit := compe.Capacity
  for _,v := range form.Transaction {
    if v.Status == 4 {
      data := model.Participant{
        ID:v.ID,
      }
      result := d.Delete(&data)
      if result.Error != nil {
        c.JSON(http.StatusBadRequest,gin.H{"cause":"データの削除に失敗しました"})
        logger.Error.Println(result.Error.Error())
        return
      } 
    } else {
      data := model.Participant{
        ID:v.ID,
        CompetitionID:v.CompetitionID,
        UserID:v.UserID,
        Status:v.Status,
        UpdateAt:time.Now(),
      }
      result := d.Save(&data)
      if result.Error != nil {
        c.JSON(http.StatusBadRequest,gin.H{"cause":"データの更新に失敗しました"})
        logger.Error.Println(result.Error.Error())
        return
      }
    }
  }
  //参加人数チェック
  var nowParticipant int
  var mes string
  d.Model(model.Participant{}).Where("competition_id = ?",form.Transaction[0].CompetitionID).Where("status = ?",1).Count(&nowParticipant)
  if limit != nil {
    if nowParticipant > *limit {
      mes ="正常に処理は完了しましたが、定員を超過しています。"
    } else {
      mes ="正常に処理は完了しました"
    }
  } else {
    mes ="正常に処理は完了しました"
  }
  pa := make([]*model.Participant,0)
  d.Where("competition_id = ?",form.Transaction[0].CompetitionID).Preload("User",func(db *gorm.DB)*gorm.DB {
    return db.Select([]string{"id","screen_name","name","avatar"})
  }).Find(&pa)
  c.JSON(http.StatusOK,gin.H{"payload":pa,"cause":mes})
}

//post combination
func PostCombination(c *gin.Context) {
  //test
  session := sessions.Default(c)
  uid := session.Get("user_id") 
  cid := c.Param("cid")
  d := db.GetConn()
  //check
  compe := model.Competition{}
  if d.Where("id = ?",cid).Where("user_id = ?",uid).Find(&compe).RecordNotFound() {
    c.JSON(http.StatusBadRequest,gin.H{"cause":"データに不備があります"})
    return
  }
  form := model.BundleCombination{}
  if err := c.Bind(&form); err != nil {
    c.JSON(http.StatusBadRequest,gin.H{"cause":"データに不備があります"})
    return
  }
  //組み合わせ公開更新
  {
    compe.CombinationOpen = form.Open
    compe.UpdateAt = time.Now()
    result := d.Save(&compe)
    if result.Error != nil {
      logger.Error.Println(result.Error.Error())
      c.JSON(http.StatusInternalServerError,gin.H{"cause":"予期せぬエラーが発生しました"})
      return
    }
  }
  //削除対象検索
  all := make([]*model.Combination,0)
  var delete_data = make([]uint,0)
  var flag bool
  d.Select("id").Where("competition_id = ?",cid).Find(&all)
  for _,val1 := range all {
    flag = false
    for _,val2 := range form.Transaction {
      if val1.ID == val2.ID {
        flag = true
      }
    }
    if flag == false {
      delete_data = append(delete_data,val1.ID)
    }
  }
  //削除
  if len(delete_data) > 0 {
    result := d.Debug().Delete(&model.Combination{},delete_data)
    if result.Error != nil {
      logger.Error.Println(result.Error.Error())
      c.JSON(http.StatusInternalServerError,gin.H{"cause":"データの削除に失敗しました"})
      return
    } 
  }
  for _,val := range form.Transaction {
   if d.Where("id = ?",val.ID).Find(&model.Combination{}).RecordNotFound() {
     //新規作成
     data := model.Combination{
       CompetitionID:val.CompetitionID,
       StartTime: val.StartTime,
       StartInOut: val.StartInOut,
       Member1: val.Member1,
       Member2: val.Member2,
       Member3: val.Member3,
       Member4: val.Member4,
       UpdateAt: time.Now(),
       DeletedAt: nil,
     }
     result := d.Create(&data)
     if result.Error != nil {
       c.JSON(http.StatusInternalServerError,gin.H{"cause":"予期せぬエラーが発生しました"})
       logger.Error.Println(result.Error.Error())
       return
     }
   } else {
     //更新
     data := model.Combination{
       ID:val.ID,
       CompetitionID:val.CompetitionID,
       StartTime:val.StartTime,
       StartInOut:val.StartInOut,
       Member1: val.Member1,
       Member2: val.Member2,
       Member3: val.Member3,
       Member4: val.Member4,
       UpdateAt: time.Now(),
       DeletedAt:nil,
     }
     result := d.Save(&data)
     if result.Error != nil {
       c.JSON(http.StatusInternalServerError,gin.H{"cause":"予期せぬエラーが発生しました"})
       logger.Error.Println(result.Error.Error())
       return
     }
   }
 }
 resultdata := make([]*model.Combination,0)
 d.Where("competition_id = ?",cid).Order("start_time asc").Find(&resultdata)
 c.JSON(http.StatusOK,gin.H{"combination_open":compe.CombinationOpen,"payload":resultdata,"cause":"正常に終了しました"})
}


func PostDm(c *gin.Context) {
  form := model.PostDm{}
  var sendErr []string 
  uid := sessions.Default(c).Get("user_id")
  if err := c.Bind(&form); err != nil {
    c.JSON(http.StatusBadRequest,gin.H{"cause":"データに不備があります"})
    logger.Error.Println(err.Error())
    return
  }
  d := db.GetConn()
  user := model.User{}
  d.Where("id = ?",uid).Find(&user)
  token := twitter.Token{
    AccessToken:user.Token,
    AccessSecretToken:user.Secret,
  }
  client := twitter.NewClient(token)
  for _,v := range form.Users {
    info,_,err1 := client.Users.Show(&twicli.UserShowParams{
    ScreenName: v,
    })
    if err1 != nil {
      sendErr = append(sendErr,v + "への送信に失敗しました [原因] " + err1.Error())
      continue
    }
    _,_,err2 := client.DirectMessages.EventsNew(twitter.NewDmEvent(info.IDStr,form.Message))
    if err2 != nil {
      sendErr = append(sendErr, v + "への送信に失敗しました [原因] " + err2.Error())
    }
  }
  if len(sendErr) > 0 {
    var temp string = ""
    for _,v := range sendErr {
      temp += v + "\n"
    }
    c.JSON(http.StatusInternalServerError,gin.H{"cause":temp})
    return
  }
  c.JSON(http.StatusOK,gin.H{"cause":"正常に処理が完了しました"})
}

func GetTwitterAuthUrl(c *gin.Context) {
  session := sessions.Default(c)
  client := twitter.GetConnect()
  rt,err := client.RequestTemporaryCredentials(nil,config.GetConfig().Apserver.TwitterCallBackUrl,nil)
  if err != nil {
    c.JSON(http.StatusBadRequest,gin.H{"cause":err.Error()})
    return
  }
  session.Set("request_token",rt.Token)
  session.Set("request_token_secret",rt.Secret)
  if err := session.Save(); err != nil {
    c.JSON(http.StatusInternalServerError,gin.H{"cause":"予期せぬエラーが発生しました"})
    logger.Error.Println(err.Error())
    return
  }
  url := client.AuthorizationURL(rt,nil)
  c.JSON(http.StatusOK,url)
}

func TwitterCallBack(c *gin.Context) {
  session := sessions.Default(c)
  rq := model.TwitterCallBackReq{}
  rq.Token = c.Query("oauth_token")
  rq.Verifier = c.Query("oauth_verifier")
  rt := session.Get("request_token")
  rts := session.Get("request_token_secret")
  at,err := twitter.GetAccessToken(
                &oauth.Credentials{
                  Token:rt.(string),
                  Secret:rts.(string),
                 },
                 rq.Verifier,
               )
  if err != nil {
    c.JSON(http.StatusInternalServerError,gin.H{"cause":err.Error()})
    return
  }
  account := model.TwitterAccount{}
  if err := twitter.GetAccount(at,&account); err != nil {
    c.JSON(http.StatusInternalServerError,gin.H{"cause":err.Error()})
    return
  }

  //DB
  d := db.GetConn()
  user := model.User{}
  if d.Where("sns_id = ?",account.ID).Find(&user).RecordNotFound() {
    //insert table
    user.SnsId = account.ID
    user.ScreenName = account.ScreenName
    user.Name = account.Name
    user.Avatar = account.Avatar
    user.Description = account.Description
    user.LoginType = "TWITTER"
    user.Token = at.Token
    user.Secret = at.Secret
    user.UpdateAt = time.Now()
    result := d.Create(&user)
    if result.Error != nil {
      c.JSON(http.StatusInternalServerError,gin.H{"cause":"予期せぬエラーが発生しました"})
      return
    }
  } else {
    user.SnsId = account.ID
    user.ScreenName = account.ScreenName
    user.Name = account.Name
    user.Avatar = account.Avatar
    user.Description = account.Description
    user.Token = at.Token
    user.Secret = at.Secret
    user.UpdateAt = time.Now()
    result := d.Save(&user)
    if result.Error != nil {
      c.JSON(http.StatusInternalServerError,gin.H{"cause":"予期せぬエラーが発生しました"})
      return
    }
  }
  result := model.User{}
  d.Where("screen_name = ?",account.ScreenName).Select("id").Find(&result)
  session.Set("user_id",result.ID)
  if err := session.Save(); err != nil {
    c.JSON(http.StatusInternalServerError,gin.H{"cause":"予期せぬエラーが発生しました"})
    return
  }
  to := config.GetConfig().Apserver.RedirectUrl + "/" + account.ScreenName
  c.Redirect(http.StatusPermanentRedirect,to)
}

func SendOgpHtml(c *gin.Context) {
  cidstr := c.Param("cid")
  cid,err := strconv.Atoi(cidstr) 
  if err != nil {
    c.AbortWithStatus(500)
    return
  }
  ogp := twitter.MakeOgpInfo(cid)
  if ogp == nil  {
    c.AbortWithStatus(500)
    return
  }

  c.HTML(http.StatusOK,"index.tmpl",gin.H{
    "title": ogp.OgTitle,
    "description": ogp.OgDescription,
    "og_title": ogp.OgTitle,
    "og_description": ogp.OgDescription,
    "og_type": ogp.OgType,
    "og_url": ogp.OgUrl,
    "og_image": ogp.OgImage,
    "twitter_card": ogp.Twitter.Card,
    "twitter_site": ogp.Twitter.Site,
    "twitter_creater": ogp.Twitter.Creater,
  })
}
