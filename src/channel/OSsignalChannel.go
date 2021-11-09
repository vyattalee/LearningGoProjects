package main

import (
	"fmt"
	"os"
	"os/signal"
)

func cleanup() {
	fmt.Println("cleanup")
}

func main() {
	//c := make(chan os.Signal)
	//signal.Notify(c, os.Interrupt) // This traps Control-c to safely shut down the server
	////, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGUSR1
	////go func() {
	//
	//
	//	fmt.Println("OsSignal：", <-c)
	//	//srvChanIn <- d2networking.ServerEventStop
	////}()

	c := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			// sig is a ^C, handle it
			fmt.Println("OsSignal：", sig)
		}
		cleanup()
		done <- true

	}()

	fmt.Println("awaiting signal")
	<-done
	fmt.Println("exiting")

	//c := make(chan os.Signal)
	//signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	//go func() {
	//	fmt.Println("OsSignal：", <-c)
	//	cleanup()
	//	os.Exit(1)
	//}()
	//
	//for {
	//	fmt.Println("sleeping...")
	//	time.Sleep(10 * time.Second) // or runtime.Gosched() or similar per @misterbee
	//}

	//ch := make(chan os.Signal, 1)
	//signal.Notify(ch, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGUSR1)
	//for {
	//	s := <-ch
	//	switch s {
	//	case syscall.SIGQUIT:
	//		fmt.Println("SIGSTOP")
	//		return
	//	case syscall.SIGSTOP:
	//		fmt.Println("SIGSTOP")
	//		return
	//	case syscall.SIGHUP:
	//		fmt.Println("SIGHUP")
	//		return
	//	case syscall.SIGKILL:
	//		fmt.Println("SIGKILL")
	//		return
	//	case syscall.SIGUSR1:
	//		fmt.Println("SIGUSR1")
	//		return
	//	default:
	//		fmt.Println("default")
	//		return
	//	}
	//}

}
