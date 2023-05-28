package main

import (
	"context"
	"log"
	"net/http"
	"time"


	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"

	"github.com/Degoke/doc-verify/common"
	"github.com/Degoke/doc-verify/user"
	"github.com/Degoke/doc-verify/verification"
)

var ctx = context.TODO()


var (
	g errgroup.Group
)

func mainRouter() http.Handler {
	e := gin.New()
	e.Use(gin.Recovery())

	v1 := e.Group("/api/v1")

	user.UsersRouter(v1.Group("/users"))
	user.UserRouter(v1.Group("/user"))

	v1.Use(user.AuthMiddleware(true))
	verification.VerificationRouter(v1.Group("/verification"))
	verification.VerificationsRouter(v1.Group("/verifications"))

	return e
}

func fileRouter() http.Handler {
	e := gin.New()

	e.Use(gin.Recovery())
	e.MaxMultipartMemory = 8 << 20

	v1 := e.Group("/api/v1/upload")
	v1.Use(user.AuthMiddleware(true))

	verification.UploadRouter(v1.Group("verification"))

	return e
}

func main() {
	common.LoadENV()
	client := common.Init()

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	
	mainServer := &http.Server{
		Addr: ":8080",
		Handler: mainRouter(),
		ReadTimeout: 5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	fileServer := &http.Server{
		Addr: ":8081",
		Handler: fileRouter(),
		ReadTimeout: 5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	g.Go(func() error {
		return mainServer.ListenAndServe()
	})

	g.Go(func() error {
		return fileServer.ListenAndServe()
	})

	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}