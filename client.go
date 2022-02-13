package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	Send chan []byte
	UID  string
	mm   *MessageManager
	om   *ObjectManager
	tm   *TokenManager
	conn *websocket.Conn
}

func (c *Client) readPump() {
	defer func() {
		c.mm.Logout <- c
		c.tm.ReleaseTokenWithUsername(c.UID)
		c.om.CancelDraggingWithUID(c.UID)
		c.conn.Close()
	}()
	log.Printf("client %v: start listening\n", c.UID)
	for {
		messageType, message, err := c.conn.ReadMessage()
		log.Println(string(message))
		if messageType == websocket.CloseMessage {
			log.Printf("%v left\n", c.UID)
			break
		}
		if messageType != websocket.TextMessage {
			log.Printf("%v message is not text\n", c.UID)
			break
		}
		if err != nil {
			log.Println(err)
			break
		}

		// parse the message
		var env Envelope
		err2 := json.Unmarshal(message, &env)
		if err2 != nil || env.Type == "" {
			continue
		}
		ok := c.tm.ValidateToken(env.Token, c.UID) && c.UID == env.UID
		if !ok {
			return
		}
		log.Println(env)
		switch env.Type {
		case "drag_start":
			{
				var req struct {
					Envelope
					DragStartRequest
				}
				err := json.Unmarshal(message, &req)
				if err != nil {
					log.Println(err)
					continue
				}
				err2 := c.om.StartDragging(c.UID, req.CID)
				if err2 != nil {
					tosend, err := json.Marshal(DragStartFailResponse{
						Status:  "err",
						Message: "drag fail",
					})
					if err != nil {
						log.Panic(err)
					}
					c.Send <- tosend
					continue
				}
				tosend, err := json.Marshal(DragStartResponse{
					Status: "ok",
					CID:    req.CID,
					UID:    req.UID,
					Event:  "drag_start",
				})
				if err != nil {
					log.Panic(err)
				}
				c.mm.Broadcast <- tosend
			}
		case "drag_cancel":
			{

			}
		case "drag_finish":
			{
				var req struct {
					Envelope
					DragFinishRequest
				}
				err := json.Unmarshal(message, &req)
				if err != nil {
					continue
				}
				err2 := c.om.FinishDragging(c.UID, req.CID, req.Pos)
				if err2 != nil {
					log.Println(err2)
					tosend, err := json.Marshal(DragStartFailResponse{
						Status:  "err",
						Message: "drag finish fail",
					})
					if err != nil {
						log.Panic(err)
					}
					c.Send <- tosend
					continue
				}
				tosend, err := json.Marshal(DragFinishResponse{
					Status: "ok",
					CID:    req.CID,
					UID:    req.UID,
					Event:  "drag_finish",
					Pos:    req.Pos,
				})
				if err != nil {
					log.Panic(err)
				}
				c.mm.Broadcast <- tosend
			}
		}
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(time.Minute * 2)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case tosend := <-c.Send:
			c.conn.WriteMessage(websocket.TextMessage, tosend)
			log.Printf("send(%v): %v", c.UID, string(tosend))
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(time.Second * 30))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
