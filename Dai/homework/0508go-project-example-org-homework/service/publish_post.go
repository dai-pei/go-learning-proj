package service

import (
	"errors"
	"example-proj/repository"
	"fmt"
	"time"
	"unicode/utf16"
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

// Init函数一般是先于main函数执行，所以任何文件内定义的init函数都会被执行 初始化顺序：变量初始化->init()->main()
// func init() {
// 	idGen = &idworker.IdWorker{}
// 	idGen.InitIdWorker(1, 1)
// }

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
	id, err := generateIdBySnowFlake(100, 100)

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

var lastTimeStamp int64 = 0
var curTimeStamp int64 = lastTimeStamp
var sn int64 = 0

// snowflake算法
// https://blog.csdn.net/fly910905/article/details/82054196
func generateIdBySnowFlake(machineId int64, datacenterId int64) (int64, error) {
	// 如果想让时间戳范围更长，也可以减去一个日期
	curTimeStamp := time.Now().UnixNano() / 1000000

	if curTimeStamp == lastTimeStamp {
		// 2的12次方 -1 = 4095，每毫秒可产生4095个ID
		if sn > 4095 {
			time.Sleep(time.Millisecond)
			curTimeStamp = time.Now().UnixNano() / 1000000
			sn = 0
		}
	} else {
		sn = 0
	}
	sn++
	lastTimeStamp = curTimeStamp
	// 应为时间戳后面有22位，所以向左移动22位
	curTimeStamp = curTimeStamp << 22
	machineId = machineId << 17
	datacenterId = datacenterId << 12
	// 通过与运算把各个部位连接在一起
	return int64(curTimeStamp) | machineId | datacenterId | sn, nil
}
