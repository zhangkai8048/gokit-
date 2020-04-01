package main

import (
	"testing"
	"myproject.com/grpc-kit/util"
)
func TestGenRSAPubAndPri(t *testing.T){
	util.GenRSAPubAndPri(1024,"./pem")
}

