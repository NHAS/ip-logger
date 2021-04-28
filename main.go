package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"

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

func main() {

	listenAddr := flag.String("server", "0.0.0.0:8080", "Server listen address")

	flag.Parse()

	serverMode := false

	flag.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "server":
			serverMode = true
		}
	})

	if serverMode {

		err := models.OpenDataBase("db.sql")
		if err != nil {
			log.Fatal(err)
		}

		err = StartCommandHandler()
		if err != nil {
			log.Fatal(err)
		}

		http.HandleFunc("/a/", redirectionHandler)
		log.Fatal(http.ListenAndServe(*listenAddr, nil))
	}

	c, err := makeCommand(flag.Args())
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.Dial("unix", SockAddr)
	if err != nil {
		log.Fatal(err)
	}

	b, err := json.Marshal(c)
	if err != nil {
		log.Fatal(err)
	}

	conn.Write(b)

	re := bufio.NewReader(conn)

	line, _, _ := re.ReadLine()

	fmt.Println(string(line))
}
