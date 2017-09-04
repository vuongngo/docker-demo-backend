package ws_test

import (
	"app/ws"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	gWs "github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func TestWs(t *testing.T) {
	assert := assert.New(t)
	var err error
	//ws hub
	hub := ws.NewHub()
	go hub.Run()

	websocket := &ws.WebSocket{
		Hub: hub,
	}

	// Serve websocket
	r := mux.NewRouter()
	r.HandleFunc("/ws", websocket.ServeWs).Methods("GET")
	http.Handle("/", r)

	ts := httptest.NewServer(r)
	defer ts.Close()
	d := gWs.Dialer{}

	cookie := &http.Cookie{
		Name:  "token",
		Value: "123",
	}
	header := http.Header{}
	header.Set("Cookie", cookie.String())

	_, _, err = d.Dial("ws://"+ts.Listener.Addr().String()+"/ws", header)
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(time.Second * 3)
	assert.Equal(1, hub.GetClientCount(), "Number of client")
}
