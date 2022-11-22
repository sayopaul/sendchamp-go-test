package main

import (
	"github.com/sayopaul/sendchamp-go-test/bootstrap"
	"go.uber.org/fx"
)

func main() {
	//bootstrap app and inject dependencies using fx
	fx.New(bootstrap.Bootstrap).Run()
}
