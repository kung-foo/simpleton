package main

import (
	"crypto/rand"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/docopt/docopt-go"
	"github.com/julienschmidt/httprouter"
)

// VERSION is set by the makefile
var VERSION = "0.0.0"

var usage = `
Usage:
    simpleton [options] [-v ...] [--port=<port>]
    simpleton -h | --help | --version

Options:
    -h --help               Show this screen.
    --version               Show version.
    --port=<port>           Port to listen on [default: 8080].
    -v                      Enable verbose logging (-vv for very verbose).
`

type config struct {
	verbose     bool
	veryVerbose bool
}

// global config struct
var Config = &config{}

var (
	random1kBuffer = make([]byte, 1024)
	ok             = []byte("OK")
)

func main() {
	mainEx(os.Args[1:])
}

func mainEx(argv []string) {
	args, err := docopt.Parse(usage, argv, true, VERSION, false)

	if err != nil {
		log.Fatal(err)
	}

	Config.verbose = args["-v"].(int) > 0
	Config.veryVerbose = args["-v"].(int) > 1

	if Config.veryVerbose {
		log.SetLevel(log.DebugLevel)
		log.Debugf("args: %v", args)
	}

	port, err := strconv.Atoi(args["--port"].(string))

	if err != nil {
		log.Fatal(err)
	}

	if Config.verbose {
		log.Infof("Build:        %s", VERSION)
		log.Infof("NumCPU:       %d", runtime.NumCPU())
		log.Infof("GOMAXPROCS:   %d", runtime.GOMAXPROCS(0))
	}

	rand.Read(random1kBuffer)

	router := httprouter.New()
	router.GET("/", index)
	router.GET("/data/:size", dataNKB)
	router.GET("/sleep/:time", sleepNms)
	router.PUT("/data/null", dataNull)
	router.POST("/data/null", dataNull)

	log.Infof("Listening on: %d", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), router))
}

func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if Config.veryVerbose {
		log.Debugf("%+v", r)
	}
	w.Write(ok)
}

func dataNKB(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if Config.veryVerbose {
		log.Debugf("%+v", r)
	}
	if sz, err := strconv.Atoi(p.ByName("size")); !returnOnError(err, w) {
		for i := 0; i < sz; i++ {
			w.Write(random1kBuffer)
		}
	}
}

func dataNull(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	defer r.Body.Close()
	if Config.veryVerbose {
		log.Debugf("%+v", r)
	}
	if written, err := io.Copy(ioutil.Discard, r.Body); !returnOnError(err, w) {
		fmt.Fprintf(w, "%d", written)
	}
}

func sleepNms(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if Config.veryVerbose {
		log.Debugf("%+v", r)
	}
	if t, err := strconv.Atoi(p.ByName("time")); !returnOnError(err, w) {
		time.Sleep(time.Millisecond * time.Duration(t))
		w.Write(ok)
	}
}

func returnOnError(err error, w http.ResponseWriter) bool {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return true
	}
	return false
}
