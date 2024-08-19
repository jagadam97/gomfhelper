package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jagadam97/backend/apis"
)


func main() {
    router := gin.Default()
    router.GET("/wapi", apis.GetMutualFunds)

    router.Run("localhost:8080")
}