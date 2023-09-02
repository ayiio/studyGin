package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func func1(ctx *gin.Context) {
	fmt.Println("func1...")
}

func func2(ctx *gin.Context) {
	fmt.Println("func2 before")
	ctx.Next()
	fmt.Println("func2 after")
}

func func3(ctx *gin.Context) {
	fmt.Println("func3...")
	//ctx.Abort()
}

func func4(ctx *gin.Context) {
	fmt.Println("func4...")
	ctx.Set("name", "test")
}

func func5(ctx *gin.Context) {
	fmt.Println("func5...")
	value, ok := ctx.Get("name")
	if ok {
		vStr := value.(string)
		fmt.Println(vStr)
	}
}

func main() {
	r := gin.Default()

	r.GET("/hello", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "ok")
	})

	//r.Use()
	shopGroup := r.Group("/shop", func1, func2)
	shopGroup.Use(func3)
	{
		shopGroup.GET("/index", func4, func5)
	}

	r.Run(":8080")
}
