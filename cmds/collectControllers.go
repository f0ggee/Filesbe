package cmds

import (
	"Kaban/internal/Deliver/httpController"

	"github.com/gorilla/mux"
)

type Controllers struct {
	download        httpController.DownloadController
	encryptDownload httpController.DownloadEncryptController
	uploader        httpController.UploadController
	encryptUpload   httpController.UploadEncryptController
	login           httpController.LoginController
	register        httpController.RegisterController
	urlBuild        httpController.BuildUrlController
	checkUser       httpController.CheckUserAuthController
}

type CollectorsController struct {
	Apps                AppBuilders
	AuthTokensCollector TokensAuthCollector
	SessionCollector    httpController.Session
	Router              *mux.Router
	ParserCollector     ParsersCollector
}

func CollectControllers(c CollectorsController) Controllers {
	controllerDownload := GetControllerDownloadBuilder(c.Apps.download)
	controllerEncryptDownload := GetControllerDownloadEncryptBuilder(&c.Apps.encryptDownload)
	controllerUploader := GetControllerFileUploadBuilder(ControllerFileUploadBuilder{
		Token:    &c.AuthTokensCollector.Rf,
		Sessions: c.SessionCollector,
		App:      &c.Apps.upload,
	})

	controllerUploadEncrypt := GetControllerFileUploadEncryptBuilder(ControllerFileUploadEncrypt{
		Token:    &c.AuthTokensCollector.Jwt,
		Sessions: c.SessionCollector,
		App:      c.Apps.encryptUpload,
	}, c.Router)
	controllerLogin := GetControllerLoginBuilder(ControllerLoginBuilder{
		Sessions: c.SessionCollector,
		Parser:   c.ParserCollector,
		App:      &c.Apps.login,
	})
	controllerRegister := GetControllerRegisterBuilder(RegisterController{
		Sessions: c.SessionCollector,
		App:      c.Apps.register,
		Parser:   c.ParserCollector,
	})
	controllerUrlBuild := GetControllerUrlBuildBuilder()
	controllerCheckUser := GetCheckUsers()
	return Controllers{
		download:        controllerDownload,
		encryptDownload: controllerEncryptDownload,
		uploader:        controllerUploader,
		encryptUpload:   controllerUploadEncrypt,
		login:           controllerLogin,
		register:        controllerRegister,
		urlBuild:        controllerUrlBuild,
		checkUser:       controllerCheckUser,
	}
}
