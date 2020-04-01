package main

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"log"
	"testing"
	"time"
)

func TestJwt(t *testing.T)  {
	type UserClaim struct {
		Uname string `json:"username"`
        jwt.StandardClaims
	}
	priBytes, err := ioutil.ReadFile("./pem/private.pem")
	if err != nil {
		log.Fatal("私钥文件读取失败")
	}

	pubBytes, err := ioutil.ReadFile("./pem/public.pem")
	if err != nil {
		log.Fatal("公钥文件读取失败")
	}
	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(pubBytes)
	if err != nil {
		log.Fatal("公钥文件不正确")
	}

	priKey, err := jwt.ParseRSAPrivateKeyFromPEM(priBytes)
	if err != nil {
		log.Fatal("私钥文件不正确")
	}


	//sec := []byte("abc123")
	token_obj := jwt.NewWithClaims(jwt.SigningMethodRS256,UserClaim{Uname:"zhangkai"})
	token,err := token_obj.SignedString(priKey)
	if err != nil{
		panic(err)
	}
	fmt.Println(token)


	getToken ,err := jwt.ParseWithClaims(token, &UserClaim{},func(token *jwt.Token) (i interface{}, e error) {
		return pubKey,nil
	})
	if getToken.Valid {
		fmt.Println(getToken.Claims.(*UserClaim).Uname)
	}




}


func TestJwtExpire(t *testing.T)  {
	type UserClaim struct {
		Uname string `json:"username"`
		jwt.StandardClaims
	}
	priBytes, err := ioutil.ReadFile("./pem/private.pem")
	if err != nil {
		log.Fatal("私钥文件读取失败")
	}

	pubBytes, err := ioutil.ReadFile("./pem/public.pem")
	if err != nil {
		log.Fatal("公钥文件读取失败")
	}
	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(pubBytes)
	if err != nil {
		log.Fatal("公钥文件不正确")
	}

	priKey, err := jwt.ParseRSAPrivateKeyFromPEM(priBytes)
	if err != nil {
		log.Fatal("私钥文件不正确")
	}

   user := UserClaim{Uname:"zhangkai"}
   user.ExpiresAt = time.Now().Add(time.Second * 5).Unix()
   //UserClaim嵌套了jwt.StandardClaims，使用它的Add方法添加过期时间是5秒后，这里要使用unix()
	token_obj := jwt.NewWithClaims(jwt.SigningMethodRS256,user)
	token,err := token_obj.SignedString(priKey)
	if err != nil{
		panic(err)
	}
	fmt.Println(token)

   for {
	   getToken, err := jwt.ParseWithClaims(token, &UserClaim{}, func(token *jwt.Token) (i interface{}, e error) {
		   return pubKey, nil
	   })
	   if getToken.Valid {
		   fmt.Println(getToken.Claims.(*UserClaim).Uname)
	   }else if ve, ok := err.(*jwt.ValidationError); ok { //官方写法招抄就行
		   if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			   fmt.Println("错误的token")
		   } else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			   fmt.Println("token过期或未启用")
		   } else {
			   fmt.Println("无法处理这个token", err)
		   }

	   }
	   if err != nil{
	   	fmt.Println(err)
	   	break
	   }
	   time.Sleep(time.Second)
   }



}