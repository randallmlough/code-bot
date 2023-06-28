package main

import (
	"fmt"
	"runtime/debug"
)

//nolint:all
func (app *application) backgroundTask(fn func()) {
	app.wg.Add(1)

	go func() {
		defer app.wg.Done()

		defer func() {
			err := recover()
			if err != nil {
				app.logger.Error(fmt.Errorf("%s", err), debug.Stack())
			}
		}()

		fn()
	}()
}
