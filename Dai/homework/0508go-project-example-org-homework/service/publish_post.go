package service

import (
	"errors"
	"example-proj/repository"
	"fmt"
	"time"
	"unicode/utf16"

	idworker "github.com/gitstliu/go-id-worker"
)

// service中query_page_info的功能
// 1. check params
// 2. prepare Info
// 3. pack date

// 在service层的query_page_info.go中，没有直接写对repository的增删改查操作接口，而是
// 声明了一个flow类，所有操作都是类内函数，这样能够更好地封装？

// type Post struct {
// 	Id         int64  `json:"id"`
// 	ParentId   int64  `json:"parent_id"`
// 	Content    string `json:"content"`
// 	CreateTime int64  `json:"create_time"`
// }

// 本service要完成的工作，封装在一个类内：
// 验证参数正确性
// 生成id
// 生成时间戳
// 封装数据
// 调用dao写入文件

var idGen *idworker.IdWorker

func init() {
	idGen = &idworker.IdWorker{}
	idGen.InitIdWorker(1, 1)
}

func PublishPost(topicId int64, content string) (int64, error) {
	fmt.Println("service Publish Post")
	return NewPublishPostFlow(topicId, content).Do()
}

func NewPublishPostFlow(topicId int64, content string) *PublishPostFlow {
	return &PublishPostFlow{
		content: content,
		topicId: topicId,
	}
}

type PublishPostFlow struct {
	content string
	topicId int64
	postId  int64
}

func (f *PublishPostFlow) Do() (int64, error) {
	if err := f.checkParam(); err != nil {
		return 0, err
	}
	if err := f.publish(); err != nil {
		return 0, err
	}
	return f.postId, nil
}

func (f *PublishPostFlow) checkParam() error {
	if len(utf16.Encode([]rune(f.content))) >= 500 {
		return errors.New("content length must be less than 500")
	}
	return nil
}

func (f *PublishPostFlow) publish() error {
	post := &repository.Post{
		ParentId:   f.topicId,
		Content:    f.content,
		CreateTime: time.Now().Unix(),
	}
	id, err := idGen.NextId()
	if err != nil {
		return err
	}
	post.Id = id
	if err := repository.NewPostDaoInstance().InsertPost(post); err != nil {
		return err
	}
	f.postId = post.Id
	return nil
}
