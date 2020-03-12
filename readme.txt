

go使用thrift的例子

1.环境准备
安装go
安装thrift编译器（可以执行thrift命令自动生成代码）

2.项目初始化
创建项目gen-go-thrift-tutorial
go mod init xxx
go get xxx/thrift       //安装thrift库，不同于thrift编译器

3.使用.thrift定义服务
format.thrift   //format服务
shared.thrift  //shared服务
tutorial.thrift //tutorial服务，依赖shared服务

4.thrift命令生成
thrift --gen go format.thrift
thrift -r --gen go tutorial.thrift      //有依赖，使用 -r
将生成的代码放到方便管理的目录
自动生成的代码可能有一些错误需要手动修改，比如包名错误，方法中缺少ctx参数等

5.实现server端handler
format-server.go
定义handler结构体FormatDataImpl
以handler为接收子实现.thrift中定义的方法, xxx.DoFormat(), xxx.Hello()

tutorial-server.go
定义handler结构体CalculatorHandler
以handler为接收子实现.thrift中定义的方法，xxx.Ping(), xxx.Add() ...

6.启动server端，提供服务(一直运行)
protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
transportFactory := thrift.NewTTransportFactory()
transport := thrift.NewTServerSocket(*addr)    //此处有ip,port
handler := NewCalculatorHandler()
processor := NewCalculatorProcessor(handler)
server := thrift.NewTSimpleServer4(processor, transport, transportFactory, protocolFactory)
server.Serve()

7. client端请求服务(一次结束)
protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
transportFactory := thrift.NewTTransportFactory()
socket, err = thrift.NewTSocket(*addr)           //此处有ip,port
transport, err = transportFactory.GetTransport(socket)
err := transport.Open()
defer transport.Close()
client := tutorial.NewCalculatorClientFactory(transport, protocolFactory)
client.Ping(ctx)


