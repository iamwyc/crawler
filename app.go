package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"iamwyc/crawler/crawler"
	"net/http"
)

func main() {
	gin.SetMode(gin.DebugMode)
	router := gin.Default()
	router.StaticFS("/assets", http.Dir("assets"))
	router.LoadHTMLGlob("templates/*")
	router.GET("/", index)
	router.GET("/crawler", crawl)
	router.Run()
}

func index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func crawl(c *gin.Context) {
	uri := c.Query("uri")
	err := crawler.HttpGet(uri)
	if err != nil {
		fmt.Printf("%v", err)
		c.JSON(200, "{\"ret\":\"err\"}")
	}else {
		c.JSON(200, "{\"ret\":\"0\"}")
	}
}
