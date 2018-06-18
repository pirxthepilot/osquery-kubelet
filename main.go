package main

import (
	"flag"
	"log"
	"time"

	"github.com/kolide/osquery-go"
	"github.com/kolide/osquery-go/plugin/table"

	"github.com/pirxthepilot/osquery-kubelet/internal/kubepod"
)

func main() {

	var (
		flSocketPath = flag.String("socket", "", "Path to socket")
		_            = flag.Int("timeout", 0, "Timeout")
		_            = flag.Int("interval", 0, "Interval")
		_            = flag.Bool("verbose", false, "Verbose")
	)
	flag.Parse()

	// allow for osqueryd to create the socket path
	time.Sleep(2 * time.Second)

	server, err := osquery.NewExtensionManagerServer(
		"osquery-kubelet",
		*flSocketPath,
	)
	if err != nil {
		log.Fatalf("Error creating extension: %s\n", err)
	}

	// Create and register a new table plugin with the server.
	// table.NewPlugin requires the table plugin name,
	// a slice of Columns and a Generate function.
	server.RegisterPlugin(table.NewPlugin(
		"kubelet_pods",
		kubepod.KubePodColumns(),
		kubepod.KubePodGenerate,
	))
	if err := server.Run(); err != nil {
		log.Fatalln(err)
	}

}
