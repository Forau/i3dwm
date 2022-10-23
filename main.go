package main

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"io"
	"log"
	"net"
	"os"
)

const (
	I3Magic = "i3-ipc"

	I3Command = iota
	I3GetWorkspaces
	I3Subscribe
	I3GetOutputs
	I3GetTree
	I3GetMarks
	I3GetBarConfig
	I3GetVersion
)

type I3Node struct {
	ID            int64    `json:"id"`
	Name          string   `json:"name"`
	Type          string   `json:"type"`
	Focused       bool     `json:"focused"`
	Nodes         []I3Node `json:"nodes"`
	FloatingNodes []I3Node `json:"floating_nodes"`
	Layout        string   `json:"layout"`
	Output        string   `json:"output"`
}

func main() {
	log.Printf("warning, this is still in initial state, and does not yet work")
	sockAddr := os.Getenv("I3SOCK")
	if sockAddr == "" {
		log.Panicf("no env variable I3SOCK for i3 socket")
	}

	sock, err := net.Dial("unix", sockAddr)
	if err != nil {
		log.Panic(err)
	}
	defer sock.Close()

	rsp, err := send_msg(context.Background(), sock, 4)
	if err != nil {
		log.Panicf("err: %+v", err)
	}
	var nod I3Node
	err = json.Unmarshal([]byte(rsp), &nod)
	if err != nil {
		log.Panicf("err: %+v", err)
	}
	//log.Printf("nod: %+v", nod)
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(nod)
}

func send_msg(ctx context.Context, sock net.Conn, typ uint32, data ...byte) (string, error) {
	err := Message{Type: 4}.Write(sock)
	if err != nil {
		return "", err
	}

	var msg Message
	err = msg.Read(sock)

	log.Printf("read: %s %v", msg.Payload, err)
	return string(msg.Payload), nil
}

type Message struct {
	Type    uint32
	Payload []byte
}

func (m Message) Write(w io.Writer) error {
	_, err := w.Write([]byte(I3Magic))
	if err != nil {
		return err
	}
	if err = binary.Write(w, binary.LittleEndian, uint32(len(m.Payload))); err != nil {
		return err
	}
	return binary.Write(w, binary.LittleEndian, m.Type)
}

func (m *Message) Read(r io.Reader) error {
	header := make([]byte, 6)
	if _, err := r.Read(header); err != nil {
		return err
	}
	var l uint32
	if err := binary.Read(r, binary.LittleEndian, &l); err != nil {
		return err
	}
	if err := binary.Read(r, binary.LittleEndian, &m.Type); err != nil {
		return err
	}
	m.Payload = make([]byte, l)
	_, err := r.Read(m.Payload)
	return err
}
