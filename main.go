package main

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"flag"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"text/template"
)

const (
	I3Command = iota
	I3GetWorkspaces
	I3Subscribe
	I3GetOutputs
	I3GetTree
	I3GetMarks
	I3GetBarConfig
	I3GetVersion

	I3Magic = "i3-ipc"
)

var (
	perMonitorNumWS = flag.Int("numws", 10, "number of workspaces per monitor")
	startNumber     = flag.Int("numstart", 1, "first workspace number (0 or 1)")
	i3sock          = flag.String("i3sock", "I3SOCK", "envvar for the unix socket")
)

func main() {
	ctx := context.Background()
	flag.Parse()
	sockAddr := os.Getenv(*i3sock)
	if sockAddr == "" {
		log.Panicf("no env variable %q for i3 socket", *i3sock)
	}

	sock, err := net.Dial("unix", sockAddr)
	if err != nil {
		log.Panicf("unable to open socket: %v", err)
	}
	defer sock.Close()

	nodes, err := sendMsg(ctx, sock, 4)
	if err != nil {
		log.Panicf("err: %+v", err)
	}
	if len(nodes) != 1 {
		log.Panicf("expected 1 result, but got %d", len(nodes))
	}

	ws, curr, ok := nodes[0].Current()
	if !ok {
		log.Panicf("no focused workspace?")
	}

	if flag.NArg() == 0 {
		// No args, just debug
		log.Printf("no commands run, so just show some debug messages\nCurrent WS: %v\nCurrent con: %v", ws, curr)
		return
	}

	displayAdder, err := ws.DisplayAdder(*perMonitorNumWS, *startNumber)
	if err != nil {
		log.Panicf("unable to find current display")
	}

	templ := template.New("i3cmd")
	templ.Funcs(map[string]any{
		"ws": func(i int) int { return i + displayAdder },
	})
	tmpl, err := templ.Parse(strings.Join(flag.Args(), " "))
	if err != nil {
		log.Panicf("unable to parse template: %v", err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, map[string]any{})
	if err != nil {
		log.Panicf("unable to run template: %v", err)
	}
	_, err = sendMsg(ctx, sock, I3Command, buf.Bytes()...)
	if err != nil {
		log.Panicf("error sending i3 command: %v", err)
	}
}

func sendMsg(ctx context.Context, sock net.Conn, typ uint32, data ...byte) ([]I3Node, error) {
	err := Message{Type: typ, Payload: data}.Write(sock)
	if err != nil {
		return nil, err
	}

	var msg Message
	err = msg.Read(sock)
	if err != nil {
		return nil, err
	}
	if msg.Payload[0] == '{' {
		var nod I3Node
		err = json.Unmarshal(msg.Payload, &nod)
		return []I3Node{nod}, err
	}
	var nodes []I3Node
	return nodes, json.Unmarshal(msg.Payload, &nodes)
}

type Message struct {
	Type    uint32
	Payload []byte
}

func (m Message) Write(w io.Writer) error {
	var buf bytes.Buffer
	_, err := buf.WriteString(I3Magic)
	if err != nil {
		return err
	}
	if err = binary.Write(&buf, binary.LittleEndian, uint32(len(m.Payload))); err != nil {
		return err
	}
	if err = binary.Write(&buf, binary.LittleEndian, m.Type); err != nil {
		return err
	}
	if _, err = buf.Write(m.Payload); err != nil {
		return err
	}

	log.Printf("sending %q", buf.String())
	_, err = io.Copy(w, &buf)
	return err
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

// I3Node is a subset of what the nodes fields, but it all we need for this
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

func (n I3Node) DisplayAdder(perDisp, start int) (int, error) {
	myNum, err := strconv.Atoi(strings.Split(n.Name, " ")[0])
	if err != nil {
		return 0, err
	}
	disp := (myNum - start) / perDisp
	return disp * perDisp, nil
}

func (n I3Node) Current() (ws, con *I3Node, found bool) {
	weAreWS := n.Type == "workspace"
	if n.Focused {
		if weAreWS {
			return &n, nil, true
		}
		return nil, &n, true
	}

	nodes := append([]I3Node{}, n.Nodes...)
	nodes = append(nodes, n.FloatingNodes...)
	for _, nod := range nodes {
		ws, con, found = nod.Current()
		if found {
			if ws == nil {
				// Check if it is us..
				if weAreWS {
					ws = &n
				}
			}
			return
		}
	}
	return
}
