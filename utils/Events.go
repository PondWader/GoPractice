package utils

type EventEmitter struct {
	listeners map[string](map[string]func(interface{}))
}

func (emitter *EventEmitter) AddListener(event string, listenerName string, listener func(interface{})) {
	if emitter.listeners[event] == nil {
		emitter.listeners[event] = make(map[string]func(interface{}))
	}
	emitter.listeners[event][listenerName] = listener
}

func (emitter *EventEmitter) RemoveListener(event string, listenerName string) {
	delete(emitter.listeners[event], listenerName)
}

func (emitter *EventEmitter) Emit(event string, data interface{}) {
	if emitter.listeners[event] == nil {
		return
	}

	for _, listener := range emitter.listeners[event] {
		go listener(data)
	}
}
