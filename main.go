package main

import (
	"gitee.com/webook/internal/repository"
	"gitee.com/webook/internal/repository/dao"
	"gitee.com/webook/internal/service"
	"gitee.com/webook/internal/web"
	"gitee.com/webook/internal/web/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

func main() {

	db := initDB()
	server := initWebServer()
	u := initUser(db)

	u.RegisterUserRoutes(server)
	server.Run(":8080")
}

func initWebServer() *gin.Engine {
	server := gin.Default()
	server.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:3000"},
		//AllowMethods: []string{"POST"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"x-jwt-token"},
		AllowCredentials: true,
		//AllowOriginFunc: func(origin string) bool {
		//	return origin == "https://github.com"
		//},
		MaxAge: 12 * time.Hour,
	}))
	//store := cookie.NewStore([]byte("secret"))
	//store, err := redis.NewStore(
	//	16, "tcp",
	//	"localhost:6379",
	//	"",
	//	[]byte("VgMkqI1mTTsfR+yn"),
	//	[]byte("VgMkqI1mTTsfR+yn"),
	//)
	//if err != nil {
	//	panic(err)
	//}
	//server.Use(sessions.Sessions("mysession", store))

	//server.Use(middleware.NewLoginMiddlewareBuilder().Build())
	server.Use(middleware.NewLoginJWTMiddlewareBuilder().Build())
	return server
}

func initUser(db *gorm.DB) *web.UserHandler {
	ud := dao.NewUserDAO(db)
	repo := repository.NewUserRepository(ud)
	svc := service.NewUserService(repo)

	u := web.NewUserHandler(svc)
	return u
}

func initDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:8082)/webook"))
	if err != nil {
		// 只会再初始化过程中 panic
		// panic 相当于整个 goroutine 结束
		// 一旦初始化过程出错 应用就不要启动了
		panic(err)
	}
	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}

	return db
}
