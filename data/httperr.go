package data

// ErrPayload : object model for the error page,
// send tthis to the error page so as to elaborate an error
type ErrPayload struct {
	Code   uint   `json:"code"`
	Status string `json:"status"`
	Msg    string `json:"msg"`
	GoBack string `json:"goback"` // url to go back to the last ok page
}
