package main

import (
	"fmt"
	"net/http"
	"project3/gee"
)

func main() {
	r := gee.New()
	r.GET("/", func(ctx *gee.Context) {
		fmt.Fprintf(ctx.Writer, "URL.Path = %q\n", ctx.Path)
	})

	r.GET("/hello", func(ctx *gee.Context) {
		ctx.String(http.StatusOK, "hello %s, you're at %s\n", ctx.Query("name"), ctx.Path)
	})

	r.POST("/login", func(ctx *gee.Context) {
		ctx.JSON(http.StatusOK, gee.H{
			"username:": ctx.PostForm("username"),
			"password":  ctx.PostForm("password"),
		})
	})

	r.Run(":8080")
}
