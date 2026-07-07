package Deliver

import (
	"Kaban/internal/DomainLevel"
	"Kaban/internal/InfrastructureLayer/DeliverPackages/RepoDownloadNoEncrypt"
	"Kaban/internal/Service/Application"

	"net/http"
)

type AnswerDownloadNoEncrypt struct {
	Answ RepoDownloadNoEncrypt.NewRepoDownloadNoEncrypt
}

type UrlBuilderDownloadNoEncrypt struct {
	UrlWork RepoDownloadNoEncrypt.NewRepoDownloadNoEncrypt
}
type NetworkDownloadNoEncrypt struct {
	W http.ResponseWriter
	R *http.Request
}

type NewDownloadWithNotEncryptApplication struct {
	Application.NewDownloadNotEncrypt
}

type NewDownloadWithNotEncrypt struct {
	AnswerDownloadNoEncrypt
	UrlBuilderDownloadNoEncrypt
	NetworkDownloadNoEncrypt
	NewDownloadWithNotEncryptApplication
}

func GetNewDownloadWithNotEncrypt(answerDownloadNoEncrypt AnswerDownloadNoEncrypt, urlBuilderDownloadNoEncrypt UrlBuilderDownloadNoEncrypt, networkDownloadNoEncrypt NetworkDownloadNoEncrypt, newDownloadWithNotEncryptApplication NewDownloadWithNotEncryptApplication) *NewDownloadWithNotEncrypt {
	return &NewDownloadWithNotEncrypt{AnswerDownloadNoEncrypt: answerDownloadNoEncrypt, UrlBuilderDownloadNoEncrypt: urlBuilderDownloadNoEncrypt, NetworkDownloadNoEncrypt: networkDownloadNoEncrypt, NewDownloadWithNotEncryptApplication: newDownloadWithNotEncryptApplication}
}

func (d NewDownloadWithNotEncrypt) DownloadWithNotEncrypt(w http.ResponseWriter, r *http.Request) {

	name := d.UrlWork.GetData(r)
	if name == "" {
		d.Answ.SetBadAnswer(RepoDownloadNoEncrypt.DownloadNoEncryptIncomingData{
			W:               d.W,
			Err:             DomainLevel.ErrorCantGetFileName,
			StatusOperation: DomainLevel.Break,
		})
		return
	}

	err := d.DownloadWithNonEncrypt(name, d.R.Context())
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
