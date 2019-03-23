package handler

import (
	"net/http"
	"time"
	"fmt"
	"io/ioutil"
	"filestore_server/util"
	dblayer "filestore_server/db"
)

const(
	pwd_salt = "*#890"
)

// 处理用户注册请求
func SignupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		data,err := ioutil.ReadFile("./static/view/signup.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(data)
		return
	}
	r.ParseForm()
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	phone := r.Form.Get("phone")

	if len(username) < 3 || len(password) < 5 {
		w.Write([]byte("invalid parameter"))
		return 
	}

	enc_pwd := util.Sha256([]byte(password+pwd_salt))
	suc := dblayer.UserSignup(username, enc_pwd, phone)
	if suc {
		w.Write([]byte("SUCCESS"))
	} else {
		w.Write([]byte("FAILED"))
	}
}

// 登录接口
func SigninHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	encPwd := util.Sha256([]byte(password+pwd_salt))

	// 1.校验用户名密码
	pwdChecked := dblayer.UserSignin(username, encPwd)
	if !pwdChecked {
		w.Write([]byte("FAILED"))
		return 
	}
	// 2.生成访问token
	token := GenToken(username)
	upRes := dblayer.UpdateToken(username, token)
	if !upRes {
		w.Write([]byte("FAILED"))
		return 
	}
	// 3.重定向到首页 
	//w.Write([]byte("http://"+r.Host+"/static/view/home.html"))
	resp := util.RespMsg{
		Code : 0,
		Msg : "OK",
		Data : struct{
			Location string
			Username string
			Token string
		}{
			Location : "http://" + r.Host + "/static/view/home.html",
			Username : username,
			Token : token,
		},
	}
	w.Write(resp.JSONBytes())
}

func UserInfoHandler(w http.ResponseWriter, r *http.Request) {
	// 1. 解析请求参数
	r.ParseForm()
	username := r.Form.Get("username")
	token := r.Form.Get("token")
	// 2. 验证token是否有效
	isTokenValid := IsTokenValid(token)
	if !isTokenValid {
		w.WriteHeader(http.StatusForbidden)
		return 
	}
	// 3. 查询用户信息
	user,err := dblayer.GetUserInfo(username)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return 
	}
	// 4. 组装并相应用户数据
	resp := util.RespMsg{
		Code : 0,
		Msg : "OK",
		Data :user,
	}
	w.Write(resp.JSONBytes())
}

func GenToken(username string) string {
	// 40位  MD5(username+timestamp+token_salt) + timestamp[:8]
	ts := fmt.Sprintf("%x", time.Now().Unix())
	tokenPrefix := util.MD5([]byte(username+ts+"_tokensalt"))
	return tokenPrefix + ts[:8]
}

func IsTokenValid(token string) bool {
	if len(token) != 40 {
		return false
	}
	// TODO: 判断token的时效性，是否过期
	// TODO: 从数据库表tbl_user_token查询username对应的token信息
	// TODO: 对比两个token是否一致
	return true
}