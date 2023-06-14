package main

import (
	"log"

	rpc "github.com/TikTokTechImmersion/assignment_demo_2023/rpc-server/kitex_gen/rpc/imservice"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	etcd "github.com/kitex-contrib/registry-etcd"
)


type MyServiceImpl struct{}

func (s *MyServiceImpl) Ping(ctx context.Context, message string) (string, error) {
    // Perform any necessary processing with the message
    // In this case, simply reply with "pong"
    response := "pong"
    return response, nil
}

// Define the RPC service interface
type MyService interface {
    // Method for ping-pong interaction
    Ping(ctx context.Context, message string) (string, error)
}
// import context

func main() {
	r, err := etcd.NewEtcdRegistry([]string{"etcd:2379"}) // r should not be reused.
	if err != nil {
		log.Fatal(err)
	}

	svr := rpc.NewServer(new(IMServiceImpl), server.WithRegistry(r), server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
		ServiceName: "demo.rpc.server",
	}))

	err = svr.Run()
	if err != nil {
		log.Println(err.Error())
	}
}
