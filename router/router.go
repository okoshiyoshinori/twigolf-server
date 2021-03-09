package router

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/okoshiyoshinori/twigolf-server/config"
	"github.com/okoshiyoshinori/twigolf-server/controller"
)

func GetRouter() *gin.Engine {
  r := gin.Default()
  store := cookie.NewStore([]byte(config.GetConfig().Apserver.CookieKey))
  r.Use(sessions.Sessions("twigoluSession",store))
  r.Use(CorsMiddleware())
  
  privateApi := r.Group("/api/private/",PrivateMiddleware())
  privateRouter(privateApi)
  publicApi := r.Group("/api/public/")
  publicRouter(publicApi)

  return r
}

func privateRouter(group *gin.RouterGroup) {
}

func publicRouter(group *gin.RouterGroup) {
  group.GET("/prefs",controller.GetPref)
  group.GET("/user/:uid",controller.GetUser)
  group.GET("/competition",controller.GetCompetition)
  group.GET("/competition/:id",controller.GetCompetitonDetail)
  group.GET("/comments/:cid",controller.GetComment)
  group.GET("/participants/:cid",controller.GetPaticipants)
  group.GET("/search",controller.SearchCompetition)
  group.GET("/clubs",controller.GetClubs)
  //private
  group.POST("/competition",controller.PostCompetiton)
}

func CorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", config.GetConfig().Apserver.Origin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func PrivateMiddleware() gin.HandlerFunc {
  return func(c *gin.Context) {
    //ここに認証関連の処理を追加
    c.Next()
  }
}
