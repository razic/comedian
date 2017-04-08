package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	uinamespb "github.com/razic/comedian/api/services/uinames"
	"github.com/urfave/cli"
	"google.golang.org/grpc"
)

const usage = `
                              _ _
  ___ ___  _ __ ___   ___  __| (_) __ _ _ __  
 / __/ _ \| '_   _ \ / _ \/ _  | |/ _  | '_ \ 
| (_| (_) | | | | | |  __/ (_| | | (_| | | | |
 \___\___/|_| |_| |_|\___|\__,_|_|\__,_|_| |_|

http server combining uinames + icndb apis
`

func main() {
	app := cli.NewApp()

	app.Name = "comedian"
	app.Usage = usage
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "socket, s",
			Value: "/var/run/uinames.sock",
			Usage: "socket path to grpc server",
		},
	}
	app.Action = func(c *cli.Context) error {
		dialOpts := append(
			[]grpc.DialOption{grpc.WithInsecure()},
			grpc.WithDialer(func(addr string, timeout time.Duration) (net.Conn, error) {
				return net.DialTimeout("unix", addr, timeout)
			}),
		)
		conn, err := grpc.Dial(c.String("socket"), dialOpts...)

		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}

		defer conn.Close()

		uinames := uinamespb.NewUinamesClient(conn)

		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			nameResp, err := uinames.GetName(context.Background(), &uinamespb.GetNameRequest{})

			if err != nil {
				fmt.Fprintf(w, "Error: %v", err)
				return
			}

			fmt.Fprintf(w, "Hi there, I love %s", nameResp.Name)
		})

		http.ListenAndServe(":8080", nil)

		if err != nil {
			log.Fatalf("could not get name: %v", err)
			os.Exit(1)
		}

		return nil
	}

	app.Run(os.Args)
}
