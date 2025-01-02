package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/thanhphuocnguyen/go-eshop/internal/cmd"
)

var InterruptSignals = []os.Signal{
	os.Interrupt,
	syscall.SIGTERM,
	syscall.SIGINT,
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), InterruptSignals...)
	defer stop()
	cmd.Execute(ctx)
}
