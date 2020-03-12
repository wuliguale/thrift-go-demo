package main

import (
	"context"
	"fmt"
	"gen-go-thrift-tutorial/common"
	"gen-go-thrift-tutorial/format"
	"github.com/apache/thrift/lib/go/thrift"
	"log"
	"strings"
)

//1. handler
type FormatDataImpl struct {}

//3. method
func (fdi *FormatDataImpl) DoFormat(ctx context.Context, data *format.Data) (r *format.Data, err error){
	var rData format.Data
	rData.Text = strings.ToUpper(data.Text)

	return &rData, nil
}

func (fdi *FormatDataImpl) Hello(ctx context.Context, data *format.Data) (r *format.Data, err error) {
	data.Text += " hello"
	return data, err
}


func main() {
	//1. protocolFactory
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()

	//2. transportFactory
	transportFactory := thrift.NewTTransportFactory()

	//3. socket
	serverTransport, err := thrift.NewTServerSocket(common.HOST + ":" + common.PORT)
	if err != nil {
		log.Fatalln("Error:", err)
	}

	//4. xxx handler
	handler := &FormatDataImpl{}

	//5. xxx processor
	processor := format.NewFormatDataProcessor(handler)

	//6. NewTSimpleServer4
	server := thrift.NewTSimpleServer4(processor, serverTransport, transportFactory, protocolFactory)
	fmt.Println("Running at:", common.HOST + ":" + common.PORT)

	//7. serve
	server.Serve()
}
