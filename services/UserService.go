package services

import (
	"errors"
	"myproject.com/grpc-kit/util"
	"strconv"
)
type IUserservice interface {
	GetName(userid int ) string
	DelUser(userid int ) error
}

type UserService struct {
	
}

func (this *UserService)GetName(userid int)string  {
	if userid == 101{
		return "zhangkai"+util.ServiceName+":"+strconv.Itoa(util.ServicePort)
	}
	return "guester"
}

func (this *UserService)DelUser(userid int) error  {
	if userid == 101{
		return errors.New("没有权限")
	}
	return nil
}
