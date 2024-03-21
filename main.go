package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"slices"
	"strings"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	config "github.com/roychowdhuryrohit-dev/projectmeer/lib"
	"github.com/roychowdhuryrohit-dev/projectmeer/lib/algos"
	"github.com/roychowdhuryrohit-dev/projectmeer/lib/routes"
)

func main() {
	config.Config()
	nodeListValue, _ := config.ConfigMap.Load(config.NodeListValue)
	nodeList := strings.Split(nodeListValue.(string), ",")
	replicaID, _ := config.ConfigMap.Load(config.ReplicaID)
	port, _ := config.ConfigMap.Load(config.Port)

	curNode := slices.Index(nodeList, replicaID.(string))
	if curNode == -1 {
		log.Panicln("invalid replica id")
	}
	config.NodeList = nodeList
	fg := algos.NewFugueMax[rune](curNode, nodeList)
	r := setupServer(fg, replicaID.(string))

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	srv := &http.Server{
		Addr:    ":" + port.(string),
		Handler: r,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Panicln(err.Error())
		}
	}()
	log.Printf("node#%d server started at :%s ", curNode, port)

	<-exit
	log.Println("shutting down server")
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
	defer func() {
		cancel()
	}()
	srv.Shutdown(ctx)
}

func setupServer(fg *algos.FugueMax[rune], origin string) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://" + origin},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	fs := http.FileServer(http.Dir("./assets/build"))
	r.Handle("/*", fs)
	r.Mount("/web", webRoutes(fg))
	r.Mount("/p2p", p2pRoutes(fg))
	return r
}

func webRoutes(fg *algos.FugueMax[rune]) chi.Router {
	r := chi.NewRouter()
	r.Post("/insertText", routes.InsertText(fg))
	r.Post("/deleteText", routes.DeleteText(fg))
	r.Get("/getText", routes.GetText(fg))
	r.Get("/getNodeList", routes.GetNodeList())
	return r
}

func p2pRoutes(fg *algos.FugueMax[rune]) chi.Router {
	r := chi.NewRouter()
	r.Post("/receivePrimitive", routes.ReceivePrimitiveHandler(fg))
	return r
}
