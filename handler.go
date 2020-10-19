package junebao_top

import (
	"JuneGoBlog/src/junebao.top/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"reflect"
)

type CheckFunc func(ctx *gin.Context, req BaseReqInter) (BaseRespInter, error)
type LogicFunc func(ctx *gin.Context, req BaseReqInter) BaseRespInter

// 解析请求，整合检查请求参数，响应逻辑，并响应
func EasyHandler(cf CheckFunc, lf LogicFunc, req interface{}) gin.HandlerFunc {
	// EasyHandler 只会执行一次， 每次请求过来真正执行的是 EasyHandler 返回的这个 HandlerFunc
	// 所以从 routes 中传过来的参数 req 并不会与上下文绑定，HandlerFunc 会根据 req 的类型
	// 反射获得一个新的 request, 避免两次请求的参数相互叠加
	return func(context *gin.Context) {
		t := reflect.TypeOf(req)
		request := reflect.New(t).Interface().(BaseReqInter)
		var resp interface{}
		if err := request.JSON(context); err != nil {
			msg := fmt.Sprintf("Request binding failed，type of req is %v, context is %v",
				reflect.TypeOf(req), context)
			utils.ExceptionLog(err, msg)
			resp = ParamErrorRespHeader
		} else {
			if checkResp, err := cf(context, request); err != nil {
				resp = checkResp
			} else {
				resp = lf(context, request)
			}
		}
		context.Set("resp", resp)
		context.JSON(http.StatusOK, resp)
	}
}
