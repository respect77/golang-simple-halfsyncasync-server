package process_data

import (
	"net"
	packet_data "../packet_data"
)
const (
	/*1~9999 packet_data*/
	P_Connect packet_data.PacketType = 10000
	P_Close  packet_data.PacketType = 10001
)

type ProcessData struct {
	Process_id packet_data.PacketType
	Param interface{}
	User_context interface{
		Send([]byte)
	}
}


type ConnectParam struct {
	Conn *net.Conn
}

type CloseParam struct {
	User_index int
}

func CreateProcessData(packet_type packet_data.PacketType, param interface{}, userCtx interface{Send([]byte)}) *ProcessData {
	return &ProcessData{Process_id : packet_type, Param : param, User_context : userCtx}
}

func CreateConnectParam(conn *net.Conn) *ConnectParam {
	return &ConnectParam{Conn : conn}
}

func CreateCloseParam(user_index int) *CloseParam {
	return &CloseParam{User_index : user_index}
}
