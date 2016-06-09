package protocol

import (
	"gopkg.in/mgo.v2/bson"
)

// 请求的消息类型
const (
	// 登陆
	IM_LOGIN = iota
	IM_GET_USER_INFO
	IM_GET_USER_FRIENDS
	IM_CHAT_P2P
	IM_GET_OFFLINE_MSG
	IM_MODIFY_PWD
	IM_FIND_USER
	IM_ADD_FRIEND
	IM_DELETE_FRIEND
	IM_GET_GROUP_LIST
	IM_GET_GROUP_USERS
	IM_CHAT_GROUP
	IM_PUSH
	IM_PUSH_REP
	IM_GET_NOTIFICATION_LIST
	IM_GET_NOTIFICATION

	IM_UPLOAD_FILE
	IM_GET_USER_LIST
	IM_DOWNLOND_FILE

	// 查询
	IM_USER_STATUS
	IM_USER_INFO

	IM_PUSH_GET_REPLY

	// 登出
	IM_EXIT
)

// 通用反馈数据结构
type Rep struct {
	Cmd int    `json:"cmd"`
	Ack string `json:"ack"`
	Msg string `json:"msg"`
}

// 获取常用联系人列表
type FriendItem struct {
	Id   string `bson:"id"`
	Name string `bson:"name"`
}

type FriendRep struct {
	Cmd     int          `json:"cmd"`
	Ack     string       `json:"ack"`
	Friends []FriendItem `json:"friends"`
}

// 查找用户
type UserRep struct {
	Cmd   int          `json:"cmd"`
	Ack   string       `json:"ack"`
	Users []FriendItem `json:"users"`
}

// 添加常用联系人
type AddFriRep struct {
	Cmd  int    `json:"cmd"`
	Ack  string `json:"ack"`
	Id   string `json:"id"`
	Name string `json:"name"`
}

// 删除常用联系人
type DelFriRep struct {
	Cmd  int    `json:"cmd"`
	Ack  string `json:"ack"`
	Id   string `json:"id"`
	Name string `json:"name"`
}

// 获取指定用户的群列表
type GetGroupListItem struct {
	Id   string `bson:"id"`
	Name string `bson:"name"`
}
type GetGroupListRep struct {
	Cmd    int                `json:"cmd"`
	Ack    string             `json:"ack"`
	Groups []GetGroupListItem `json:"groups"`
}

// 获取指定群的所有成员
type GetGroupUserItem struct {
	Id   string `bson:"id"`
	Name string `bson:"name"`
}
type GetGroupUsersRep struct {
	Cmd   int                `json:"cmd"`
	Ack   string             `json:"ack"`
	Id    string             `json:"id"`
	Users []GetGroupUserItem `json:"users"`
}

// 推送消息结构
type PushMsg struct {
	ID        bson.ObjectId `bson:"_id"`
	Fromid    string        `bson:"id"`
	Fromname  string        `bson:"name"`
	School    string        `bson:"school"`
	Specialty string        `bson:"specialty"`
	Grade     string        `bson:"grade"`
	Class     string        `bson:"class"`
	Title     string        `bson:"title"`
	Content   string        `bson:"content"`
	Time      string        `bson:"time"`
}
type NotificationMsg struct {
	Fromid   string `json:"id"`
	Fromname string `json:"name"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	Time     string `json:"time"`
}
