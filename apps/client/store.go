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

func newStore() gen.ProcessBehavior {
	return &store{}
}

type store struct {
	act.Pool
	router    *mux.Router
	getRouter *mux.Router
}

func (s *store) Init(args ...any) (opts act.PoolOptions, err error) {
	s.initRouter()

	if err = s.registerStoreHandler(); err != nil {
		return opts, err
	}

	if err = s.registerFeed(); err != nil {
		return opts, err
	}

	s.registerStatic()

	if err = s.startWebServer(args[0]); err != nil {
		return opts, err
	}

	opts.WorkerFactory = newStoreHandler
	return opts, nil
}

func (s *store) startWebServer(args any) error {
	cfg, ok := args.(*config)
	if !ok {
		s.Log().Error("invalid config: %#v", args)
		return errors.New("client.store: invalid client config")
	}

	var webOptions meta.WebServerOptions
	webOptions.Port = uint16(cfg.Port)
	webOptions.Host = cfg.Host
	webOptions.Handler = s.router

	webserver, err := meta.CreateWebServer(webOptions)
	if err != nil {
		s.Log().Error("client.store: failed to start web server %s", err)
		return err
	}

	if _, err = s.SpawnMeta(webserver, gen.MetaOptions{}); err != nil {
		webserver.Terminate(err)
		return err
	}

	s.Log().Info("client.store: web server started [//%s:%d/]", webOptions.Host, webOptions.Port)
	return nil
}

func (s *store) initRouter() {
	s.router = mux.NewRouter()
	s.router.Use(middleware)
	s.getRouter = s.router.Methods("GET").Subrouter()
}

func (s *store) registerFeed() error {
	feedHandler := sse.CreateHandler(sse.HandlerOptions{
		ProcessPool: []gen.Atom{"client_feed"},
		Heartbeat:   15 * time.Second,
	})

	if _, err := s.SpawnMeta(feedHandler, gen.MetaOptions{}); err != nil {
		s.Log().Error("client.store: failed to register feed %s", err)
		return err
	}

	s.getRouter.Handle("/feed", feedHandler)
	return nil
}

func (s *store) registerStoreHandler() error {
	storeHandler := meta.CreateWebHandler(meta.WebHandlerOptions{})
	_, err := s.SpawnMeta(storeHandler, gen.MetaOptions{})

	if err != nil {
		s.Log().Error("client.store: failed to register store handler %s", err)
	} else {
		s.getRouter.Handle("/store", storeHandler)
	}
	return err
}

func (s *store) registerStatic() {
	assetDir, _ := filepath.Abs("./client/dist")
	indexFile := filepath.Join(assetDir, "index.html")

	assetHandler := http.FileServer(http.Dir(assetDir))
	s.getRouter.PathPrefix("/assets/").Handler(assetHandler)

	indexHandler := func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, indexFile)
	}

	s.getRouter.NotFoundHandler = http.HandlerFunc(indexHandler)
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
