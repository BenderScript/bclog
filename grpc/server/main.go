package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/golang/protobuf/ptypes"
	"github.com/google/uuid"
	"github.com/repenno/bclog/grpc/bclogpb"
	"google.golang.org/grpc"
	"log"
	"net"
)

var (
	port = flag.Int("port", 10000, "The server port")
)

type BcLogServer struct {
	m map[string]string
}

func newServer() *BcLogServer {
	s := &BcLogServer{
		m: make(map[string]string),
	}
	return s
}

func (bs *BcLogServer) CreateEntry(ctx context.Context, req *bclogpb.CreateEntryReq) (*bclogpb.CreateEntryResp, error) {
	fmt.Println("Received entry with text: ", req.Entry.Text)
	newId := uuid.New().String()
	entryid := bclogpb.EntryId{
		Id: newId,
	}
	resp := bclogpb.CreateEntryResp{
		Entryid:   &entryid,
		Timestamp: ptypes.TimestampNow(),
	}
	bs.m[newId] = req.Entry.Text
	return &resp, nil
}

func (bs *BcLogServer) GetEntry(ctx context.Context, req *bclogpb.GetEntryReq) (*bclogpb.GetEntryResp, error) {
	fmt.Println("Received get entry with id: ", req.Id)

	entryid := bclogpb.EntryId{
		Id: req.Id,
	}
	text := bs.m[req.Id]
	baseEntry := bclogpb.BaseEntry{
		Text:    text,
		Entryid: &entryid,
	}
	resp := bclogpb.GetEntryResp{
		Entry: &baseEntry,
	}
	return &resp, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	bclogpb.RegisterBlockChainLogServer(grpcServer, newServer())
	log.Println("Starting GRPC Server")
	grpcServer.Serve(lis)
}
