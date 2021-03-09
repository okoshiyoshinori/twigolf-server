package util

import (
	"strconv"

	"github.com/okoshiyoshinori/twigolf-server/config"
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

