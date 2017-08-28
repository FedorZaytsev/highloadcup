package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/ogier/pflag"
)

var (
	Version   string
	BuildTime string
	ts        time.Time
	DB        *Database
)

func init() {
	var versReq bool
	pflag.StringVarP(&configPath, "config", "c", "config.toml", "Used for set path to config file.")
	pflag.BoolVarP(&versReq, "version", "v", false, "Use for build time and version print")
	var err error
	pflag.Parse()
	if versReq {
		fmt.Println("Version: ", Version)
		fmt.Println("Build time:", BuildTime)
		os.Exit(0)
	}
	Config, err = configure()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	Log, err = initLogger()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	DB, err = DatabaseInit()
	if err != nil {
		fmt.Println("Error", err)
		os.Exit(1)
	}
	ts = time.Now()
}

func writeAnswer(w http.ResponseWriter, code int, body string) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Cache-Control", "application/json")
	w.WriteHeader(code)
	fmt.Fprintf(w, "%s", body)
}

func generateError(data string) string {
	return fmt.Sprintf("{\"error\": \"%s\"}", data)
}

func main() {
	Log.Infof("Started\n")

	go loadToServer()
	http.HandleFunc("/users/new", newUser)
	http.HandleFunc("/users/", processUser)
	http.HandleFunc("/locations/new", newLocation)
	http.HandleFunc("/locations/", processLocation)
	http.HandleFunc("/visits/new", newVisit)
	http.HandleFunc("/visits/", processVisit)
	http.ListenAndServe(":80", nil)

	//uncomment if it is a demon
	//sgnl := make(chan os.Signal, 1)
	//signal.Notify(sgnl, os.Interrupt, syscall.SIGTERM)
	//<-sgnl
}
