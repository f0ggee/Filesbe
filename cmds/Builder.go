//This file contains special functions that build objects.
//It's a second part of building objects.

package cmds

import (
	"Kaban/internal/Deliver/httpController"
	"Kaban/internal/DomainLevel"
	"Kaban/internal/Service/Application"

	"github.com/gorilla/mux"
)

//WARNING! Every module's name must follow this syntax: SomethingBuilder
//There are application builder functions

type DownloadBuilderApplication struct {
	Trans FileTransferringCollector
	Red   RedisCollector
	Parse ParsersCollector
}

func GetDownloadApplicationBuilder(Transferring DownloadBuilderApplication) Application.NewDownload {
	return Application.NewDownload{
		DownloadDelivery: Application.DownloadDelivery{
			Reader: Transferring.Red.Read,
		},
		DownloadFileControl: Application.DownloadFileControl{
			Downloader: &Transferring.Trans.Download,
			Uploader:   &Transferring.Trans.HttpFileTransferringCollector.Upload,
			Deleter:    &Transferring.Trans.Delete,
			Decoder:    Transferring.Parse.Parser,
		},
	}
}

type DownloadEncryptApplication struct {
	Red          RedisCollector
	Crypt        CryptoCollector
	Keys         SessionKeysCollector
	Transferring FileTransferringCollector
}

func GetDownloadEncryptApplicationBuilder(D DownloadEncryptApplication) Application.NewDownloadEncrypt {
	return Application.NewDownloadEncrypt{
		NewDownloadEncryptDelivery: Application.NewDownloadEncryptDelivery{
			ReaderRedis:  D.Red.Read,
			DeleterRedis: D.Red.Delete,
		},
		NewDownloadEncryptCrypto: Application.NewDownloadEncryptCrypto{
			Crypto:        &D.Crypt.AesCtrRealization,
			EncrypterKeys: D.Keys.Keys,
		},
		NewDownloadEncryptFileControl: Application.NewDownloadEncryptFileControl{
			Uploader:   &D.Transferring.HttpFileTransferringCollector.Upload,
			Downloader: &D.Transferring.S3FileTransferringCollector.Download,
			Deleter:    &D.Transferring.Delete,
		},
	}
}

type LoginApplication struct {
	DB     DatabaseCollector
	Crypt  CryptoCollector
	Tokens TokensAuthCollector
}

func GetLoginApplicationBuilder(L LoginApplication) Application.NewLogin {

	return Application.NewLogin{
		ReaderDatabase: L.DB.Read,
		Validate:       &L.Crypt.Validate,
		Rf:             &L.Tokens.Rf,
		Jwt:            &L.Tokens.Rf,
	}
}

type RegisterApplication struct {
	DB     DatabaseCollector
	Tokens TokensAuthCollector
	Crypt  CryptoCollector
}

func GetRegisterBuilder(R RegisterApplication) Application.NewRegisterApplication {
	return Application.NewRegisterApplication{
		CheckingDb: nil,
		WriterDb:   nil,
		Generator:  nil,
		Rft:        nil,
		Jwt:        nil,
	}
}

type UploadApplication struct {
	Crypt CryptoCollector
	Parse ParsersCollector
	S3    S3FileTransferringCollector
	Red   RedisCollector
}

func GetUploadBuilder(U UploadApplication) Application.NewUpload {
	return Application.NewUpload{
		NewFileUploaderCrypto: Application.NewFileUploaderCrypto{
			Generator: U.Crypt.Generate,
		},
		NewFileUploaderDataMange: Application.NewFileUploaderDataMange{
			Encode: U.Parse.Parser,
		},
		NewFileUploaderDelivery: Application.NewFileUploaderDelivery{
			Uploader:   &U.S3.Upload,
			WriteRedis: &U.Red.Write,
		},
	}
}

type UploadEncrypt struct {
	Parse        ParsersCollector
	Crypt        CryptoCollector
	Transferring FileTransferringCollector
	Red          RedisCollector
}

func GetUploadEncryptBuilder(d UploadEncrypt) Application.NewUploadEncrypt {

	return Application.NewUploadEncrypt{
		NewUploadEncryptDataManage: Application.NewUploadEncryptDataManage{
			Encode: d.Parse.Parser,
		},
		NewUploadEncryptCrypto: Application.NewUploadEncryptCrypto{
			Generate: d.Crypt.Generate,
			Encrypt:  &d.Crypt.AesGcmRealization,
		},
		NewUploadEncryptDelivery: Application.NewUploadEncryptDelivery{
			Uploader:     &d.Transferring.S3FileTransferringCollector.Upload,
			RedisWriter:  &d.Red.Write,
			RedisChecker: &d.Red.Check,
			RedisDeleter: d.Red.Delete,
			Deleter:      &d.Transferring.S3FileTransferringCollector.Delete,
		},
	}
}

//There are controller builders

func GetControllerDownloadBuilder(app Application.NewDownload) httpController.DownloadController {
	return httpController.DownloadController{
		DownloadApplication: &app,
	}
}

func GetControllerDownloadEncryptBuilder(app Application.DownloadEncryptApplication) httpController.DownloadEncryptController {
	return httpController.DownloadEncryptController{
		DownloadEncryptApplication: app,
		GetDataRequest:             httpController.GetDataRequest,
	}
}

type ControllerFileUploadBuilder struct {
	Token    DomainLevel.AuthMaker
	Sessions httpController.Session
	App      Application.UploadApplication
}

func GetControllerFileUploadBuilder(C ControllerFileUploadBuilder) httpController.UploadController {

	return httpController.UploadController{
		FileUploaderSessions: httpController.FileUploaderSessions{
			Rf:       C.Token,
			Jwt:      C.Token,
			Sessions: C.Sessions,
		},
		UploadApplication: C.App,
		UrlData:           httpController.GetUploadData,
	}
}

type ControllerFileUploadEncrypt struct {
	Token    DomainLevel.AuthMaker
	Sessions httpController.Session
	App      Application.NewUploadEncrypt
}

func GetControllerFileUploadEncryptBuilder(C ControllerFileUploadEncrypt, r *mux.Router) httpController.UploadEncryptController {
	return httpController.UploadEncryptController{
		NewFileUploaderEncryptSession: httpController.NewFileUploaderEncryptSession{
			ReadSession: C.Sessions,
			Jwt:         C.Token,
			Rf:          C.Token,
		},
		UploaderEncryptApplication: &C.App,
		Rout:                       r,
		UrlUploadData:              httpController.GetUploadEncryptData,
	}
}

type ControllerLoginBuilder struct {
	Sessions httpController.Session
	Parser   ParsersCollector
	App      Application.LoginApplication
}

func GetControllerLoginBuilder(C ControllerLoginBuilder) httpController.LoginController {
	return httpController.LoginController{
		Sess:   C.Sessions,
		Parses: C.Parser.Parser,
		App:    C.App,
	}
}

type RegisterController struct {
	Sessions httpController.Session
	App      Application.NewRegisterApplication
	Parser   ParsersCollector
}

func GetControllerRegisterBuilder(c RegisterController) httpController.RegisterController {

	return httpController.RegisterController{
		RegisterNet: httpController.RegisterNet{},
		Session:     c.Sessions,
		Decoder:     c.Parser.Parser,
		App:         &c.App,
	}
}

func GetControllerUrlBuildBuilder() httpController.BuildUrlController {
	return httpController.BuildUrlController{
		Url: httpController.UrlBuilder,
	}
}

func GetCheckUsers() httpController.CheckUserAuthController {
	return httpController.CheckUserAuthController{
		CheckUserAuthNetWork: httpController.CheckUserAuthNetWork{},
	}
}
