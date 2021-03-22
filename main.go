package main

import (
	"github.com/gin-gonic/gin"
	"github.com/okoshiyoshinori/twigolf-server/config"
	"github.com/okoshiyoshinori/twigolf-server/logger"
	"github.com/okoshiyoshinori/twigolf-server/router"
)

var route *gin.Engine 

func init() {
  route = router.GetRouter()
}


func main(){
  logger.Info.Println("Listening and serving HTTP on",config.GetConfig().Apserver.Port)
  route.Run(config.GetConfig().Apserver.Port)
}
