package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ayang64/grdc/asciiart"
	"github.com/ayang64/grdc/render"
)

func GRDC() {
	t := time.NewTicker(500 * time.Millisecond)

	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT)

	for {
		asciiart.Encode(os.Stdout, render.TextToBitmap(time.Now().Format("03:04:05")))
		select {
		case <-t.C:
		case <-sig:
			return
		}
	}
}
func main() {
	GRDC()
	fmt.Printf("\x1bc")
}
