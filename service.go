package main

import (
	"flag"
	"fmt"
	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"golang.org/x/time/rate"

	. "myproject.com/grpc-kit/services"
	"myproject.com/grpc-kit/util"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)
var logger log.Logger
func init()  {

	{
		logger = log.NewLogfmtLogger(os.Stdout)
		logger = log.WithPrefix(logger, "mykit", "1.0")
		logger = log.WithPrefix(logger, "time", log.DefaultTimestampUTC) //加上前缀时间
		logger = log.WithPrefix(logger, "caller", log.DefaultCaller)     //加上前缀，日志输出时的文件和第几行代码

	}
	port := flag.Int("port", 8080, "http listen port")
	name := flag.String("name", "userservice", "name")
	flag.Parse()
	if *name == "" {
		logger.Log("请指定服务名")
	}
	if *port == 0 {
		logger.Log("请指定端口")
	}
	util.SetServiceNameAndPort(*name,*port)
}

func main()  {
	//限流
	limit := rate.NewLimiter(1, 1)
	//统一错误处理
	options := []httptransport.ServerOption{ //生成ServerOtion切片，传入我们自定义的错误处理函数
		httptransport.ServerErrorEncoder(util.MyErrorEncoder),
	}
    //用户服务
	user := &UserService{}
	//这里添加中间件
	endp := util.RateLimit(limit)((CheckTokenMiddleware())((UserServiceLogMiddleware(logger))(GenUserEnpointMiddleware(user))))
	serverHander := httptransport.NewServer(endp, DecodeUserRequest, EncodeUserResponse, options...)
    //权限服务
	accessService := &AccessService{}
	endp_access := AccessEndpoint(accessService)
	accessHandler := httptransport.NewServer(endp_access,DecodeAccessRequest,EncodeAccessResponse,options...)

     //生成路由
     router :=  mux.NewRouter()
	{
		router.Methods("POST").Path("/access-token").Handler(accessHandler)            //注册token获取的handler
		router.Methods("GET", "DELETE").Path(`/user/{uid:\d+}`).Handler(serverHander)
		router.Methods("GET").Path("/health").HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			writer.Header().Set("Content-type","application/json")
			writer.Write([]byte(`{"status":"OK"}`))
		})
	}

    errChan := make(chan error)
     //开启协程
     go func() {
		 //注册中心注册服务
		 util.RegService()
		 err := http.ListenAndServe(":"+strconv.Itoa(util.ServicePort),router)
		 if err != nil{
			 fmt.Println(err.Error())
			 //如果有异常，清除注册中心的服务
			 errChan <- err
		 }
	 }()

     go func() {
     	sigChan := make(chan os.Signal)
     	signal.Notify(sigChan,syscall.SIGINT,syscall.SIGTERM)
		 errChan <- fmt.Errorf("%s", <- sigChan)
	 }()

    //从通道读取信号
	getErr  := <- errChan
    util.UnRegService()
   logger.Log(getErr)




}
