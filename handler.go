package main

import (
	"html"
	"html/template"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/go-chi/chi/v5"
	"golang.org/x/exp/errors/fmt"
)

type Page struct {
	Title string
}

func EscapeHTML(s string) template.HTML {
	return template.HTML(html.UnescapeString(s))
}

func (server *FServer) RegisterHandlers() {
	r := server.Base.Handler.(*chi.Mux)

	wd, err := os.Getwd()

	if err != nil {
		log.Fatal(err)
	}

	views := path.Join(wd, "www", "views")

	if server.Config.InContainer {
		views = path.Join(wd, "app", "www", "views")
	}

	r.Route("/", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			page := Page{
				Title: "fhir-kit",
			}

			indexHtml := path.Join(views, "index.html")

			tmpl, err := template.New("index.html").Funcs(template.FuncMap{"escape": EscapeHTML}).ParseFiles(indexHtml)

			if err != nil {
				fmt.Errorf("error parsing template %s: %s", indexHtml, err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("server error"))
			}

			err = tmpl.Execute(w, page)

			if err != nil {
				fmt.Errorf("error executing template %s: %s", indexHtml, err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("server error"))
			}
			w.WriteHeader(http.StatusOK)
		})

		r.Get("/dashboard", func(w http.ResponseWriter, r *http.Request) {
			page := Page{
				Title: "Dashboard",
			}

			dashHtml := path.Join(views, "dashboard.html")
			tmpl, err := template.New("dashboard.html").Funcs(template.FuncMap{"escape": EscapeHTML}).ParseFiles(dashHtml)

			if err != nil {
				fmt.Errorf("error parsing template %s: %s", dashHtml, err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
			}

			err = tmpl.Execute(w, page)

			if err != nil {
				fmt.Errorf("error executing template %s: %s", dashHtml, err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
			}

			w.WriteHeader(http.StatusOK)
		})
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		_, err := w.Write([]byte("Not Found"))
		if err != nil {
			fmt.Errorf("error writing response: %s", err.Error())
			log.Fatal(err)
		}
	})

	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(405)
		_, err := w.Write([]byte("Not Allowed"))
		if err != nil {
			log.Fatal(err)
		}
	})
}

func HeartbeatHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	_, err := w.Write([]byte("pong"))
	if err != nil {
		log.Fatal(err)
	}
}

func CreateProjectHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	_, err := w.Write([]byte("Project Created"))
	if err != nil {
		log.Fatal(err)
	}
}
