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

func StartCommandHandler(domain string) (err error) {
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

			go runCommands(domain, conn)
		}
	}()

	return
}

func makeCommand(args []string) (c Command, err error) {
	if len(args) > 0 {

		u, err := url.Parse(args[len(args)-1])
		if err != nil {
			return c, err
		}

		switch u.Scheme {
		case "http", "https":
			c.Cmd = "create"
			c.Args = args
		default:
			c.Cmd = args[0]
			c.Args = args[1:]
		}

		return c, err
	}

	return c, fmt.Errorf("Unable to create command")
}

func runCommands(domain string, conn net.Conn) {
	defer conn.Close()

	var cmd Command
	err := json.NewDecoder(conn).Decode(&cmd)
	if err != nil {
		return
	}

	switch cmd.Cmd {
	case "ls":

		var urls []models.Url

		if len(cmd.Args) == 0 {
			urls, err = models.GetAllUrls()
			if err != nil {
				fmt.Fprintf(conn, "Loading all entries failed: %s", err.Error())
				return
			}
		} else {
			u, err := models.GetUrl(cmd.Args[0])
			if err != nil {
				fmt.Fprintf(conn, "Loading entry failed: %s", err.Error())
				return
			}

			urls = []models.Url{u}

		}

		b, err := json.MarshalIndent(&urls, "", "    ")
		if err != nil {
			fmt.Fprintf(conn, "Marshalling entries failed: %s", err.Error())
			return
		}
		conn.Write(b)
		return

	case "create":
		label := ""
		if len(cmd.Args) > 1 {
			label = cmd.Args[0]
		}

		if len(cmd.Args) > 0 {
			id, err := models.NewUrl(cmd.Args[len(cmd.Args)-1], label)
			if err != nil {
				fmt.Fprintf(conn, "Creatng encountered an error %s", err.Error())
				return
			}

			fmt.Fprintf(conn, "%s/a/%s", domain, id)
			return
		}

		fmt.Fprintf(conn, "Not enough arguments for create")

	case "rm":
		if len(cmd.Args) != 1 {
			conn.Write([]byte("Not enough arguments"))
			return
		}
		err = models.DeleteUrl(cmd.Args[0])
		if err != nil {
			log.Println(err)
			fmt.Fprintf(conn, "Error deleting %s : %s", cmd.Args[0], err.Error())
			return
		}

		fmt.Fprintf(conn, "Deleted %s", cmd.Args[0])

	default:
		conn.Write([]byte("Unknown command " + cmd.Cmd + "\n"))
	}

}

type Command struct {
	Cmd  string
	Args []string
}
