package util

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"

)

//日志中间件,每一个service都应该有自己的日志中间件
func UserServiceLogMiddleware(logger log.Logger) endpoint.Middleware { //Middleware type Middleware func(Endpoint) Endpoint
	return func(next endpoint.Endpoint) endpoint.Endpoint { //Endpoint type Endpoint func(ctx context.Context, request interface{}) (response interface{}, err error)
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			//r := request.(UserRequest) //通过类型断言获取请求结构体
			logger.Log("method","event", "get user", "userid")
			return next(ctx, request)
		}
	}
}
