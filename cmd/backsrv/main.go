package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/chfanghr/WTFCarProject/raspi"
	"github.com/chfanghr/cleanuphandler"
	"log"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"
	"runtime/pprof"
)

var ServiceName = "backendService"
var Logger *log.Logger
var fca string

func main() {
	pidFilePath := flag.String("pidFile", fmt.Sprint(""), "path to pid file")
	configFile := flag.String("configFile", "", "path to config file")
	logFilePath := flag.String("logFile", "", "path to log file")
	daemonize := flag.Bool("daemonize", false, "daemonize or not")
	closeStdio := flag.Bool("closeStdio", false, "daemon close stdio or not")
	listenNetwork := flag.String("listenNetwork", "tcp", "rpc server listen to which type of network")
	listenAddress := flag.String("listenAddress", "0.0.0.0:8888", "rpc server listen to which address")
	fakeCarAddress := flag.String("fakeCarAddress", "0.0.0.0:8887", "address of fakeCar's http service ")
	//networkTimeout := flag.Duration("networkTimeout", time.Second*5, "connection timeout")
	cpuProfile := flag.String("cpuProfile", "", "path to cpu profile")
	flag.StringVar(&ServiceName, "serviceName", ServiceName, "name of rpc service")
	flag.Parse()

	fca = *fakeCarAddress

	if *cpuProfile != "" {
		profile, err := os.Create(*cpuProfile)
		if err != nil {
			log.Fatalln("error occur while creating profile :", err)
		}
		_ = pprof.StartCPUProfile(profile)
		defer pprof.StopCPUProfile()
	}

	var err error
	Logger, err = SetupLogger(*logFilePath, !*closeStdio)
	if err != nil {
		log.SetFlags(log.LstdFlags)
		log.Fatalln("error occur while setting up Logger :", err)
	}
	cleanuphandler.SetLogger(Logger)

	if *daemonize {
		Logger.Println("setting up daemon")
		_, err := daemon(0, 0)
		if err != nil {
			Logger.Fatalln("error occur while setting up daemon :", err)
		}
	}

	Logger.Println("setting up pid file")
	err = SetupPidFile(*pidFilePath)
	if err != nil {
		Logger.Fatalln("error occur while setting up pid file :", err)
	}

	cleanuphandler.AddCleanupHandlers(raspi.CleanUpHandler)

	Logger.Println("loading carService")
	carService, err := LoadCarService(*configFile)
	//carService := car.Service(car.NewGeneralServiceHandler(NewFakeCar(Logger)))
	if err != nil {
		Logger.Fatalln("error occur while loading carService :", err)
	}
	if carService == nil {
		Logger.Fatalln("car.Service should not be nil,panic")
	}

	Logger.Println("register rpc service")
	err = rpc.RegisterName(ServiceName, carService)
	if err != nil {
		Logger.Fatalln("error occur while register rpc service :", err)
	}

	Logger.Println("setting up listener")
	l, err := SetupListener(*listenNetwork, *listenAddress)
	if err != nil {
		Logger.Fatalln("error occur while setting up listener :", err)
	}

	Logger.Println("service is up")
	Logger.Println("listen on :", l.Addr().String())

	defer func() {
		os.Exit(0)
		//TODO
		//cleanuphandler.Wait()
	}()

	go func() {
		defer func() { _ = l.Close() }()
		ctx, cancel := context.WithCancel(context.Background())
		cleanuphandler.AddCleanupHandlers(func(logger *log.Logger) {
			cancel()
		})
		for {
			select {
			case <-ctx.Done():
				return
			case <-func() chan struct{} {
				ch := make(chan struct{})
				go func() {
					conn, err := l.Accept()
					if err != nil {
						select {
						case <-ctx.Done():
							return
						default:
						}
						Logger.Println("error occur while accepting connection :", err)
						if conn != nil {
							_ = conn.Close()
							return
						}
					}
					select {
					case <-ctx.Done():
						return
					default:
					}
					Logger.Println("serve connection :", conn.RemoteAddr())
					rpc.ServeCodec(jsonrpc.NewServerCodec(conn))
					ch <- struct{}{}
				}()
				return ch
			}():
			}
		}
	}()
}
