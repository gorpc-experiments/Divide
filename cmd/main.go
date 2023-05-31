package main

import (
	"errors"
	"github.com/gorpc-experiments/ServiceCore"
	"log"
	"os"
)

type Args struct {
	A, B int
}

type Quotient struct {
	Quo, Rem int
	Err      string
}

type Arith struct {
	ServiceCore.CoreHealth
}

func (t *Arith) Divide(args *Args, quo *Quotient) error {
	if args.B == 0 {
		return errors.New("divide by zero")
	}
	quo.Quo = args.A / args.B
	quo.Rem = args.A % args.B

	hostname, _ := os.Hostname()
	quo.Err = hostname

	log.Printf("%d/ %d = %d rem %d\n", args.A, args.B, quo.Quo, quo.Rem)

	return nil
}

func main() {
	ServiceCore.SetupLogging()

	arith := new(Arith)

	ServiceCore.PublishMicroService(arith, true)
}