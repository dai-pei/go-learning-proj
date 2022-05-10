package cotroller

import (
	"example-proj/service"
	"fmt"
	"strconv"
)

type PostResponse struct {
	Code int64       `json:code`
	Msg  string      `json:msg`
	Data interface{} `json:data`
}

// 这里Data返回topic title、post content，创建的时间戳

func PublishPost(parentIdStr string, content string) *PostResponse {
	parentId, _ := strconv.ParseInt(parentIdStr, 10, 64)

	fmt.Println(parentId)

	// if err != nil {
	// 	return &PostResponse{
	// 		Code: -1,
	// 		Msg:  err.Error(),
	// 	}
	// }

	// return nil

	postId, err := service.PublishPost(parentId, content)
	if err != nil {
		return &PostResponse{
			Code: -1,
			Msg:  err.Error(),
		}
	}
	return &PostResponse{
		Code: 0,
		Msg:  "success",
		Data: map[string]int64{
			"PostId": postId,
		},
	}
}
