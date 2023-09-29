package middleware

import (
	"encoding/gob"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type LoginMiddlewareBuilder struct {
}

func NewLoginMiddlewareBuilder() *LoginMiddlewareBuilder {
	return &LoginMiddlewareBuilder{}
}

func (l *LoginMiddlewareBuilder) Build() gin.HandlerFunc {
	gob.Register(time.Now())
	return func(context *gin.Context) {
		println("in login middleware")
		if context.Request.URL.Path == "/users/login" || context.Request.URL.Path == "/users/signup" {
			return
		}
		s := sessions.Default(context)
		println("s", s)
		id := s.Get("userId")

		if id == nil {
			context.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		println("id", id)
		updateTime := s.Get("update_time")
		s.Set("userId", id)
		s.Options(sessions.Options{
			MaxAge: 10,
		})
		println("updateTime", updateTime)
		now := time.Now().UnixMilli()
		// 说明第一次登录
		if updateTime == nil {

			s.Set("update_time", now)

			if err := s.Save(); err != nil {
				println("first save error")
				fmt.Printf("Found error %v\n", err)
			}
			return
		}
		// updateTime 有
		updateTimeVal, ok := updateTime.(int64)
		if !ok {
			context.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if now-updateTimeVal > 5*1000 {
			s.Set("update_time", now)
			println("save error")
			if err := s.Save(); err != nil {
				fmt.Printf("Found error %v\n", err)
			}

		}
	}
}
