package grpc

import (
	"fmt"
	"net"

	"SDK/lib/processcontrol"
	grpcexample "SDK/proto/grpcexample"
	grpcgcpsdk "SDK/proto/grpcgcpsdk"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var log = logrus.WithFields(logrus.Fields{"package": "grpc-SdkServer"})

func ExampleStart(ip string, port string, tls bool, certFile string, keyFile string) {
	lis, err := net.Listen("tcp", fmt.Sprintf("%v:%v", ip, port))
	if err != nil {
		log.Fatal("ExampleStart failed to listen: %v", err)
		return
	}
	var opts []grpc.ServerOption
	if tls {
		if certFile != "" && keyFile != "" {
			creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
			if err != nil {
				log.Fatal("ExampleStart Failed to generate credentials %v", err)
			} else {
				opts = []grpc.ServerOption{grpc.Creds(creds)}
			}
		}
	}
	grpcServer := grpc.NewServer(opts...)
	grpcexample.RegisterExampleServiceServer(grpcServer, newServer())
	log.Info("ExampleServer grpcServer start")
	grpcServer.Serve(lis)
}

func SdkStart(ip string, port string, tls bool, certFile string, keyFile string) {
	lis, err := net.Listen("tcp", fmt.Sprintf("%v:%v", ip, port))
	if err != nil {
		log.Fatal("SdkStart failed to listen: %v", err)
		processcontrol.ProcessExit <- 1
		return
	}
	var opts []grpc.ServerOption
	if tls {
		if certFile != "" && keyFile != "" {
			creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
			if err != nil {
				log.Fatal("SdkStart Failed to generate credentials %v", err)
			} else {
				opts = []grpc.ServerOption{grpc.Creds(creds)}
			}
		}
	}
	grpcServer := grpc.NewServer(opts...)
	grpcgcpsdk.RegisterSDKServiceServer(grpcServer, newSdkServer())
	log.Info("SdkServer grpcServer start")
	grpcServer.Serve(lis)
}
