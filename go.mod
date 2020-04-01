module myproject.com/grpc-kit

go 1.13

replace myproject.com/grpc-kit/services => E:\\src\\myproject.com\\grpc-kit\\services

replace myproject.com/grpc-kit/cliens => E:\\src\\myproject.com\\grpc-kit\\cliens

replace myproject.com\grpc-kit\util => E:\src\myproject.com\grpc-kit\util

require (
	github.com/afex/hystrix-go v0.0.0-20180502004556-fa1af6a1f4f5
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/go-kit/kit v0.10.0
	github.com/google/uuid v1.1.1
	github.com/gorilla/mux v1.7.3
	github.com/hashicorp/consul/api v1.3.0
	github.com/hashicorp/golang-lru v0.5.3 // indirect
	github.com/miekg/dns v1.1.22 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/tidwall/gjson v1.6.0
	github.com/tidwall/pretty v1.0.1 // indirect
	golang.org/x/crypto v0.0.0-20191108234033-bd318be0434a // indirect
	golang.org/x/net v0.0.0-20191109021931-daa7c04131f5 // indirect
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0
)
