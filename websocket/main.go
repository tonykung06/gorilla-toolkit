package main

import (
	"io"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	// basicWS()
	// wsReaderWriter()
	wsJSON()
}

func wsReaderWriter() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("testing"))
	})
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, _ := upgrader.Upgrade(w, r, nil)
		go func() {
			for {
				messageType, reader, _ := conn.NextReader()
				writer, _ := conn.NextWriter(messageType)
				io.Copy(writer, reader)
				writer.Close()
			}
		}()
	})
	http.ListenAndServe(":3000", nil)
}

func basicWS() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("testing"))
	})
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, _ := upgrader.Upgrade(w, r, nil)
		go func() {
			for {
				messageType, p, _ := conn.ReadMessage()
				conn.WriteMessage(messageType, p)
			}
		}()
	})
	http.ListenAndServe(":3000", nil)

	//on browser console,
	//var ws = new WebSocket('ws://localhost:3000/ws')
	//ws.addEventListener('message', function(e) {console.log(e);})
	//ws.send("OMG")
}

func wsJSON() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("testing"))
	})
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, _ := upgrader.Upgrade(w, r, nil)
		go func() {
			for {
				var author Author
				conn.ReadJSON(&author)
				conn.WriteMessage(websocket.TextMessage, []byte("Author's Name: "+author.Name))
			}
		}()
	})
	http.ListenAndServe(":3000", nil)

	//on browser console,
	//var author = {name: 'Tony Kung', books: ['book1', 'book2']};
	//var ws = new WebSocket('ws://localhost:3000/ws');
	//ws.addEventListener('message', function(e) {console.log(e.data);});
	//ws.send(JSON.stringify(author));
}

type Author struct {
	Name  string   `json:"name"`
	Books []string `json:"books"`
}
