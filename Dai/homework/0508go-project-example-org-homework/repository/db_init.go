package repository

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

// 批量声明
// 使用var关键字和小括号()可以将一组变量定义放在一起，用于批量声明。

// var (
//   id int
//   name string
//   score []float32
//   job struct {
//     post string
//   }
//   lock func() bool
// )

var (
	topicIndexMap map[int64]*Topic
	postIndexMap  map[int64][]*Post
	rwMutex       sync.RWMutex
)

func Init(filePath string) error {
	if err := initTopicIndexMap(filePath); err != nil {
		return err
	}
	if err := initPostIndexMap(filePath); err != nil {
		return err
	}
	return nil
}

func initTopicIndexMap(filePath string) error {
	open, err := os.Open(filePath + "topic")
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(open)
	topicTmpMap := make(map[int64]*Topic)
	for scanner.Scan() {
		text := scanner.Text()
		var topic Topic
		if err := json.Unmarshal([]byte(text), &topic); err != nil {
			return err
		}
		topicTmpMap[topic.Id] = &topic
	}
	topicIndexMap = topicTmpMap
	// printMap(topicIndexMap)
	return nil
}
func printMap(topicIndexMap map[int64]*Topic) {
	for _, v := range topicIndexMap {
		fmt.Println(v)
	}
}

func initPostIndexMap(filePath string) error {
	open, err := os.Open(filePath + "post")
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(open)
	postTmpMap := make(map[int64][]*Post)
	for scanner.Scan() {
		text := scanner.Text()
		var post Post
		if err := json.Unmarshal([]byte(text), &post); err != nil {
			return err
		}
		posts, ok := postTmpMap[post.ParentId]
		if !ok {
			postTmpMap[post.ParentId] = []*Post{&post}
			continue
		}
		posts = append(posts, &post)
		postTmpMap[post.ParentId] = posts
	}
	postIndexMap = postTmpMap
	return nil
}
