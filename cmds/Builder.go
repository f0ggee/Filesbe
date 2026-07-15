package cmds

import (
	"Kaban/internal/Deliver"
	"Kaban/internal/DomainLevel"
	"Kaban/internal/InfrastructureLayer/AuthTokensManage"
	"Kaban/internal/InfrastructureLayer/DeliverPackages/RepoDownloadEncryptRepo"
	"Kaban/internal/InfrastructureLayer/DeliverPackages/RepoDownloadNoEncrypt"
	"Kaban/internal/InfrastructureLayer/DeliverPackages/RepoLoginRealizations"
	"Kaban/internal/InfrastructureLayer/DeliverPackages/RepoRegisterRepository"
	"Kaban/internal/InfrastructureLayer/DeliverPackages/RepoSessionHandle"
	"Kaban/internal/InfrastructureLayer/DeliverPackages/RepoUsersCheckAuth"
	"Kaban/internal/InfrastructureLayer/DeliverPackages/RepofileUploaderEncryptRepo"
	"Kaban/internal/InfrastructureLayer/DeliverPackages/RepofileUploaderNoEncryptRepo"
	"Kaban/internal/InfrastructureLayer/DeliverPackages/RepourlBuilder"
	"Kaban/internal/InfrastructureLayer/RepoParsers"
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

func GetDownloadEncryptApplicationBuilder(a *RedisCollector, s *S3Collector, z *KeysCollector, c *CollectorCrypto, f *FileControlCollector) *Application.NewDownloadEncrypt {

	delivery := Application.NewDownloadEncryptDelivery{
		ReaderRedis:  &a.Read,
		DownloadS3:   s.S3Download,
		DeleterRedis: &a.Delete,
		DeleterS3:    s.Deleter,
	}

	crypto := Application.NewDownloadEncryptCrypto{
		Decrypt:       c.Decrypt,
		EncrypterKeys: z.Keys,
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
func GetUploadApplication(c *CollectorCrypto, f *FileControlCollector, z *RepoParsersCollector, s3 *S3Collector, r *RedisCollector) *Application.NewFileUploader {

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
	f          *FileControlCollector
	s3         *S3Collector
	z          *RepoParsersCollector
	c          *CollectorCrypto
	serverKeys *DomainLevel.NewServerKeys
	r          *RedisCollector
}

func GetUploadEncryptApplicationBuilder(d UploadEncryptBuilderIncomeData) *Application.NewUploadEncrypt {

	data := Application.NewUploadEncryptDataManage{
		FileManger: d.f.FileSettings,
		Encode:     d.z.Encode,
	}

	crypto := Application.NewUploadEncryptCrypto{
		Generate:   d.c.Generate,
		ServerKeys: *d.serverKeys,
		Encrypt:    d.c.Encrypt,
	}

	delivery := Application.NewUploadEncryptDelivery{
		UploaderS3:   d.s3.Uploader,
		RedisWriter:  &d.r.Write,
		RedisChecker: &d.r.Check,
		RedisDeleter: &d.r.Delete,
		DeleterS3:    d.s3.Deleter,
	}
	return Application.GetNewUploadEncrypt(data, crypto, delivery)
}

// Delivery builders

type RegisterControllerBuilder struct {
	Answ    RepoDownloadNoEncrypt.Answer
	UrlWork RepoDownloadNoEncrypt.UrlWork
	s       *Application.NewDownload
}

func GetControllerDownloadBuilder(d *RegisterControllerBuilder) *Deliver.NewDownloadWithNotEncrypt {

	answ := Deliver.AnswerDownloadNoEncrypt{
		Answ: d.Answ,
	}

	url := Deliver.NewDownloadWithNotEncryptUrlBuilder{
		UrlWork: d.UrlWork,
	}

	NetWork := Deliver.NetworkDownloadNoEncrypt{}
	App := Deliver.NewDownloadWithNotEncryptApplication{
		NewDownload: *d.s,
	}

	return Deliver.GetNewDownloadWithNotEncrypt(answ, url, NetWork, App)
}

type EncryptDownloadControllerIncomeData struct {
	Answ     RepoDownloadEncryptRepo.AnswersDownloadEncrypt
	UrlBuild RepoDownloadEncryptRepo.UrlWork
	d        *Application.NewDownloadEncrypt
}

func GetEncryptDownloadControllerBuilder(d EncryptDownloadControllerIncomeData) *Deliver.NewDownloadEncrypt {
	answ := Deliver.AnswerDownloadEncrypt{
		S: d.Answ,
	}
	url := Deliver.UrlBuilderDownloadEncrypt{
		UrlBuild: d.UrlBuild,
	}
	net := Deliver.NetworkDownloadEncrypt{}
	App := Deliver.NewDownloadWithEncryptApplication{
		NewDownloadEncrypt: *d.d,
	}
	return Deliver.GetNewDownloadEncrypt(answ, url, net, App)
}

type UploaderBuilderIncomeData struct {
	R           *mux.Router
	ReadSession RepoSessionHandle.Session
	AuthCheck   AuthTokensManage.AuthCheck
	Answers     RepofileUploaderEncryptRepo.AnswersUploadEncrypt
	Build       RepofileUploaderEncryptRepo.UrlUploadEncrypt
	A           Application.NewUploadEncrypt
}

func GetUploaderEncrypterControllerBuilder(d *UploaderBuilderIncomeData) *Deliver.NewFileUploaderEncrypt {
	net := Deliver.FileUploaderEncryptNetwork{}
	sess := Deliver.NewFileUploaderEncryptSession{
		ReadSession: d.ReadSession,
		AuthCheck:   d.AuthCheck,
	}
	details := Deliver.NewFileUploaderEncryptDetails{
		Answers: d.Answers,
		Build:   d.Build,
		Rout:    d.R,
	}

	App := Deliver.NewFileUploaderEncryptApplication{
		NewUploadEncrypt: d.A,
	}
	return Deliver.GetNewFileUploaderEncrypt(net, sess, details, App)
}

type UploaderEncryptBuilderIncomeData struct {
	S       RepofileUploaderNoEncryptRepo.Answers
	Builder RepofileUploaderNoEncryptRepo.NewUploaderNoEncrypt
	Session RepoSessionHandle.Session
	Auth    AuthTokensManage.AuthCheck
}

func GetFileUploaderControllerBuilder(data UploaderEncryptBuilderIncomeData) *Deliver.NewFileUploader {
	net := Deliver.NewFileUploaderNet{}
	details := Deliver.NewFileUploaderWorkDetails{
		S:       nil,
		Builder: RepofileUploaderNoEncryptRepo.NewUploaderNoEncrypt{},
	}
	session := Deliver.NewFileUploaderSessions{
		Session: nil,
		Auth:    nil,
	}

	app := Deliver.NewFileUploaderApp{
		NewFileUploader: Application.NewFileUploader{},
	}
	return Deliver.GetNewFileUploader(net, details, session, app)
}

type NewLoginIncomeData struct {
	S      *RepoLoginRealizations.LoginAnswers
	Sess   RepoSessionHandle.Session
	Parses *RepoParsers.Parsing
	app    *Application.NewLogin
}

func GetNewLoginBuilder(data NewLoginIncomeData) *Deliver.NewLogin {
	net := Deliver.LoginNet{}
	depends := Deliver.LoginDepends{
		S:    data.S,
		Sess: data.Sess,
	}
	parse := Deliver.ParseLogin{
		Parses: data.Parses,
	}
	app := Deliver.LoginApplication{
		NewLogin: *data.app,
	}
	return Deliver.GetNewLogin(net, depends, parse, app)
}

type CheckUserBuilderIncomeData struct {
	answers *RepoUsersCheckAuth.SetUsersChecker
	Session RepoSessionHandle.Session
	Auth    AuthTokensManage.NewAuthChecker
}

func GetCheckUserBuilder(data CheckUserBuilderIncomeData) *Deliver.NewCheckUserAuth {
	net := Deliver.GetNewCheckUserAuthNetWork(nil, nil)
	auth := Deliver.GetNewCheckUserAuthWorkDetails(data.answers)
	session := Deliver.GetNewCheckUserAuthSessions(data.Auth, data.Session)
	return Deliver.GetNewCheckUserAuth(*net, *auth, *session)
}

type RegisterBuilderIncomeData struct {
	Answ    RepoRegisterRepository.RegisterAnswers
	Session RepoSessionHandle.Session
	D       RepoParsers.Decode
	App     Application.NewRegisterApplication
}

func GetRegisterBuilder(data RegisterBuilderIncomeData) *Deliver.NewRegister {

	Details := Deliver.NewRegisterDetails{
		Answ:    data.Answ,
		Session: data.Session,
		D:       data.D,
	}

	app := Deliver.NewRegisterApp{
		App: data.App,
	}

	net := Deliver.RegisterNet{}
	return Deliver.GetNewRegister(net, Details, app)
}

func GetUrlUploaderBuilder(Url RepourlBuilder.UrlBuilderAnswer) *Deliver.NewBuildUrl {
	url := Deliver.UrlSettings{
		Url: Url,
	}
	net := Deliver.UrlNetwork{}
	return Deliver.GetNewBuildUrl(url, net)
}
