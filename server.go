package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type FServer struct {
	Base   *http.Server
	Config Config
}

func NewServer(c Config) *FServer {
	server := &FServer{
		Base: &http.Server{
			Handler: chi.NewRouter(),
			Addr:    ":" + c.Port,
		},
		Config: c,
	}
	return server
}

func (server *FServer) Start() error {

	err := server.RegisterMiddlewares()

	if err != nil {
		return err
	}

	server.RegisterHandlers()

	router := server.Base.Handler.(*chi.Mux)

	if server.Config.Verbose {
		fmt.Println()
		theader := fmt.Sprintf("%-6s | %-6s\n", "METHOD", "ROUTE")
		tsep := fmt.Sprintf("%-6s + %-6s\n", "------", "------")
		fmt.Print(theader)
		fmt.Print(tsep)

		walker := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
			if route != "/" && route[len(route)-1:] == "/" {
				route = route[:len(route)-1]
			}

			fmt.Printf("%-6s | %-6s\n", method, route)
			return nil
		}

		if err := chi.Walk(router, walker); err != nil {
			return err
		}
		fmt.Print(tsep + "\n")

	}

	server.Printf("Listening on port %s\n", server.Base.Addr)

	if server.Config.InContainer {
		server.Printf("Server is running in a container\n")
	}

	err = server.Base.ListenAndServe()

	if err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (server *FServer) Printf(str string, args ...interface{}) {
	if server.Config.Verbose {
		fmt.Printf(str, args...)
	}
}
