package clients

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

func EncodeRequestGetUserInfo (c context.Context, r *http.Request, response interface{}) error{
     user_request := response.(UserRequest)
     r.URL.Path += "/user/"+strconv.Itoa(user_request.Uid)
     return nil
}

func DecodeResponseGetUserInfo (c context.Context,r  *http.Response) (response interface{}, err error){
	if r.StatusCode > 400 {
		return nil,errors.New("404 error page")
	}
	//将json 转化为对象
	var user_response UserResponse
	err = json.NewDecoder(r.Body).Decode(&user_response)
	if err != nil{
		return nil,err
	}
	return user_response ,nil
}