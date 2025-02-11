package main

import (
	"fmt"
	"strconv"
	"hash/fnv"
	"github.com/asynkron/protoactor-go/actor"
)

type NodeActor struct {
	nodeID		uint64 //Chord identifier (a hash)
	address 	string
	port		int
	data 		map[string]string //piece of the overall hashtable this node holds
	//actorPID	*actor.PID //Proto.Actor identifier
	//successor 	*actor.PID //actor that would be queried by this node
	//predecessor *actor.PID //actor that would query this node
	//fingerTable []*actor.PID //list of actors at the predetermined locations on the ring
}

func (state *NodeActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
		case *Message:
			fmt.Printf("Got a message: %v", msg.Text)
	}
}

func makeNode(address string, port int) *NodeActor {

	//get the hash of the address to serve as node ID
	s := address + strconv.Itoa(port)
	hash := fnv.New64a()
	hash.Write([]byte(s)) 

	return &NodeActor {
		nodeID: hash.Sum64(),
		address: address,
		port: port,
		data: make(map[string]string),
	}
}
