package repository

import (
	"encoding/json"
	"os"
	"sync"
)

type Post struct {
	Id         int64  `json:"id"`
	ParentId   int64  `json:"parent_id"`
	Content    string `json:"content"`
	CreateTime int64  `json:"create_time"`
}
type PostDao struct {
}

var (
	postDao  *PostDao
	postOnce sync.Once
)

// 只需要实例化一个Dao对象，本项目中每次调用了NewPostDaoInstance函数
// 如果没有创建过PostDao对象，则会创建，如果已经创建了，则直接返回它
// 这里是单例模式的一个应用
func NewPostDaoInstance() *PostDao {
	postOnce.Do(
		func() {
			postDao = &PostDao{}
		})
	return postDao
}
func (*PostDao) QueryPostsByParentId(parentId int64) []*Post {
	return postIndexMap[parentId]
}

// 注意这两个变量位置为db_init.go
// var (
// 	topicIndexMap map[int64]*Topic
// 	postIndexMap  map[int64][]*Post
// )

func (*PostDao) InsertPost(post *Post) error {
	f, err := os.OpenFile("./data/post", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}

	defer f.Close()
	marshal, _ := json.Marshal(post)
	if _, err = f.WriteString(string(marshal) + "\n"); err != nil {
		return err
	}
	// map需要读写锁添加，但是写入file不需要是吗
	rwMutex.Lock()
	postList, ok := postIndexMap[post.ParentId]
	if !ok {
		postIndexMap[post.ParentId] = []*Post{post}
	} else {
		postList = append(postList, post)
		postIndexMap[post.ParentId] = postList
	}
	rwMutex.Unlock()
	return nil
}
