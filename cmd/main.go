package main

import (
	"github.com/mislavperi/jafa/cmd/api/bootstrap"
	"github.com/mislavperi/jafa/utils"
)

func main() {
	api := bootstrap.Api()
	utils.Run(api)
}
