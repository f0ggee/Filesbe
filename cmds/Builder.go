package cmds

import (
	"Kaban/internal/Deliver"
	"Kaban/internal/DomainLevel"
	"Kaban/internal/InfrastructureLayer/DeliverPackages/RepoDownloadNoEncrypt"
	"Kaban/internal/Service/Application"
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
func GetControllerDownloadBuilder(DeliverPackagesCollector *DeliverPackagesCollector) *Deliver.NewDownloadWithNotEncrypt {
	return Deliver.GetNewDownloadWithNotEncrypt(Deliver.AnswerDownloadNoEncrypt{Answ: DeliverPackagesCollector.DownloadCollector.Answ}, Deliver.UrlBuilderDownloadNoEncrypt{UrlWork: RepoDownloadNoEncrypt.NewRepoDownloadNoEncrypt(DeliverPackagesCollector.UrlBuilderCollector.Answ)}, Deliver.NetworkDownloadNoEncrypt{}, Deliver.NewDownloadWithNotEncryptApplication{})
}
