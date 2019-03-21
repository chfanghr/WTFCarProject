package main

import (
	"flag"
	"fmt"
	"github.com/chfanghr/WTFCarProject/raspi"
	"github.com/felixge/tcpkeepalive"
	"log"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"
	"runtime/pprof"
	"time"
)

var ServiceName = "backendService"
var Logger *log.Logger
var CleanUpFuncs *CleanUpHandlerArray

func main() {
	pidFilePath := flag.String("pidFile", fmt.Sprint(""), "path to pid file")
	configFile := flag.String("configFile", "", "path to config file")
	logFilePath := flag.String("logFile", "", "path to log file")
	daemonize := flag.Bool("daemonize", false, "daemonize or not")
	closeStdio := flag.Bool("closeStdio", false, "daemon close stdio or not")
	listenNetwork := flag.String("listenNetwork", "tcp", "rpc server listen to which type of network")
	listenAddress := flag.String("listenAddress", ":8888", "rpc server listen to which address")
	networkTimeout := flag.Duration("networkTimeout", time.Second*5, "connection timeout")
	cpuProfile := flag.String("cpuProfile", "", "path to cpu profile")
	flag.StringVar(&ServiceName, "serviceName", ServiceName, "name of rpc service")
	flag.Parse()

	if *cpuProfile != "" {
		profile, err := os.Create(*cpuProfile)
		if err != nil {
			log.Fatalln("error occur while creating profile :", err)
		}
		pprof.StartCPUProfile(profile)
		defer pprof.StopCPUProfile()
	}
	var err error
	Logger, err = SetupLogger(*logFilePath, !*closeStdio)
	if err != nil {
		log.SetFlags(log.LstdFlags)
		log.Fatalln("error occur while setting up Logger :", err)
	}
	CleanUpFuncs = NewCleanUpHandlerArray(Logger)

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

	CleanUpFuncs.Add(raspi.CleanUpHandler)

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
	//getDeadline := func() time.Time { return time.Now().Add(*networkTimeout) }

	go func() {
		ch := make(chan int)
		CleanUpFuncs.Add(func(i *log.Logger) {
			close(ch)
			i.Println("close listener")
			l.Close()
		})
		for {
			select {
			case <-ch:
				return
			default:
			}
			conn, err := l.Accept()
			select {
			case <-ch:
				return
			default:
			}
			if err != nil {
				Logger.Println("error occur while accepting connection :", err)
				if conn != nil {
					conn.Close()
				}
				continue
			}
			Logger.Println("connection accepted :", conn.RemoteAddr())
			select {
			case <-ch:
				return
			default:
			}
			connKeepAlive, err := conn, tcpkeepalive.SetKeepAlive(conn, *networkTimeout, 4, time.Second)
			if err != nil {
				Logger.Println("error occur while setting up connection :", err)
				if connKeepAlive != nil {
					connKeepAlive.Close()
				}
				continue
			}
			go func() {
				defer func() {
					if connKeepAlive != nil {
						Logger.Println("connection closed :", connKeepAlive.RemoteAddr())
						connKeepAlive.Close()
					}
				}()
				select {
				case <-ch:
					return
				default:
				}
				Logger.Println("serve connection :", connKeepAlive.RemoteAddr())
				rpc.ServeCodec(jsonrpc.NewServerCodec(connKeepAlive))
			}()
		}
	}()
	defer func() {
		if CleanUpFuncs != nil {
			CleanUpFuncs.Wait()
		}
	}()
}
