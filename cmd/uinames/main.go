package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	pb "github.com/razic/comedian/api/services/uinames"
	"github.com/urfave/cli"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const usage = `
       _ 
 _   _(_)_ __   __ _ _ __ ___   ___  ___ 
| | | | | '_ \ / _  | '_   _ \ / _ \/ __|
| |_| | | | | | (_| | | | | | |  __/\__ \
 \__,_|_|_| |_|\__,_|_| |_| |_|\___||___/

grpc wrapper for uinames api
`

const endpoint = "https://uinames.com/api/"

var (
	httpClient = http.Client{
		Timeout: time.Duration(10 * time.Second),
	}
)

type uinames struct{}

// GetName gets a name from the uinames api
func (u *uinames) GetName(ctx context.Context, in *pb.GetNameRequest) (*pb.GetNameResponse, error) {
	res, err := httpClient.Get(endpoint)

	if err != nil {
		log.Fatalf("failed to get api response from uinames: %v", err)
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		log.Fatalf("failed read response body: %v", err)
		return nil, err
	}

	resp := new(pb.GetNameResponse)

	err = json.Unmarshal(body, &resp)

	if err != nil {
		log.Fatalf("failed unmarshal response body: %v", err)
		return nil, err
	}

	return resp, nil
}

func main() {
	app := cli.NewApp()

	app.Name = "uinames"
	app.Usage = usage
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "socket, s",
			Value: "/var/run/uinames.sock",
			Usage: "socket path to grpc server",
		},
	}
	app.Action = func(c *cli.Context) error {
		//TODO: find a better place for this
		os.Remove(c.String("socket"))

		listen, err := net.Listen("unix", c.String("socket"))

		if err != nil {
			log.Fatalf("failed to listen: %v", err)
			os.Exit(1)
		}

		sigChan := make(chan os.Signal, 1)
		server := grpc.NewServer()

		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		pb.RegisterUinamesServer(server, &uinames{})

		go func() {
			<-sigChan
			os.Remove(c.String("socket"))
			os.Exit(0)
		}()

		if err := server.Serve(listen); err != nil {
			log.Fatalf("failed to serve: %v", err)
			os.Exit(1)
		}

		return nil

	}

	app.Run(os.Args)
}
