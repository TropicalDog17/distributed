package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"golang.org/x/exp/slices"
)

type Message struct {
	Id   int    `json:"id"`
	Src  string `json:"src"`
	Dest string `json:"dest"`
	Body any    `json:"body"`
}
type Response struct {
	Src  string `json:"src"`
	Dest string `json:"dest"`
	Body any    `json:"body"`
}

var node_id string
var next_msg_id int
var messages = make([]int, 0)
var adjacents = make([]string, 0)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		msg := Message{}
		err := json.Unmarshal([]byte(line), &msg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
			continue
		}
		fmt.Fprintf(os.Stderr, "Received \"%v\\n\"\n", msg)

		body := msg.Body.(map[string]any)

		switch body["type"].(string) {
		case "init":
			node_id = body["node_id"].(string)
			fmt.Fprintf(os.Stderr, "Initialized node %s\n", node_id)
			body["type"] = "init_ok"
			reply(msg, body)
		case "broadcast":
			fmt.Fprintf(os.Stderr, "Receive broadcast message %v\n", msg)
			reply_broadcast(msg, body)
		case "read":
			fmt.Fprintf(os.Stderr, "Receive read message %v\n", msg)
			body["type"] = "read_ok"
			reply_read(msg, body)
		case "topology":
			fmt.Fprintf(os.Stderr, "Receive topology message %v\n", msg)
			body["type"] = "topology_ok"
			reply_topology(msg, body)
		}

	}
}
func reply(request Message, body map[string]any) {
	next_msg_id++
	id := next_msg_id

	res := Response{}
	body["in_reply_to"] = body["msg_id"]
	body["msg_id"] = id
	res.Src = node_id
	res.Dest = request.Src
	res.Body = body

	jsonBytes, err := json.Marshal(res)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	fmt.Fprintln(os.Stdout, string(jsonBytes))
}
func broadcast(body map[string]any, dest string) {
	body["type"] = "broadcast"

	res := Response{}
	res.Src = node_id
	res.Dest = dest
	res.Body = body
	jsonBytes, err := json.Marshal(res)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	fmt.Fprintln(os.Stdout, string(jsonBytes))
}
func reply_broadcast(msg Message, body map[string]any) {
	body["type"] = "broadcast_ok"
	float64_val := body["message"].(float64)
	broadcast_value := int(float64_val)
	fmt.Fprintf(os.Stderr, "Receive broadcast value %v \n", broadcast_value)
	if !slices.Contains(messages, broadcast_value) {
		messages = append(messages, broadcast_value)
	}
	temp := body["message"]
	delete(body, "message")
	reply(msg, body)
	body["message"] = temp
	for _, element := range adjacents {
		if element != msg.Src {
			broadcast(body, element)
		}
	}

}
func reply_read(msg Message, body map[string]any) {
	body["messages"] = messages
	reply(msg, body)
}
func reply_topology(msg Message, body map[string]any) {
	a := body["topology"]
	b := a.(map[string]any)[node_id]
	for _, e := range b.([]any) {
		f, ok := e.(string)
		if ok {
			adjacents = append(adjacents, f)
		} else {
			return
		}
	}
	delete(body, "topology")
	reply(msg, body)
}
