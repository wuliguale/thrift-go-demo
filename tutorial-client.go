package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"gen-go-thrift-tutorial/common"
	"gen-go-thrift-tutorial/tutorial"
	"github.com/apache/thrift/lib/go/thrift"
	"os"
)

var defaultCtx = context.Background()
func main() {
	flag.Usage = common.Usage
	//server := flag.Bool("server", false, "Run server")
	protocol := flag.String("P", "binary", "Specify the protocol (binary, compact, json, simplejson)")
	framed := flag.Bool("framed", false, "Use framed transport")
	buffered := flag.Bool("buffered", false, "Use buffered transport")
	addr := flag.String("addr", "localhost:9090", "Address to listen to")
	secure := flag.Bool("secure", false, "Use tls secure transport")

	flag.Parse()


	//1.根据需要的传输协议创建protocolFactory
	var protocolFactory thrift.TProtocolFactory
	switch *protocol {
	case "compact":
		protocolFactory = thrift.NewTCompactProtocolFactory()
	case "simplejson":
		protocolFactory = thrift.NewTSimpleJSONProtocolFactory()
	case "json":
		protocolFactory = thrift.NewTJSONProtocolFactory()
	case "binary", "":
		protocolFactory = thrift.NewTBinaryProtocolFactoryDefault()
	default:
		fmt.Fprint(os.Stderr, "Invalid protocol specified", protocol, "\n")
		common.Usage()
		os.Exit(1)
	}

	//2.根据需要创建transportFactory
	var transportFactory thrift.TTransportFactory
	if *buffered {
		transportFactory = thrift.NewTBufferedTransportFactory(8192)
	} else {
		transportFactory = thrift.NewTTransportFactory()
	}
	if *framed {
		transportFactory = thrift.NewTFramedTransportFactory(transportFactory)
	}

	//3.根据需要创建socket (socket, ssl socket)
	var transport thrift.TTransport
	var err error
	if *secure {
		cfg := new(tls.Config)
		cfg.InsecureSkipVerify = true
		transport, err = thrift.NewTSSLSocket(*addr, cfg)
	} else {
		transport, err = thrift.NewTSocket(*addr)
	}
	if err != nil {
		fmt.Println("Error opening socket:", err)
	}

	//4. 根据transportFactory和socket得到transport
	transport, err = transportFactory.GetTransport(transport)
	if err != nil {
		fmt.Println(err)
	}

	//5. transport open
	if err := transport.Open(); err != nil {
		fmt.Println(err)
	}
	//6. transport close
	defer transport.Close()

	//7. 创建client
	client := tutorial.NewCalculatorClientFactory(transport, protocolFactory)

	//8. client调用服务方法
	client.Ping(defaultCtx)
	fmt.Println("ping()")

	sum, _ := client.Add(defaultCtx, 1, 1)
	fmt.Print("1+1=", sum, "\n")


	work := tutorial.NewWork()
	work.Op = tutorial.Operation_DIVIDE
	work.Num1 = 1
	work.Num2 = 0
	quotient, err := client.Calculate(defaultCtx, 1, work)
	if err != nil {
		switch v := err.(type) {
		case *tutorial.InvalidOperation:
			fmt.Println("Invalid operation:", v)
		default:
			fmt.Println("Error during operation:", err)
		}
	} else {
		fmt.Println("Whoa we can divide by 0 with new value:", quotient)
	}


	work.Op = tutorial.Operation_SUBTRACT
	work.Num1 = 15
	work.Num2 = 10
	diff, err := client.Calculate(defaultCtx, 1, work)
	if err != nil {
		switch v := err.(type) {
		case *tutorial.InvalidOperation:
			fmt.Println("Invalid operation:", v)
		default:
			fmt.Println("Error during operation:", err)
		}
	} else {
		fmt.Print("15-10=", diff, "\n")
	}


	log, err := client.GetStruct(defaultCtx, 1)
	if err != nil {
		fmt.Println("Unable to get struct:", err)
	} else {
		fmt.Println("Check log:", log.Value)
	}

}



