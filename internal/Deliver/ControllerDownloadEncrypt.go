package Deliver

import (
	"Kaban/internal/DomainLevel"
	"Kaban/internal/InfrastructureLayer/DeliverPackages/RepoDownloadEncryptRepo"
	"Kaban/internal/Service/Application"
	"net/http"
)

type AnswerDownloadEncrypt struct {
	S RepoDownloadEncryptRepo.AnswersDownloadEncrypt
}
type UrlBuilderDownloadEncrypt struct {
	UrlBuild RepoDownloadEncryptRepo.UrlWork
}

type NetworkDownloadEncrypt struct {
	W http.ResponseWriter
	R *http.Request
}
type NewDownloadWithEncryptApplication struct {
	Application.NewDownloadEncrypt
}
type NewDownloadEncrypt struct {
	Answer AnswerDownloadEncrypt
	Url    UrlBuilderDownloadEncrypt
	Net    NetworkDownloadEncrypt
	NewDownloadWithEncryptApplication
}

func GetNewDownloadEncrypt(answer AnswerDownloadEncrypt, url UrlBuilderDownloadEncrypt, net NetworkDownloadEncrypt, newDownloadWithEncryptApplication NewDownloadWithEncryptApplication) *NewDownloadEncrypt {
	return &NewDownloadEncrypt{Answer: answer, Url: url, Net: net, NewDownloadWithEncryptApplication: newDownloadWithEncryptApplication}
}

func (d NewDownloadEncrypt) DownloadWithEncrypt() {

	fileName := d.Url.UrlBuild.GetDataRequest(d.Net.R)
	if fileName == "" {

		d.Answer.S.SetBadAnswers(RepoDownloadEncryptRepo.IncomingDataAnswer{
			W:     d.Net.W,
			Error: DomainLevel.ErrorCantGetFileName,
		})
		return
	}
	err := d.DownloadEncrypt(Application.NewDownloadEncryptIncomingData{
		NewDownloadEncryptNetwork: Application.NewDownloadEncryptNetwork{
			d.Net.W,
		},
		Ctx:          d.Net.R.Context(),
		EncryptedURl: fileName,
	})
	if err != nil {
		d.Answer.S.SetBadAnswers(RepoDownloadEncryptRepo.IncomingDataAnswer{
			W:     d.Net.W,
			Error: err.Error(),
		})
		return
	}
	return

}
