package core

import "sync"

var (
	listeners = make(map[string][]any)
	mu        sync.RWMutex
)

func On[T Event](eventName string, fn func(T)) {
	mu.Lock()
	defer mu.Unlock()
	listeners[eventName] = append(listeners[eventName], fn)
}

func Dispatch[T Event](e T) {
	mu.RLock()
	defer mu.RUnlock()
	for _, fn := range listeners[e.Name()] {
		if handler, ok := fn.(func(T)); ok {
			go handler(e)
		}
	}
}
