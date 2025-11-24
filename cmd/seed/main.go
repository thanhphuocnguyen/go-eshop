package main

import (
	"context"
	"os"

	"github.com/thanhphuocnguyen/go-eshop/internal/cmd"
)

func main() {
	ret := cmd.ExecuteSeed(context.Background())
	os.Exit(ret)
}
