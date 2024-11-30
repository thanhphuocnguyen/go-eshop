package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/thanhphuocnguyen/go-eshop/internal/cmd"
)

var interruptSignals = []os.Signal{
	os.Interrupt,
	syscall.SIGTERM,
	syscall.SIGINT,
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), interruptSignals...)
	defer stop()
	cmd.Execute(ctx)
}
