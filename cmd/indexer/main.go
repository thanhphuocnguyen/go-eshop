package main

import (
	"context"
	"os"

	"github.com/thanhphuocnguyen/go-eshop/internal/cmd"
)

func main() {
	ctx := context.Background()
	ret := cmd.ExecuteIndexer(ctx)
	os.Exit(ret)
}
