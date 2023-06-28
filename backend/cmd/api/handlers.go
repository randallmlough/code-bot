package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/randallmlough/code-bot/internal/openai"
	"github.com/randallmlough/code-bot/internal/request"
	"github.com/randallmlough/code-bot/internal/response"
)

type ConvertInput struct {
	To   string `json:"to"`
	From string `json:"from,omitempty"`
}

func (app *application) convert(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	input := ConvertInput{}
	if err := request.DecodeJSON(w, r, &input); err != nil {
		app.logger.Error(err, nil)
		app.badRequest(w, r, errors.New("bad request"))
		return
	}
	ctx, cancel := context.WithTimeout(ctx, 3*time.Minute)
	defer cancel()

	prompt := openai.ConvertPromptGenerator(input.From, input.To)
	result, err := app.gpt.Convert(ctx, prompt)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			app.logger.Error(fmt.Errorf("openAI took too long: %w", err), nil)
			app.contextTimeout(w, r)
			return
		}
		app.serverError(w, r, errors.New("request failed"))
		return
	}

	data := response.Envelope{
		"data": result,
	}

	if err := response.JSON(w, http.StatusOK, data); err != nil {
		app.serverError(w, r, err)
	}
}
