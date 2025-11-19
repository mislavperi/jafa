package main

import (
	"github.com/mislavperi/jafa/server/cmd/server/bootstrap"
	"github.com/mislavperi/jafa/server/utils"
)

func main() {
	server := bootstrap.Server()
	utils.Run(server)
}
