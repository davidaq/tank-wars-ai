// Autogenerated by Thrift Compiler (0.10.0)
// DO NOT EDIT UNLESS YOU ARE SURE THAT YOU KNOW WHAT YOU ARE DOING

package main

import (
        "flag"
        "fmt"
        "math"
        "net"
        "net/url"
        "os"
        "strconv"
        "strings"
        "lib/thrift"
        "player"
)


func Usage() {
  fmt.Fprintln(os.Stderr, "Usage of ", os.Args[0], " [-h host:port] [-u url] [-f[ramed]] function [arg1 [arg2...]]:")
  flag.PrintDefaults()
  fmt.Fprintln(os.Stderr, "\nFunctions:")
  fmt.Fprintln(os.Stderr, "  void uploadMap( gamemap)")
  fmt.Fprintln(os.Stderr, "  void uploadParamters(Args arguments)")
  fmt.Fprintln(os.Stderr, "  void assignTanks( tanks)")
  fmt.Fprintln(os.Stderr, "  void latestState(GameState state)")
  fmt.Fprintln(os.Stderr, "   getNewOrders()")
  fmt.Fprintln(os.Stderr)
  os.Exit(0)
}

func main() {
  flag.Usage = Usage
  var host string
  var port int
  var protocol string
  var urlString string
  var framed bool
  var useHttp bool
  var parsedUrl url.URL
  var trans thrift.TTransport
  _ = strconv.Atoi
  _ = math.Abs
  flag.Usage = Usage
  flag.StringVar(&host, "h", "localhost", "Specify host and port")
  flag.IntVar(&port, "p", 9090, "Specify port")
  flag.StringVar(&protocol, "P", "binary", "Specify the protocol (binary, compact, simplejson, json)")
  flag.StringVar(&urlString, "u", "", "Specify the url")
  flag.BoolVar(&framed, "framed", false, "Use framed transport")
  flag.BoolVar(&useHttp, "http", false, "Use http")
  flag.Parse()
  
  if len(urlString) > 0 {
    parsedUrl, err := url.Parse(urlString)
    if err != nil {
      fmt.Fprintln(os.Stderr, "Error parsing URL: ", err)
      flag.Usage()
    }
    host = parsedUrl.Host
    useHttp = len(parsedUrl.Scheme) <= 0 || parsedUrl.Scheme == "http"
  } else if useHttp {
    _, err := url.Parse(fmt.Sprint("http://", host, ":", port))
    if err != nil {
      fmt.Fprintln(os.Stderr, "Error parsing URL: ", err)
      flag.Usage()
    }
  }
  
  cmd := flag.Arg(0)
  var err error
  if useHttp {
    trans, err = thrift.NewTHttpClient(parsedUrl.String())
  } else {
    portStr := fmt.Sprint(port)
    if strings.Contains(host, ":") {
           host, portStr, err = net.SplitHostPort(host)
           if err != nil {
                   fmt.Fprintln(os.Stderr, "error with host:", err)
                   os.Exit(1)
           }
    }
    trans, err = thrift.NewTSocket(net.JoinHostPort(host, portStr))
    if err != nil {
      fmt.Fprintln(os.Stderr, "error resolving address:", err)
      os.Exit(1)
    }
    if framed {
      trans = thrift.NewTFramedTransport(trans)
    }
  }
  if err != nil {
    fmt.Fprintln(os.Stderr, "Error creating transport", err)
    os.Exit(1)
  }
  defer trans.Close()
  var protocolFactory thrift.TProtocolFactory
  switch protocol {
  case "compact":
    protocolFactory = thrift.NewTCompactProtocolFactory()
    break
  case "simplejson":
    protocolFactory = thrift.NewTSimpleJSONProtocolFactory()
    break
  case "json":
    protocolFactory = thrift.NewTJSONProtocolFactory()
    break
  case "binary", "":
    protocolFactory = thrift.NewTBinaryProtocolFactoryDefault()
    break
  default:
    fmt.Fprintln(os.Stderr, "Invalid protocol specified: ", protocol)
    Usage()
    os.Exit(1)
  }
  client := player.NewPlayerServerClientFactory(trans, protocolFactory)
  if err := trans.Open(); err != nil {
    fmt.Fprintln(os.Stderr, "Error opening socket to ", host, ":", port, " ", err)
    os.Exit(1)
  }
  
  switch cmd {
  case "uploadMap":
    if flag.NArg() - 1 != 1 {
      fmt.Fprintln(os.Stderr, "UploadMap requires 1 args")
      flag.Usage()
    }
    arg18 := flag.Arg(1)
    mbTrans19 := thrift.NewTMemoryBufferLen(len(arg18))
    defer mbTrans19.Close()
    _, err20 := mbTrans19.WriteString(arg18)
    if err20 != nil { 
      Usage()
      return
    }
    factory21 := thrift.NewTSimpleJSONProtocolFactory()
    jsProt22 := factory21.GetProtocol(mbTrans19)
    containerStruct0 := player.NewPlayerServerUploadMapArgs()
    err23 := containerStruct0.ReadField1(jsProt22)
    if err23 != nil {
      Usage()
      return
    }
    argvalue0 := containerStruct0.Gamemap
    value0 := argvalue0
    fmt.Print(client.UploadMap(value0))
    fmt.Print("\n")
    break
  case "uploadParamters":
    if flag.NArg() - 1 != 1 {
      fmt.Fprintln(os.Stderr, "UploadParamters requires 1 args")
      flag.Usage()
    }
    arg24 := flag.Arg(1)
    mbTrans25 := thrift.NewTMemoryBufferLen(len(arg24))
    defer mbTrans25.Close()
    _, err26 := mbTrans25.WriteString(arg24)
    if err26 != nil {
      Usage()
      return
    }
    factory27 := thrift.NewTSimpleJSONProtocolFactory()
    jsProt28 := factory27.GetProtocol(mbTrans25)
    argvalue0 := player.NewArgs_()
    err29 := argvalue0.Read(jsProt28)
    if err29 != nil {
      Usage()
      return
    }
    value0 := argvalue0
    fmt.Print(client.UploadParamters(value0))
    fmt.Print("\n")
    break
  case "assignTanks":
    if flag.NArg() - 1 != 1 {
      fmt.Fprintln(os.Stderr, "AssignTanks requires 1 args")
      flag.Usage()
    }
    arg30 := flag.Arg(1)
    mbTrans31 := thrift.NewTMemoryBufferLen(len(arg30))
    defer mbTrans31.Close()
    _, err32 := mbTrans31.WriteString(arg30)
    if err32 != nil { 
      Usage()
      return
    }
    factory33 := thrift.NewTSimpleJSONProtocolFactory()
    jsProt34 := factory33.GetProtocol(mbTrans31)
    containerStruct0 := player.NewPlayerServerAssignTanksArgs()
    err35 := containerStruct0.ReadField1(jsProt34)
    if err35 != nil {
      Usage()
      return
    }
    argvalue0 := containerStruct0.Tanks
    value0 := argvalue0
    fmt.Print(client.AssignTanks(value0))
    fmt.Print("\n")
    break
  case "latestState":
    if flag.NArg() - 1 != 1 {
      fmt.Fprintln(os.Stderr, "LatestState requires 1 args")
      flag.Usage()
    }
    arg36 := flag.Arg(1)
    mbTrans37 := thrift.NewTMemoryBufferLen(len(arg36))
    defer mbTrans37.Close()
    _, err38 := mbTrans37.WriteString(arg36)
    if err38 != nil {
      Usage()
      return
    }
    factory39 := thrift.NewTSimpleJSONProtocolFactory()
    jsProt40 := factory39.GetProtocol(mbTrans37)
    argvalue0 := player.NewGameState()
    err41 := argvalue0.Read(jsProt40)
    if err41 != nil {
      Usage()
      return
    }
    value0 := argvalue0
    fmt.Print(client.LatestState(value0))
    fmt.Print("\n")
    break
  case "getNewOrders":
    if flag.NArg() - 1 != 0 {
      fmt.Fprintln(os.Stderr, "GetNewOrders requires 0 args")
      flag.Usage()
    }
    fmt.Print(client.GetNewOrders())
    fmt.Print("\n")
    break
  case "":
    Usage()
    break
  default:
    fmt.Fprintln(os.Stderr, "Invalid function ", cmd)
  }
}
