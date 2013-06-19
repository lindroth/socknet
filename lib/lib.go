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

	go func() {
		for mess := range input {
			if _, err := ws.Write([]byte(mess)); err != nil {
				log.Fatal(err)
			}
		}
	}()

	go func() {
		var msg = make([]byte, 512)
		var n int

		for n, err = ws.Read(msg); err == nil; n, err = ws.Read(msg) {
			output <- string(msg[:n])

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
