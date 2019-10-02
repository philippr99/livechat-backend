package auth

import (
	"context"
	"github.com/gorilla/websocket"
	"net/http"
)
var upgrader = websocket.Upgrader{}
func Middleware() func (http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// Headers to allow CROSS-Origin
			w.Header().Add("Access-Control-Allow-Origin", "*")
			w.Header().Add("Access-Control-Allow-Methods", "POST, GET")
			w.Header().Add("Access-Control-Max-Age", "3600")
			w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With, Token",);

			// and call the next with our new context
			ctx := context.WithValue(r.Context(), "token", r.Header.Get("Token"))
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}
