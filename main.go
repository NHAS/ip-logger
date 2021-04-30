package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"

	"github.com/NHAS/ip-logger/models"
	"github.com/NHAS/ip-logger/util"
)

func redirectionHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Path)

	u, err := models.GetUrl(models.GetId(r.URL.Path))
	if err != nil {
		http.NotFound(w, r)
		log.Println(err)
		return
	}

	err = models.NewVisit(u.Identifier, r.RemoteAddr, r.Header.Get("User-Agent"))
	if err != nil {
		log.Fatal(err)
	}

	http.Redirect(w, r, u.Destination, 302)

}

type config struct {
	Domain string
}

func main() {

	listenAddr := flag.String("server", "0.0.0.0:8080", "Server listen address")
	configPath := flag.String("config", "config.json", "Configuration file for server location")

	flag.Parse()

	serverMode := false

	flag.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "server":
			serverMode = true
		}
	})

	configBytes, err := ioutil.ReadFile(*configPath)
	if err != nil {
		log.Fatal(err)
	}

	var conf config
	err = json.Unmarshal(configBytes, &conf)
	if err != nil {
		log.Fatal(err)
	}

	if serverMode {

		err = models.OpenDataBase("db.sql")
		if err != nil {
			log.Fatal(err)
		}

		err = StartCommandHandler(conf.Domain)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("Listening on ", *listenAddr)
		http.HandleFunc("/a/", redirectionHandler)

		log.Fatal(http.ListenAndServe(*listenAddr, nil))
	}

	c, err := makeCommand(flag.Args())
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.Dial("unix", SockAddr)
	if err != nil {
		fmt.Println("Unable to contact server")
		return
	}

	b, err := json.Marshal(c)
	if err != nil {
		log.Fatal(err)
	}

	conn.Write(b)

	buf := get(conn)

	switch c.Cmd {
	case "ls":
		var u []models.Url

		err := json.Unmarshal(buf, &u)
		if err != nil {
			fmt.Printf("Error: %s", buf)
			return
		}

		util.PrintTable(conf.Domain, u)

		if len(u) == 1 {
			util.PrintVisits(u[0])
		}
	default:
		fmt.Printf("%s\n", buf)
	}

}

func get(conn net.Conn) (result []byte) {
	for {
		b := make([]byte, 1)
		_, err := conn.Read(b)
		if err != nil {
			break
		}
		result = append(result, b[0])
	}

	return result
}
