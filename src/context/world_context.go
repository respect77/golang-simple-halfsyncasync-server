package context

import (
	"fmt"
	"net"
	"sync"

	process_data "../process_data"
	packet_data "../packet_data"
	procedure "../procedure"
)

/*WorldContext : world context*/
type WorldContext struct {
	world_index     int
	user_map        map[int]*UserContext
	last_user_index int
	Channel         chan *process_data.ProcessData
	process_map map[packet_data.PacketType]func(*process_data.ProcessData)
}

func (self *WorldContext) Init(world_index int) {
	self.last_user_index = 1
	self.world_index = world_index
	self.process_map = make(map[packet_data.PacketType]func(*process_data.ProcessData))
	self.user_map = make(map[int]*UserContext)
	self.Channel = make(chan *process_data.ProcessData, 1000)
	
	self.process_map[packet_data.PacketType_LoginType] = procedure.Login
	self.process_map[packet_data.PacketType_MoveType] = procedure.Move
	self.process_map[process_data.P_Connect] = self.UserAdd
	self.process_map[process_data.P_Close] = self.UserClose
	//self.process_map[process_data.P_Connect] = procedure.Login
}

func (self *WorldContext) Svc(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		fmt.Println("channel wait : ", self.world_index)
		data := <-(self.Channel)
		self.process_map[data.Process_id](data)
	}
	fmt.Println("world done")
}

func (self *WorldContext) UserCount() int {
	return int(len(self.user_map))
}

func (self *WorldContext) UserAdd(data *process_data.ProcessData) {
	fmt.Printf("connected : %d", self.world_index)
	conn := data.Param.(*net.Conn)
	
	user_ctx := new(UserContext)
	user_ctx.Init(self.last_user_index, conn, self.Channel)
	self.user_map[self.last_user_index] = user_ctx

	self.last_user_index++
	
	user_ctx.Send([]byte("login"))
}

func (self *WorldContext) UserClose(data *process_data.ProcessData) {
	fmt.Println("closed")
}
