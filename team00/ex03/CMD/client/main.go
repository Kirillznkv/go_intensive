package main

import (
	"context"
	pb "ex01/API/nlo"
	"flag"
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"gonum.org/v1/gonum/stat"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
	"time"
)

type psqlLog struct {
	db *pg.DB
}

func (l *psqlLog) forecastLog(i int, data *pb.Response, m, s float64) {
	fmt.Printf("%d: (mean:%f std:%f) ", i, m, s)
	fmt.Printf("%s %f %s\n", data.SessionId, data.Frequency, data.Timestamp.AsTime())
}

func (l *psqlLog) anomalyLog(data *pb.Response) {
	a := &anomalyData{
		SessionId: data.SessionId,
		Frequency: data.Frequency,
		Timestamp: data.Timestamp.AsTime(),
	}
	_, err := l.db.Model(a).Insert()
	if err != nil {
		log.Fatal(err)
	}
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

func findAnomaly(conn *grpc.ClientConn, db *pg.DB) {
	l := psqlLog{db: db}
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
			l.forecastLog(i+1, data, mean, std)
			if i == 1000 {
				forecastArr = nil
			}
		} else if data.Frequency < mean-(std*(*kFlag)) || data.Frequency > mean+(std*(*kFlag)) {
			l.anomalyLog(data)
		}
		i++
	}
}

var kFlag *float64

func init() {
	kFlag = flag.Float64("k", 1, "STD anomaly coefficient")
	flag.Parse()
}

type anomalyData struct {
	tableName struct{} `pg:"anomaly"`
	SessionId string   `pg:"type:text"`
	Frequency float64
	Timestamp time.Time `pg:"type:timestamp"`
}

func createAnomalyTable(db *pg.DB) error {
	model := (*anomalyData)(nil)
	err := db.Model(model).CreateTable(&orm.CreateTableOptions{})
	if err != nil && err.Error() == "42P07" {
		return err
	}
	return nil
}

func connectDB() *pg.DB {
	db := pg.Connect(&pg.Options{
		Addr:     ":5432",
		User:     "kshanti",
		Database: "postgres",
	})
	return db
}

func main() {
	db := connectDB()
	defer db.Close()
	err := createAnomalyTable(db)
	if err != nil {
		log.Fatal(err)
	}
	conn := connectToServ("127.0.0.1:8080")
	defer conn.Close()
	findAnomaly(conn, db)
}
