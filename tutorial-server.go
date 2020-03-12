package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"gen-go-thrift-tutorial/common"
	"gen-go-thrift-tutorial/shared"
	"gen-go-thrift-tutorial/tutorial"
	"github.com/apache/thrift/lib/go/thrift"
	"os"
	"strconv"
)

//1. handler
type CalculatorHandler struct {
	log map[int]*shared.SharedStruct
}

//2. new handler
func NewCalculatorHandler() *CalculatorHandler {
	return &CalculatorHandler{log: make(map[int]*shared.SharedStruct)}
}

//3. method
func (p *CalculatorHandler) Ping(ctx context.Context) (err error) {
	fmt.Print("ping()\n")
	return nil
}

func (p *CalculatorHandler) Add(ctx context.Context, num1 int32, num2 int32) (retval17 int32, err error) {
	fmt.Print("add(", num1, ",", num2, ")\n")
	return num1 + num2, nil
}

func (p *CalculatorHandler) Calculate(ctx context.Context, logid int32, w *tutorial.Work) (val int32, err error) {
	fmt.Print("calculate(", logid, ", {", w.Op, ",", w.Num1, ",", w.Num2, "})\n")
	switch w.Op {
	case tutorial.Operation_ADD:
		val = w.Num1 + w.Num2
		break
	case tutorial.Operation_SUBTRACT:
		val = w.Num1 - w.Num2
		break
	case tutorial.Operation_MULTIPLY:
		val = w.Num1 * w.Num2
		break
	case tutorial.Operation_DIVIDE:
		if w.Num2 == 0 {
			ouch := tutorial.NewInvalidOperation()
			ouch.WhatOp = int32(w.Op)
			ouch.Why = "Cannot divide by 0"
			err = ouch
			return
		}
		val = w.Num1 / w.Num2
		break
	default:
		ouch := tutorial.NewInvalidOperation()
		ouch.WhatOp = int32(w.Op)
		ouch.Why = "Unknown operation"
		err = ouch
		return
	}
	entry := shared.NewSharedStruct()
	entry.Key = logid
	entry.Value = strconv.Itoa(int(val))
	k := int(logid)
	/*
	   oldvalue, exists := p.log[k]
	   if exists {
	     fmt.Print("Replacing ", oldvalue, " with ", entry, " for key ", k, "\n")
	   } else {
	     fmt.Print("Adding ", entry, " for key ", k, "\n")
	   }
	*/
	p.log[k] = entry
	return val, err
}

func (p *CalculatorHandler) GetStruct(ctx context.Context, key int32) (*shared.SharedStruct, error) {
	fmt.Print("getStruct(", key, ")\n")
	v, _ := p.log[int(key)]
	return v, nil
}

func (p *CalculatorHandler) Zip(ctx context.Context) (err error) {
	fmt.Print("zip()\n")
	return nil
}


func main() {
	flag.Usage = common.Usage
	//server := flag.Bool("server", true, "Run server")
	protocol := flag.String("P", "binary", "Specify the protocol (binary, compact, json, simplejson)")
	framed := flag.Bool("framed", false, "Use framed transport")
	buffered := flag.Bool("buffered", false, "Use buffered transport")
	addr := flag.String("addr", "localhost:9090", "Address to listen to")
	secure := flag.Bool("secure", false, "Use tls secure transport")

	flag.Parse()

	//1. protocolFactory
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

	//2. transportFactory
	var transportFactory thrift.TTransportFactory
	if *buffered {
		transportFactory = thrift.NewTBufferedTransportFactory(8192)
	} else {
		transportFactory = thrift.NewTTransportFactory()
	}
	if *framed {
		transportFactory = thrift.NewTFramedTransportFactory(transportFactory)
	}

	//3. socket
	var transport thrift.TServerTransport
	var err error
	if *secure {
		cfg := new(tls.Config)
		if cert, err := tls.LoadX509KeyPair("server.crt", "server.key"); err == nil {
			cfg.Certificates = append(cfg.Certificates, cert)
		} else {
			fmt.Println(err)
		}
		transport, err = thrift.NewTSSLServerSocket(*addr, cfg)
	} else {
		transport, err = thrift.NewTServerSocket(*addr)
	}

	if err != nil {
		fmt.Println(err)
	}

	//4. xxxHandler
	handler := NewCalculatorHandler()

	//5. xxxProcessor
	processor := tutorial.NewCalculatorProcessor(handler)

	//6. NewTSimpleServer4
	server2 := thrift.NewTSimpleServer4(processor, transport, transportFactory, protocolFactory)
	fmt.Println("Starting the simple server... on ", *addr)

	//7. serve
	server2.Serve()
}

