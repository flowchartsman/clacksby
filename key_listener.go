package main

import (
	hook "github.com/robotn/gohook"
)

type keyboardEventListener struct {
	enabled   bool
	stopChan  chan struct{}
	eventChan chan hook.Event
}

func newListener() *keyboardEventListener {
	return &keyboardEventListener{}
}

func (l *keyboardEventListener) start() chan hook.Event {
	l.enabled = true
	l.eventChan = hook.Start()
	return l.eventChan
}

func (l *keyboardEventListener) stop() {
	if l.enabled {
		hook.End()
		l.enabled = false
	}
}

func (l *keyboardEventListener) events() chan hook.Event {
	return l.eventChan
}
