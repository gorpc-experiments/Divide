package main

import (
	"errors"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/gorpc-experiments/GalaxyClient"
	"log"
	"net/http"
	"net/rpc"
	"os"
	"strings"
)

type Args struct {
	A, B int
}

type Quotient struct {
	Quo, Rem int
	Err      string
}

type Arith int

func (t *Arith) Divide(args *Args, quo *Quotient) error {
	if args.B == 0 {
		return errors.New("divide by zero")
	}
	quo.Quo = args.A / args.B
	quo.Rem = args.A % args.B

	hostname, _ := os.Hostname()
	quo.Err = hostname

	spew.Dump(args, quo)

	return nil
}

func main() {
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		fmt.Println(pair[0], "=", pair[1])
	}

	arith := new(Arith)
	err := rpc.Register(arith)

	client, err := GalaxyClient.NewGalaxyClient()

	if err != nil {
		log.Println(err.Error())
		return
	}

	client.RegisterToGalaxy(arith)

	rpc.HandleHTTP()

	err = http.ListenAndServe(":2345", nil)
	if err != nil {
		log.Println(err.Error())
	}
}
