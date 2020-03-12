package main

import (
	"context"
	"fmt"
	"gen-go-thrift-tutorial/common"
	"gen-go-thrift-tutorial/format"
	"github.com/apache/thrift/lib/go/thrift"
	"log"
	"net"
)



func main()  {
	//1. protocolFactory
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()

	//2. transportFactory
	transportFactory := thrift.NewTTransportFactory()

	//3. socket
	tSocket, err := thrift.NewTSocket(net.JoinHostPort(common.HOST, common.PORT))
	if err != nil {
		log.Fatalln("tSocket error:", err)
	}

	//4. transport
	transport, err := transportFactory.GetTransport(tSocket)

	//5. transport.Open
	if err := transport.Open(); err != nil {
		log.Fatalln("Error opening:", common.HOST + ":" + common.PORT)
	}

	//6. transport.Close
	defer transport.Close()

	//7. client
	client := format.NewFormatDataClientFactory(transport, protocolFactory)


	//8. call
	data := format.Data{Text:"hello,world!"}
	d, err := client.DoFormat(context.Background(), &data)
	fmt.Println(d.Text)

	data2 := format.Data{"abc"}
	d2, err := client.Hello(context.Background(), &data2)
	fmt.Println(d2, err)
}

