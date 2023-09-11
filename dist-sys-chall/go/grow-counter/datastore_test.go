package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandleReadData(t *testing.T) {
	ds := DataStore{}
	ds.kv = map[string]int{
		"a": 12,
		"b": 34,
		"c": 56,
		"d": 78,
		"e": 910,
	}
	// Id   int            `json:"id"`
	// Src  string         `json:"src"`
	// Dest string         `json:"dest"`
	// Body map[string]any `json:"body"`
	msg := Message{
		Id:   1,
		Src:  "n0",
		Dest: "n2",
		Body: map[string]any{
			"type":  "read",
			"value": "a",
		},
	}
	response := ds.handle(msg)
	assert.Equal(t, response.Body["type"], "read_ok", "Should be ok")
	assert.Equal(t, response.Body["value"], 12, "Should pass")

}
