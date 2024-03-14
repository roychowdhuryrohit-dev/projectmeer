package algos

import (
	"bytes"
	"compress/gzip"
	"encoding/gob"
	"fmt"
	"slices"

	config "github.com/roychowdhuryrohit-dev/projectmeer/node-user/lib"
)

type ID struct {
	Sender  string
	Counter int
}

type Side string

const (
	LeftSide  Side = "L"
	RightSide Side = "R"
)

type Node[T any] struct {
	Id                          *ID
	Value                       T
	IsDeleted                   bool
	Parent, RightOrigin         *Node[T]
	LeftChildren, RightChildren []*Node[T]
	Side                        Side
	Size                        int
}

type InsertMessage[T any] struct {
	Id          ID
	Value       T
	Parent      ID
	Side        Side
	RightOrigin ID
}

type DeleteMessage struct {
	Id ID
}

type NodeSave[T any] struct {
	Value               T
	IsDeleted           bool
	Parent, RightOrigin *ID
	Side                Side
	Size                int
}

type Tree[T any] struct {
	Root      *Node[T]
	nodesByID map[string][]*Node[T]
}

func NewTree[T any]() *Tree[T] {
	root := &Node[T]{
		Id: &ID{
			Sender:  "",
			Counter: 0,
		},
		Side:          RightSide,
		IsDeleted:     true,
		LeftChildren:  make([]*Node[T], 0),
		RightChildren: make([]*Node[T], 0),
	}
	return &Tree[T]{
		Root: root,
		nodesByID: map[string][]*Node[T]{
			"": {root},
		},
	}
}

func (tr *Tree[T]) AddNode(id *ID, value T, parent *Node[T], side Side, rightOriginID *ID) error {
	node := &Node[T]{
		Id:            id,
		Side:          side,
		IsDeleted:     false,
		Value:         value,
		Parent:        parent,
		LeftChildren:  make([]*Node[T], 0),
		RightChildren: make([]*Node[T], 0),
	}
	if rightOriginID != nil {
		rightOrigin, err := tr.GetByID(rightOriginID)
		if err != nil {
			return err
		}
		node.RightOrigin = rightOrigin
	}

	if _, ok := tr.nodesByID[id.Sender]; !ok {
		tr.nodesByID[id.Sender] = []*Node[T]{}
	}
	tr.nodesByID[id.Sender] = append(tr.nodesByID[id.Sender], node)

	tr.InsertIntoSiblings(node)
	tr.UpdateSize(node, 1)
	return nil
}

func (tr *Tree[T]) GetByID(id *ID) (*Node[T], error) {
	bySender, ok := tr.nodesByID[id.Sender]
	if !ok {
		return nil, fmt.Errorf("unknown ID (%+v)", *id)
	}
	if 0 > id.Counter || id.Counter > len(bySender) {
		return nil, fmt.Errorf("unknown ID (%+v)", *id)
	}
	return bySender[id.Counter], nil
}

func (tr *Tree[T]) InsertIntoSiblings(node *Node[T]) {
	parent := node.Parent
	if node.Side == RightSide {
		i := 0
		for ; i < len(parent.RightChildren); i++ {
			if node.RightOrigin != nil && parent.RightChildren[i].RightOrigin != nil && !(tr.isLess(node.RightOrigin, parent.RightChildren[i].RightOrigin) || (node.RightOrigin == parent.RightChildren[i].RightOrigin && node.Id.Sender > parent.RightChildren[i].Id.Sender)) {
				break
			}
		}
		parent.RightChildren = slices.Insert(parent.RightChildren, i, node)

	} else {
		i := 0
		for ; i < len(parent.LeftChildren); i++ {
			if !(node.Id.Sender > parent.LeftChildren[i].Id.Sender) {
				break
			}
		}
		parent.LeftChildren = slices.Insert(parent.LeftChildren, i, node)
	}
}

func (tr *Tree[T]) UpdateSize(node *Node[T], delta int) {
	for anc := node; anc != nil; anc = anc.Parent {
		anc.Size += delta
	}
}

func (tr *Tree[T]) isLess(a *Node[T], b *Node[T]) bool {
	if a == b {
		return false
	}
	if a == nil {
		return false
	}
	if b == nil {
		return true
	}
	aDepth := tr.depth(a)
	bDepth := tr.depth(b)

	aAnc := a
	bAnc := b
	if aDepth > bDepth {
		var lastSide Side
		for i := aDepth; i > bDepth; i-- {
			lastSide = aAnc.Side
			aAnc = aAnc.Parent
		}
		if aAnc == b {
			return lastSide != LeftSide
		}
	}
	if bDepth > aDepth {
		var lastSide Side
		for i := bDepth; i > aDepth; i-- {
			lastSide = bAnc.Side
			bAnc = bAnc.Parent
		}
		if bAnc == a {
			return lastSide != RightSide
		}
	}
	for aAnc.Parent != bAnc.Parent {
		aAnc = aAnc.Parent
		bAnc = bAnc.Parent
	}
	if aAnc.Side != bAnc.Side {
		return aAnc.Side == LeftSide
	} else if aAnc.Parent != nil {
		siblings := aAnc.Parent.RightChildren
		if aAnc.Side == LeftSide {
			siblings = aAnc.Parent.LeftChildren
		}
		return slices.Index(siblings, aAnc) < slices.Index(siblings, bAnc)
	}
	return false
}

func (tr *Tree[T]) depth(node *Node[T]) int {
	depth := 0
	for current := node; current.Parent != nil; current = current.Parent {
		depth++
	}
	return depth
}

func (tr *Tree[T]) GetByIndex(node *Node[T], index int) (*Node[T], error) {
	if index < 0 || index >= node.Size {
		return nil, fmt.Errorf("index out of range")
	}

	remaining := index
recurse:
	for {
		for _, child := range node.LeftChildren {
			if remaining < child.Size {
				node = child
				continue recurse
			}
			remaining -= child.Size
		}
		if !node.IsDeleted {
			if remaining == 0 {
				return node, nil
			}
			remaining--
		}
		for _, child := range node.RightChildren {
			if remaining < child.Size {
				node = child
				continue recurse
			}
			remaining -= child.Size
		}
		return nil, fmt.Errorf("index in range but not found")
	}
}

func (tr *Tree[T]) LeftmostDescendent(node *Node[T]) *Node[T] {
	desc := node
	for len(desc.LeftChildren) != 0 {
		desc = desc.LeftChildren[0]
	}
	return desc
}

func (tr *Tree[T]) NextNonDescendent(node *Node[T]) *Node[T] {
	current := node
	for current.Parent != nil {
		siblings := current.Parent.RightChildren
		if current.Side == LeftSide {
			siblings = current.Parent.LeftChildren
		}
		index := slices.Index(siblings, current)
		if index < (len(siblings) - 1) {
			nextSibling := siblings[index+1]
			return tr.LeftmostDescendent(nextSibling)
		} else if current.Side == LeftSide {
			return current.Parent
		}
		current = current.Parent
	}
	return nil
}

func (tr *Tree[T]) Traverse(node *Node[T]) []T {
	var result []T
	current := node
	type Child struct {
		side       Side
		childIndex int
	}
	stack := []*Child{{side: LeftSide, childIndex: 0}}
	for {
		top := stack[len(stack)-1]
		children := current.RightChildren
		if top.side == LeftSide {
			children = current.LeftChildren
		}
		if top.childIndex == len(children) {
			if top.side == LeftSide {
				if !current.IsDeleted {
					result = append(result, current.Value)
				}
				top.side = RightSide
				top.childIndex = 0
			} else {
				if current.Parent == nil {
					return result
				}
				current = current.Parent
				if len(stack) > 0 {
					stack = stack[:len(stack)-1]
				}
			}
		} else {
			child := children[top.childIndex]
			top.childIndex++
			if child.Size > 0 {
				current = child
				stack = append(stack, &Child{side: LeftSide, childIndex: 0})
			}
		}
	}
}

func (tr *Tree[T]) Save() ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	save := make(map[string][]*NodeSave[T], 0)
	for sender, bySender := range tr.nodesByID {
		nodeSaves := make([]*NodeSave[T], len(bySender))
		for i, node := range bySender {
			nodeSaves[i].Value = node.Value
			nodeSaves[i].IsDeleted = node.IsDeleted
			if node.Parent == nil {
				nodeSaves[i].Parent = nil
			} else {
				nodeSaves[i].Parent = node.Parent.Id
			}
			nodeSaves[i].Side = node.Side
			nodeSaves[i].Size = node.Size
		}
		save[sender] = nodeSaves
	}

	if err := enc.Encode(save); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (tr *Tree[T]) Load(saveData []byte) error {
	buf := bytes.NewBuffer(saveData)
	dec := gob.NewDecoder(buf)
	save := make(map[string][]*NodeSave[T], 0)

	if err := dec.Decode(&save); err != nil {
		return err
	}

	for sender, bySenderSave := range save {
		if sender == "" {
			tr.Root.Size = bySenderSave[0].Size
			continue
		}

		for counter, nodeSave := range bySenderSave {
			var node *Node[T]
			node.Id.Sender = sender
			node.Id.Counter = counter
			node.Parent = nil
			node.Value = nodeSave.Value
			node.IsDeleted = nodeSave.IsDeleted
			node.Side = nodeSave.Side
			node.Size = nodeSave.Size
			node.LeftChildren = make([]*Node[T], 0)
			node.RightChildren = make([]*Node[T], 0)
			tr.nodesByID[sender] = append(tr.nodesByID[sender], node)
		}
	}

	for sender, bySender := range tr.nodesByID {
		if sender == "" {
			continue
		}
		if save[sender] != nil {
			bySenderSave := save[sender]
			for i := 0; i < len(bySender); i++ {
				node := bySender[i]
				nodeSave := bySenderSave[i]
				if nodeSave.Parent != nil {
					var err error
					node.Parent, err = tr.GetByID(nodeSave.Parent)
					if err != nil {
						return err
					}
				}
				if nodeSave.RightOrigin != nil {
					var err error
					node.RightOrigin, _ = tr.GetByID(nodeSave.RightOrigin)
					if err != nil {
						return err
					}
				} else {
					node.RightOrigin = nil
				}
			}
		}
	}
	var readyNodes []*Node[T]
	pendingNodes := make(map[*Node[T]][]*Node[T])

	for sender, bySender := range tr.nodesByID {
		if sender == "" {
			continue
		}
		for i := 0; i < len(bySender); i++ {
			node := bySender[i]
			if node.RightOrigin == nil {
				readyNodes = append(readyNodes, node)
			} else {
				_, ok := pendingNodes[node.RightOrigin]
				if !ok {
					pendingNodes[node.RightOrigin] = make([]*Node[T], 0)
				}
				pendingNodes[node.RightOrigin] = append(pendingNodes[node.RightOrigin], node)
			}
		}
	}

	for len(readyNodes) != 0 {
		node := readyNodes[len(readyNodes)-1]
		readyNodes = readyNodes[:len(readyNodes)-1]
		tr.InsertIntoSiblings(node)
		deps := pendingNodes[node]
		if deps != nil {
			readyNodes = append(readyNodes, deps...)
		}
		delete(pendingNodes, node)
	}
	if len(pendingNodes) != 0 {
		return fmt.Errorf("failed to validate all nodes")
	}
	return nil
}

type FugueMax[T any] struct {
	CausalBroadcast[T]
	counter int
	tree    *Tree[T]
}

func NewFugueMax[T any](curNode int, nodeList []string) *FugueMax[T] {
	fg := &FugueMax[T]{
		counter: 0,
		tree:    NewTree[T](),
		CausalBroadcast: CausalBroadcast[T]{
			sendSeq:   0,
			delivered: make([]int, len(nodeList)),
			buffer:    make([]Message, 0),
			nodeList:  nodeList,
			curNode:   curNode,
		},
	}
	*fg.CausalBroadcast.cp = fg
	return fg
}

func (fg *FugueMax[T]) Insert(index int, values ...T) error {
	for i := 0; i < len(values); i++ {
		err := fg.insertOne(index+i, values[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func (fg *FugueMax[T]) insertOne(index int, value T) error {
	rid, ok := config.ConfigMap.Load(config.ReplicaID)
	if !ok {
		return fmt.Errorf("missing replica ID")
	}
	id := &ID{
		Sender:  rid.(string),
		Counter: fg.counter,
	}
	fg.counter++
	leftOrigin := fg.tree.Root
	if index != 0 {
		leftOrigin, _ = fg.tree.GetByIndex(fg.tree.Root, index-1)
	}
	var msg InsertMessage[T]
	if len(leftOrigin.RightChildren) == 0 {
		msg = InsertMessage[T]{
			Id:     *id,
			Value:  value,
			Parent: *leftOrigin.Id,
			Side:   RightSide,
		}
		rightOrigin := fg.tree.NextNonDescendent(leftOrigin)
		msg.RightOrigin = ID{}
		if rightOrigin != nil {
			msg.RightOrigin = *rightOrigin.Id
		}

	} else {
		rightOrigin := fg.tree.LeftmostDescendent(leftOrigin.RightChildren[0])
		msg = InsertMessage[T]{
			Id:     *id,
			Value:  value,
			Parent: *rightOrigin.Id,
			Side:   LeftSide,
		}
	}
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(msg); err != nil {
		return err
	}
	fg.SendPrimitive(buf.Bytes(), "insert")
	return nil
}

func (fg *FugueMax[T]) Delete(startIndex int, count int) error {
	for i := 0; i < count; i++ {
		err := fg.deleteOne(startIndex)
		if err != nil {
			return err
		}
	}
	return nil
}

func (fg *FugueMax[T]) deleteOne(index int) error {
	node, err := fg.tree.GetByIndex(fg.tree.Root, index)
	if err != nil {
		return err
	}
	msg := DeleteMessage{
		Id: *node.Id,
	}

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(msg); err != nil {
		return err
	}
	fg.SendPrimitive(buf.Bytes(), "delete")
	return nil
}

func (fg *FugueMax[T]) ReceivePrimitive(message []byte, messageType string) error {
	buf := bytes.NewBuffer(message)
	dec := gob.NewDecoder(buf)
	switch messageType {
	case "insert":

		iMsg := InsertMessage[T]{}

		if err := dec.Decode(&iMsg); err != nil {
			return err
		}
		parent, err := fg.tree.GetByID(&iMsg.Parent)
		if err != nil {
			return err
		}
		fg.tree.AddNode(&iMsg.Id, iMsg.Value, parent, iMsg.Side, &iMsg.RightOrigin)
	case "delete":
		dMsg := DeleteMessage{}

		if err := dec.Decode(&dMsg); err != nil {
			return err
		}
		node, err := fg.tree.GetByID(&dMsg.Id)
		if err != nil {
			return err
		}
		if !node.IsDeleted {
			node.Value = *new(T)
			node.IsDeleted = true
			fg.tree.UpdateSize(node, -1)
		}
	default:
		return fmt.Errorf("bad message type")
	}
	return nil
}

func (fg *FugueMax[T]) Length() int {
	return fg.tree.Root.Size
}

func (fg *FugueMax[T]) Get(index int) (T, error) {
	if index < 0 || index >= fg.Length() {
		return *new(T), fmt.Errorf("index out of bounds")
	}
	node, err := fg.tree.GetByIndex(fg.tree.Root, index)
	if err != nil {
		return *new(T), err
	}

	return node.Value, nil
}

func (fg *FugueMax[T]) Values() []T {
	return fg.tree.Traverse(fg.tree.Root)
}

func (fg *FugueMax[T]) SavePrimitive() ([]byte, error) {
	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	s, err := fg.tree.Save()
	if err != nil {
		return nil, err
	}
	w.Write(s)
	w.Close()
	return b.Bytes(), nil
}

func (fg *FugueMax[T]) LoadPrimitives(savedState []byte) error {
	if savedState == nil {
		return fmt.Errorf("saved state is null")
	}
	b := bytes.NewBuffer(savedState)
	r, err := gzip.NewReader(b)
	if err != nil {
		return err
	}
	defer r.Close()
	var s bytes.Buffer
	_, err = s.ReadFrom(r)
	if err != nil {
		return err
	}
	return fg.tree.Load(s.Bytes())
}
