package router

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/tmdgo/dependencies"
)

type Router struct {
	routes  []Route
	manager *dependencies.Manager
}

func (router *Router) Init(manager *dependencies.Manager) {
	router.routes = make([]Route, 0)
	router.manager = manager
}

func (router *Router) AddRoute(route Route) {
	router.routes = append(router.routes, route)
}

func (router *Router) AddController(controller Controller) {
	router.routes = append(router.routes, controller.GetRoutes()...)
}

func (router *Router) ListenAndServe(address string) {
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	r := mux.NewRouter()

	for _, route := range router.routes {
		r.HandleFunc(route.Path, routerHandler(router, route, route.HandleFunc)).Methods(route.Method)
	}

	srv := &http.Server{
		Addr:         address,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		Handler:      handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(r),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	srv.Shutdown(ctx)
	log.Println("shutting down")
	os.Exit(0)
}
