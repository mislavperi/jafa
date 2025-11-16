package main

import (
	"github.com/mislavperi/jafa/cmd/server/bootstrap"
	"github.com/mislavperi/jafa/utils"
)

func main() {
	server := bootstrap.Server()
	utils.Run(server)
}
