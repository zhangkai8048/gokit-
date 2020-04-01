package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/consul"
	"github.com/go-kit/kit/sd/lb"
	httptransport "github.com/go-kit/kit/transport/http"
	consulapi "github.com/hashicorp/consul/api"
	"io"
	. "myproject.com/grpc-kit/clients"
	"net/url"
	"os"
	"time"
)

func callback(e error)error{
    fmt.Println("降级回调")
    return errors.New("降级回调")
}
func Myhystrix(run func()error)error{
	configA := hystrix.CommandConfig{ //创建一个hystrix的config
		Timeout:                3000,                  //command运行超过3秒就会报超时错误，并且在一个统计窗口内处理的请求数量达到阈值会调用我们传入的降级回调函数
		MaxConcurrentRequests:  5,                     //控制最大并发数为5，并且在一个统计窗口内处理的请求数量达到阈值会调用我们传入的降级回调函数
		RequestVolumeThreshold: 20,                     //判断熔断的最少请求数，默认是5；只有在一个统计窗口内处理的请求数量达到这个阈值，才会进行熔断与否的判断
		ErrorPercentThreshold:  5,                     //判断熔断的阈值，默认值5，表示在一个统计窗口内有50%的请求处理失败，比如有20个请求有10个以上失败了会触发熔断器短路直接熔断服务
		SleepWindow:            int(time.Second * 10), //熔断器短路多久以后开始尝试是否恢复，这里设置的是10
	}
	hystrix.ConfigureCommand("getUser", configA) //hystrix绑定command
	err:=hystrix.Do("getUser", func() error {
		err:=run()
		return err
	},callback)
	if err != nil{
		return err
	}
	return nil
}
func run() error {
	//通过注册中心获取服务
	{
		config := consulapi.DefaultConfig()
		config.Address = "localhost:8500"
		api_client, _ := consulapi.NewClient(config)
		client:= consul.NewClient(api_client)

		var logger log.Logger
		{
			logger = log.NewLogfmtLogger(os.Stdout)
			var Tag = []string{"primary"}
			name := flag.String("name", "userservice", "name")
			flag.Parse()
			if *name == "" {
				//panic("请指定服务名")
				return errors.New("请指定服务名")
			}

			//循环调用服务
			for true {
              time.Sleep(1*time.Second)
			//第二部创建一个consul的实例
			instancer := consul.NewInstancer(client, logger, *name, Tag, true)
			//最后的true表示只有通过健康检查的服务才能被得到
			{
				factory := func(service_url string) (endpoint.Endpoint, io.Closer, error) {
					//factory定义了如何获得服务端的endpoint,这里的service_url是从consul中读取到的service的address我这里是192.168.3.14:8000
					tart, _:= url.Parse("http://" + service_url)
					//server ip +8080真实服务的地址
					return httptransport.NewClient("GET", tart, EncodeRequestGetUserInfo, DecodeResponseGetUserInfo).Endpoint(), nil, nil
					//我再GetUserInfo_Request里面定义了访问哪一个api把url拼接成了http://192.168.3.14:8000/v1/user/{uid}的形式
				}
				endpointer := sd.NewEndpointer(instancer, factory, logger)
				endpoints, _ := endpointer.Endpoints()
				fmt.Println("服务有", len(endpoints), "条")
				//mylb := lb.NewRoundRobin(endpointer) //使用go-kit自带的轮询
				mylb := lb.NewRandom(endpointer, time.Now().UnixNano()) //使用go-kit自带的轮询
				getUserInfo,err  := mylb.Endpoint() //轮询获取服务
				if err != nil {
					os.Exit(1)
					return  err
				}
				ctx := context.Background() //第三步：创建一个context上下文对象

				//第四步：执行
				res, err := getUserInfo(ctx, UserRequest{Uid: 101})
				if err != nil {
					os.Exit(1)
					return  err
				}
				//第五步：断言，得到响应值
				userinfo := res.(UserResponse)
				fmt.Println(userinfo.Result)
			}

			}
		}
	}
	return nil

}

func main()  {
	err := Myhystrix(run)
	if err != nil{
		fmt.Println(err)
	}
}