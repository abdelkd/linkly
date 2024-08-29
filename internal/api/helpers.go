package api

import (
	"encoding/json"
	"io"
	"net/http"
)

func (app *Application) writeJson(w http.ResponseWriter, data any, status int) {
	w.WriteHeader(status)

	jsonData, err := json.Marshal(data)
	if err != nil {
		app.ErrorLog.Println(err.Error())
		app.serverError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.Write(jsonData)
	w.Write([]byte("\n"))
}

func (app *Application) serverError(w http.ResponseWriter, err error) {
	app.ErrorLog.Println(err.Error())
	app.jsonMessage(w, "internal server error", false, http.StatusInternalServerError)
}

type ResponseMessage struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
}

func (app *Application) jsonMessage(w http.ResponseWriter, message string, ok bool, code int) {
	var response ResponseMessage

	response.Ok = ok
	response.Message = message

	app.writeJson(w, response, code)
}

func (app *Application) GetRequestBody(r *http.Request, v any) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, v)
}
