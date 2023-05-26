package main

import (
	"errors"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/gorpc-experiments/ServiceCore"
	"log"
	"net/http"
	"net/rpc"
	"os"
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
	arith := new(Arith)
	err := rpc.Register(arith)

	client, err := ServiceCore.NewGalaxyClient()

	if err != nil {
		log.Println(err.Error())
		return
	}

	client.RegisterToGalaxy(arith)

	rpc.HandleHTTP()
	port := ServiceCore.GetRPCPort()

	println("Divide is running on port", port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		log.Println(err.Error())
	}
}
