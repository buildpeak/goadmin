package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"goadmin-backend/internal/cmd/api"
)

//nolint:cyclop // This is a command line tool
func main() {
	app := flag.String("app", "api", "the app to run")

	flag.Parse()

	if app == nil || *app == "" {
		*app = "api"
	}

	var cfg map[string]interface{}

	switch *app {
	case "api":
		cfg = parseAPIConfig()
	default:
		log.Fatalf("unknown app: %s", *app)
	}

	prog := os.Args[0]

	//nolint:gomnd // This is a command line tool
	if len(os.Args) < 2 {
		//nolint:forbidigo // This is a command line tool
		fmt.Printf("usage: %s <command> [args]\n", path.Base(prog))
		os.Exit(1)
	}

	cmd := os.Args[1]

	switch cmd {
	case "get":
		//nolint:gomnd // This is a command line tool
		if len(os.Args) < 3 {
			Printf("usage: %s get <key>", prog)
		}

		key := os.Args[2]

		// find value of key
		var value interface{}

		value = cfg

		keys := strings.Split(key, ".")

		//nolint:varnamelen // We need the index
		for i, k := range keys {
			//nolint:forcetypeassert // We know that value is a map
			if i == len(keys)-1 {
				//nolint:forcetypeassert,errcheck // We know that value is a string
				value = value.(map[string]interface{})[k].(string)
			} else {
				//nolint:forcetypeassert,errcheck // We know that value is a map
				value = value.(map[string]interface{})[k].(map[string]interface{})
			}
		}

		//nolint:forcetypeassert // We know that value is a string
		Printf("%s\n", value.(string))

		return
	case "show":
		Printf("%s\n", jsonfy(cfg))
	}
}

func parseAPIConfig() map[string]interface{} {
	// Parse the config file
	cfg, err := api.NewConfig()
	if err != nil {
		log.Fatalf("error parsing config: %v", err)
	}

	jsonStr := jsonfy(cfg)

	var mcfg map[string]interface{}

	err = json.Unmarshal([]byte(jsonStr), &mcfg)
	if err != nil {
		log.Fatalf("error unmarshalling: %v", err)
	}

	return mcfg
}

func Printf(format string, a ...any) {
	//nolint:forbidigo // This is a command line tool
	fmt.Printf(format, a...)
}

func jsonfy(v interface{}) string {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Sprintf("error marshalling: %v", err)
	}

	return string(b)
}
