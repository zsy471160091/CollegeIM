package process

import (
	"encoding/json"
	"github.com/donnie4w/go-logger/logger"
	"github.com/donnie4w/json4g"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"protocol"
	"time"
)

type cmd_process map[int](func(msg []byte, sendChan chan<- []byte))

var g_DB_URL string

var g_File_Address string

var CMD_PROCESS cmd_process

var g_clientList *clientList

func InitDB(db_url string) error {
	session, err := mgo.Dial(db_url)
	if err != nil {
		logger.Fatal(err.Error())
		return err
	}
	defer session.Close()
	g_DB_URL = db_url
	return nil
}

func InitFileAddress(address string) {
	g_File_Address = address
}

func init() {
	CMD_PROCESS = cmd_process{
		protocol.IM_LOGIN:            do_login,
		protocol.IM_GET_USER_INFO:    do_get_user_info,
		protocol.IM_GET_USER_FRIENDS: do_get_user_friends,
		protocol.IM_CHAT_P2P:         do_chat_p2p,
		protocol.IM_UPLOAD_FILE:      do_upload_file,
		protocol.IM_MODIFY_PWD:       do_modify_pwd,
		protocol.IM_FIND_USER:        do_find_user,
		protocol.IM_ADD_FRIEND:       do_add_friend,
		protocol.IM_DELETE_FRIEND:    do_delete_friend,
		// IM_USER_STATUS:Do_User_STATUS,
		// IM_USER_INFO:Do
		// IM_CHAT_GROUP
		// IM_PUSH
		// IM_PUSH_REPLY
		// IM_PUSH_GET_REPLY
		// protocol.IM_OFF_LINE: do_offLine,
	}
	g_clientList = newClientList(30000)
}

// 登陆验证
func do_login(msg []byte, sendChan chan<- []byte) {
	// 初始化反馈数据
	rep := protocol.Rep{}
	rep.Cmd = protocol.IM_LOGIN

	// 反馈数据
	defer func() {
		msg, err := json.Marshal(rep)
		if err != nil {
			logger.Warn("编码成json数据时出错, err:", err)
		}
		sendChan <- msg
	}()

	// 从请求的数据中取出用户名和密码（32位小写MD5的值）
	jroot, err := json4g.LoadByString(string(msg))
	if err != nil {
		logger.Warn("json解析失败:", err)
		rep.Ack = "error"
		return
	}
	jnode := jroot.GetNodeByName("id")
	reqId := jnode.ValueString
	jnode = jroot.GetNodeByName("passwd")
	reqPasswd := jnode.ValueString
	logger.Debug("do_login:", "id", reqId, "passwd", reqPasswd)

	// 从数据库中取出相应用户密码的MD5字符串
	type PWD struct {
		Name   string `bson:"name"`
		Passwd string `bson:"passwd"`
	}
	session, err := mgo.Dial(g_DB_URL)
	if err != nil {
		logger.Warn("查询用户", reqId, "密码时连接数据库失败，err:", err)
		rep.Ack = "error"
		return
	}
	c := session.DB("D_USER").C("C_USER_INFO")
	defer session.Close()

	userpwd := PWD{}
	err = c.Find(bson.M{"id": reqId}).One(&userpwd)
	if err != nil {
		logger.Warn("查询用户", reqId, "密码时出错，err:", err)
		rep.Ack = "error"
		rep.Msg = "用户名或密码错误"
		return
	}

	// 比对结果
	if reqPasswd == userpwd.Passwd {
		//修改用户状态为在线
		c = session.DB("D_USER").C("C_USER_STATUS")
		err = c.Update(bson.M{"id": reqId}, bson.M{"$set": bson.M{
			"status": "online",
		}})

		if err != nil {
			logger.Warn("更改用户", reqId, "状态为在线时出错，err:", err)
			rep.Ack = "error"
			return
		}

		rep.Ack = "success"
		rep.Msg = userpwd.Name
		logger.Info("用户", reqId, "登录成功")
		g_clientList.add(reqId, sendChan)
		// go do_send_offline_msg(reqId)
	} else {
		rep.Ack = "error"
		rep.Msg = "用户名或密码错误"
	}
	return
}

// 获取指定用户信息
func do_get_user_info(msg []byte, sendChan chan<- []byte) {
	rep := protocol.Rep{}
	rep.Cmd = protocol.IM_GET_USER_INFO

	defer func() {
		msg, err := json.Marshal(rep)
		if err != nil {
			logger.Warn("编码成json数据时出错, err:", err)
		}
		sendChan <- msg
	}()

	// 获取指定的用户id
	reqId := getId(string(msg), "id")
	logger.Debug("do_get_user_info:", "id", reqId)

	// 从数据库中取出相应数据
	session, err := mgo.Dial(g_DB_URL)
	if err != nil {
		rep.Ack = "error"
		rep.Msg = "Error: 服务器端无法连接到数据库"
		return
	}
	c := session.DB("D_USER").C("C_USER_INFO")

	defer session.Close()

	// 获取指定用户的信息
	type userInfo struct {
		// ID        bson.ObjectId `json:"_id"`
		StuID     string `bson:"id"`
		Name      string `bson:"name"`
		Age       string `bson:"age"`
		Grade     string `bson:"grade"`
		Specialty string `bson:"specialty"`
		Class     string `bson:"class"`
		Identity  string `bson:"identity"`
	}
	userinfo := userInfo{}
	err = c.Find(bson.M{"id": reqId}).One(&userinfo)

	if err != nil {
		rep.Ack = "error"
		rep.Msg = "Error: 没有找到请求的用户"
		return
	} else {
		tmp, _ := json.Marshal(userinfo)
		rep.Ack = "success"
		rep.Msg = string(tmp)
	}

}

// 获取指定用户的常用联系人列表
func do_get_user_friends(msg []byte, sendChan chan<- []byte) {
	rep := protocol.FriendRep{}
	rep.Cmd = protocol.IM_GET_USER_FRIENDS

	defer func() {
		msg, err := json.Marshal(rep)
		if err != nil {
			logger.Warn("编码成json数据时出错, err:", err)
		}
		sendChan <- msg
	}()

	// 获取指定用户id
	reqId := getId(string(msg), "id")
	logger.Debug("do_get_user_info:", "id", reqId)

	// 连接数据库服务器并指定数据库和集合
	session, err := mgo.Dial(g_DB_URL)
	if err != nil {
		rep.Ack = "error"
		return
	}
	c := session.DB("D_USER").C("C_USER_STATUS")
	defer session.Close()

	// 从数据库中取出相应数据
	// type friendItem struct {
	// 	Id   string `bson:"id"`
	// 	Name string `bson:"name"`
	// }

	type friend struct {
		Friends []protocol.FriendItem `bson:"friends"`
	}

	fri := friend{}
	err = c.Find(bson.M{"id": reqId}).One(&fri)

	if err != nil {
		rep.Ack = "error"
		return
	} else {
		rep.Friends = fri.Friends
		rep.Ack = "success"
	}

}

// 点对点聊天
func do_chat_p2p(msg []byte, sendChan chan<- []byte) {
	type chatMsg struct {
		// ID        bson.ObjectId `json:"_id"`
		From_id string `bson:"from_id"`
		To_id   string `bson:"to_id"`
		Msg     string `bson:"msg"`
		Time    string `bson:"time"`
		OffLine string `bson:"offline"`
	}
	chatmsg := chatMsg{}
	// 从请求信息中取出目标用户id
	fromId := getId(string(msg), "from_id")
	toId := getId(string(msg), "to_id")
	logger.Debug("do_get_user_info:", "from_id", fromId, "to_id", toId)

	chatmsg.From_id = fromId
	chatmsg.To_id = toId
	chatmsg.Msg = getId(string(msg), "msg")
	chatmsg.Time = time.Now().Format("2006-01-02 15:04:05")

	// 判断目标用户是否在线
	if g_clientList.contains(toId) {
		logger.Debug("目标用户", toId, "在线")

		jroot, err := json4g.LoadByString(string(msg))
		if err != nil {
			logger.Warn("json解析失败:", err)
		}
		jroot.AddNode(json4g.NowJsonNode("time", chatmsg.Time))

		targetChan := g_clientList.get(toId)
		targetChan <- []byte(jroot.ToString())
		chatmsg.OffLine = "N"
	} else {
		chatmsg.OffLine = "Y"
	}

	// 消息记录插入数据库
	session, err := mgo.Dial(g_DB_URL)
	if err != nil {
		logger.Warn("Error: 服务器端无法连接到数据库")
		return
	}
	c := session.DB("D_" + toId).C("C_ALL_MSG")

	defer session.Close()

	c.Insert(&chatmsg)

	c = session.DB("D_" + fromId).C("C_ALL_MSG")
	chatmsg.OffLine = "N"
	c.Insert(&chatmsg)

	return
}

// 上传文件
func do_upload_file(msg []byte, sendChan chan<- []byte) {

}

// 修改密码
func do_modify_pwd(msg []byte, sendChan chan<- []byte) {
	// 初始化反馈数据
	rep := protocol.Rep{}
	rep.Cmd = protocol.IM_MODIFY_PWD

	// 反馈数据
	defer func() {
		msg, err := json.Marshal(rep)
		if err != nil {
			logger.Warn("编码成json数据时出错, err:", err)
		}
		sendChan <- msg
	}()

	// 从请求的数据中取出新旧密码（32位小写MD5的值）
	jroot, err := json4g.LoadByString(string(msg))
	if err != nil {
		logger.Warn("json解析失败:", err)
		rep.Ack = "error"
		return
	}
	jnode := jroot.GetNodeByName("id")
	reqId := jnode.ValueString
	jnode = jroot.GetNodeByName("new_pwd")
	reqNewPWD := jnode.ValueString
	jnode = jroot.GetNodeByName("old_pwd")
	reqOldPWD := jnode.ValueString
	logger.Debug("do_modify_pwd:", "jnode", reqId, "new_pwd", reqNewPWD, "old_pwd", reqOldPWD)

	// 从数据库中取出相应用户密码的MD5字符串
	type PWD struct {
		Passwd string `bson:"passwd"`
	}
	session, err := mgo.Dial(g_DB_URL)
	if err != nil {
		logger.Warn("查询用户", reqId, "密码时连接数据库失败，err:", err)
		rep.Ack = "error"
		return
	}
	c := session.DB("D_USER").C("C_USER_INFO")
	defer session.Close()

	userpwd := PWD{}
	err = c.Find(bson.M{"id": reqId}).One(&userpwd)
	if err != nil {
		logger.Warn("查询用户", reqId, "密码时出错，err:", err)
		rep.Ack = "error"
		rep.Msg = "修改密码错误"
		return
	}

	// 比对结果
	if reqOldPWD == userpwd.Passwd {
		//修改用户密码

		err = c.Update(bson.M{"id": reqId}, bson.M{"$set": bson.M{
			"passwd": reqNewPWD,
		}})

		if err != nil {
			logger.Warn("更改用户", reqId, "密码时出错，err:", err)
			rep.Ack = "error"
			rep.Msg = "修改密码错误"
			return
		}

		rep.Ack = "success"
		rep.Msg = "密码修改成功"

		// go do_send_offline_msg(reqId)
	} else {
		rep.Ack = "error"
		rep.Msg = "密码错误"
	}
	return
}

// 查找用户
func do_find_user(msg []byte, sendChan chan<- []byte) {
	rep := protocol.UserRep{}
	rep.Cmd = protocol.IM_FIND_USER

	defer func() {
		msg, err := json.Marshal(rep)
		if err != nil {
			logger.Warn("编码成json数据时出错, err:", err)
		}
		sendChan <- msg
	}()

	// 获取请求信息
	reqId := getId(string(msg), "id")
	reqName := getId(string(msg), "name")
	logger.Debug("do_find_user:", "id", reqId, "name", reqName)

	// 连接数据库服务器并指定数据库和集合
	session, err := mgo.Dial(g_DB_URL)
	if err != nil {
		rep.Ack = "error"
		return
	}
	c := session.DB("D_USER").C("C_USER_INFO")
	defer session.Close()

	// 从数据库中取出相应数据
	fri := []protocol.FriendItem{}
	if reqId != "" {
		err = c.Find(bson.M{"id": reqId}).All(&fri)
	} else {
		err = c.Find(bson.M{"name": reqName}).All(&fri)
	}

	if len(fri) != 0 {
		rep.Ack = "success"
		rep.Users = fri
	} else {
		rep.Ack = "error"
	}
}

// 添加常用联系人
func do_add_friend(msg []byte, sendChan chan<- []byte) {
	rep := protocol.AddFriRep{}
	rep.Cmd = protocol.IM_ADD_FRIEND

	defer func() {
		msg, err := json.Marshal(rep)
		if err != nil {
			logger.Warn("编码成json数据时出错, err:", err)
		}
		sendChan <- msg
	}()

	// 获取请求信息
	reqId := getId(string(msg), "id")
	reqFriId := getId(string(msg), "fri_id")
	reqFriName := getId(string(msg), "fri_name")
	logger.Debug("do_add_friend:", "id", reqId, "fri_id", reqFriId, "fri_name", reqFriName)

	// 连接数据库服务器并指定数据库和集合
	session, err := mgo.Dial(g_DB_URL)
	if err != nil {
		rep.Ack = "error"
		return
	}
	c := session.DB("D_USER").C("C_USER_STATUS")
	defer session.Close()

	// 在数据库中修改相应数据
	err = c.Update(bson.M{"id": reqId},
		bson.M{"$push": bson.M{
			"friends": bson.M{"id": reqFriId, "name": reqFriName},
		}})
	if err != nil {
		rep.Ack = "error"
	} else {
		rep.Ack = "success"
		rep.Id = reqFriId
		rep.Name = reqFriName
	}

}

// 删除常用联系人
func do_delete_friend(msg []byte, sendChan chan<- []byte) {
	rep := protocol.AddFriRep{}
	rep.Cmd = protocol.IM_DELETE_FRIEND

	defer func() {
		msg, err := json.Marshal(rep)
		if err != nil {
			logger.Warn("编码成json数据时出错, err:", err)
		}
		sendChan <- msg
	}()

	// 获取请求信息
	reqId := getId(string(msg), "id")
	reqFriId := getId(string(msg), "fri_id")
	reqFriName := getId(string(msg), "fri_name")
	logger.Debug("do_delete_friend:", "id", reqId, "fri_id", reqFriId, "fri_name", reqFriName)

	// 连接数据库服务器并指定数据库和集合
	session, err := mgo.Dial(g_DB_URL)
	if err != nil {
		rep.Ack = "error"
		return
	}
	c := session.DB("D_USER").C("C_USER_STATUS")
	defer session.Close()

	// 在数据库中修改相应数据
	err = c.Update(bson.M{"id": reqId},
		bson.M{"$pull": bson.M{
			"friends": bson.M{"id": reqFriId, "name": reqFriName},
		}})
	if err != nil {
		rep.Ack = "error"
	} else {
		rep.Ack = "success"
		rep.Id = reqFriId
		rep.Name = reqFriName
	}

}

// 推送离线信息
func do_send_offline_msg(id string) {
	type chatMsg struct {
		ID      bson.ObjectId `bson:"_id"`
		From_id string        `bson:"from_id"`
		To_id   string        `bson:"to_id"`
		Msg     string        `bson:"msg"`
		Name    string        `bson:"name"`
		Time    string        `bson:"time"`
		OffLine string        `bson:"offline"`
	}

	//
	type chatRep struct {
		Cmd     int    `json:"cmd"`
		From_id string `json:"from_id"`
		To_id   string `json:"to_id"`
		Msg     string `json:"msg"`
		Name    string `json:"name"`
		Time    string `json:"time"`
	}
	chatrep := chatRep{
		Cmd: protocol.IM_CHAT_P2P,
	}
	chatmsg := []chatMsg{}

	// 查询离线消息
	session, err := mgo.Dial(g_DB_URL)
	if err != nil {
		logger.Warn("Error: 服务器端无法连接到数据库")
		return
	}
	c := session.DB("D_" + id).C("C_ALL_MSG")

	defer session.Close()

	c.Find(bson.M{"offline": "Y"}).All(&chatmsg)

	targetChan := g_clientList.get(id)
	for _, v := range chatmsg {
		chatrep.From_id = v.From_id
		chatrep.To_id = v.To_id
		chatrep.Name = v.Name
		chatrep.Time = v.Time
		chatrep.Msg = v.Msg

		msg, err := json.Marshal(chatrep)
		if err != nil {
			logger.Warn("编码成json数据时出错, err:", err)
		}
		targetChan <- msg
	}
	// c = session.DB("D_" + fromId).C("C_ALL_MSG")
	// chatmsg.OffLine = "N"
	// c.Insert(&chatmsg)

	return
}

// 用户离线
func OffLine(reqId string) {
	//从在线用户列表中删除离线用户
	g_clientList.remove(reqId)
	//修改用户状态为离线
	session, err := mgo.Dial(g_DB_URL)
	if err != nil {
		logger.Warn("查询用户", reqId, "密码时连接数据库失败，err:", err)
		return
	}
	c := session.DB("D_USER").C("C_USER_STATUS")
	c.Update(bson.M{"id": reqId}, bson.M{"$set": bson.M{
		"status": "offline",
	}})

	defer session.Close()
}

// 打印用户列表
func ShowClientList() {
	logger.Debug("在线用户数", g_clientList.len())
}

// 从json请求中获取id
func getId(msg, nodeName string) string {
	jroot, err := json4g.LoadByString(string(msg))
	if err != nil {
		logger.Warn("json解析失败:", err)
		return "nil"
	}
	jnode := jroot.GetNodeByName(nodeName)
	return jnode.ValueString
}
