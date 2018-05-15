package main

import (
	"encoding/json"
	"fmt"
)

type MessageType string

const (
	HelloMsgType = "Hello"
	StateMsgType = "State"
	ChatMsgType  = "Chat"
	ErrorMsgType = "Error"
)

type Message map[MessageType]interface{}

func (m *Message) UnmarshalJSON(b []byte) error {
	msgSub := make(map[MessageType]interface{})
	if err := json.Unmarshal(b, &msgSub); err != nil {
		return err
	}

	for key, val := range msgSub {
		data, err := json.Marshal(val)
		if err != nil {
			return err
		}
		msgSub[key] = data
	}
	*m = msgSub
	return nil
}

type ErrorMsg struct {
	Message string `json:"username"`
}

type HelloMsg struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Room     Room   `json:"room"`
	Version  string `json:"version"`
}

type Room struct {
	Name string `json:"name"`
}

type ClientChatMsg string

type ServerChatMsg struct {
	Username string `json:"username"`
	Message  string `json:"message"`
}

type StateMsg struct{}

func MarshalMessage(msg interface{}) ([]byte, error) {
	retMsg := make(map[MessageType]interface{})

	switch msg.(type) {
	case HelloMsg:
		retMsg[HelloMsgType] = msg
	case ClientChatMsg:
		retMsg[ChatMsgType] = msg
	case ServerChatMsg:
		retMsg[ChatMsgType] = msg
	case StateMsg:
		retMsg[StateMsgType] = msg
	default:
		return nil, fmt.Errorf("Message type not recognized")
	}
	return json.Marshal(retMsg)
}

func UnmarshalMessage(data []byte) (MessageType, interface{}, error) {
	var msg Message
	var retMsg interface{}
	if err := json.Unmarshal(data, &msg); err != nil {
		return "", nil, err
	}

	for msgType, payload := range msg {

		switch msgType {
		case HelloMsgType:
			retMsg = &HelloMsg{}
		case ChatMsgType:
			retMsg = &ServerChatMsg{}
		case StateMsgType:
			retMsg = &StateMsg{}
		default:
			return msgType, nil, nil
		}
		if err := json.Unmarshal(payload.([]byte), retMsg); err != nil {
			return "", nil, err
		}
		return msgType, retMsg, nil
	}
	return "", nil, fmt.Errorf("bad message")
}
