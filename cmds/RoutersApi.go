package cmds

import (
	"Kaban/internal/Deliver"
	"net/http"

	"github.com/gorilla/mux"
)

func GetLoginApi(postRequest *mux.Router, app *Deliver.NewLoginController) *mux.Route {
	return postRequest.HandleFunc("/login/api", func(writer http.ResponseWriter, request *http.Request) {
		app.LoginNet.W = writer
		app.LoginNet.R = request
		app.Login()
	}).Methods(http.MethodPost)
}
func GetRegisterApiRouter(postRequest *mux.Router, app *Deliver.NewRegister) *mux.Route {
	return postRequest.HandleFunc("/register/api", func(writer http.ResponseWriter, request *http.Request) {
		app.RegisterNet = Deliver.RegisterNet{
			W: writer,
			R: request,
		}
		app.Register()
	}).Methods(http.MethodPost)
}

func GetUploaderApiRouter(postRequest *mux.Router, app *Deliver.NewUploader) *mux.Route {
	return postRequest.HandleFunc("/downloader/api", func(writer http.ResponseWriter, request *http.Request) {
		app.NewFileUploaderNet = Deliver.NewFileUploaderNet{
			W: writer,
			R: request,
		}
		app.FileUploaderNoEncrypt(postRequest)
	}).Methods(http.MethodPost)
}

func GetMainApiRouter(router *mux.Router, app *Deliver.NewCheckUserAuth) *mux.Route {
	return router.HandleFunc("/maine/api", func(writer http.ResponseWriter, request *http.Request) {
		app.W = writer
		app.R = request
		app.CheckUserAuth()

	}).Methods("GET")
}
func GetDoUrlApiRouter(router *mux.Router, app *Deliver.NewBuildUrl) *mux.Route {
	return router.HandleFunc("/doUrl/api", func(writer http.ResponseWriter, request *http.Request) {
		app.Net.W = writer
		app.SetUrl()
	}).Methods(http.MethodGet)
}

func GetEncryptUploaderApiRouter(postRequest *mux.Router, app *Deliver.NewUploaderEncrypt) *mux.Route {
	return postRequest.HandleFunc("/downloader2/api", func(writer http.ResponseWriter, request *http.Request) {
		app.NewFileUploaderEncryptNetwork.W = writer
		app.NewFileUploaderEncryptNetwork.R = request
		app.FileUploaderEncrypt()

	}).Methods(http.MethodPost)
}

func GetDownloadApi(getRequest *mux.Router, app *Deliver.NewDownloadWithNotEncrypt) *mux.Route {
	return getRequest.HandleFunc("/d/{name}", func(writer http.ResponseWriter, request *http.Request) {
		app.NetworkDownloadNoEncrypt.W = writer
		app.NetworkDownloadNoEncrypt.R = request
		app.DownloadWithNotEncrypt()
	}).Methods(http.MethodGet)
}
func GetEncryptDownloadApi(getRequest *mux.Router, app *Deliver.NewDownloadEncrypt) *mux.Route {
	return getRequest.HandleFunc("/d2/{name}", func(writer http.ResponseWriter, request *http.Request) {
		app.Net.W = writer
		app.Net.R = request
		app.DownloadWithEncrypt()
	}).Methods(http.MethodGet)
}
