package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/abdelkd/linkly/internal/data"
	"github.com/abdelkd/linkly/internal/util"
	"github.com/abdelkd/linkly/internal/validator"
	"github.com/gorilla/mux"
)

const API_VERSION = 1

type HealthCheckResponse struct {
	Database   string `json:"database"`
	ApiVersion int    `json:"api_version"`
}

func (app *Application) handleHealthCheck(w http.ResponseWriter, _ *http.Request) {
	var response HealthCheckResponse

	databaseStatus := app.DB.QueryRow("SELECT 'Success'")
	if err := databaseStatus.Scan(&response.Database); err != nil {
		app.ErrorLog.Println(err)
		response.Database = "Failed"
	}

	response.ApiVersion = API_VERSION
	app.writeJson(w, response, http.StatusOK)
}

func (app *Application) handleNewLink(w http.ResponseWriter, r *http.Request) {
	var request data.CreateLinkRequest
	var response data.CreateLinkResponse

	err := app.GetRequestBody(r, &request)
	if err != nil {
		app.serverError(w, err)
		return
	}

	isUrl := validator.IsValidUrl(request.Location)
	if !isUrl {
		app.InfoLog.Println("invalid url")
		app.jsonMessage(w, "location must be a valid url", false, http.StatusBadRequest)
		return
	}

	linkByHash, err := app.Models.Links.GetByHashCode(util.HashString(request.Location))
	if !errors.Is(err, sql.ErrNoRows) {
		if err != nil {
			app.ErrorLog.Fatal(err)
		}

		if linkByHash != nil {
			response.Link = linkByHash.Url
			response.Code = linkByHash.Code

			app.writeJson(w, response, http.StatusFound)
			return
		}
	}

	var link data.Link
	link.Url = request.Location

	err = app.Models.Links.Add(&link)
	if err != nil {
		app.InfoLog.Println(err)
		app.jsonMessage(w, "failed to create a new link", false, http.StatusInternalServerError)
	}

	response.Link = fmt.Sprintf("%sv%d/link/%s", app.Env.BaseURL, API_VERSION, link.Code)
	response.Code = link.Code

	app.writeJson(w, response, http.StatusCreated)
}

func (app *Application) handleGetLink(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	linkCode := vars["id"]
	if linkCode == "" {
		app.jsonMessage(w, "please provide a valid link code", false, http.StatusNotFound)
		return
	}

	url, err := app.Models.Links.Get(linkCode)
	if err != nil {
		app.jsonMessage(w, "link could not be found", false, http.StatusNotFound)
		return
	}

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}
