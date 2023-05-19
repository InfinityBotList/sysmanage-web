package actions

import "net/http"

type Action struct {
	Name          string
	Description   string
	ConfirmDialog string                                        // If unset, no confirm dialog will be shown
	Handler       func(*ActionContext) (*ActionResponse, error) `json:"-"`
}

type ActionContext struct {
	Request *http.Request
}

type ActionResponse struct {
	StatusCode int
	Resp       string
	TaskID     string
}
