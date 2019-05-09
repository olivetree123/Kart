package middleware

import (
	"fmt"
	"kart/config"
	"kart/global"
	"kart/utils"
	"net/http"
	"strings"
)

// Auth 认证中间件
func Auth(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("method = ", r.Method)
		if r.Method == "OPTIONS" {
			h.ServeHTTP(w, r)
			utils.JSONResponse(nil, w)
			return
		}
		auth := r.Header.Get("Authorization")
		fmt.Println("Authorization = ", auth)
		rs := strings.Split(auth, " ")
		if len(rs) != 2 || rs[0] != config.Config.GetString("TokenScheme") || global.GetToken(rs[1]) == nil {
			utils.PermitCORS(w)
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		user := global.GetToken(rs[1])
		r.Header.Set("userID", user.ID)
		h.ServeHTTP(w, r)
	})
}
