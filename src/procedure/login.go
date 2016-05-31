package procedure

import (
    "fmt"
    process_data "../process_data"
)

/*Login :로그인*/
func Login(data *process_data.ProcessData) {
    fmt.Println("Login")
    fmt.Println(data.Param)
}
