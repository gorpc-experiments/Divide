package main

import (
	"errors"
	"fmt"
	"github.com/AliceDiNunno/KubernetesUtil"
	"github.com/davecgh/go-spew/spew"
	"github.com/gorpc-experiments/GalaxyClient"
	"log"
	"net/http"
	"net/rpc"
	"os"
	"strconv"
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

func getPort() int {
	port := 0
	if KubernetesUtil.IsRunningInKubernetes() {
		port = KubernetesUtil.GetInternalServicePort()
	}
	if port == 0 {
		env_port := os.Getenv("PORT")
		if env_port == "" {
			log.Fatalln("PORT env variable isn't set")
		}
		envport, err := strconv.Atoi(env_port)
		if err != nil {
			log.Fatalln(err.Error())
		}
		port = envport
	}

	return port
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
	port := getPort()

	println("Divide is running on port", port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		log.Println(err.Error())
	}
}
