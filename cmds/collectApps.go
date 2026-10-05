package cmds

import "Kaban/internal/Service/Application"

type AppBuilders struct {
	download        Application.NewDownload
	encryptDownload Application.NewDownloadEncrypt
	login           Application.NewLogin
	register        Application.NewRegisterApplication
	upload          Application.NewUpload
	encryptUpload   Application.NewUploadEncrypt
}
type CollectorsApp struct {
	FileTransferringCollector FileTransferringCollector
	RedisCollector            RedisCollector
	ParserCollector           ParsersCollector
	CryptoCollector           CryptoCollector
	SessionKeys               SessionKeysCollector
	DatabaseCollector         DatabaseCollector
	AuthTokensCollector       TokensAuthCollector
}

func CollectApplicationBuilders(c CollectorsApp) AppBuilders {
	downloadApplication := GetDownloadApplicationBuilder(DownloadBuilderApplication{
		Trans: c.FileTransferringCollector,
		Red:   c.RedisCollector,
		Parse: c.ParserCollector,
	})
	downloadEncryptApplication := GetDownloadEncryptApplicationBuilder(DownloadEncryptApplication{
		Red:          c.RedisCollector,
		Crypt:        c.CryptoCollector,
		Keys:         c.SessionKeys,
		Transferring: c.FileTransferringCollector,
	})
	loginApplication := GetLoginApplicationBuilder(LoginApplication{
		DB:     c.DatabaseCollector,
		Crypt:  c.CryptoCollector,
		Tokens: c.AuthTokensCollector,
	})
	registerApplication := GetRegisterBuilder(RegisterApplication{
		DB:     c.DatabaseCollector,
		Tokens: c.AuthTokensCollector,
		Crypt:  c.CryptoCollector,
	})
	uploadApplication := GetUploadBuilder(UploadApplication{
		Crypt: c.CryptoCollector,
		Parse: c.ParserCollector,
		S3:    c.FileTransferringCollector.S3FileTransferringCollector,
		Red:   c.RedisCollector,
	})
	uploadEncryptApplication := GetUploadEncryptBuilder(UploadEncrypt{
		Parse:        c.ParserCollector,
		Crypt:        c.CryptoCollector,
		Transferring: c.FileTransferringCollector,
		Red:          c.RedisCollector,
	})
	return AppBuilders{
		download:        downloadApplication,
		encryptDownload: downloadEncryptApplication,
		login:           loginApplication,
		register:        registerApplication,
		upload:          uploadApplication,
		encryptUpload:   uploadEncryptApplication,
	}
}
