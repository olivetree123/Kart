package auth

import (
	"kart/config"
	"kart/global"
	"net/http"
	"strings"
)

// Middleware 认证中间件
func Middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		rs := strings.Split(auth, " ")
		if len(rs) != 2 || rs[0] != config.Config.GetString("TokenScheme") || global.GetToken(rs[1]) == nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		h.ServeHTTP(w, r)
	})
}
