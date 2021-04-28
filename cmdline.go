package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/url"
	"os"

	"github.com/NHAS/ip-logger/models"
)

const SockAddr = "/tmp/iplogcontrol.sock"

func StartCommandHandler() (err error) {
	if err = os.RemoveAll(SockAddr); err != nil {
		log.Fatal(err)
	}

	go func() {
		l, err := net.Listen("unix", SockAddr)
		if err != nil {
			log.Fatal("listen error:", err)
		}
		defer l.Close()

		for {
			conn, err := l.Accept()
			if err != nil {
				log.Fatal("accept error:", err)
			}

			go runCommands(conn)
		}
	}()

	return
}

func makeCommand(args []string) (c Command, err error) {
	if len(args) > 0 {

		switch args[0] {
		case "ls":
			c.Cmd = "ls"
			c.Args = args[1:]
			return
		case "rm":
			c.Cmd = "rm"
			c.Args = args[1:]
			return
		}

		_, err = url.Parse(args[len(args)-1])
		if err == nil {
			c.Cmd = "create"
			c.Args = args

		}

		return

	}

	return c, fmt.Errorf("Unable to create command")
}

func runCommands(conn net.Conn) {
	defer conn.Close()

	var cmd Command
	err := json.NewDecoder(conn).Decode(&cmd)
	if err != nil {
		return
	}

	switch cmd.Cmd {
	case "ls":
		urls, err := models.GetAllUrls()
		log.Println(err)
		b, _ := json.MarshalIndent(urls, "", "    ")
		conn.Write(b)

	case "create":
		label := ""
		if len(cmd.Args) > 1 {
			label = cmd.Args[0]
		}

		if len(cmd.Args) > 0 {
			id, err := models.NewUrl(cmd.Args[len(cmd.Args)-1], label)
			if err != nil {
				conn.Write([]byte(err.Error()))
				return
			}

			conn.Write([]byte(id))
			return
		}

	case "rm":
		if len(cmd.Args) != 1 {
			log.Println("Nothing")
			return
		}
		models.DeleteUrl(cmd.Args[0])

	}

}

type Command struct {
	Cmd  string
	Args []string
}
