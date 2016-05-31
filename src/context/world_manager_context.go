package context

import (
	"net"
	"sync"

	process_data "../process_data"
)

type WorldManagerContext struct {
	world_list []WorldContext
}

func (self *WorldManagerContext) Init(world_count int, wg *sync.WaitGroup) {
	//var wg sync.WaitGroup
	self.world_list = make([]WorldContext, world_count)
	for i := 0; i < world_count; i++ {
		wg.Add(1)
		self.world_list[i].Init(i + 1)
		go self.world_list[i].Svc(wg)
	}
}

func (self *WorldManagerContext) UserAdd(conn *net.Conn) {
	index := 0
	for i := 1; i < len(self.world_list); i++ {
		if self.world_list[i].UserCount() < self.world_list[index].UserCount() {
			index = i
		}
	}
	data := process_data.CreateProcessData(process_data.P_Connect, conn, new(UserContext))
	self.world_list[index].Channel <- data
}
