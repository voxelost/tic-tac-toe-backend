package utils

import (
	"fmt"
	"sync"
)

type ModifiableQueue[T Identifiable] struct {
	sync.Mutex
	chan_    chan ID
	mappings map[ID]T
}

// Return a new ModifiableQueue[T]
func NewModifiableQueue[T Identifiable](queueSize int) *ModifiableQueue[T] {
	return &ModifiableQueue[T]{
		chan_:    make(chan ID, queueSize),
		mappings: make(map[ID]T),
	}
}

// Push an element on the top of the queue
func (mq *ModifiableQueue[T]) Push(obj T) {
	mq.Lock()
	defer mq.Unlock()

	id := obj.GetId()
	if _, ok := mq.mappings[id]; ok {
		return
	}

	mq.chan_ <- id
	mq.mappings[id] = obj
}

// Block until there's at least one element in the queue, pop it
func (mq *ModifiableQueue[T]) PopBlocking() T {
	id := <-mq.chan_

	mq.Lock()
	defer mq.Unlock()

	obj, ok := mq.mappings[id]
	if !ok {
		panic(fmt.Sprintf("ModifiableQueue Queue and Cache mismatch! Couldn't find object of ID: %s", id))
	}

	delete(mq.mappings, id)

	return obj
}

// Pop the first element in queue if one exists and return ok=true, if not, return ok=false
func (mq *ModifiableQueue[T]) PopNonBlocking() (obj T, ok bool) {
	if len(mq.chan_) < 1 {
		return obj, false
	}

	id := <-mq.chan_
	obj, ok = mq.mappings[id]
	if !ok {
		panic(fmt.Sprintf("ModifiableQueue Queue and Cache mismatch! Couldn't find object of ID: %s", id))
	}

	delete(mq.mappings, id)

	return obj, ok
}

// Return an object saved under given ID.
func (mq *ModifiableQueue[T]) Get(id ID) T {
	mq.Lock()
	defer mq.Unlock()

	return mq.mappings[id]
}

func (mq *ModifiableQueue[T]) Delete(Id ID) {
	mq.RemoveMultiple([]ID{Id})
}

// Remove Multiple IDs from the queue
func (mq *ModifiableQueue[T]) RemoveMultiple(Ids []ID) {
	mq.Lock()
	defer mq.Unlock()

	idsToPreserve := []ID{}
	for len(mq.chan_) > 0 {
		id_ := <-mq.chan_
		shouldPreserve := true
		for _, idToDelete := range Ids {
			if id_ == idToDelete {
				shouldPreserve = false
				break
			}
		}

		_, exists := mq.mappings[id_]
		if shouldPreserve && exists {
			idsToPreserve = append(idsToPreserve, id_)
		} else {
			delete(mq.mappings, id_)
		}
	}

	for _, id := range idsToPreserve {
		mq.chan_ <- id
	}
}

// Get the list if IDs in the same order as they are present in the queue
func (mq *ModifiableQueue[T]) GetIds() []ID {
	mq.Lock()
	defer mq.Unlock()

	ids := []ID{}
	for len(mq.chan_) > 0 {
		ids = append(ids, <-mq.chan_)
	}

	for _, id := range ids {
		mq.chan_ <- id
	}

	return ids
}

// Get current queue length
func (mq *ModifiableQueue[T]) Len() int {
	return len(mq.chan_)
}
