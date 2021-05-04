package util

import (
	"fmt"
	"strconv"
	"time"

	"github.com/okoshiyoshinori/twigolf-server/config"
	"github.com/okoshiyoshinori/twigolf-server/model"
	"github.com/speps/go-hashids"
)

var (
	h *hashids.HashID
)

func init() {
	hd := hashids.NewData()
	hd.Salt = config.GetConfig().Apserver.Salt 
	hd.MinLength = 8
	h, _ = hashids.NewWithData(hd)
}


func Encode(s int) string {
	e, _ := h.Encode([]int{s})
	return e
}

func Decode(s string) (int, error) {
	d, err := h.DecodeWithError(s)
	if err != nil {
		return 0, err
	}
	return d[0], nil
}

func LimitNum(page string) (offset int,num int,err error) {
  if page == "" || page == "0" {
    page = "1"
  }
  p,err := strconv.Atoi(page)
  if err != nil {
    return 0,0,err
  }
  off := (p - 1) * config.GetConfig().Apserver.NumPerPage
  return off,config.GetConfig().Apserver.NumPerPage,nil
}

func GetSexToString(num uint) string {
  if num == 1 {
    return "男性"
  } else {
    return "女性"
  }
}

func GetInOutToString(num uint) string {
  if num == 1 {
    return "IN"
  } else {
    return "OUT"
  }
}

func DateToString(t *time.Time) string {
  return fmt.Sprintf("%04d年%02d月%02d日 %02d:%02d",t.Year(),int(t.Month()),t.Day(),t.Hour(),t.Minute()) 
}

func TimeToString(t time.Time) string {
  return fmt.Sprintf("%02d:%02d",t.Hour(),t.Minute())
}

func GetUserRealName(d []*model.Participant,uid uint) string {
  for _,v := range d {
    if v.UserID == uid {
      if v.User.RealName != nil {
        return *v.User.RealName
      } else {
        return v.User.ScreenName
      }
    }
  }
  return ""
}

func CalcAge(t time.Time) int {
    dateFormatOnlyNumber := "20060102"

    now := time.Now().Format(dateFormatOnlyNumber)
    birthday := t.Format(dateFormatOnlyNumber)

    nowInt, err := strconv.Atoi(now)
    if err != nil {
        return 0
    }
    birthdayInt, err := strconv.Atoi(birthday)
    if err != nil {
        return 0
    }

    age := (nowInt - birthdayInt) / 10000
    return age
}

