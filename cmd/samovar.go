package main

import (
	"flag"
	"time"

	"samovar/apps/api"
	"samovar/apps/bus"
	"samovar/apps/device"
	"samovar/apps/mqtt"

	"ergo.services/application/observer"
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
	flag.StringVar(&OptionNodeName, "name", "Samovar@localhost", "node name")
	flag.StringVar(&OptionNodeCookie, "cookie", lib.RandomString(16), "a secret cookie for the network messaging")
}

func main() {
	var options gen.NodeOptions

	flag.Parse()

	options.Applications = []gen.ApplicationBehavior{
		observer.CreateApp(observer.Options{}),
		bus.CreateBusApp(),
		api.CreateApiApp(),
		mqtt.CreateMqttApp(),
		device.CreateDeviceApp(),
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
