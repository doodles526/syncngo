package main

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"net"
	"sync"
	"time"
)

type Client struct {
	serverConn net.Conn

	config *ClientArgs

	handlers map[MessageType]MessageHandler
	handlerM sync.Mutex
}

type MessageHandler func(payload interface{})

type ClientArgs struct {
	Host     string
	Port     string
	Username string
	Password string
	Room     string
	Version  string
}

func NewClient(args *ClientArgs) (*Client, error) {
	c := &Client{
		config:   args,
		handlers: make(map[MessageType]MessageHandler),
	}

	sc, err := net.Dial("tcp", net.JoinHostPort(args.Host, args.Port))
	if err != nil {
		return nil, err
	}
	c.serverConn = sc

	go c.handleMessages()

	helloC, hFunc, errFunc := c.helloHandlerGen()
	c.registerHandler("Hello", hFunc)
	c.registerHandler("Error", errFunc)

	if err := c.negotiateHello(); err != nil {
		return nil, err
	}

	_, ok := <-helloC
	if !ok {
		return nil, fmt.Errorf("Failed Hello")
	}
	c.unregisterHandler("Hello")
	c.unregisterHandler("Error")

	go c.heartbeat()

	return c, err
}

func (c *Client) handleMessages() error {
	connbuf := bufio.NewReader(c.serverConn)
	for {
		msgData, err := connbuf.ReadBytes(0x0a)
		fmt.Println(msgData)
		if err != nil {
			return err
		}

		msgType, msg, err := UnmarshalMessage(msgData)
		if err != nil {
			return err
		}

		c.handlerM.Lock()

		if handler, ok := c.handlers[msgType]; ok {
			handler(msg)
		}

		c.handlerM.Unlock()
	}
	return nil
}

func (c *Client) registerHandler(name MessageType, handler MessageHandler) {
	c.handlerM.Lock()
	defer c.handlerM.Unlock()

	c.handlers[name] = handler
}

func (c *Client) unregisterHandler(name MessageType) {
	c.handlerM.Lock()
	defer c.handlerM.Unlock()

	delete(c.handlers, name)
}

func (c *Client) negotiateHello() error {
	h := HelloMsg{
		Username: c.config.Username,
		Password: fmt.Sprintf("%x", md5.Sum([]byte(c.config.Password))),
		Room:     Room{Name: c.config.Room},
		Version:  c.config.Version,
	}

	if err := c.sendMessage(h); err != nil {
		return err
	}
	return nil
}

func (c *Client) helloHandlerGen() (<-chan *HelloMsg, MessageHandler, MessageHandler) {
	resultCh := make(chan *HelloMsg)
	errHandler := func(payload interface{}) {
		close(resultCh)
	}
	helloHandler := func(payload interface{}) {
		helloMsg := payload.(*HelloMsg)
		resultCh <- helloMsg
	}
	return resultCh, helloHandler, errHandler
}

func (c *Client) sendMessage(msg interface{}) error {
	msgData, err := MarshalMessage(msg)
	if err != nil {
		return err
	}
	_, err = c.serverConn.Write(msgData)
	_, err = c.serverConn.Write([]byte{0x0d, 0x0a})
	return err
}

func (c *Client) heartbeat() error {
	for _ = range time.Tick(1 * time.Second) {
		c.sendMessage(StateMsg{})
		c.sendMessage(ClientChatMsg("Kendall +7 - Makeout"))
	}

	return nil
}
