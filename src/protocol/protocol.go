package protocol

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

	IM_UPLOAD_FILE
	IM_GET_USER_LIST
	IM_DOWNLOND_FILE

	// 查询
	IM_USER_STATUS
	IM_USER_INFO

	// 聊天
	IM_CHAT_P2P_REP
	IM_CHAT_GROUP

	// 推送
	IM_PUSH
	IM_PUSH_REPLY
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
