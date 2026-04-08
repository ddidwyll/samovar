package api

import (
	"net/http"
	"path/filepath"
	"time"

	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"
	"ergo.services/ergo/meta"
	"ergo.services/meta/sse"

	"github.com/gorilla/mux"
)

func createWeb() gen.ProcessBehavior {
	return &web{}
}

type web struct {
	act.Pool
}

// Init invoked on a start this process.
func (w *web) Init(args ...any) (act.PoolOptions, error) {
	var webOptions meta.WebServerOptions
	var poolOptions act.PoolOptions

	router := mux.NewRouter()
	getRouter := router.Methods("GET").Subrouter()

	// sse handler
	sseHandler := sse.CreateHandler(sse.HandlerOptions{
		Heartbeat: 15 * time.Second,
	})
	_, err := w.SpawnMeta(sseHandler, gen.MetaOptions{})
	if err != nil {
		w.Log().Error("unable to spawn SSEHandler meta-process: %s", err)
		return poolOptions, err
	}

	getRouter.Handle("/events", sseHandler)

	err = w.registerStores(getRouter)
	if err != nil {
		return poolOptions, err
	}

	registerStatic(getRouter)

	webOptions.Port = 4000
	webOptions.Host = "localhost"

	webOptions.Handler = router

	webserver, err := meta.CreateWebServer(webOptions)
	if err != nil {
		w.Log().Error("unable to create web server meta-process: %s", err)
		return poolOptions, err
	}
	webserverid, err := w.SpawnMeta(webserver, gen.MetaOptions{})
	if err != nil {
		// invoke Terminate to close listening socket
		webserver.Terminate(err)
		return poolOptions, err
	}

	proto := "http"
	if webOptions.CertManager != nil {
		proto = "https"
	}
	w.Log().Info("started web server %s: use %s://%s:%d/", webserverid, proto, webOptions.Host, webOptions.Port)

	poolOptions.WorkerFactory = createWebWorker
	return poolOptions, nil
}

func (w *web) registerStores(router *mux.Router) error {
	storesHandler := meta.CreateWebHandler(meta.WebHandlerOptions{})
	_, err := w.SpawnMeta(storesHandler, gen.MetaOptions{})

	if err != nil {
		w.Log().Error("unable to spawn WebHandler meta-process: %s", err)
	} else {
		router.Handle("/stores/{store}", storesHandler)
	}
	return err
}

func registerStatic(router *mux.Router) {
	assetDir, _ := filepath.Abs("./client/dist")
	indexFile := filepath.Join(assetDir, "index.html")

	assetHandler := http.FileServer(http.Dir(assetDir))
	router.PathPrefix("/assets/").Handler(assetHandler)

	indexHandler := func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, indexFile)
	}

	router.NotFoundHandler = http.HandlerFunc(indexHandler)
}

const htmlPage = `<!DOCTYPE html>
<html>
<head>
    <title>SSE Example</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, monospace;
            max-width: 800px;
            margin: 50px auto;
            padding: 20px;
            background: #1a1a2e;
            color: #eee;
        }
        h1 { color: #00d9ff; }
        .status {
            padding: 10px 20px;
            border-radius: 5px;
            margin: 20px 0;
            font-weight: bold;
        }
        .connected { background: #0a3d0a; color: #4caf50; }
        .disconnected { background: #3d0a0a; color: #f44336; }
        .info {
            background: #16213e;
            padding: 20px;
            border-radius: 10px;
            margin: 20px 0;
        }
        .info h3 { margin-top: 0; color: #00d9ff; }
        .counter { font-size: 48px; color: #00d9ff; }
        #events {
            background: #0f0f23;
            border: 1px solid #333;
            border-radius: 5px;
            padding: 15px;
            height: 300px;
            overflow-y: auto;
            font-family: monospace;
            font-size: 13px;
        }
        .event { margin: 5px 0; padding: 5px; border-left: 3px solid #00d9ff; padding-left: 10px; }
        .event-welcome { border-color: #4caf50; }
        .event-time { border-color: #ff9800; }
        .timestamp { color: #666; }
    </style>
</head>
<body>
    <h1>SSE Server Example</h1>

    <div id="status" class="status disconnected">Disconnected</div>

    <div class="info">
        <h3>Live Data</h3>
        <div>Counter: <span id="counter" class="counter">-</span></div>
        <div>Server Time: <span id="time">-</span></div>
        <div>Connected Clients: <span id="clients">-</span></div>
    </div>

    <h3>Event Log</h3>
    <div id="events"></div>

    <script>
        const statusEl = document.getElementById('status');
        const eventsEl = document.getElementById('events');
        const counterEl = document.getElementById('counter');
        const timeEl = document.getElementById('time');
        const clientsEl = document.getElementById('clients');

        function connect() {
            const es = new EventSource('/events');

            es.onopen = () => {
                statusEl.textContent = 'Connected';
                statusEl.className = 'status connected';
                addEvent('system', 'Connection established');
            };

            es.onerror = () => {
                statusEl.textContent = 'Disconnected (reconnecting...)';
                statusEl.className = 'status disconnected';
                addEvent('system', 'Connection lost');
            };

            es.addEventListener('welcome', (e) => {
                addEvent('welcome', e.data);
            });

            es.addEventListener('time', (e) => {
                const data = JSON.parse(e.data);
                counterEl.textContent = data.counter;
                timeEl.textContent = data.time;
                clientsEl.textContent = data.clients;
                addEvent('time', e.data);
            });
        }

        function addEvent(type, data) {
            const now = new Date().toLocaleTimeString();
            const div = document.createElement('div');
            div.className = 'event event-' + type;
            div.innerHTML = '<span class="timestamp">[' + now + ']</span> <strong>' + type + ':</strong> ' + data;
            eventsEl.insertBefore(div, eventsEl.firstChild);

            // Keep only last 50 events
            while (eventsEl.children.length > 50) {
                eventsEl.removeChild(eventsEl.lastChild);
            }
        }

        connect();
    </script>
</body>
</html>`
