package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	pb "github.com/razic/comedian/api/services/icndb"
	"github.com/urfave/cli"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const usage = `
  _                 _ _
 (_) ___ _ __   __| | |__ 
 | |/ __| '_ \ / _  | '_ \ 
 | | (__| | | | (_| | |_) |
 |_|\___|_| |_|\__,_|_.__/ 

grpc wrapper for icndb api
`

const endpoint = "https://api.icndb.com/jokes/random?firstName=%s&lastName=%s&limitTo=nerdy"

var (
	httpClient = http.Client{
		Timeout: time.Duration(10 * time.Second),
	}
)

type icndb struct{}

// GetJoke gets a name from the icndb api
func (u *icndb) GetJoke(ctx context.Context, in *pb.GetJokeRequest) (*pb.GetJokeResponse, error) {
	res, err := httpClient.Get(fmt.Sprintf(endpoint, in.FirstName, in.LastName))

	if err != nil {
		log.Fatalf("failed to get joke from icndb api: %v", err)
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		log.Fatalf("failed read response body: %v", err)
		return nil, err
	}

	resp := new(pb.GetJokeResponse)

	err = json.Unmarshal(body, &resp)

	if err != nil {
		return nil, err
	}

	return resp, err
}

func main() {
	app := cli.NewApp()

	app.Name = "icndb"
	app.Usage = usage
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "socket, s",
			Value: "/var/run/icndb.sock",
			Usage: "socket path to grpc server",
		},
	}
	app.Action = func(c *cli.Context) error {
		listen, err := net.Listen("unix", c.String("socket"))

		if err != nil {
			log.Fatalf("failed to listen: %v", err)
			os.Exit(1)
		}

		sigChan := make(chan os.Signal, 1)
		server := grpc.NewServer()

		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		pb.RegisterIcndbServer(server, &icndb{})

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
