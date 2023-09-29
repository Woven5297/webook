package middleware

import (
	"encoding/gob"
	"gitee.com/webook/internal/web"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
	"time"
)

type LoginJWTMiddlewareBuilder struct {
}

func NewLoginJWTMiddlewareBuilder() *LoginJWTMiddlewareBuilder {
	return &LoginJWTMiddlewareBuilder{}
}

func (l *LoginJWTMiddlewareBuilder) Build() gin.HandlerFunc {
	gob.Register(time.Now())
	return func(context *gin.Context) {
		println("in login middleware")
		if context.Request.URL.Path == "/users/login" || context.Request.URL.Path == "/users/signup" {
			return
		}
		tokenStr := context.GetHeader("Authorization")
		if tokenStr == "" {
			context.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		s := strings.Split(tokenStr, " ")
		if len(s) != 2 {
			context.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		token := s[1]
		claims := &web.UserClaims{}
		t, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0"), nil
		})
		if err != nil {
			context.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if t == nil || t.Valid {
			context.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}
}
