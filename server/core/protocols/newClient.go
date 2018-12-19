package protocols

import "fmt"

func GetClientIdMessage(id int) []byte {
	str := fmt.Sprintf("{\"clientId\":%d}", id)

	return []byte(str)
}
