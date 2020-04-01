"# gokit-" 
post  http://127.0.0.1:8082/access-token   {"username":"jerry","userpass":"123"}
get http://127.0.0.1:8082/user/101?token=带上上一步的token参数
 
