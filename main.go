package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/royalbhati/yaml-conf-go/config"
)

func main() {
	log := log.New(os.Stdout, "URL SHORTENER : ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	cfgFlags, err := config.ParseFlags()
	if err != nil {
		log.Println("Config Parse error:", err)
		os.Exit(1)
	}
	cfg, err := config.ParseConfig(cfgFlags.Path)
	if err != nil {
		log.Println("Config Parse error:", err)
		os.Exit(1)
	}

	if err := run(cfg, log); err != nil {
		log.Println("main: error:", err)
		os.Exit(1)
	}
}

type App struct {
	mux *chi.Mux
}

func NewApp() *App {
	return &App{
		mux: chi.NewRouter(),
	}
}

func (h *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}

func run(cfg *config.Config, log *log.Logger) error {
	log.Println("heyyy", cfg.Web)
	log.Println("heyyy", cfg.Web.Host)

	app := NewApp()
	app.mux.MethodFunc(http.MethodGet, "/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hey wassup?"))
	})

	api := http.Server{
		Addr:         cfg.Web.Host,
		Handler:      app,
		ReadTimeout:  cfg.Web.Timeout.Read * time.Second,
		WriteTimeout: cfg.Web.Timeout.Write * time.Second,
	}

	if err := api.ListenAndServe(); err != nil {
		log.Println("ERR", err)

	}
	return nil
}
