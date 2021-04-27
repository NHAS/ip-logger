package main

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/NHAS/ip-logger/models"
	"github.com/NHAS/ip-logger/util"
)

func redirectionHandler(w http.ResponseWriter, r *http.Request) {
	u, err := models.GetUrl(util.GetId(r.URL.Path))
	if err != nil {

		err = models.NewVisit(u.Identifier, r.RemoteAddr)
		if err != nil {
			log.Fatal(err)
		}

		http.Redirect(w, r, u.Destination, 302)

		return
	}

	http.NotFound(w, r)

	return
}

const SockAddr = "/tmp/iplogcontrol.sock"

func main() {
	if err := os.RemoveAll(SockAddr); err != nil {
		log.Fatal(err)
	}

	err := models.OpenDataBase("db.sql")
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		l, err := net.Listen("unix", SockAddr)
		if err != nil {
			log.Fatal("listen error:", err)
		}
		defer l.Close()

		for {
			// Accept new connections, dispatching them to echoServer
			// in a goroutine.
			conn, err := l.Accept()
			if err != nil {
				log.Fatal("accept error:", err)
			}

			go commands(conn)
		}
	}()

	http.HandleFunc("/a/", redirectionHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func commands(conn net.Conn) {
	var cmd Command
	err := json.NewDecoder(conn).Decode(&cmd)
	if err != nil {
		return
	}

	switch cmd.Cmd {
	case "ls":
	case "create":
		label := ""
		if len(cmd.Args) > 1 {
			label = cmd.Args[0]
		}

		if len(cmd.Args) > 0 {
			id, err := models.NewUrl(cmd.Args[len(cmd.Args)-1], label)
			if err != nil {
				return
			}

			conn.Write([]byte(id))
		}
	case "rm":

	}
}

type Command struct {
	Cmd  string
	Args []string
}
