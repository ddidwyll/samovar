package client

import (
	"errors"
	"net/http"
	"path/filepath"
	"time"

	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"
	"ergo.services/ergo/meta"
	"ergo.services/meta/sse"

	"github.com/gorilla/mux"
)

func newApi() gen.ProcessBehavior {
	return &api{}
}

type api struct {
	act.Pool
	router    *mux.Router
	getRouter *mux.Router
}

func (a *api) Init(args ...any) (opts act.PoolOptions, err error) {
	a.initRouter()

	if err = a.registerApiHandler(); err != nil {
		return opts, err
	}

	if err = a.registerFeed(); err != nil {
		return opts, err
	}

	a.registerStatic()

	if err = a.startWebServer(args[0]); err != nil {
		return opts, err
	}

	opts.WorkerFactory = newApiHandler
	return opts, nil
}

func (a *api) startWebServer(args any) error {
	cfg, ok := args.(*config)
	if !ok {
		a.Log().Error("invalid config: %#v (%v)", args, ok)
		return errors.New("invalid client api config")
	}

	var webOptions meta.WebServerOptions

	webOptions.Port = uint16(cfg.Port)
	webOptions.Host = cfg.Host
	webOptions.Handler = a.router

	webserver, err := meta.CreateWebServer(webOptions)
	if err != nil {
		a.Log().Error("client.api: failed to start web server %s", err)
		return err
	}

	if _, err = a.SpawnMeta(webserver, gen.MetaOptions{}); err != nil {
		webserver.Terminate(err)
		return err
	}

	a.Log().Info("client.api: web server started [//%s:%d/]", webOptions.Host, webOptions.Port)
	return nil
}

func (a *api) initRouter() {
	a.router = mux.NewRouter()
	a.router.Use(middleware)
	a.getRouter = a.router.Methods("GET").Subrouter()
}

func (a *api) registerFeed() error {
	feedHandler := sse.CreateHandler(sse.HandlerOptions{
		ProcessPool: []gen.Atom{"client_feed"},
		Heartbeat:   15 * time.Second,
	})

	if _, err := a.SpawnMeta(feedHandler, gen.MetaOptions{}); err != nil {
		a.Log().Error("client.api: failed to register feed %s", err)
		return err
	}

	a.getRouter.Handle("/feed", feedHandler)
	return nil
}

func (a *api) registerApiHandler() error {
	apiHandler := meta.CreateWebHandler(meta.WebHandlerOptions{})
	_, err := a.SpawnMeta(apiHandler, gen.MetaOptions{})

	if err != nil {
		a.Log().Error("client.api: failed to register api handler %s", err)
	} else {
		a.getRouter.Handle("/api", apiHandler)
	}
	return err
}

func (a *api) registerStatic() {
	assetDir, _ := filepath.Abs("./client/dist")
	indexFile := filepath.Join(assetDir, "index.html")

	assetHandler := http.FileServer(http.Dir(assetDir))
	a.getRouter.PathPrefix("/assets/").Handler(assetHandler)

	indexHandler := func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, indexFile)
	}

	a.getRouter.NotFoundHandler = http.HandlerFunc(indexHandler)
}

func middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")

		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Methods", "GET, PUT, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			w.WriteHeader(http.StatusOK)
			return
		}

		h.ServeHTTP(w, r)
	})
}
