package main

import (
	"flag"

	"time"

	"samovar/apps/api"

	"ergo.services/application/observer"

	"ergo.services/logger/colored"

	"ergo.services/ergo"
	"ergo.services/ergo/gen"
	"ergo.services/ergo/lib"
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

	// create applications that must be started
	apps := []gen.ApplicationBehavior{
		observer.CreateApp(observer.Options{}),
		api.CreateApi(),
	}
	options.Applications = apps

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
