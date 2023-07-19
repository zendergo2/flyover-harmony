package main

import (
	"bytes"
	"log"
	"net/http"
	"os"
	"os/exec"

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
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		wd, _ := os.Getwd()
		cmd := exec.Command(wd + "/test.sh")
		var outb, errb bytes.Buffer
		cmd.Stdout = &outb
		cmd.Stderr = &errb
		// var pipe, _ = cmd.StdoutPipe()
		// var errpipe, _ = cmd.StderrPipe()
		inbuf := bytes.Buffer{}
		inbuf.Write([]byte(message))
		cmd.Stdin = &inbuf
		cmd.Run()
		log.Println("ERR: " + errb.String())
		var out = outb.Bytes()
		if bytes.HasSuffix(out, []byte("\n")) {
			outb.Truncate(len(out) - 1)
			outb.Write([]byte{'\r', '\n'})
		}
		err = c.WriteMessage(websocket.TextMessage, outb.Bytes())
		if err != nil {
			log.Println("write:", err)
			break
		}
		// scanner := bufio.NewScanner(pipe)
		// for scanner.Scan() {
		// 	m := scanner.Text()
		// 	log.Println(m)
		// 	// err = c.WriteMessage(websocket.TextMessage, []byte(m))
		// 	// if err != nil {
		// 	// 	log.Println("write:", err)
		// 	// 	break
		// 	// }
		// }
		// log.Println("Started")
		// errscanner := bufio.NewScanner(errpipe)
		// for errscanner.Scan() {
		// 	m := errscanner.Text()
		// 	log.Println("ERR: " + m)
		// 	// err = c.WriteMessage(websocket.TextMessage, []byte("ERR: "+m))
		// 	// if err != nil {
		// 	// 	log.Println("write:", err)
		// 	// 	break
		// 	// }
		// }
	}

}

func main() {
	fs := http.FileServer(http.Dir("./client/"))
	http.Handle("/client/", http.StripPrefix("/client/", fs))
	http.HandleFunc("/echo", echo)
	http.HandleFunc("/process", process)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
