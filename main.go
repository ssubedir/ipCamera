package main

import (
	"fmt"
	proto "ipCamera/proto/ipcamera/proto"
	"ipCamera/service"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {

	// grpc
	gs := grpc.NewServer()
	c := service.NewService()
	proto.RegisterIpCameraServer(gs, c)
	reflection.Register(gs)

	go func() {
		// create a TCP socket for inbound server connections
		l, err := net.Listen("tcp", ":9000")
		if err != nil {
			os.Exit(1)
		}

		fmt.Println("Starting ipCamera on port 9000")

		// listen for requests
		gs.Serve(l)
	}()

	sigs := make(chan os.Signal, 1)
	sigdone := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Signals
	go func() {
		sig := <-sigs
		fmt.Println(sig)
		sigdone <- true
	}()

	<-sigdone
	gs.GracefulStop()
	fmt.Println("Gracefulstop server ")

}
