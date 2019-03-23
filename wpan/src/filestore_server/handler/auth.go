package handler

import (
	"net/http"
	"fmt"
)

// http请求拦截器
func HTTPInterceptor(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		username := r.Form.Get("username")
		token := r.Form.Get("token")
		if len(username) < 3 || !IsTokenValid(token) {
			fmt.Println(len(username))
			w.WriteHeader(http.StatusForbidden)
			return
		}
		h(w,r)
	})
}