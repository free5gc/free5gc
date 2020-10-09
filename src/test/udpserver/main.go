package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/calee0219/fatal"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "UDP echo server"
	app.Usage = "./udpecho"
	app.Action = action
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "addr, a",
			Value: "0.0.0.0",
			Usage: "Set local UDP address to listen",
		},
		cli.StringFlag{
			Name:  "port, p",
			Value: "30000",
			Usage: "Set local UDP port to listen",
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func serve(listener net.PacketConn, addr net.Addr, buf []byte) {
	fmt.Printf("%s\t: %s\n", addr, buf)
	_, err := listener.WriteTo(buf, addr)
	if err != nil {
		fatal.Fatalf("listener WriteTo error in serve: %+v", err)
	}
}

func action(c *cli.Context) error {
	src := c.String("addr") + ":" + c.String("port")
	listener, err := net.ListenPacket("udp", src)
	if err != nil {
		fatal.Fatalf("ListenPacket: %v", err)
	}
	defer func() {
		errLis := listener.Close()
		if errLis != nil {
			fatal.Fatalf("listener Close error in action: %+v", errLis)
		}
	}()

	fmt.Printf("UDP server start and listening on %s.\n", src)

	for {
		buf := make([]byte, 1024)
		n, addr, err1 := listener.ReadFrom(buf)
		if err1 != nil {
			err = fmt.Errorf("ReadFrom: %v", err1)
			return err
		}
		go serve(listener, addr, buf[:n])

	}
}
