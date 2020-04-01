package util

import (
	"context"
	"net/http"
)

type MyError struct {
	Code int
	Message string
}

func NewMyError(code int ,msg string)error  {
	return &MyError{Code:code,Message:msg}
}

func (this *MyError)Error()string{
	return this.Message
}

func MyErrorEncoder(ctx context.Context, err error, w http.ResponseWriter) {
	contentType, body := "text/plain; charset=utf-8", []byte(err.Error())
	w.Header().Set("Content-type", contentType) //设置请求头
	if myerr ,ok:= err.(*MyError);ok {
		w.WriteHeader(myerr.Code) //写入返回码
	}else{
		w.WriteHeader(500)
	}
	w.Write(body)

}
