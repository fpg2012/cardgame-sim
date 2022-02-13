package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func genRootHandler(upgrader *websocket.Upgrader, tm *TokenManager, mm *MessageManager, om *ObjectManager) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		if messageType != websocket.TextMessage {
			log.Println("login message must be text message")
			return
		}
		var s LoginRequest
		err2 := json.Unmarshal(p, &s)
		if err2 != nil {
			return
		}
		token, err := tm.AllocateToken(s.UID)
		if err != nil {
			conn.WriteJSON(LoginResponse{
				Status:  "err",
				Message: "token allocation failed",
			})
		} else {
			conn.WriteJSON(LoginResponse{
				Status:  "ok",
				Token:   fmt.Sprint(token),
				Message: "welcome!",
			})
		}
		client := Client{
			Send: make(chan []byte),
			UID:  s.UID,
			mm:   mm,
			om:   om,
			tm:   tm,
			conn: conn,
		}
		mm.Login <- &client
		log.Printf("%v login, token: %v", s.UID, token)
		go client.readPump()
		go client.writePump()
	})
}

func genGetTokenHandler(tm *TokenManager) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		username := r.URL.Query().Get("username")
		if username == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		token, err := tm.AllocateToken(username)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write([]byte(fmt.Sprint(token)))
	})
}

func main() {
	log.Println("Hello!")

	tm := NewTokenManager()
	mm := NewMessageManager()
	om := NewObjectManager()

	go mm.run()

	http.Handle("/", genRootHandler(
		&websocket.Upgrader{
			ReadBufferSize: 1048576, WriteBufferSize: 1048576,
		},
		tm,
		mm,
		om,
	))
	http.Handle("/token", genGetTokenHandler(tm))
	http.ListenAndServe("0.0.0.0:6673", nil)
}
