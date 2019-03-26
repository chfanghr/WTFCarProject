package simplecli

import (
	"flag"
	"fmt"
	"github.com/chfanghr/WTFCarProject/car"
	"github.com/chfanghr/WTFCarProject/location"
	"log"
	"strconv"
)

var (
	host = flag.String("host", "localhost", "hostname/ip of backend")
	port = flag.Uint64("port", 8888, "port of backend")
	ns   = flag.String("ns", "backendService", "name of rpc service")
)

func main() {
	flag.Parse()
	if flag.NArg() == 0 {
		log.Fatalln("no command specific")
	}
	cli := car.NewGeneralClient("tcp", fmt.Sprintf("%s:%d", *host, *port), *ns)
	if err := cli.IsServiceAvailable(); err != nil {
		log.Fatalln(err)
	}

	switch flag.Arg(0) {
	case "GetLocation":
		p, err := cli.GetLocation()
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("location:x=%f,y=%f\n", p.GetX(), p.GetY())
		return
	case "MoveTo":
		if flag.NArg() < 3 {
			log.Fatalln("no enough parameters")
		}

		x, err1 := strconv.ParseFloat(flag.Arg(1), 64)
		y, err2 := strconv.ParseFloat(flag.Arg(2), 64)

		if err1 != nil || err2 != nil {
			log.Fatalln(err1, err2)
		}

		dest := location.NewPoint2D(x, y)
		if err := cli.MoveTo(*dest); err != nil {
			log.Fatalln(err)
		}
		log.Printf("move to:x=%f,y=%f\n", x, y)
		return
	case "StopMovement":
		if err := cli.StopMovement(); err != nil {
			log.Fatalln(err)
		}
		log.Printf("movement stopped\n")
		return
	default:
		log.Fatalln("unknown command")
	}

}
