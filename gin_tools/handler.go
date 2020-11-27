package gin_tools

import (
	"fmt"
	"github.com/520MianXiangDuiXiang520/GinTools/check_tools"
	"github.com/520MianXiangDuiXiang520/GinTools/log_tools"
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

		resp := struct {
			Header BaseRespInter `json:"header"`
		}{
			Header: ParamErrorRespHeader,
		}

		// 请求数据绑定
		if err := request.JSON(context); err != nil {
			msg := fmt.Sprintf("Request binding failed，type of req is %v, context is %v",
				reflect.TypeOf(req), context)
			utils.ExceptionLog(err, msg)
			resp.Header = ParamErrorRespHeader
			context.Set("resp", resp)
			context.JSON(http.StatusOK, resp)
		} else {
			// 标签检查请求参数
			if !check_tools.CheckRequest(request) {
				resp.Header = ParamErrorRespHeader
				context.Set("resp", resp)
				context.JSON(http.StatusOK, resp)
			} else {
				// 自定义方法检查请求参数
				if checkResp, err := cf(context, request); err != nil {
					context.Set("resp", checkResp)
					context.JSON(http.StatusOK, checkResp)
				} else {
					// 执行业务逻辑
					r := lf(context, request)
					context.Set("resp", r)
					context.JSON(http.StatusOK, r)
				}
			}

		}

		context.Set("resp", resp)
		context.JSON(http.StatusOK, resp)
	}
}
