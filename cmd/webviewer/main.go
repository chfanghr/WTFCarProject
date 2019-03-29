package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"github.com/chfanghr/WTFCarProject/location"
	"github.com/chfanghr/WTFCarProject/map2d"
	"github.com/chfanghr/WTFCarProject/rpcprotocal"
	"github.com/chfanghr/cleanuphandler"
	"github.com/gorilla/websocket"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"
)

var (
	logFile       = flag.String("log", "", "path to log file")
	serveAddr     = flag.String("addr", "localhost:8886", "address to serve on")
	mapFile       = flag.String("map", "map.json", "path to <map>.json")
	fakeCarWSHost = flag.String("fakeCarWS", "localhost:8887", "hostname/IP to fakeCarService")
	backsrvAddr   = flag.String("backsrvAddr", "localhost:8888", "hostname/IP address to backsrv")
	backsrvName   = flag.String("backsrvName", "backendService", "name of backend rpc service")
	logger        *log.Logger
)

const webPage = "index.html"

var (
	webPageData []byte
	backsrv     *backend
	mapData     *map2d.Map2D
)

type backend struct {
	addr    string
	cli     *rpc.Client
	service string
	mapData *map2d.Map2D
}

func newBackend(addr string, name string, m *map2d.Map2D) (*backend, error) {
	cli, err := jsonrpc.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &backend{
		addr:    addr,
		cli:     cli,
		service: name,
		mapData: m,
	}, err
}
func (b *backend) MoveTo(p location.Point2D) error {
	tmp := rpcprotocal.Point2D{}
	err := b.cli.Call(b.service+"GetLocation", 0, &tmp)
	if err != nil {
		return err
	}
	path := b.mapData.ComputePathTo(*rpcprotocal.Point2DToLocationPoint2D(tmp), p)
	for _, v := range path {
		if err := b.cli.Call(b.service+".MoveTo", *rpcprotocal.Point2DFromLocationPoint2D(v), new(int)); err != nil {
			return err
		}
	}
	return nil
}
func (b *backend) reconnect() error {
	cli, err := jsonrpc.Dial("tcp", b.addr)
	if err != nil {
		return err
	}
	b.cli = cli
	return nil
}

type connection struct {
	*websocket.Conn
	back *backend
}
type message struct {
	MoveTo *rpcprotocal.Point2D `json:"move_to"`
}

func (c *connection) worker() {
	defer func() { _ = c.Close() }()
	for {
		_, r, err := c.NextReader()
		if err != nil {
			return
		}
		m := &message{}
		err = json.NewDecoder(r).Decode(m)
		if err != nil {
			return
		}
		if m.MoveTo != nil {
			_ = c.back.MoveTo(*location.NewPoint2D(m.MoveTo.X, m.MoveTo.Y))
		}
	}
}
func newConnection(b *backend, ws *websocket.Conn) *connection {
	tmp := &connection{
		ws, b,
	}
	go tmp.worker()
	return tmp
}

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
func setupBacksrv() {
	tmp, err := newBackend(*backsrvAddr, *backsrvName, mapData)
	if err != nil {
		logger.Fatalln(err)
	}
	backsrv = tmp
}
func serve() {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		_, err := writer.Write(webPageData)
		if err != nil {
			writer.WriteHeader(http.StatusServiceUnavailable)
		}
	})
	http.HandleFunc("/ws/", func(w http.ResponseWriter, r *http.Request) {
		up := websocket.Upgrader{}
		ws, err := up.Upgrade(w, r, nil)
		if err != nil {
			logger.Println("error upgrade ws", err)
			return
		}
		logger.Println("ws connected ", ws.RemoteAddr())
		c := newConnection(backsrv, ws)
		go c.worker()
		return
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
func processWebPage() {
	t, err := template.New("webpage").Parse(string(webPageData))
	if err != nil {
		logger.Fatalln(err)
	}
	buf := bytes.NewBuffer([]byte{})
	mapData, _ := json.Marshal(mapData.Map)
	if err = t.Execute(buf, struct {
		MapData      string
		FakeCarWSURL string
		ClientWSURL  string
	}{
		string(mapData),
		"ws://" + *fakeCarWSHost + "/",
		"ws://" + *serveAddr + "/ws/",
	}); err != nil {
		logger.Fatalln(err)
	}
	webPageData = buf.Bytes()
}
func main() {
	processWebPage()
	setupBacksrv()
	go serve()
	cleanuphandler.Wait()
}
