package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	wsservice "app.wsservice"
)

func main() {
	fmt.Println("Start service")
	done := setupWSServerCloseHandler()

	var configpath string
	argconfigpath := flag.String("path", "", "a string")
	flag.Parse()
	if *argconfigpath == "" {
		mainpath, _ := os.Getwd()
		configpath = fmt.Sprintf("%s/%s", mainpath, ".conf/wsservice.yaml")
	} else {
		configpath = *argconfigpath
	}
	ctx, cancel := context.WithCancel(context.Background())

	wssrv, err := wsservice.NewWsService(ctx, configpath)
	if err != nil {
		fmt.Println("NewWsService can't create ws service engine:", err)
		cancel()
		os.Exit(1)
	}

	if err := wssrv.Run(); err != nil {
		fmt.Println("Ws service start error:", err)
		wssrv.Shutdown()
	}

	go func() {
		for {
			select {
			case <-wssrv.Complete():
				cancel()
			case <-done:
				go wssrv.Shutdown()
				<-wssrv.Complete()
				cancel()
			}
		}
	}()
	<-ctx.Done()
	fmt.Println("Finish ws service")
}

func setupWSServerCloseHandler() chan bool {
	c := make(chan os.Signal)
	d := make(chan bool)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-c
		d <- true
	}()
	return d
}
