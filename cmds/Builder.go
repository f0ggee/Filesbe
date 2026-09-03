package cmds

import (
	"Kaban/internal/Deliver/httpController"
	"Kaban/internal/Deliver/httpController/DeliverPackages/RepoDownloadEncryptRepo"
	"Kaban/internal/Deliver/httpController/DeliverPackages/RepoDownloadNoEncrypt"
	"Kaban/internal/Deliver/httpController/DeliverPackages/RepoLoginRealizations"
	"Kaban/internal/Deliver/httpController/DeliverPackages/RepoRegisterRepository"
	"Kaban/internal/Deliver/httpController/DeliverPackages/RepoUsersCheckAuth"
	"Kaban/internal/Deliver/httpController/DeliverPackages/RepofileUploaderEncryptRepo"
	"Kaban/internal/Deliver/httpController/DeliverPackages/RepofileUploaderNoEncryptRepo"
	"Kaban/internal/Deliver/httpController/DeliverPackages/RepourlBuilder"
	"Kaban/internal/DomainLevel"
	"Kaban/internal/InfrastructureLayer/AuthTokensManage"
	"Kaban/internal/InfrastructureLayer/ProtocolManage"
	"Kaban/internal/InfrastructureLayer/RepoEncrypterKeys"
	"Kaban/internal/InfrastructureLayer/RepoParsers"
	"Kaban/internal/InfrastructureLayer/RepoSessionHandle"
	"Kaban/internal/Service/Application"

	"github.com/gorilla/mux"
)

// application builders
func GetDownloadApplicationBuilder(s *S3Collector, a *RedisCollector, f *FileControlCollector) *Application.NewDownload {
	delivery := Application.DownloadDelivery{
		Reader:     &a.Read,
		DeleterS3:  s.Deleter,
		S3Download: s.S3Download,
	}
	fileControl := Application.DownloadFileControl{
		Transfer:     f.Transfer,
		FileManaging: f.FileSettings,
	}
	return Application.GetNewDownload(delivery, fileControl, Application.DownloadNetwork{})
}

func GetDownloadEncryptApplicationBuilder(a *RedisCollector, s *S3Collector, z *RepoEncrypterKeys.Keys, c *CollectorCrypto, f *FileControlCollector) *Application.NewDownloadEncrypt {

	delivery := Application.NewDownloadEncryptDelivery{
		ReaderRedis:  &a.Read,
		DownloadS3:   s.S3Download,
		DeleterRedis: &a.Delete,
		DeleterS3:    s.Deleter,
	}

	crypto := Application.NewDownloadEncryptCrypto{
		Decrypt:       c.Decrypt,
		EncrypterKeys: *z,
	}
	file := Application.NewDownloadEncryptFileControl{
		Transfer:     f.Transfer,
		FileManaging: f.FileSettings,
	}
	return Application.GetNewDownloadEncrypt(delivery, crypto, file)
}

func GetLoginApplicationBuilder(c *CollectorCrypto, d *CollectorDatabaseManage, x *CollectorAuthTokensManage) *Application.NewLogin {

	data := Application.NewLoginData{
		ReaderDatabase: d.Reader,
	}
	crypto := Application.NewLoginCrypto{
		Validate: c.Validate,
	}

	auth := Application.NewLoginAuth{
		GeneratingTokens: x.Creating,
	}

	return Application.GetNewLogin(data, crypto, auth)
}

func GetRegisterApplicationBuilder(d *CollectorDatabaseManage, x *CollectorAuthTokensManage, c *CollectorCrypto) *Application.NewRegisterApplication {
	data := Application.NewRegisterDataMange{
		CheckingDb:      &d.Checker,
		WriterDb:        &d.Writer,
		GeneratorTokens: x.Creating,
	}
	crypto := Application.NewRegisterCrypto{
		Generator: c.Generate,
	}
	return Application.GetNewRegisterApplication(data, crypto)
}
func GetUploadApplicationBuilder(c *CollectorCrypto, f *FileControlCollector, z *RepoParsersCollector, s3 *S3Collector, r *RedisCollector) *Application.NewUpload {

	crypto := Application.NewFileUploaderCrypto{
		Generator: c.Generate,
	}

	data := Application.NewFileUploaderDataMange{
		FileSettings: f.FileSettings,
		Encode:       z.Encode,
	}
	delivery := Application.NewFileUploaderDelivery{
		UploadS3:   s3.Uploader,
		WriteRedis: &r.Write,
	}
	return Application.GetNewFileUploader(crypto, data, delivery)
}

type UploadEncryptBuilderIncomeData struct {
	F          *FileControlCollector
	S3         *S3Collector
	Z          *RepoParsersCollector
	C          *CollectorCrypto
	ServerKeys *DomainLevel.NewServerKeys
	R          *RedisCollector
}

func GetUploadEncryptApplicationBuilder(d UploadEncryptBuilderIncomeData) *Application.NewUploadEncrypt {

	data := Application.NewUploadEncryptDataManage{
		FileManger: d.F.FileSettings,
		Encode:     d.Z.Encode,
	}

	crypto := Application.NewUploadEncryptCrypto{
		Generate:   d.C.Generate,
		ServerKeys: *d.ServerKeys,
		Encrypt:    d.C.Encrypt,
	}

	delivery := Application.NewUploadEncryptDelivery{
		UploaderS3:   d.S3.Uploader,
		RedisWriter:  &d.R.Write,
		RedisChecker: &d.R.Check,
		RedisDeleter: &d.R.Delete,
		DeleterS3:    d.S3.Deleter,
	}
	return Application.GetNewUploadEncrypt(data, crypto, delivery)
}

// Delivery builders

type RegisterControllerBuilderIncomeData struct {
	Answ    RepoDownloadNoEncrypt.Answer
	UrlWork RepoDownloadNoEncrypt.UrlWork
	App     *Application.NewDownload
}

func GetControllerDownloadBuilder(d *RegisterControllerBuilderIncomeData) *httpController.DownloadNew {

	answ := httpController.AnswerDownloadNoEncrypt{
		Answ: d.Answ,
	}

	url := httpController.NewDownloadWithNotEncryptUrlBuilder{
		UrlWork: d.UrlWork,
	}

	NetWork := httpController.DownloadNetwork{}
	App := httpController.DownloadApp{
		NewDownload: *d.App,
	}

	return httpController.GetNewDownloadWithNotEncrypt(answ, url, NetWork, App)
}

type EncryptDownloadControllerIncomeData struct {
	Answ     RepoDownloadEncryptRepo.AnswersDownloadEncrypt
	UrlBuild RepoDownloadEncryptRepo.UrlWork
	App      *Application.NewDownloadEncrypt
}

func GetEncryptDownloadControllerBuilder(d EncryptDownloadControllerIncomeData) *httpController.NewDownloadEncrypt {
	answ := httpController.AnswerDownloadEncrypt{
		S: d.Answ,
	}
	url := httpController.UrlBuilderDownloadEncrypt{
		UrlBuild: d.UrlBuild,
	}
	net := httpController.NetworkDownloadEncrypt{}
	App := httpController.NewDownloadWithEncryptApplication{
		NewDownloadEncrypt: *d.App,
	}
	return httpController.GetNewDownloadEncrypt(answ, url, net, App)
}

type UploaderBuilderIncomeData struct {
	R           *mux.Router
	ReadSession RepoSessionHandle.Session
	AuthCheck   AuthTokensManage.AuthCheck
	Answers     RepofileUploaderEncryptRepo.AnswersUploadEncrypt
	Build       RepofileUploaderEncryptRepo.UrlUploadEncrypt
	A           Application.NewUploadEncrypt
}

func GetUploaderEncrypterControllerBuilder(d *UploaderBuilderIncomeData) *httpController.NewUploaderEncrypt {
	net := httpController.NewFileUploaderEncryptNetwork{}
	sess := httpController.NewFileUploaderEncryptSession{
		ReadSession: d.ReadSession,
		AuthCheck:   d.AuthCheck,
	}
	details := httpController.NewFileUploaderEncryptDetails{
		Answers: d.Answers,
		Build:   d.Build,
		Rout:    d.R,
	}

	App := httpController.NewFileUploaderEncryptApplication{
		NewUploadEncrypt: d.A,
	}
	return httpController.GetNewFileUploaderEncrypt(net, sess, details, App)
}

type UploaderEncryptBuilderIncomeData struct {
	S       RepofileUploaderNoEncryptRepo.Answers
	Builder RepofileUploaderNoEncryptRepo.NewUploaderNoEncrypt
	Session RepoSessionHandle.Session
	Auth    AuthTokensManage.AuthCheck
	App     Application.NewUpload
}

func GetUploaderControllerBuilder(data UploaderEncryptBuilderIncomeData) *httpController.NewFileUploader {
	net := httpController.FileUploaderNet{}
	details := httpController.NewFileUploaderWorkDetails{
		S:       data.S,
		Builder: data.Builder,
	}
	session := httpController.FileUploaderSessions{
		Session: data.Session,
		Auth:    data.Auth,
	}

	app := httpController.NewFileUploaderApp{
		NewUpload: data.App,
	}
	return httpController.GetNewFileUploader(net, details, session, app)
}

type NewLoginIncomeData struct {
	S      *RepoLoginRealizations.LoginAnswers
	Sess   RepoSessionHandle.Session
	Parses RepoParsers.Decode
	App    *Application.NewLogin
}

func GetControllerLoginBuilder(data NewLoginIncomeData) *httpController.NewLoginController {
	net := httpController.LoginNet{}
	depends := httpController.LoginDepends{
		S:    data.S,
		Sess: data.Sess,
	}
	parse := httpController.ParseLogin{
		Parses: data.Parses,
	}
	app := httpController.LoginApplication{
		NewLogin: *data.App,
	}
	return httpController.GetNewLogin(net, depends, parse, app)
}

type CheckUserBuilderIncomeData struct {
	answers *RepoUsersCheckAuth.SetUsersChecker
	Session RepoSessionHandle.Session
	Auth    AuthTokensManage.AuthCheck
}

func GetControllerCheckAuthBuilder(data CheckUserBuilderIncomeData) *httpController.CheckUserAuth {
	net := httpController.GetNewCheckUserAuthNetWork(nil, nil)
	auth := httpController.GetNewCheckUserAuthWorkDetails(data.answers)
	session := httpController.GetNewCheckUserAuthSessions(data.Auth, data.Session)
	return httpController.GetNewCheckUserAuth(*net, *auth, *session)
}

type RegisterBuilderIncomeData struct {
	Answ    RepoRegisterRepository.RegisterAnswers
	Session RepoSessionHandle.Session
	D       RepoParsers.Decode
	App     Application.NewRegisterApplication
}

func GetControllerRegisterBuilder(data RegisterBuilderIncomeData) *httpController.NewRegister {

	Details := httpController.NewRegisterDetails{
		Answ:    data.Answ,
		Session: data.Session,
		D:       data.D,
	}

	app := httpController.NewRegisterApp{
		App: data.App,
	}

	net := httpController.RegisterNet{}
	return httpController.GetNewRegister(net, Details, app)
}

func GetControllerUrlUploaderBuilder(Url RepourlBuilder.UrlBuilderAnswer) *httpController.NewBuildUrl {
	url := httpController.UrlSettings{
		Url: Url,
	}
	net := httpController.UrlNetwork{}
	return httpController.GetNewBuildUrl(url, net)
}

type ProtocolManageBuilder struct {
	EncrypterKeys RepoEncrypterKeys.Keys
	ServerKeys    DomainLevel.NewServerKeys
	C             *CollectorCrypto
	Parser        *RepoParsersCollector
	GrpcConn      GrpcCollector
	Red           RedisCollector
}

type ProtocolExchanges struct {
	Start      ProtocolManage.NewExchangeInitializer
	Processing ProtocolManage.Exchanger
}

func GetProtocolManageBuilder(data ProtocolManageBuilder) *ProtocolExchanges {
	key := ProtocolManage.NewExchangeInitializerKey{
		Keys:       data.EncrypterKeys,
		ServerKeys: data.ServerKeys,
	}

	crypto := ProtocolManage.NewExchangeInitializerCrypto{
		CryptoGenerating: data.C.Generate,
		CryptoEncrypt:    data.C.Encrypt,
		CryptoDecrypt:    data.C.Decrypt,
		CryptoValidate:   data.C.Validate,
	}
	Parser := ProtocolManage.NewExchangeInitializerParsers{
		Encode: data.Parser.Encode,
		Decode: data.Parser.Decode,
	}

	Del := ProtocolManage.NewExchangeInitializerDeliver{
		Grcp: data.GrpcConn.Sender,
	}
	delProcess := ProtocolManage.NewExchangerDeliver{
		Redis: &data.Red.Read,
	}
	cryptoProcess := ProtocolManage.NewExchangerCrypto{
		Decrypter:  data.C.Decrypt,
		Validation: data.C.Validate,
	}
	keyProcess := ProtocolManage.NewExchangerKeys{
		ServerKeys:    data.ServerKeys,
		EncrypterKeys: data.EncrypterKeys,
	}
	parser := ProtocolManage.NewExchangerParsers{
		Decoder: data.Parser.Decode,
	}
	return &ProtocolExchanges{
		Start:      *ProtocolManage.GetNewExchangeInitializer(key, crypto, Parser, Del),
		Processing: ProtocolManage.GetNewExchanger(delProcess, parser, cryptoProcess, keyProcess),
	}
}
