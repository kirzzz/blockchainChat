package network

import (
	network "blockchainStorage/common"
	"encoding/json"
	"fmt"
	"net"
)

type Message struct {
	Command string `json:"command"`
	Data    []byte `json:"data"`
}

type Network struct {
	NodeList []network.Node
}

func (n *Network) Broadcast(command string, data []byte) error {
	msg := &Message{command, data}
	payload, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	for _, node := range n.NodeList {
		conn, err := net.Dial("tcp", node.Address)
		if err != nil {
			continue
		}

		_, err = conn.Write(payload)
		if err != nil {
			continue
		}

		err = conn.Close()
		if err != nil {
			return fmt.Errorf("failed to close connection: %w", err)
		}
	}

	return nil
}

func (n *Network) StartServer(port int, handler func(msg *Message, conn net.Conn)) error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("failed to start listener: %w", err)
	}
	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			_ = fmt.Errorf("failed to close connection: %w", err)
		}
	}(listener)

	for {
		conn, err := listener.Accept()
		if err != nil {
			// Handle error
			continue
		}

		go func() {
			defer conn.Close()

			buf := make([]byte, 1024)
			n, err := conn.Read(buf)
			if err != nil {
				// Handle error
				return
			}

			var msg Message
			err = json.Unmarshal(buf[:n], &msg)
			if err != nil {
				// Handle error
				return
			}

			handler(&msg, conn)
		}()
	}
}
