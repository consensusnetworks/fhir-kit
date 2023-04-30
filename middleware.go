package main

import (
	"errors"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (server *FServer) RegisterMiddlewares() error {
	r := server.Base.Handler.(*chi.Mux)

	// keep logger first
	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(middleware.StripSlashes)
	r.Use(SetDefaultTimeZone)
	// TODO: add application/xml with fhir+json and fhir+xml for backwards compatibility
	r.Use(middleware.AllowContentType("application/json"))
	// r.Use(SetWeakETag)
	r.Use(middleware.Heartbeat("/ping"))

	wd, err := os.Getwd()

	if err != nil {
		return err
	}

	publicDir := http.Dir(path.Join(wd, "www/public"))

	if server.Config.InContainer {
		publicDir = http.Dir(path.Join(wd, "app", "www/public"))
	}

	err = FileServer(r, "/public", publicDir)
	// err := FileServer(r, "/public", http.FS(WebContent))

	if err != nil {
		return err
	}
	return nil
}

func ClientContentTypePrefer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Prefer") != "return=representation" {
			// TODO: set the context based on the Accept header
		}
		next.ServeHTTP(w, r)
	})
}
func FileServer(r chi.Router, p string, root http.FileSystem) error {
	if strings.ContainsAny(p, "{}*") {
		return errors.New("file server error: URL path cannot contain variables")
	}

	if p != "/" && p[len(p)-1] != '/' {
		r.Get(p, http.RedirectHandler(p+"/", http.StatusMovedPermanently).ServeHTTP)
		p += "/"
	}

	p += "*"

	r.Get(p, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})

	return nil
}

func SetDefaultTimeZone(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Header.Set("Date", time.Now().Format(time.RFC1123))
		next.ServeHTTP(w, r)
	})
}

func SetWeakETag(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("ETag", `W/"weak"`)
		next.ServeHTTP(w, r)
	})
}
