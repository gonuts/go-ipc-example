package main

// +build ignore

import (
	"io"
	"log"
	"net/rpc/jsonrpc"
	"os"
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

func main() {
	var err error
	adder := Adder(1)
	var ret int

	rwc := new(ReadWriteCloser)

	rwc.ReadCloser = os.Stdin
	rwc.WriteCloser = os.Stdout

	client := jsonrpc.NewClient(rwc)

	for i := 0; i < 3; i++ {
		err = client.Call("Adder.Add", &adder, &ret)
		if err != nil {
			log.Fatal(err)
		}
	}
}
