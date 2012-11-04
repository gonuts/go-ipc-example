package main

// +build ignore

import (
	"fmt"
	"io"
	"log"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os/exec"
)

type ReadWriteCloser struct {
	io.ReadCloser
	io.WriteCloser
}

func (rw *ReadWriteCloser) Close() error {
	rw.ReadCloser.Close()
	rw.WriteCloser.Close()
	return nil
}

// Defined for the RPC package
type Adder int

var no int

// Returns how many times the function has been called.
func (p *Adder) Add(in *Adder, ret *int) error {
	no++
	*p += *in + 2
	*in = *p
	return nil
}

func main() {
	var err error

	clientApp := exec.Command("./client-example", "")

	rwc := new(ReadWriteCloser)

	rwc.WriteCloser, err = clientApp.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}

	rwc.ReadCloser, err = clientApp.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	serv := rpc.NewServer()
	codec := jsonrpc.NewServerCodec(rwc)
	fmt.Println("Made RPC server")
	m := new(Adder)
	serv.Register(m)
	fmt.Println("Registered adder service")

	err = clientApp.Start()
	if err != nil {
		log.Fatal(err)
	}

	go serv.ServeCodec(codec)

	clientApp.Wait()

	fmt.Printf("Adder has been called %d times and is now: %d\n", no, *m)
}
