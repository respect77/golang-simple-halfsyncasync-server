package context

import (
	"bytes"
	"fmt"
	"net"
	//"unsafe"

	packet_data "../packet_data"
	process_data "../process_data"
	proto "github.com/golang/protobuf/proto"
)

type UserContext struct {
	user_index             int
	conn                   *net.Conn
	send_channel           chan []byte
	linked_process_channel chan *process_data.ProcessData
}

func (self *UserContext) RecvExec() {
	var receive_buffer bytes.Buffer
	d := make([]byte, 4096) // 4096 크기의 바이트 슬라이스 생성

	header_exists := false
	
	header := &packet_data.Header{}
	
	temp_header, temp_err := proto.Marshal(&packet_data.Header{
		BodySize:   proto.Int32(1),
		PacketType: packet_data.PacketType_MoveType.Enum(),})
	if temp_err != nil {
		fmt.Println("temp_err : ", temp_err)
	}

	header_len := int(len(temp_header))

	for {
		n, recv_err := (*self.conn).Read(d) // 클라이언트에서 받은 데이터를 읽음
		if recv_err != nil {
			self.linked_process_channel <- process_data.CreateProcessData(process_data.P_Close, process_data.CreateCloseParam(self.user_index), self)
			return
		}
		//fmt.Println(n)

		_, w_err := receive_buffer.Write(d[:n])

		if w_err != nil {
			self.linked_process_channel <- process_data.CreateProcessData(process_data.P_Close, process_data.CreateCloseParam(self.user_index), self)
			return
		}

		if header_exists == false {
			//fmt.Println("receive_buffer.Len() : ", receive_buffer.Len())
			//fmt.Println("header_len : ", header_len)
			if receive_buffer.Len() < header_len {
				//헤더크기만큼 안들어왔으면 계속받음
				continue
			}
			header_data := receive_buffer.Next(header_len)
			header_exists = true

			//fmt.Println("header start")

			unmarshal_err := proto.Unmarshal(header_data, header)
			if unmarshal_err != nil {
				fmt.Println("unmarshaling error: ", unmarshal_err)
			}
			//fmt.Println(header.GetBodySize())
		}

		body_len := int(header.GetBodySize())
		if receive_buffer.Len() < body_len {
			//바디크기만큼 안들어왔으면 계속받음
			fmt.Println("Not Body")
			continue
		}

		body_data := receive_buffer.Next(body_len)
		
		self.linked_process_channel <- process_data.CreateProcessData(header.GetPacketType(), body_data, self)
		
		header_exists = false
	}
}

func (self *UserContext) SendExec() {
	for {
		p := <-(self.send_channel)
		_, err := (*self.conn).Write(p) // 클라이언트로 데이터를 보냄
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Sended")
	}
}

func (self *UserContext) Init(i int, conn *net.Conn, linked_world_ctx_channel chan *process_data.ProcessData) {
	self.send_channel = make(chan []byte, 1000)
	self.user_index = i
	self.conn = conn
	//this.Conn.SetNoDelay(false)
	self.linked_process_channel = linked_world_ctx_channel
	go self.RecvExec()
	go self.SendExec()
}

func (self *UserContext) Send(data []byte) {
	self.send_channel <- data
}