package common

import "encoding/json"

type Node struct {
	Address string
}

func (pc *Node) UnmarshalJSON(p []byte) error {
	if string(p) == `""` {
		// empty string, do nothing
		return nil
	}
	// Prevent recursion to this method by declaring a new
	// type with same underlying type as PrimaryContact and
	// no methods.
	type x Node
	return json.Unmarshal(p, &(*x)(pc).Address)
}
