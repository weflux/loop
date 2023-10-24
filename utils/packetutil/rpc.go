package packetutil

import (
	"github.com/google/uuid"
	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/packets"
	proxypb "github.com/weflux/loopify/protocol/proxy"
	shared "github.com/weflux/loopify/protocol/shared"
	"github.com/weflux/loopify/utils/clientutil"
	"github.com/weflux/loopify/utils/topicutil"
)

func ToConnectRequest(cl *mqtt.Client, pk *packets.Packet) *proxypb.ConnectRequest {
	req := &proxypb.ConnectRequest{
		Id: uuid.NewString(),
		Metadata: &shared.Metadata{
			Topic:   "",
			Session: cl.ID,
			User:    string(cl.Properties.Username),
		},

		Payload: &proxypb.ConnectRequest_Payload{
			Username: string(pk.Connect.Username),
			Password: string(pk.Connect.Password),
		},
	}

	return req
}

func FromConnectRequest(req *proxypb.ConnectRequest) (packets.Packet, error) {

	return packets.Packet{}, nil
}

func ToConnectReply(pk *packets.Packet) *proxypb.ConnectReply {
	return nil
}

func FromConnectReply(rep *proxypb.ConnectReply) (packets.Packet, error) {
	return packets.Packet{}, nil
}

func ToRPCRequest(pk *packets.Packet) *proxypb.RPCRequest {

	return nil
}

func ToRPCReply(pk *packets.Packet) *proxypb.RPCReply {
	return nil
}

func FromRPCRequest(req *proxypb.RPCRequest) (packets.Packet, error) {
	return packets.Packet{}, nil
}

func FromRPCReply(cl *mqtt.Client, rep *proxypb.RPCReply) (packets.Packet, error) {
	ct := clientutil.GetContentType(cl)
	data, err := ct.Marshal(rep)
	if err != nil {
		return packets.Packet{}, nil
	}
	pk := packets.Packet{
		FixedHeader: packets.FixedHeader{Type: packets.Publish},
		TopicName:   topicutil.ReplyTopic(),
		Ignore:      true,
		Payload:     data,
	}
	return pk, nil
}
