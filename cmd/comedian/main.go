package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/razic/comedian/api/services/icndb"
	icndbpb "github.com/razic/comedian/api/services/icndb"
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
			Name:  "uinames-socket",
			Value: "/var/run/uinames.sock",
			Usage: "socket path to uinames grpc server",
		},
		cli.StringFlag{
			Name:  "icndb-socket",
			Value: "/var/run/icndb.sock",
			Usage: "socket path to icndb grpc server",
		},
	}
	app.Action = func(c *cli.Context) error {
		dialOpts := append(
			[]grpc.DialOption{grpc.WithInsecure()},
			grpc.WithDialer(func(addr string, timeout time.Duration) (net.Conn, error) {
				return net.DialTimeout("unix", addr, timeout)
			}),
		)

		uinamesConn, err := grpc.Dial(c.String("uinames-socket"), dialOpts...)

		if err != nil {
			log.Fatalf("could not connect to uinames grpc server: %v", err)
		}

		defer uinamesConn.Close()

		icndbConn, err := grpc.Dial(c.String("icndb-socket"), dialOpts...)

		if err != nil {
			log.Fatalf("could not connect to icndb grpc server: %v", err)
		}

		defer icndbConn.Close()

		uinames := uinamespb.NewUinamesClient(uinamesConn)
		icndb := icndb.NewIcndbClient(icndbConn)

		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			nameResp, err := uinames.GetName(context.Background(), &uinamespb.GetNameRequest{})

			if err != nil {
				fmt.Fprintf(w, "error getting name from uinames grpc: %v", err)
				return
			}

			jokeResp, err := icndb.GetJoke(context.Background(), &icndbpb.GetJokeRequest{
				FirstName: nameResp.Name,
				LastName:  nameResp.Surname,
			})

			if err != nil {
				fmt.Fprintf(w, "error getting joke from icndb grpc: %v", err)
				return
			}

			fmt.Fprintf(w, jokeResp.Value.Joke)
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
