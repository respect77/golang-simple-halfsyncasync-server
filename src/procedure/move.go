package procedure

import (
    "fmt"
	proto "github.com/golang/protobuf/proto"
    process_data "../process_data"
	packet_data "../packet_data"
)

/*Move :Move*/
func Move(data *process_data.ProcessData) {
    fmt.Println("Move")
	body := &packet_data.MoveRequest{}

	unmarshal_err := proto.Unmarshal(data.Param.([]byte), body)
	if unmarshal_err != nil {
		fmt.Println("unmarshaling error: ", unmarshal_err)
	}
			
	fmt.Println(body.GetPos().Enum())
	
	data.User_context.Send([]byte("data"))
}
