package main

import (
	//	"fmt"
	"os"
	//	"os/signal"
)

func main() {
	app := NewApp()
	//ctrl + c
	// c := make(chan os.Signal, 1)
	// signal.Notify(c, os.Interrupt)
	// go func() {
	// 	sig := <-c
	// 	fmt.Printf("receive signal %s\n", sig.String())
	// 	app.Stop()
	// }()
	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}
