package bot

import (
	"encoding/json"
	"errors"

	"github.com/BaiMeow/SimpleBot/message"
	"github.com/google/uuid"
)

// ErrJsonUnmarshal json序列化中出错
var ErrJsonUnmarshal = errors.New("JsonUnmarshalError")

var waitReply = make(map[uuid.UUID]func([]byte, bool))

type preUnmarshalReply struct {
	Echo    uuid.UUID `json:"echo"`
	RetCode int       `json:"retcode"`
	Status  string    `json:"status"`
}

type apiCallFramework struct {
	Action string      `json:"action"`
	Params interface{} `json:"params"`
	Echo   uuid.UUID   `json:"echo"`
}

type groupMsg struct {
	GroupID    int64       `json:"group_id"`
	Message    interface{} `json:"message"`
	AutoEscape bool        `json:"auto_escape"`
}

type groupMsgReplyDetails struct {
	Data struct {
		MessageID int32 `json:"message_id"`
	} `json:"data"`
}

type privateMsg struct {
	UserID     int64       `json:"user_id"`
	Message    interface{} `json:"message"`
	AutoEscape bool        `json:"auto_escape"`
}

type privateMsgReplyDetails struct {
	MessageID int32 `json:"message_id"`
}

func handleAPIReply(data []byte) {
	reply := new(preUnmarshalReply)
	if waitReply[reply.Echo] != nil {
		waitReply[reply.Echo](data, reply.Status == "ok")
		delete(waitReply, reply.Echo)
	}
}

//SendGroupMsg 发送群聊消息(不含匿名消息)
func (b *Bot) SendGroupMsg(group int64, msg *message.Msg) (int32, error) {
	id := uuid.New()
	bytes, err := json.Marshal(&apiCallFramework{
		Action: "send_group_msg",
		Params: groupMsg{
			GroupID:    group,
			Message:    msg.ToArrayMessage(),
			AutoEscape: false,
		},
		Echo: id,
	})
	if err != nil {
		return 0, err
	}
	b.driver.Write(bytes)
	msgID := make(chan int32, 1)
	waitReply[id] = func(data []byte, ok bool) {
		if !ok {
			msgID <- 0
			return
		}
		details := new(groupMsgReplyDetails)
		if err := json.Unmarshal(data, &details); err != nil {
			msgID <- 0
			return
		}
		msgID <- details.Data.MessageID
	}
	recMsgID := <-msgID
	if recMsgID == 0 {
		return 0, ErrJsonUnmarshal
	}
	return recMsgID, nil
}

//SendPrivateMsg 发送私聊消息
func (b *Bot) SendPrivateMsg(qq int64, msg *message.Msg) (int32, error) {
	id := uuid.New()
	bytes, err := json.Marshal(&apiCallFramework{
		Action: "send_private_msg",
		Params: privateMsg{
			UserID:     qq,
			Message:    msg.ToArrayMessage(),
			AutoEscape: false,
		},
		Echo: id,
	})
	if err != nil {
		return 0, err
	}
	b.driver.Write(bytes)
	msgID := make(chan int32, 1)
	waitReply[id] = func(data []byte, ok bool) {
		if !ok {
			msgID <- 0
		}
		details := new(privateMsgReplyDetails)
		if err := json.Unmarshal(data, details); err != nil {
			msgID <- 0
		}
		msgID <- details.MessageID
	}
	recMsgID := <-msgID
	if recMsgID == 0 {
		return 0, ErrJsonUnmarshal
	}
	return recMsgID, nil
}
