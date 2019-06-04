package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/iamwyc/crawler/crawler"
	"github.com/iamwyc/crawler/global"
	"net/http"
)

func main() {
	gin.SetMode(gin.DebugMode)
	router := gin.Default()
	router.GET("/", index)
	router.GET("/assets/css/bootstrap.min.css", bootstrap)
	router.GET("/assets/js/jquery-3.4.1.min.js", jq)
	router.GET("/crawler", crawl)
	router.Run()
}

func index(c *gin.Context) {
	c.Data(http.StatusOK, "text/html", global.IndexHTML)
}

func bootstrap(c *gin.Context) {
	c.Data(http.StatusOK, "text/css", global.BootstrapCSS)
}
func jq(c *gin.Context) {
	c.Data(http.StatusOK, "application/x-javascript", global.JqJS)
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
