package utils

import "sync"

type EventEmitter struct {
	listeners map[string](map[string]func(interface{}))
	mu        *sync.RWMutex
}

func NewEventEmitter() *EventEmitter {
	return &EventEmitter{
		mu: &sync.RWMutex{},
	}
}

func (emitter *EventEmitter) AddListener(event string, listenerName string, listener func(interface{})) {
	emitter.mu.Lock()
	if emitter.listeners[event] == nil {
		emitter.listeners[event] = make(map[string]func(interface{}))
	}
	emitter.listeners[event][listenerName] = listener
	emitter.mu.Unlock()
}

func (emitter *EventEmitter) RemoveListener(event string, listenerName string) {
	emitter.mu.Lock()
	delete(emitter.listeners[event], listenerName)
	emitter.mu.Unlock()
}

func (emitter *EventEmitter) Emit(event string, data interface{}) {
	emitter.mu.RLock()
	if emitter.listeners[event] == nil {
		emitter.mu.RUnlock()
		return
	}

	for _, listener := range emitter.listeners[event] {
		go listener(data)
	}
	emitter.mu.RLock()
}
