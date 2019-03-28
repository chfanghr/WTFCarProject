package main

import (
	"flag"
	"github.com/chfanghr/WTFCarProject/map2d"
	"github.com/chfanghr/cleanuphandler"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var (
	logFile   = flag.String("log", "", "path to log file")
	serveAddr = flag.String("addr", "localhost:8886", "address to serve on")
	mapFile   = flag.String("map", "map.json", "path to <map>.json")
	logger    *log.Logger
)

const webPage = "index.html"

var (
	webPageData []byte
	mapData     *map2d.Map2d
)

func init() {
	flag.Parse()
	logger = log.New(os.Stdout, "", log.LstdFlags)
	if *logFile != "" {
		file, err := os.OpenFile(*logFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
		if err != nil {
			logger.Fatalln(err)
		}
		logger = log.New(file, "", log.LstdFlags)
		cleanuphandler.AddCleanupHandlers(func(i *log.Logger) {
			_ = file.Close()
		})
	}
	if *mapFile == "" {
		logger.Fatalln("map file required")
	} else {
		mapRawData, err := ioutil.ReadFile(*mapFile)
		if err != nil {
			logger.Fatalln(err)
		}
		mapData, err = map2d.NewMap2d(mapRawData)
		if err != nil {
			logger.Fatalln(err)
		}
	}
	_webPageData, err := ioutil.ReadFile(webPage)
	if err != nil {
		logger.Fatalln(err)
	}
	webPageData = _webPageData
}

func serve() {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		_, err := writer.Write(webPageData)
		if err != nil {
			writer.WriteHeader(http.StatusServiceUnavailable)
		}
	})
	err := make(chan error)
	logger.Println("service start on", *serveAddr)
	go func() {
		err <- http.ListenAndServe(*serveAddr, http.DefaultServeMux)
	}()
	if e := <-err; e != nil {
		logger.Fatalln(e)
	}
}

func processWebPage(d *map2d.Map2d, h []byte) {
	//TODO HANDLE WEB PAGE HERE
}

func main() {
	processWebPage(mapData, webPageData)
	serve()
}
