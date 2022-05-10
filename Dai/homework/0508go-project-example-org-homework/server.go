package main

import (
	"os"

	"example-proj/cotroller"
	"example-proj/repository"

	"gopkg.in/gin-gonic/gin.v1"
)

func main() {
	if err := Init("./data/"); err != nil {
		os.Exit(-1)
	}
	r := gin.Default()
	r.GET("/community/page/get/:id", func(c *gin.Context) {
		topicId := c.Param("id")
		data := cotroller.QueryPageInfo(topicId)
		c.JSON(200, data)
	})
	r.POST("/community/page/post", func(c *gin.Context) {
		topicId, _ := c.GetPostForm("topicId")
		content, _ := c.GetPostForm("content")
		data := cotroller.PublishPost(topicId, content)
		c.JSON(200, data)

	})
	err := r.Run()
	if err != nil {
		return
	}
}

func Init(filePath string) error {
	if err := repository.Init(filePath); err != nil {
		return err
	}
	return nil
}
