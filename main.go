package main

import (
	"bytes"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{} // use default options

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func process(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	vp := &VirtualProcess{
		stderr: bytes.Buffer{},
		proc:   nil,
	}
	vp.Start()
	currentMessage := make([]byte, 1024)
	cmLength := 0
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		cmLength += len(message)
		// currentMessage[cmLength:len(message)]message
		currentMessage = append(currentMessage[cmLength:], message...)
		vp.stdin.Write(currentMessage)

		result := make([]byte, 1024)
		i, e := vp.Read(result)
		log.Printf("Read: %d %s", i, e)
		log.Printf("recv/send: %s / %s", currentMessage, result)
		err = c.WriteMessage(websocket.TextMessage, result)
		currentMessage = []byte{}
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
	vp.Wait()

}

func main() {
	fs := http.FileServer(http.Dir("./client/"))
	http.Handle("/client/", http.StripPrefix("/client/", fs))
	http.HandleFunc("/echo", echo)
	http.HandleFunc("/process", process)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
