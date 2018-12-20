package main

import (
	"flag"
	"fmt"
	"github.com/chfanghr/backend/raspi"
	"github.com/felixge/tcpkeepalive"
	"log"
	"net/rpc"
	"net/rpc/jsonrpc"
	"time"
)

var serviceName = "backendService"
var logger *log.Logger
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
	flag.StringVar(&serviceName, "serviceName", serviceName, "name of rpc service")
	flag.Parse()

	var err error
	logger, err = SetupLogger(*logFilePath, !*closeStdio)
	if err != nil {
		log.SetFlags(log.LstdFlags)
		log.Fatalln("error occur while setting up logger :", err)
	}
	CleanUpFuncs = NewCleanUpHandlerArray(logger)

	if *daemonize {
		logger.Println("setting up daemon")
		_, err := daemon(0, 0)
		if err != nil {
			logger.Fatalln("error occur while setting up daemon :", err)
		}
	}

	logger.Println("setting up pid file")
	err = SetupPidFile(*pidFilePath)
	if err != nil {
		logger.Fatalln("error occur while setting up pid file :", err)
	}

	CleanUpFuncs.Add(raspi.CleanUpHandler)

	logger.Println("loading carService")
	carService, err := LoadCarService(*configFile)
	//carService := car.Service(car.NewGeneralServiceHandler(NewFakeCar(logger)))
	if err != nil {
		logger.Fatalln("error occur while loading carService :", err)
	}
	if carService == nil {
		logger.Fatalln("car.Service should not be nil,panic")
	}

	logger.Println("register rpc service")
	err = rpc.RegisterName(serviceName, carService)
	if err != nil {
		logger.Fatalln("error occur while register rpc service :", err)
	}

	logger.Println("setting up listener")
	l, err := SetupListener(*listenNetwork, *listenAddress)
	if err != nil {
		logger.Fatalln("error occur while setting up listener :", err)
	}

	logger.Println("service is up")
	logger.Println("listen on :", l.Addr().String())
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
				logger.Println("error occur while accepting connection :", err)
				if conn != nil {
					conn.Close()
				}
				continue
			}
			logger.Println("connection accepted :", conn.RemoteAddr())
			select {
			case <-ch:
				return
			default:
			}
			connKeepAlive, err := conn, tcpkeepalive.SetKeepAlive(conn, *networkTimeout, 4, time.Second)
			if err != nil {
				logger.Println("error occur while setting up connection :", err)
				if connKeepAlive != nil {
					connKeepAlive.Close()
				}
				continue
			}
			go func() {
				defer func() {
					if connKeepAlive != nil {
						logger.Println("connection closed :", connKeepAlive.RemoteAddr())
						connKeepAlive.Close()
					}
				}()
				select {
				case <-ch:
					return
				default:
				}
				logger.Println("serve connection :", connKeepAlive.RemoteAddr())
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
