package lib

import (
	"code.google.com/p/go.net/websocket"
	"log"
	"net/http"
	"net/url"
)

type Socknet struct {
}

func (self *Socknet) Connect(origin string, location string, header http.Header) (input chan string, output chan string, err error) {

	config := &websocket.Config{
		Location: parseUrl(location),
		Origin:   parseUrl(origin),
		Version:  13,
		Header:   header,
	}

	var ws *websocket.Conn
	if ws, err = websocket.DialConfig(config); err != nil {
		return
	}
	input = make(chan string)
	output = make(chan string)

	closer := func() {
		defer func() {
			recover()
		}()
		close(output)
		ws.Close()
	}

	go func() {
		defer closer()
		for mess := range input {
			if err := websocket.Message.Send(ws, mess); err != nil {
				log.Fatal(err)
				break
			}
		}
	}()

	go func() {
		defer closer()
		var msg string
		for err := websocket.Message.Receive(ws, &msg); err == nil; err = websocket.Message.Receive(ws, &msg) {
			output <- msg
		}
	}()

	return input, output, nil

}

func parseUrl(location string) *url.URL {
	locationUrl, err := url.Parse(location)
	if err != nil {
		log.Fatal(err)
	}
	return locationUrl
}
