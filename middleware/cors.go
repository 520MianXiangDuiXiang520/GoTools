package middleware

import (
	junebaotop "JuneGoBlog/src/junebao.top"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// 允许白名单
var AccessControlAllowOrigin []string = []string{
	"http://localhost:8081",
	"http://localhost:8082",
	"http://127.0.0.1:8889",
	"http://localhost:8889",
	"http://localhost:81",
	"http://39.106.168.39:80",
	"http://39.106.168.39:81",
}

func CorsHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		origin := context.GetHeader("Origin")
		log.Println(origin)
		method := context.Request.Method
		context.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		for _, allow := range AccessControlAllowOrigin {
			if allow == origin {
				context.Header("Access-Control-Allow-Origin", origin)
				break
			}
		}
		context.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		context.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,"+
			"session,X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT,"+
			" X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type,"+
			" Pragma,token,openid,opentoken, cookies, Cookies, cookie, Cookies")
		context.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, "+
			"Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar")
		context.Header("Access-Control-Max-Age", "172800")
		context.Header("Access-Control-Allow-Credentials", "true")
		context.Set("content-type", "application/json")
		//设置返回格式是json

		if method == "OPTIONS" {
			context.Abort()
			context.JSON(http.StatusOK, junebaotop.SuccessRespHeader)
		}
		context.Next()
	}
}
