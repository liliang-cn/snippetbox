package main

import (
	"fmt"
	"net/http"
)

// 中间件就是一个函数，接受调用链中下一个处理函数作为参数，执行完某些逻辑之后并调用下一个处理函数，然后返回一个处理函数
// servemux 之前的中间件函数将作用于所有的请求
// servemux 之后的中间件函数将作用于特定的请求处理函数
func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("X-Frame-Options", "deny")

		next.ServeHTTP(w, r)
	})
}

// logRequest 的签名与普通的中间件函数一样，所以可以调用，此外由于是 application 的方法，可以调用 infoLog 的方法
func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())

		next.ServeHTTP(w, r)
	})
}

// 从 panic 中恢复并发送错误的响应
func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// defer 函数在 go 清理时最终会被调用
		defer func() {
			// 检查是否有 panic 发生
			if err := recover(); err != nil {
				// 关闭连接
				w.Header().Set("Connection", "close")
				// 写入 错误输出
				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// 需要认证的中间件
func (app *application) requireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !app.isAuthenticated(r) {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}

		w.Header().Add("Cache-Control", "no-store")
		next.ServeHTTP(w, r)
	})
}
