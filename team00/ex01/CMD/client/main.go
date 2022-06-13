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

func forecastLog(i int, data *pb.Response, m, s float64) {
	fmt.Printf("%d: (mean:%f std:%f) ", i, m, s)
	fmt.Printf("%s %f %s\n", data.SessionId, data.Frequency, data.Timestamp.AsTime())
}

func anomalyLog(data *pb.Response) {
	fmt.Printf("Anomaly: ")
	fmt.Printf("%s %f %s\n", data.SessionId, data.Frequency, data.Timestamp.AsTime())
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

func execDo(ch chan *pb.Response, conn *grpc.ClientConn) {
	stream := getDoStream(conn)
	for {
		data, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		ch <- data
	}
	close(ch)
}

func findForecast(conn *grpc.ClientConn) {
	forecastArr := make([]float64, 1000)
	var mean, std float64
	ch := make(chan *pb.Response)
	i := 0
	go execDo(ch, conn)
	for {
		data := <-ch
		if i < 1000 {
			forecastArr[i] = data.Frequency
			mean, std = stat.MeanStdDev(forecastArr[:i+1], nil)
			forecastLog(i+1, data, mean, std)
			if i == 1000 {
				forecastArr = nil
			}
		} else if data.Frequency < mean-(std*(*kFlag)) || data.Frequency > mean+(std*(*kFlag)) {
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
