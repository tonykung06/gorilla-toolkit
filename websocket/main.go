package main

import (
	"io"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var analyticsConns []*websocket.Conn
var logCh chan string

func main() {
	// basicWS()
	// wsReaderWriter()
	// wsJSON()
	dispatching()
}

func init() {
	logCh = make(chan string)
	go func() {
		for msg := range logCh {
			for _, c := range analyticsConns {
				w, _ := c.NextWriter(websocket.TextMessage)
				w.Write([]byte(time.Now().Format("2006-01-02 15:04:05")))
				w.Write([]byte(" - " + msg))
				w.Close()
			}
		}
	}()
}

func dispatching() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("testing"))
	})
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, _ := upgrader.Upgrade(w, r, nil)
		analyticsConns = append(analyticsConns, conn)
		go func(conn *websocket.Conn) {
			for {
				if _, _, err := conn.NextReader(); err != nil {
					conn.Close()
					for i := range analyticsConns {
						if analyticsConns[i] == conn {
							analyticsConns = append(analyticsConns[:i], analyticsConns[i+1:]...)
						}
					}
					break
				}
			}
		}(conn)
	})
	http.ListenAndServe(":3000", nil)
}

func wsReaderWriter() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("testing"))
	})
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, _ := upgrader.Upgrade(w, r, nil)
		go func(conn *websocket.Conn) {
			for {
				messageType, reader, _ := conn.NextReader()
				writer, _ := conn.NextWriter(messageType)
				io.Copy(writer, reader)
				writer.Close()
			}
		}(conn)
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
		go func(conn *websocket.Conn) {
			for {
				var author Author
				conn.ReadJSON(&author)
				conn.WriteMessage(websocket.TextMessage, []byte("Author's Name: "+author.Name))
			}
		}(conn)
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
