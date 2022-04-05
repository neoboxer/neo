package neo

import (
	"fmt"
	"log"
	"net/http"
	"testing"
)

func TestEngine(t *testing.T) {
	engine := NewEngine()
	engine.Use(func(ctx *Context) {
		log.Println("[INFO].....")
		ctx.Next()
	}).Use(func(ctx *Context) {
		log.Println("[DEBUG]......")
		ctx.Next()
	})

	engine.GET("/", func(ctx *Context) {
		ctx.JSON(http.StatusOK, M{})
	})

	engine.GET("/hello", func(ctx *Context) {
		ctx.Next()
		fmt.Println("hello")
	}, func(ctx *Context) {
		ctx.Next()
		fmt.Println("world")
	}, func(ctx *Context) {
		ctx.String(http.StatusOK, "hello %s, you're at %s\n", ctx.Query("name"), ctx.Path)
	})

	engine.GET("/hello/*a/b", func(ctx *Context) {
		ctx.JSON(http.StatusOK, M{
			"path": ctx.Param("a"),
		})
	})
	engine.Run(":8888")

	t.Log("hello")
}
