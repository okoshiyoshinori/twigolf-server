package main

import (
	"github.com/gin-gonic/gin"
	"github.com/okoshiyoshinori/twigolf-server/config"
	"github.com/okoshiyoshinori/twigolf-server/router"
)

var route *gin.Engine 

func init() {
  route = router.GetRouter()
}


func main(){
  route.Run(config.GetConfig().Apserver.Port)
}
