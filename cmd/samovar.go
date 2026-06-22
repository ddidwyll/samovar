package main

import (
	"flag"
	"time"

	"samovar/apps/bus"
	"samovar/apps/client"
	"samovar/apps/device"
	"samovar/apps/mqtt"
	"samovar/apps/script"
	"samovar/apps/session"
	"samovar/apps/telemetry"

	// "ergo.services/application/observer"
	"ergo.services/ergo"
	"ergo.services/ergo/gen"
	"ergo.services/ergo/lib"
	"ergo.services/logger/colored"
)

var (
	OptionNodeName   string
	OptionNodeCookie string
)

func init() {
	flag.StringVar(&OptionNodeName, "name", "samovar@localhost", "node name")
	flag.StringVar(&OptionNodeCookie, "cookie", lib.RandomString(16), "a secret cookie for the network messaging")
}

func main() {
	var options gen.NodeOptions

	flag.Parse()

	options.Applications = []gen.ApplicationBehavior{
		telemetry.CreateApp(),
		bus.CreateApp(),
		client.CreateApp(),
		script.CreateApp(),
		session.CreateApp(),
		device.CreateApp(),
		mqtt.CreateApp(),
	}

	// disable default logger to get rid of multiple logging to the os.Stdout
	options.Log.DefaultLogger.Disable = true

	// add logger "colored".
	loggercolored, err := colored.CreateLogger(colored.Options{TimeFormat: time.DateTime})
	if err != nil {
		panic(err)
	}
	options.Log.Loggers = append(options.Log.Loggers, gen.Logger{Name: "colored", Logger: loggercolored})

	// set network options
	options.Network.Cookie = OptionNodeCookie

	// starting node
	node, err := ergo.StartNode(gen.Atom(OptionNodeName), options)
	if err != nil {
		panic(err)
	}

	node.Wait()
}
