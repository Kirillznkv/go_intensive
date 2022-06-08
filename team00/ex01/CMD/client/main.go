package main

import (
	"context"
	pb "ex01/API/nlo"
	"flag"
	"fmt"
	"gonum.org/v1/gonum/stat"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
)

func forecastLog(i int, f, m, s float64) {
	fmt.Println(i, ":", "Frequency:", f, "Mean:", m, "STD:", s)
}

func anomalyLog(f float64) {
	fmt.Println("Anomaly:", f)
}

func connectToServ(addr string) *grpc.ClientConn {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	opts = append(opts, grpc.WithBlock())
	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	return conn
}

func getDoStream(conn *grpc.ClientConn) pb.Nlo_DoClient {
	c := pb.NewNloClient(conn)
	stream, err := c.Do(context.Background(), &pb.Request{})
	if err != nil {
		log.Fatalf("%v.Do(_) = _, %v", c, err)
	}
	return stream
}

func execDo(ch chan float64, conn *grpc.ClientConn) {
	stream := getDoStream(conn)
	for {
		data, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		ch <- data.Frequency
	}
	close(ch)
}

func findForecast(conn *grpc.ClientConn) {
	forecastArr := make([]float64, 1000)
	var mean, std float64
	ch := make(chan float64)
	i := 0
	go execDo(ch, conn)
	for {
		data := <-ch
		if i < 1000 {
			forecastArr[i] = data
			mean, std = stat.MeanStdDev(forecastArr[:i+1], nil)
			forecastLog(i+1, data, mean, std)
			if i == 1000 {
				forecastArr = nil
			}
		} else if data < mean-(std**kFlag) || data > mean+(std**kFlag) {
			anomalyLog(data)
		}
		i++
	}
}

var kFlag *float64

func init() {
	kFlag = flag.Float64("k", 1, "STD anomaly coefficient")
	flag.Parse()
}

func main() {
	conn := connectToServ("127.0.0.1:8080")
	defer conn.Close()
	findForecast(conn)
}
