package Deliver

import (
	"Kaban/internal/DomainLevel"
	"Kaban/internal/InfrastructureLayer/DeliverPackages/RepoDownloadNoEncrypt"
	"Kaban/internal/Service/Application"

	"net/http"
)

type AnswerDownloadNoEncrypt struct {
	Answ RepoDownloadNoEncrypt.Answer
}

type NewDownloadWithNotEncryptUrlBuilder struct {
	UrlWork RepoDownloadNoEncrypt.UrlWork
}
type NetworkDownloadNoEncrypt struct {
	W http.ResponseWriter
	R *http.Request
}

type NewDownloadWithNotEncryptApplication struct {
	Application.NewDownload
}

type NewDownloadWithNotEncrypt struct {
	AnswerDownloadNoEncrypt
	NewDownloadWithNotEncryptUrlBuilder
	NetworkDownloadNoEncrypt
	NewDownloadWithNotEncryptApplication
}

func GetNewDownloadWithNotEncrypt(answerDownloadNoEncrypt AnswerDownloadNoEncrypt, urlBuilderDownloadNoEncrypt NewDownloadWithNotEncryptUrlBuilder, networkDownloadNoEncrypt NetworkDownloadNoEncrypt, newDownloadWithNotEncryptApplication NewDownloadWithNotEncryptApplication) *NewDownloadWithNotEncrypt {
	return &NewDownloadWithNotEncrypt{AnswerDownloadNoEncrypt: answerDownloadNoEncrypt, NewDownloadWithNotEncryptUrlBuilder: urlBuilderDownloadNoEncrypt, NetworkDownloadNoEncrypt: networkDownloadNoEncrypt, NewDownloadWithNotEncryptApplication: newDownloadWithNotEncryptApplication}
}

func (d NewDownloadWithNotEncrypt) DownloadWithNotEncrypt() {

	name := d.UrlWork.GetData(d.R)
	if name == "" {
		d.Answ.SetBadAnswer(RepoDownloadNoEncrypt.DownloadNoEncryptIncomingData{
			W:               d.W,
			Err:             DomainLevel.ErrorCantGetFileName,
			StatusOperation: DomainLevel.Break,
		})
		return
	}

	err := d.Download(name, d.R.Context())
	if err != nil {
		d.Answ.SetBadAnswer(RepoDownloadNoEncrypt.DownloadNoEncryptIncomingData{
			W:               d.W,
			Err:             err.Error(),
			StatusOperation: "",
		})
		return
	}
	return

}
