package protocols

import "fmt"
import "net/http"

func GetClientIdMessage(w http.ResponseWriter, id int) {
	str := fmt.Sprintf("{\"clientId\":%d}", id)

	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(str))
}
