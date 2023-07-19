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
		stdin:  bytes.Buffer{},
		stdout: bytes.Buffer{},
		stderr: bytes.Buffer{},
		proc:   nil,
	}
	for {
		vp.stdout.Reset()
		vp.stderr.Reset()
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		vp.stdin.Write(message)

		vp.Start()
		vp.Wait()

		err = c.WriteMessage(websocket.TextMessage, vp.stdout.Bytes())
		if err != nil {
			log.Println("write:", err)
			break
		}
	}

}

func main() {
	fs := http.FileServer(http.Dir("./client/"))
	http.Handle("/client/", http.StripPrefix("/client/", fs))
	http.HandleFunc("/echo", echo)
	http.HandleFunc("/process", process)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
