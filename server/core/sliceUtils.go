package core

import intf "server/interfaces"

func InsertAt(s []*intf.GameClient, i int, val *intf.GameClient) []*intf.GameClient {
	s = append(s, nil)
	copy(s[i+1:], s[i:])
	s[i] = val
	return s
}

func DeleteAt(s []*intf.GameClient, i int) []*intf.GameClient {
	copy(s[i:], s[i+1:])
	s[len(s)-1] = nil // or the zero value of T
	s = s[:len(s)-1]
	return s
}
