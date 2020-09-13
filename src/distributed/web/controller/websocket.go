package controller

import (
	"bytes"
	"encoding/gob"
	"github.com/gorilla/websocket"
	"github.com/streadway/amqp"
	"go-distributed/src/distributed/dto"
	"go-distributed/src/distributed/qutils"
	"go-distributed/src/distributed/web/model"
	"net/http"
	"sync"
)

const url = "amqp://guest:guest@localhost:5672"

type websocketController struct {
	conn     *amqp.Connection
	ch       *amqp.Channel
	sockets  []*websocket.Conn
	mutex    sync.Mutex
	upgrader websocket.Upgrader
}

func newWebsocketController() *websocketController {
	wsc := new(websocketController)

	wsc.conn, wsc.ch = qutils.GetChannel(url)
	wsc.upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	go wsc.listenForSources()
	go wsc.listenForMessages()

	return wsc
}

func (wsc *websocketController) handleMessage(w http.ResponseWriter, r *http.Request) {
	socket, _ := wsc.upgrader.Upgrade(w, r, nil)

	wsc.addSocket(socket)
	go wsc.listenForDiscoveryRequests(socket)
}

func (wsc *websocketController) addSocket(socket *websocket.Conn) {
	wsc.mutex.Lock()
	defer wsc.mutex.Unlock()

	wsc.sockets = append(wsc.sockets, socket)
}

func (wsc *websocketController) removeSocket(socket *websocket.Conn) {
	wsc.mutex.Lock()
	defer wsc.mutex.Unlock()
	_ = socket.Close()

	for i := range wsc.sockets {
		if wsc.sockets[i] == socket {
			wsc.sockets = append(wsc.sockets[:i], wsc.sockets[i+1:]...)
		}
	}

}

func (wsc *websocketController) sendMessage(msg message) {
	wsc.mutex.Lock()
	defer wsc.mutex.Unlock()

	socketsToRemove := make([]*websocket.Conn, 0)
	for _, socket := range wsc.sockets {
		err := socket.WriteJSON(msg)

		if err != nil {
			socketsToRemove = append(socketsToRemove, socket)
		}
	}

	for _, socket := range socketsToRemove {
		wsc.removeSocket(socket)
	}
}

func (wsc *websocketController) listenForDiscoveryRequests(socket *websocket.Conn) {
	for {
		msg := message{}
		err := socket.ReadJSON(&msg)

		if err != nil {
			wsc.removeSocket(socket)
			break
		}

		if msg.Type == "discover" {
			_ = wsc.ch.Publish(
				"",
				qutils.WebappDiscoveryQueue,
				false,
				false,
				amqp.Publishing{},
			)
		}
	}
}

func (wsc *websocketController) listenForSources() {
	q := qutils.GetQueue(
		"",
		wsc.ch,
		true,
	)
	_ = wsc.ch.QueueBind(
		q.Name,
		"",
		qutils.WebappSourceExchange,
		false,
		nil,
	)

	msgs, _ := wsc.ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	for msg := range msgs {
		sensor, _ := model.GetSensorByName(string(msg.Body))
		wsc.sendMessage(message{
			Type: "source",
			Data: sensor,
		})
	}
}

func (wsc *websocketController) listenForMessages() {
	q := qutils.GetQueue(
		"",
		wsc.ch,
		true,
	)
	_ = wsc.ch.QueueBind(
		q.Name,
		"",
		qutils.WebappReadingsExchange,
		false,
		nil,
	)

	msgs, _ := wsc.ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	for msg := range msgs {
		buf := bytes.NewBuffer(msg.Body)
		dec := gob.NewDecoder(buf)
		sm := dto.SensorMessage{}

		_ = dec.Decode(&sm)
		wsc.sendMessage(message{
			Type: "reading",
			Data: sm,
		})
	}
}

type message struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}
