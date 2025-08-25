package bootstrap

import "github.com/mislavperi/jafa/internal/api"

func Api() *api.API {
	api := api.NewAPI()
	return api
}
