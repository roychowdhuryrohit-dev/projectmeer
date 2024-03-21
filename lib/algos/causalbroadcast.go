package algos

import (
	"bytes"
	"encoding/gob"
	"log"
	"net/http"
	"slices"
	"strings"

	config "github.com/roychowdhuryrohit-dev/projectmeer/lib"
)

type Message struct {
	Sender  int
	Deps    []int
	Msg     []byte
	MsgType string
}

type CPrimitive[T any] interface {
	SendPrimitive([]byte, string) error
	ReceivePrimitive([]byte, string) error
	Send([]byte) error
	Receive([]byte) error
}

type CausalBroadcast[T any] struct {
	sendSeq   int
	delivered []int
	buffer    []Message
	cp        *CPrimitive[T]
	nodeList  []string
	curNode   int
}

func (cb *CausalBroadcast[T]) SendPrimitive(msg []byte, msgType string) error {
	deps := cb.delivered
	deps[cb.curNode] = cb.sendSeq
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(Message{
		Sender:  cb.curNode,
		Deps:    deps,
		Msg:     msg,
		MsgType: msgType,
	})
	if err != nil {
		return err
	}

	return cb.Send(buf.Bytes())
}

func (cb *CausalBroadcast[T]) Send(encMsg []byte) error {
	prefix := "http://"
	prefix_secure := "https://"
	path := "/p2p/receivePrimitive"
	for _, node := range config.NodeList {
		buf := bytes.NewBuffer(encMsg)
		if !strings.HasPrefix(node, prefix) && !strings.HasPrefix(node, prefix_secure) {
			node = prefix + node
		}
		_, err := http.Post(node+path, "application/octet-stream", buf)
		if err != nil {
			return err
		}

	}
	cb.sendSeq += 1
	return nil
}

func (cb *CausalBroadcast[T]) Receive(encMsg []byte) error {
	buf := bytes.NewBuffer(encMsg)
	dec := gob.NewDecoder(buf)
	msg := Message{}
	if err := dec.Decode(&msg); err != nil {
		log.Println(err.Error())
		return err
	}
	cb.buffer = append(cb.buffer, msg)
	for i, m := range cb.buffer {
		if cb.isCausal(m.Deps, cb.delivered) {
			err := (*(cb.cp)).ReceivePrimitive(m.Msg, m.MsgType)
			if err != nil {
				log.Println(err.Error())
			}
			cb.buffer = slices.Delete(cb.buffer, i, i+1)
			cb.delivered[m.Sender] += 1

		}
	}

	return nil
}

func (cb *CausalBroadcast[T]) isCausal(v1 []int, v2 []int) bool {
	l1 := len(v1)
	l2 := len(v2)
	if l1 != l2 {
		return false
	}

	for i := 0; i < l1; i++ {
		if v1[i] > v2[i] {
			return false
		}
	}
	return true
}
