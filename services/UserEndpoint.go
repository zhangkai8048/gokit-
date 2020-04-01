package services

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"golang.org/x/time/rate"
	"myproject.com/grpc-kit/util"
)

type UserRequest struct {
  Uid int `json:"uid"`
  Method  string
  Token  string
}

type UserResponse struct {
 Result string `json:"result"`
}

func GenUserEnpointMiddleware(userservice IUserservice)endpoint.Endpoint  {
	return func(ctx context.Context, request interface{}) (response interface{}, err error){
		r:= request.(UserRequest)
		fmt.Println("当前登录用户名：",ctx.Value("LoginUser"))
		result := "nothing"
		if r.Method == "GET" {
			result = userservice.GetName(r.Uid)
		}else if r.Method=="DELETE"{
          err := userservice.DelUser(r.Uid)
          if err != nil{
          	 result = err.Error()
		  }else{
			  result =  fmt.Sprint("userid 为 %d 的用户删除成功",r.Uid)
		  }
		}

      return  UserResponse{Result:result},nil
	}
}

//token验证中间件
func CheckTokenMiddleware() endpoint.Middleware { //Middleware type Middleware func(Endpoint) Endpoint
	return func(next endpoint.Endpoint) endpoint.Endpoint { //Endpoint type Endpoint func(ctx context.Context, request interface{}) (response interface{}, err error)
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			r := request.(UserRequest) //通过类型断言获取请求结构体
			uc := UserClaim{}
			//下面的r.Token是在代码DecodeUserRequest那里封装进去的
			getToken, err := jwt.ParseWithClaims(r.Token, &uc, func(token *jwt.Token) (i interface{}, e error) {
				return []byte(secKey), err
			})
			fmt.Println(err, 123)
			if getToken != nil && getToken.Valid { //验证通过
				newCtx := context.WithValue(ctx, "LoginUser", getToken.Claims.(*UserClaim).Uname)
				return next(newCtx, request)
			} else {
				return nil, util.NewMyError(403, "error token")
			}

			//logger.Log("method", r.Method, "event", "get user", "userid", r.Uid)

		}
	}
}

//日志中间件,每一个service都应该有自己的日志中间件
func UserServiceLogMiddleware(logger log.Logger) endpoint.Middleware { //Middleware type Middleware func(Endpoint) Endpoint
	return func(next endpoint.Endpoint) endpoint.Endpoint { //Endpoint type Endpoint func(ctx context.Context, request interface{}) (response interface{}, err error)
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			r := request.(UserRequest) //通过类型断言获取请求结构体
			logger.Log("method", r.Method, "event", "get user", "userid", r.Uid)
			return next(ctx, request)
		}
	}
}

//加入限流功能中间件
func RateLimitMiddleware(limit *rate.Limiter) endpoint.Middleware { //Middleware type Middleware func(Endpoint) Endpoint
	return func(next endpoint.Endpoint) endpoint.Endpoint { //Endpoint type Endpoint func(ctx context.Context, request interface{}) (response interface{}, err error)
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			if !limit.Allow() {
				return nil, util.NewMyError(429, "toot many request")
			}
			return next(ctx, request) //执行endpoint
		}
	}
}