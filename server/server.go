package server

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
)

func Run() {
	port := ":4141"
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Use(middleware.Heartbeat("/ping"))

	r.With(formatCtx).Route("/v1", func(r chi.Router) {
		r.With(patientCtx).Route("/Patient", func(r chi.Router) {
			r.Post("/", NewPatientHandler)
		})
		r.Get("/metadata", CapabilityStmt)
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		w.Write([]byte("Resource not found"))
	})

	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(405)
		w.Write([]byte("Noop"))
	})

	fmt.Println("Server: http://127.0.0.1:4141/")

	if os.Getenv("PORT") != "" {
		port = ":" + os.Getenv("PORT")
	}

	http.ListenAndServe(port, r)
}

func CapabilityStmt(w http.ResponseWriter, r *http.Request) {
	cap, err := GetCapabilityStatement(r)
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		w.Write([]byte("Internal Server Error"))
		return
	}

	format := r.Context().Value("format").(string)

	if format == "xml" {
		xml, err := xml.Marshal(cap)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("%s", err)))
			return
		}
		w.Header().Set("Content-Type", "application/fhir+xml")
		w.WriteHeader(http.StatusCreated)
		w.Write(xml)
		return
	}

	json, err := json.Marshal(cap)

	if err != nil {
		fmt.Print("beep boop")
		log.Println(err)
		w.WriteHeader(500)
		w.Write([]byte("Internal Server Error"))
		return
	}

	w.Header().Set("Content-Type", "application/fhir+json")
	w.WriteHeader(http.StatusCreated)
	w.Write(json)
}

func formatCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		param := r.URL.Query().Get("_format")
		format := "json"
		if param == "xml" {
			format = "xml"
		}
		ctx = context.WithValue(ctx, "format", format)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func patientCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Body == nil {
			http.Error(w, "Please send a request body", 400)
			return
		}

		if r.Method == "GET" {
			id := chi.URLParam(r, "id")

			if id == "" {
				http.Error(w, "Please provide an id", 400)
				return
			}

			fmt.Println(id)

			ctx := r.Context()
			ctx = context.WithValue(ctx, "id", id)
		}
		next.ServeHTTP(w, r)
	})
}

func NewPatientHandler(w http.ResponseWriter, r *http.Request) {
	patient, err := CreatePatient(r)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("%s", err)))
		return
	}

	json, err := json.Marshal(patient)

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	format := r.Context().Value("format").(string)

	if format == "xml" {
		xml, err := xml.Marshal(patient)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("%s", err)))
			return
		}
		w.Header().Set("Content-Type", "application/fhir+xml")
		w.WriteHeader(http.StatusCreated)
		w.Write(xml)
		return
	}

	w.Header().Set("Content-Type", "application/fhir+json")
	w.WriteHeader(http.StatusCreated)
	w.Write(json)
}
