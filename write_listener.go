package main

import (
  "log"
  "reflect"
  "sync"
  "time"
  "github.com/fsnotify/fsnotify"
)

func deduplicatedWriteListener(writeFunction interface{}, writeFunctionArgs ...interface{}) {
	var (
		waitFor = 100 * time.Millisecond // Wait 100ms for new events; each new event resets the timer.
		mu     sync.Mutex
		timer  *time.Timer
	)

  // Restructure additional arguments for writeFunction
  writeFunctionArgVals := make([]reflect.Value, len(writeFunctionArgs))
    for i, arg := range writeFunctionArgs {
		writeFunctionArgVals[i] = reflect.ValueOf(arg)
	}

  // File system watcher event loop
	for {
		select {
		case err, ok := <-watcher.Errors: // Read from Errors
			if !ok { // Channel was closed
				return
			}
			log.Printf("ERROR: %s", err)

		case event, ok := <-watcher.Events: // Read from Events
			if !ok { // Channel was closed
				return
			}

			// Ignore everything outside of Create and Write
			if !event.Has(fsnotify.Create) && !event.Has(fsnotify.Write) {
				continue
			}

			mu.Lock()
			if timer == nil {
				timer = time.AfterFunc(2 * time.Second, func() {
				  // Call the writefunction supplied as argument
          reflect.ValueOf(writeFunction).Call(writeFunctionArgVals)
				 })
			}

			// Reset the timer, so it will start from 100ms again
			timer.Reset(waitFor)
      mu.Unlock()
		}
	}
}