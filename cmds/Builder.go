package cmds

import (
	"Kaban/internal/Deliver"
	"Kaban/internal/InfrastructureLayer/DeliverPackages/RepoDownloadNoEncrypt"
	"Kaban/internal/Service/Application"
)

// application builders
func GetS3Builder(s *S3Collector) *Application.S3Controlling {
	return Application.GetNewS3Controlling(s.Deleter, s.Uploader, s.S3Download)
}

func GetTransfersBuilder(s *FileControlCollector) *Application.Transfers {
	return Application.GetNewTransfers(s.Transfer)
}

func GetEncrypterKeysBuilder(s *EncrypterKeysCollector) *Application.EncrypterKeys {
	return Application.GetEncrypterKeys(s.Keys)
}
func GetCryptoBuilder(s *CollectorCrypto) *Application.GetCrypto {
	return Application.GetNewCrypto(s.Validate, s.Encrypt, s.Decrypt, s.Generate, s.Keys)
}

func GetFileManagerBuilder(s *FileControlCollector) *Application.NewFileManager {
	return Application.GetNewFileManager(s.FileSettings)
}

func GetNewControlKeysBuilder(s *CollectorCrypto) *Application.GetControlKeys {
	return Application.GetNewControlKeys(s.Keys)
}
func GetAuthTokensBuilder(s *CollectorAuthTokensManage) *Application.AuthTokens {
	return Application.GetNewAuthTokens(s.Creating, s.Validate)
}
func GetHandlerGrpcBuilder(s *GrpcCollector) *Application.HandlerGrpc {
	return Application.GetNewHandlerGrpc(s.Sender, s.Checking)
}

func GetDatabaseBuilder(s *CollectorDatabaseManage) *Application.DatabaseControlling {
	return Application.GetDatabaseControlling(&s.Writer, s.Reader, &s.Checker)
}
func GetGetRedisBuilder(s *RedisCollector) *Application.RedisControlling {
	return Application.GetRedisControlling(&s.Write, &s.Read, &s.Delete, &s.Check)
}
func GetParserBuilder(s *RepoParsersCollector) *Application.Parser {
	return Application.GetParser(s.Decode, s.Encode)
}

func GetDownloadApplicationBuilder(redis *Application.RedisControlling, s3 *Application.S3Controlling, file *Application.NewFileManager, transfers *Application.Transfers) *Application.NewDownloadNotEncrypt {

	return Application.GetNewDownloadNotEncrypt(*redis, *s3, *file, *transfers, Application.DownloadNotEncryptNetwork{})
}

func GetNewDownloadEncryptBuilder(redis *Application.RedisControlling, crypto *Application.GetCrypto, s3 *Application.S3Controlling, file *Application.NewFileManager, transfers *Application.Transfers, keys *Application.EncrypterKeys) *Application.NewDownloadEncrypt {
	return Application.GetNewNewDownloadEncrypt(*redis, *crypto, *s3, *file, *transfers, *keys)
}
func GetLoginBuilder(Crypto *Application.GetCrypto, database *Application.DatabaseControlling, auth *Application.AuthTokens) *Application.NewLogin {
	return Application.GetNewNewLogin(*database, *Crypto, *auth)
}

func GetRegisterBuilder(database *Application.DatabaseControlling, Crypto *Application.GetCrypto, auth *Application.AuthTokens) *Application.NewRegisterApplication {
	return Application.GetNewNewRegisterApplication(*database, *Crypto, *auth)
}

func GetUploadBuilder(Crypto *Application.GetCrypto, file *Application.NewFileManager, s3 *Application.S3Controlling, Parse *Application.Parser, redis *Application.RedisControlling) *Application.NewFileUploader {
	return Application.GetNewFileUploader(*Crypto, *file, *s3, *Parse, *redis)
}
func GetUploadEncryptBuilder(file *Application.NewFileManager, Crypto *Application.GetCrypto, Parse *Application.Parser, keys *Application.GetControlKeys, s3 *Application.S3Controlling, redis *Application.RedisControlling) *Application.NewUploadEncrypt {
	return Application.GetNewUploadEncrypt(*file, *Crypto, *Parse, *keys, *s3, *redis)
}

func GetControllerDownloadBuilder(DeliverPackagesCollector *deliverPackagesCollector) *Deliver.NewDownloadWithNotEncrypt {
	return Deliver.GetNewDownloadWithNotEncrypt(Deliver.AnswerDownloadNoEncrypt{Answ: DeliverPackagesCollector.DownloadCollector.Answ}, Deliver.UrlBuilderDownloadNoEncrypt{UrlWork: RepoDownloadNoEncrypt.NewRepoDownloadNoEncrypt(DeliverPackagesCollector.UrlBuilderCollector.Answ)}, Deliver.NetworkDownloadNoEncrypt{}, Deliver.NewDownloadWithNotEncryptApplication{})
}
