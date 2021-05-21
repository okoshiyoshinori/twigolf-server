package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/okoshiyoshinori/twigolf-server/config"
	"github.com/okoshiyoshinori/twigolf-server/controller"
  "time"
)

func GetRouter() *gin.Engine {
  r := gin.Default()
  r.LoadHTMLGlob("templates/*.tmpl")
  store := cookie.NewStore([]byte(config.GetConfig().Apserver.CookieKey))
  store.Options(sessions.Options{
    Path:"/",
    Secure:false,
    HttpOnly: true,
    MaxAge: 20 * 60 * 60 * 24,
  })
  r.Use(sessions.Sessions("twigolfSession",store))
  r.Use(cors.New(cors.Config{
    AllowOrigins:config.GetConfig().Apserver.Origin,
    AllowMethods: []string{
        "POST",
        "GET",
        "DELETE",
        "OPTIONS",
    },
    AllowHeaders: []string{
        "Access-Control-Allow-Credentials",
        "Access-Control-Allow-Origin",
        "Access-Control-Allow-Headers",
        "Content-Type",
        "Content-Length",
        "Accept-Encoding",
        "Authorization",
    },
    AllowCredentials: true,
    MaxAge: 24 * time.Hour,
  }))
 
  
  privateApi := r.Group("/api/private/",PrivateMiddleware())
  privateRouter(privateApi)
  publicApi := r.Group("/api/public/")
  publicRouter(publicApi)

  return r
}

func privateRouter(group *gin.RouterGroup) {
  group.POST("/competition",controller.PostCompetiton)
  group.POST("/participant",controller.PostParticipant)
  group.POST("/comments",controller.PostComment)
  group.POST("/user_basic_info",controller.PostUserBasicInfo)
  group.POST("/bundle_participant",controller.BundleParticipant)
  group.POST("/combination/:cid",controller.PostCombination)
  group.DELETE("/logout",controller.Logout)
  group.DELETE("/competition/:id",controller.DeleteCompetiton)
  group.GET("/participants_with_name/:cid",controller.GetPaticipantsWithRealName)
  group.GET("/get_combination_excel/:cid",controller.GetCombinationExcel)
  group.POST("/send_dm",controller.PostDm)
}

func publicRouter(group *gin.RouterGroup) {
  group.GET("/combination/:cid",controller.GetCombination)
  group.GET("/user/:screen_name",controller.GetUser)
  group.GET("/competition",controller.GetCompetition)
  group.GET("/competition/:id",controller.GetCompetitonDetail)
  group.GET("/comments/:cid",controller.GetComment)
  group.GET("/participants/:cid",controller.GetPaticipants)
  group.GET("/search",controller.SearchCompetition)
  group.GET("/clubs",controller.GetClubs)
  group.GET("/user_competitions/:screen_name",controller.GetUserCompetitions)
  group.GET("/session",controller.GetSession)
  group.GET("/twitterlogin",controller.GetTwitterAuthUrl)
  group.GET("/twittercallback",controller.TwitterCallBack)
  group.GET("/twitter_bot/:cid",controller.SendOgpHtml)
}
/*
func CorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
    oo
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
*/

func PrivateMiddleware() gin.HandlerFunc {
  return func(c *gin.Context) {
    //ここに認証関連の処理を追加
    session := sessions.Default(c)
    if _,ok := session.Get("user_id").(uint); !ok {
      c.AbortWithStatus(401)
      return
    }
    c.Next()
  }
}
