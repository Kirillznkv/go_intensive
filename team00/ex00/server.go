package main

import (
	"fmt"
	"github.com/google/uuid"
	"gonum.org/v1/gonum/stat/distuv"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	pb "main/api/nlo"
	"math/rand"
	"net"
	"time"
)

type logData struct {
	mean, std, freq float64
}

type GRPCServer struct {
}

func (t *GRPCServer) Do(r *pb.Request, srv pb.Nlo_DoServer) error {
	rand.Seed(time.Now().Unix())
	var ld logData
	ld.mean = rand.Float64()*21 - 10
	ld.std = rand.Float64()*1.2 + 0.3
	var res pb.Response
	res.SessionId = uuid.New().String()
	dist := distuv.Normal{
		Mu:    ld.mean,
		Sigma: ld.std,
	}
	for {
		res.Frequency = dist.Rand()
		res.Timestamp = timestamppb.Now()
		ld.freq = res.Frequency
		if err := srv.Send(&res); err != nil {
			return nil
		}
		t.myLog(ld)
	}
	return nil
}

func (t *GRPCServer) myLog(ld logData) {
	fmt.Printf("m: %f s: %f\t\t->\t%f\n", ld.mean, ld.std, ld.freq)
}

func main() {
	grpcServer := grpc.NewServer()
	srv := &GRPCServer{}
	pb.RegisterNloServer(grpcServer, srv)
	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	if err = grpcServer.Serve(l); err != nil {
		log.Fatal(err)
	}
}
