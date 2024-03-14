package algos

import (
	"bytes"
	"encoding/gob"
	"slices"
)

type Message struct {
	sender  int
	deps    []int
	msg     []byte
	msgType string
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
		sender:  cb.curNode,
		deps:    deps,
		msg:     msg,
		msgType: msgType,
	})
	if err != nil {
		return err
	}
	cb.Send(buf.Bytes())
	return nil
}

// TODO
func (cb *CausalBroadcast[T]) Send(encMsg []byte) error {
	return nil
}

func (cb *CausalBroadcast[T]) Receive(encMsg []byte) error {
	buf := bytes.NewBuffer(encMsg)
	dec := gob.NewDecoder(buf)
	msg := Message{}
	if err := dec.Decode(&msg); err != nil {
		return err
	}
	cb.buffer = append(cb.buffer, msg)
	for i, m := range cb.buffer {
		if cb.isCausal(m.deps, cb.delivered) {
			(*(cb.cp)).ReceivePrimitive(m.msg, m.msgType)
			cb.buffer = slices.Delete(cb.buffer, i, i+1)
			cb.delivered[m.sender] += 1
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
