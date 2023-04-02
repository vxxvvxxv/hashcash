package helpers

import (
	"fmt"
	"sync"
)

type ErrorList struct {
	errorList map[string]int
	mu        sync.RWMutex
}

func NewErrorList() *ErrorList {
	return &ErrorList{
		errorList: make(map[string]int),
	}
}

func (e *ErrorList) AddError(err error) {
	if err == nil {
		return
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	if _, ok := e.errorList[err.Error()]; !ok {
		e.errorList[err.Error()] = 1
	} else {
		e.errorList[err.Error()]++
	}
}

func (e *ErrorList) ParseError() error {
	e.mu.RLock()
	defer e.mu.RUnlock()

	var errParsed error

	for errDesc, counter := range e.errorList {
		if errParsed == nil {
			errParsed = fmt.Errorf("%s, count: %d", errDesc, counter)
			continue
		} else {
			errParsed = fmt.Errorf("%w, %s, count: %d", errParsed, errDesc, counter)
		}
	}

	return errParsed
}
