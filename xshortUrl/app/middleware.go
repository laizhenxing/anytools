package app

import (
	"net/http"
)

func LogMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 记录一次请求信息
		// fmt.Println("一次请求：", r.Form)
		h(w, r)
	}
}