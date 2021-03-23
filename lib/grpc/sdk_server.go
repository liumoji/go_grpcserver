package grpc

import (
	"context"
	"fmt"
	"io"
	"sync"

	grpcgcpsdk "SDK/proto/grpcgcpsdk"

	config "SDK/config_managers"
	myweb "SDK/lib/myweb"

	"google.golang.org/protobuf/types/known/structpb"
)

type SdkServer struct {
	grpcgcpsdk.UnimplementedSDKServiceServer

	mu          sync.Mutex // protects routeNotes
	streamNotes map[int]grpcgcpsdk.SDKService_LongConnectionChannelServer
	uuid        int
}

func (s *SdkServer) RegisterWithThePlatform(ctx context.Context, req *grpcgcpsdk.GegRequest) (*grpcgcpsdk.GegResponse, error) {
	log.Info("received client GegRequest token: " + req.Token + " deviceId: " + req.DeviceId)
	var res *grpcgcpsdk.GegResponse
	res = &grpcgcpsdk.GegResponse{}
	//去平台注册
	httpUrl := config.GetValue("gcp", "url")
	httpUrl = httpUrl + "/v1/gcp"
	httpClient := myweb.NewHttpsClient(httpUrl, false, "")
	if httpClient == nil {
		res.Code = int32(grpcgcpsdk.ResponseCode_SSLFAIL)
		res.Message = "loading ssl config error"
	}

	return res, nil
}

func (s *SdkServer) LongConnectionChannel(stream grpcgcpsdk.SDKService_LongConnectionChannelServer) error {
	var thisId int
	s.mu.Lock()
	s.uuid++
	thisId = s.uuid
	s.streamNotes[thisId] = stream
	s.mu.Unlock()
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			log.Print("LongConnectionChannel io.EOF uuid: " + fmt.Sprint(thisId))
			delete(s.streamNotes, thisId)
			return nil
		}
		if err != nil {
			log.Printf("LongConnectionChannel err : %s , uuid: %d end", err, thisId)
			delete(s.streamNotes, thisId)
			return err
		}

		log.Info("Receive client message: cmd: " + fmt.Sprint(in.Cmd))
		if in.Data != nil && in.Data.Fields != nil {
			for v := range in.Data.Fields {
				tt := in.Data.Fields[v].GetKind()
				log.Info("Receive client message: data [" + v + "]: " + fmt.Sprint(tt))
				log.Info("Receive client message: data [" + v + "]: " + fmt.Sprint(in.Data.Fields[v].GetKind()))
			}
		}

		//sendMessage := &grpcgcpsdk.ChaResponse{Data: &structpb.Struct{Fields: make(map[string]*structpb.Value)}}
		sendMessage := &grpcgcpsdk.ChaResponse{}
		sendMessage.Data = &structpb.Struct{}
		sendMessage.Data.Fields = make(map[string]*structpb.Value)
		sendMessage.Cmd = int32(grpcgcpsdk.ResponseCmd_SECURITYPOLICY.Number())
		sendMessage.MessageJson = "{\"test_number\": 123, \"test_string\": \"qwerty\"}"
		feilds := make(map[string]*structpb.Value)
		feilds["sdk_test"] = &structpb.Value{Kind: &structpb.Value_StringValue{StringValue: "sdk server send test message!"}}
		sendMessage.Data.Fields = feilds
		// sendMessage := &grpcgcpsdk.ChaResponse{}
		// sendMessage.Cmd = in.Cmd
		// sendMessage.Data = in.Data
		if s.streamNotes[1] != nil {
			s_err := s.streamNotes[1].Send(sendMessage)
			if s_err != nil {
				log.Printf("streamNotes[1] Send end")
				delete(s.streamNotes, 1)
			}
		}
		serr := stream.Send(sendMessage)
		if serr != nil {
			log.Printf("Send end")
			delete(s.streamNotes, thisId)
			return serr
		}
	}
}

func newSdkServer() *SdkServer {
	s := &SdkServer{streamNotes: make(map[int]grpcgcpsdk.SDKService_LongConnectionChannelServer)}
	return s
}
