package main

import (
	"context"

	"github.com/thanhphuocnguyen/go-eshop/internal/cmd"
)

func main() {
	ctx := context.Background()
	cmd.ExecuteMigrate(ctx)
}
