package bot

import (
	"sync"
)

type groupMsgHeap struct {
	heap []GroupMsgHandler
	lock sync.Mutex
}

type privateMsgHeap struct {
	heap []PrivateMsgHandler
	lock sync.Mutex
}

type listenerHeap struct {
	heap []listener
	lock sync.Mutex
}

//准备排序

func (h *listenerHeap) Len() int { return len(h.heap) }
func (h *listenerHeap) Less(i, j int) bool {
	return h.heap[i].getPriority() > h.heap[j].getPriority()
}
func (h *listenerHeap) Swap(i, j int) {
	h.heap[i], h.heap[j] = h.heap[j], h.heap[i]
}
func (h *listenerHeap) Push(l *listener) {
	h.heap = append(h.heap, *l)
}

func (h *groupMsgHeap) Len() int { return len(h.heap) }
func (h *groupMsgHeap) Less(i, j int) bool {
	return h.heap[i].getPriority() > h.heap[j].getPriority()
}
func (h *groupMsgHeap) Swap(i, j int) {
	h.heap[i], h.heap[j] = h.heap[j], h.heap[i]
}
func (h *groupMsgHeap) Push(l *GroupMsgHandler) {
	h.heap = append(h.heap, *l)
}

func (h *privateMsgHeap) Len() int { return len(h.heap) }
func (h *privateMsgHeap) Less(i, j int) bool {
	return h.heap[i].getPriority() > h.heap[j].getPriority()
}
func (h *privateMsgHeap) Swap(i, j int) {
	h.heap[i], h.heap[j] = h.heap[j], h.heap[i]
}
func (h *privateMsgHeap) Push(l *PrivateMsgHandler) {
	h.heap = append(h.heap, *l)
}
