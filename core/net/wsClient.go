package net

import (
	"encoding/json"
	"golang.org/x/net/websocket"
	"io"
	"log"
	"time"
)

type WsClient struct {
	url string
	ws  *websocket.Conn
}

func (w *WsClient) Init(url string) {
	w.url = url
	w.connect()
}

func (w *WsClient) connect() {
	if w.ws != nil {
		err := w.ws.Close()
		if err != nil {
			return
		}
	}
	origin := "http://localhost/"
	ws, err := websocket.Dial(w.url, "", origin)
	if err != nil {
		log.Println("ws建立失败", err)
		time.AfterFunc(360*time.Second, w.connect)
	} else {
		w.ws = ws
	}
}

func (w *WsClient) Start() {
	go w.ReceiveMsg()
	go w.KeepAlive()
}

func (w *WsClient) ReceiveMsg() {
	if w.ws != nil {
		defer func(ws *websocket.Conn) {
			err := ws.Close()
			if err != nil {
				log.Println("ws close 异常", err)
			}
		}(w.ws)
		for {
			var msg string
			if err := websocket.Message.Receive(w.ws, &msg); err != nil {
				if err == io.EOF {
					log.Println("ws 被服务的关闭,重新连接")
					w.connect()
				}
			} else {
				log.Printf("接收到消息: %s\n", msg)
			}
		}
	}

}

func (w *WsClient) SendMsg(msg string) {

	if w.ws != nil {
		_, err := w.ws.Write([]byte(msg))
		if err != nil {
			log.Println("ws 发送msg失败", msg)
			return
		} else {
			log.Println("ws 发送msg", msg)
		}
	}
}

type HeartbeatMessage struct {
	Command string      `json:"command"`
	Content interface{} `json:"content"`
	From    string      `json:"from"`
	To      string      `json:"to"`
	Biz     string      `json:"biz"`
}

func (w *WsClient) KeepAlive() {
	ticker := time.NewTicker(time.Duration(10) * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			heartbeat := HeartbeatMessage{
				Command: "heartbeat-ping",
				Content: struct{}{},
				From:    "30",
				To:      "0",
				Biz:     "rt-service-interactive",
			}
			message, err := json.Marshal(heartbeat)
			if err != nil {
				log.Printf("Failed to marshal heartbeat message: %v\n", err)
				continue
			}
			w.SendMsg(string(message))

		}
	}
}
