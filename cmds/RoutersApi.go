package cmds

import (
	"net/http"

	"github.com/gorilla/mux"
)

func GetLoginApi(postRequest *mux.Router) *mux.Route {
	return postRequest.HandleFunc("/login/api", func(writer http.ResponseWriter, request *http.Request) {
		Controller2.Login(writer, request, Sa)

	}).Methods("POST")
}
func GetRegisterApi(postRequest *mux.Router) *mux.Route {
	return postRequest.HandleFunc("/register/api", func(writer http.ResponseWriter, request *http.Request) {
		Controller2.Register(writer, request, Sa)

	}).Methods("POST")
}

func GetDownloadApi(postRequest *mux.Router, router *mux.Router) *mux.Route {
	return postRequest.HandleFunc("/downloader/api", func(writer http.ResponseWriter, request *http.Request) {

		Controller2.FileUploaderNoEncrypt(writer, request, router, Sa)

	}).Methods(http.MethodPost)
}

func GetMainApi(router *mux.Router) *mux.Route {
	return router.HandleFunc("/maine/api", func(writer http.ResponseWriter, request *http.Request) {
		Controller2.CheckUserAuth(writer, request, Sa)

	}).Methods("GET")
}
func GetDoUrlApi(router *mux.Router) *mux.Route {
	return router.HandleFunc("/doUrl/api", func(writer http.ResponseWriter, request *http.Request) {

		Controller2.BuildUrl(writer, request)

	}).Methods(http.MethodGet)
}

func GetEncryptDownloadApi(postRequest *mux.Router, router *mux.Router) *mux.Route {
	return postRequest.HandleFunc("/downloader2/api", func(writer http.ResponseWriter, request *http.Request) {

		Controller2.FileUploaderEncrypt(writer, request, router, Sa)

	}).Methods(http.MethodPost)
}

func GetDownloadRequest(getRequest *mux.Router) *mux.Route {
	return getRequest.HandleFunc("/d/{name}", func(writer http.ResponseWriter, request *http.Request) {

		Controller2.DownloadWithNotEncrypt(writer, request, Sa)

		//Application.Delete(ch)

	}).Methods(http.MethodGet)
}

func GetEncryptDownload(getRequest *mux.Router) *mux.Route {
	return getRequest.HandleFunc("/d2/{name}", func(writer http.ResponseWriter, request *http.Request) {

		Controller2.DownloadWithEncrypt(writer, request, Sa)

		//Application.Delete(ch)

	}).Methods(http.MethodGet)
}
