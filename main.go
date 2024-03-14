package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	config "github.com/roychowdhuryrohit-dev/projectmeer/node-user/lib"
	"github.com/roychowdhuryrohit-dev/projectmeer/node-user/lib/algos"
)

func main() {
	config.Config()
	nodeListValue, _ := config.ConfigMap.Load(config.NodeList)
	
	fg := algos.NewFugueMax[string]()
}

func setupServer() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	//r.Get("/",...) TODO: Host static homepage
	r.Mount("/web", webRoutes())
	r.Mount("/p2p", p2pRoutes())
	return r
}

func webRoutes() chi.Router {

}

func p2pRoutes() chi.Router {

}
