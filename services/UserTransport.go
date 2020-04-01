package services

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func DecodeUserRequest(c context.Context,r *http.Request) (interface{},error) {
	/*if r.URL.Query().Get("uid") != ""{
		uid,_:= strconv.Atoi(r.URL.Query().Get("uid"))
       return UserRequest{ Uid:uid},nil
	}*/
	requestMap  := mux.Vars(r)
	if uid,ok := requestMap["uid"];ok{
		uid,_:= strconv.Atoi(uid)
		return UserRequest{ Uid:uid,Method:r.Method,Token:r.URL.Query().Get("token")},nil
	}
	return nil,errors.New("参数错误")
}

func EncodeUserResponse(c context.Context,w http.ResponseWriter ,response interface{})error{
	w.Header().Set("Content-type","application/json")
	return json.NewEncoder(w).Encode(response)
}


