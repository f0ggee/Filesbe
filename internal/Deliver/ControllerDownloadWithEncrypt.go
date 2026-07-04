package Deliver

import (
	"Kaban/internal/DomainLevel"
	"Kaban/internal/InfrastructureLayer/DeliverPackages/RepoDownloadEncryptRepo"
	"Kaban/internal/Service/Application"
	"net/http"
)

type AnswerDownloadEncrypt struct {
	S RepoDownloadEncryptRepo.NewFileDownloadEncrypt
}
type UrlBuilderDownloadEncrypt struct {
	UrlBuild RepoDownloadEncryptRepo.NewFileDownloadEncrypt
}

type NetworkDownloadEncrypt struct {
	W http.ResponseWriter
	R *http.Request
}
type NewDownloadEncrypt struct {
	Answer AnswerDownloadEncrypt
	Url    UrlBuilderDownloadEncrypt
	Net    NetworkDownloadEncrypt
}

func GetNewNewDownloadEncrypt(answer AnswerDownloadEncrypt, url UrlBuilderDownloadEncrypt, net NetworkDownloadEncrypt) *NewDownloadEncrypt {
	return &NewDownloadEncrypt{Answer: answer, Url: url, Net: net}
}

func (d NewDownloadEncrypt) DownloadWithEncrypt(s *Application.HandlerPackCollect) {

	fileName := d.Url.UrlBuild.GetDataRequest(d.Net.R)
	if fileName == "" {

		d.Answer.S.SetBadAnswers(RepoDownloadEncryptRepo.IncomingDataAnswer{
			W:     d.Net.W,
			Error: DomainLevel.ErrorCantGetFileName,
		})
		return
	}
	err := s.DownloadEncrypt(d.Net.W, d.Net.R.Context(), fileName)
	if err != nil {
		d.Answer.S.SetBadAnswers(RepoDownloadEncryptRepo.IncomingDataAnswer{
			W:     d.Net.W,
			Error: err.Error(),
		})
		return
	}
	return

}
