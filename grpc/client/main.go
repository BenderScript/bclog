package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/repenno/bclog/grpc/bclogpb"
	"google.golang.org/grpc"
	"log"
	"time"
)

var (
	serverAddr = flag.String("server_addr", "127.0.0.1:10000", "The server address in the format of host:port")
)

func main() {
	flag.Parse()
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := bclogpb.NewBlockChainLogClient(conn)

	id := createEntry(client, ctx)

	getEntry(client, ctx, id)

	fmt.Println("Client finished, exiting...")

}

func createEntry(client bclogpb.BlockChainLogClient, ctx context.Context) string {

	baseEntry := bclogpb.BaseEntry{
		Text: "Store this log into blockchain",
	}
	createEntryReq := bclogpb.CreateEntryReq{
		Entry: &baseEntry,
	}
	log.Println("Sending request")
	result, err := client.CreateEntry(ctx, &createEntryReq)
	if err != nil {
		log.Fatalf("Create Entry failed %v", err)
		return ""
	}
	log.Println("Id is:", result.Entryid.Id)
	return result.Entryid.Id

}

func getEntry(client bclogpb.BlockChainLogClient, ctx context.Context, id string) int {

	getEntryReq := bclogpb.GetEntryReq{
		Id: id,
	}
	log.Println("Sending get request for id: ", id)
	result, err := client.GetEntry(ctx, &getEntryReq)
	if err != nil {
		log.Fatalf("Get Entry failed %v", err)
		return -1
	}
	log.Println("Received text is:", result.Entry.Text)
	return 0

}
