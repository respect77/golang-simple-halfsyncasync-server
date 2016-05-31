package main

import (
	"fmt"
	"net"
	"sync"
    "runtime"
	//"unsafe"

	context "./context"
	
	//proto "github.com/golang/protobuf/proto"
	packet_data "./packet_data"
)

func main() {
	fmt.Println(packet_data.Position_Down)
	
	s := packet_data.Position_Down

	fmt.Println(s.Enum())
    runtime.GOMAXPROCS(runtime.NumCPU())
	var wg sync.WaitGroup

	world_count := 2
	world_manager := new(context.WorldManagerContext)
	world_manager.Init(world_count, &wg)

	listener, listen_error := net.Listen("tcp", ":8000")
	if listen_error != nil {
		fmt.Println(listen_error)
		return
	}

	defer listener.Close()

	for {
		conn, accept_error := listener.Accept()
		fmt.Println("Connect")
		if accept_error != nil {
			fmt.Println(accept_error)
			continue
		}
		world_manager.UserAdd(&conn)
	}
	wg.Wait()
	fmt.Println("done..")
	return
}
