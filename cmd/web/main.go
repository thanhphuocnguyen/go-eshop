package main

import (
	"context"
	"os/signal"

	"github.com/thanhphuocnguyen/go-eshop/internal/cmd"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), cmd.InterruptSignals...)
	defer stop()
	cmd.Execute(ctx)
}
