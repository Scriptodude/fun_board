package protocols

import "fmt"
import i "server/interfaces"

var (
	id int
)

func NewClient(client *i.GameClient) {
	client.Id = id
	id += 1

	writeClientIdMessage(client)
}

func ExistingClient(client *i.GameClient) {
	writeClientIdMessage(client)
}

func writeClientIdMessage(client *i.GameClient) {
	str := fmt.Sprintf("{\"clientId\":%d}", client.Id)

	client.Writer.Header().Add("Content-Type", "application/json")
	client.Messages <- str
}
