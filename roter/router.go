package router

import (
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/sessions"
	"github.com/okoshiyoshinori/twigolf-server/config"
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

func privateRouter(api *gin.RouterGroup) {
}

func publicRouter(api *gin.RouterGroup) {

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
