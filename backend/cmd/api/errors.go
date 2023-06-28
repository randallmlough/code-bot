package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/randallmlough/code-bot/internal/response"
	"github.com/randallmlough/code-bot/internal/validator"
)

func (app *application) badRequest(w http.ResponseWriter, r *http.Request, err error) {
	http.Error(w, err.Error(), http.StatusBadRequest)
}

func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	trace := debug.Stack()

	app.logger.Error(err, trace)

	// if app.config.notifications.email != "" {
	// 	app.sendErrorNotification(r, err, trace)
	// }

	message := "The server encountered a problem and could not process your request"
	app.errorMessage(w, r, http.StatusInternalServerError, message, nil)
}

func (app *application) errorMessage(w http.ResponseWriter, r *http.Request, status int, message string, headers http.Header) {
	message = strings.ToUpper(message[:1]) + message[1:]

	err := response.JSONWithHeaders(w, status, map[string]string{"Error": message}, headers)
	if err != nil {
		app.logger.Error(err, debug.Stack())
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (app *application) methodNotAllowed(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("The %s method is not supported for this resource", r.Method)
	app.errorMessage(w, r, http.StatusMethodNotAllowed, message, nil)
}

func (app *application) contextTimeout(w http.ResponseWriter, r *http.Request) {
	app.errorMessage(w, r, http.StatusRequestTimeout, "Request took longer than expected", nil)
}

func (app *application) notFound(w http.ResponseWriter, r *http.Request) {
	app.errorMessage(w, r, http.StatusNotFound, "Resource not found", nil)
}

//nolint:all
func (app *application) failedValidation(w http.ResponseWriter, r *http.Request, v validator.Validator) {
	err := response.JSON(w, http.StatusUnprocessableEntity, v)
	if err != nil {
		app.serverError(w, r, err)
	}
}
