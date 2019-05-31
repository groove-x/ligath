package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
)

type TestableHandler func(w http.ResponseWriter, r *http.Request, isMock bool)
type Middleware func(next http.Handler) http.Handler

func main() {
	a := NewAPI()
	err := a.Serve()
	if err != nil {
		panic(err)
	}
}

func setMock(mock bool) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(context.WithValue(r.Context(), "mock", mock))
			next.ServeHTTP(w, r)
		})
	}
}

func toHandlerFunc(handler TestableHandler) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		isMock, ok := r.Context().Value("mock").(bool)
		if !ok {
			panic("illegal mock flag is passed")
		}

		handler(w, r, isMock)
	}
}

type API struct {
	router  *chi.Mux
	storage Storage
}

func NewAPI() *API {
	a := &API{}
	r := chi.NewRouter()

	val, ok := os.LookupEnv("MOCK")
	isMock := ok && val != ""
	r.Use(setMock(isMock))
	r.Use(middleware.Logger)

	cross := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
	})
	r.Use(cross.Handler)

	r.Route("/api", func(r chi.Router) {
		r.Get("/packages", toHandlerFunc(a.getPackages))
	})

	var s Storage
	if isMock {
		s = NewMockStorage()
	} else {
		s = NewBoltStorage()
	}

	wd, _ := os.Getwd()
	if !strings.HasSuffix(wd, "ligath") {
		fmt.Println("Please run it in repository root as working directory.")
		fmt.Printf("Current: %s\n", wd)
		os.Exit(1)
	}

	r.Get("/css/*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(wd, "frontend/dist/", r.URL.Path[1:]))
	}))
	r.Get("/js/*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(wd, "frontend/dist/", r.URL.Path[1:]))
	}))
	r.Get("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(wd, "frontend/dist/index.html"))
	}))

	a.router = r
	a.storage = s
	return a
}

func (a *API) Serve() error {
	err := http.ListenAndServe("0.0.0.0:3939", a.router)
	if err != nil {
		return fmt.Errorf("error while serving: %v", err)
	}
	return nil
}

func (a *API) getPackages(w http.ResponseWriter, r *http.Request, isMock bool) {
	render.JSON(w, r, a.storage.GetPackages())
}