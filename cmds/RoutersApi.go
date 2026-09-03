package cmds

import (
	http2 "Kaban/internal/Deliver/httpController"
	"net/http"

	"github.com/gorilla/mux"
)

func GetLoginApi(postRequest *mux.Router, app *http2.NewLoginController) *mux.Route {
	return postRequest.HandleFunc("/login/api", func(writer http.ResponseWriter, request *http.Request) {
		app.LoginNet.W = writer
		app.LoginNet.R = request
		app.Login()
	}).Methods(http.MethodPost)
}
func GetRegisterApiRouter(postRequest *mux.Router, app *http2.NewRegister) *mux.Route {
	return postRequest.HandleFunc("/register/api", func(writer http.ResponseWriter, request *http.Request) {
		app.RegisterNet = http2.RegisterNet{
			W: writer,
			R: request,
		}
		app.Register()
	}).Methods(http.MethodPost)
}

func GetUploaderApiRouter(postRequest *mux.Router, app *http2.NewFileUploader) *mux.Route {
	return postRequest.HandleFunc("/downloader/api", func(writer http.ResponseWriter, request *http.Request) {
		app.FileUploaderNet = http2.FileUploaderNet{
			W: writer,
			R: request,
		}
		app.FileUploaderNoEncrypt(postRequest)
	}).Methods(http.MethodPost)
}

func GetMainApiRouter(router *mux.Router, app *http2.CheckUserAuth) *mux.Route {
	return router.HandleFunc("/maine/api", func(writer http.ResponseWriter, request *http.Request) {
		app.W = writer
		app.R = request
		app.CheckUserAuth()

	}).Methods("GET")
}
func GetDoUrlApiRouter(router *mux.Router, app *http2.NewBuildUrl) *mux.Route {
	return router.HandleFunc("/doUrl/api", func(writer http.ResponseWriter, request *http.Request) {
		app.Net.W = writer
		app.SetUrl()
	}).Methods(http.MethodGet)
}

func GetEncryptUploaderApiRouter(postRequest *mux.Router, app *http2.NewUploaderEncrypt) *mux.Route {
	return postRequest.HandleFunc("/downloader2/api", func(writer http.ResponseWriter, request *http.Request) {
		app.NewFileUploaderEncryptNetwork.W = writer
		app.NewFileUploaderEncryptNetwork.R = request
		app.FileUploaderEncrypt()

	}).Methods(http.MethodPost)
}

func GetDownloadApi(getRequest *mux.Router, app *http2.DownloadNew) *mux.Route {
	return getRequest.HandleFunc("/d/{name}", func(writer http.ResponseWriter, request *http.Request) {
		app.DownloadNetwork.W = writer
		app.DownloadNetwork.R = request
		app.DownloadWithNotEncrypt()
	}).Methods(http.MethodGet)
}
func GetEncryptDownloadApi(getRequest *mux.Router, app *http2.NewDownloadEncrypt) *mux.Route {
	return getRequest.HandleFunc("/d2/{name}", func(writer http.ResponseWriter, request *http.Request) {
		app.Net.W = writer
		app.Net.R = request
		app.DownloadWithEncrypt()
	}).Methods(http.MethodGet)
}
