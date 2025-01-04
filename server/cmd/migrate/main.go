package main

import (
	"context"
	"os"

	"github.com/thanhphuocnguyen/go-eshop/internal/cmd"
)

func main() {
	ctx := context.Background()
	ret := cmd.ExecuteMigrate(ctx)
	os.Exit(ret)
}
