package Deliver

import (
	"Kaban/internal/DomainLevel"
	"Kaban/internal/Dto"
	"Kaban/internal/InfrastructureLayer/DeliverPackages/RepoLoginRealizations"
	"Kaban/internal/InfrastructureLayer/DeliverPackages/RepoParsers"
	"Kaban/internal/InfrastructureLayer/DeliverPackages/RepoSessionHandle"
	"net/http"
)

type Network struct {
	W http.ResponseWriter
	R *http.Request
}

type Answer struct {
	S *RepoLoginRealizations.LoginAnswers
}
type Parse struct {
	Parses *RepoParsers.Parsing
}
type NewLogin struct {
	Net           Network
	ControlAnswer Answer
	Parser        Parse
}

func (D *NewLogin) Login() {
	if D.Net.R.Method != http.MethodPost {
		D.ControlAnswer.S.SetBadAnswer(RepoLoginRealizations.AnswerDetails{
			W:         D.Net.W,
			Operation: DomainLevel.NotStart,
			Error:     DomainLevel.MethodNotAllowed,
		})
		return
	}

	DataUserLogin := &Dto.UserLoginData{}
	err := D.Parser.Parses.JsonDecode(DataUserLogin, D.Net.R.Body)
	if err != nil {
		D.ControlAnswer.S.SetBadAnswer(RepoLoginRealizations.AnswerDetails{
			W:         D.Net.W,
			Operation: DomainLevel.NotStart,
			Error:     err.Error(),
		})
		return
	}

	err = DataUserLogin.ValidateData()
	if err != nil {
		D.ControlAnswer.S.SetBadAnswer(RepoLoginRealizations.AnswerDetails{
			W:         D.Net.W,
			Operation: DomainLevel.NotStart,
			Error:     err.Error(),
		})
		return
	}

	JwtToken, RefreshToken, err := realization.LoginService(*DataUserLogin, D.Net.R.Context())
	if err != nil {
		//TODO add handling the error
		return
	}

	ReturnedData := RepoSessionHandle.SessionControl.SetNewSession(RepoSessionHandle.IncomingSessionData{
		Writer:  D.Net.W,
		Request: D.Net.R,
		Jwt:     JwtToken,
		Rt:      RefreshToken,
	})
	if ReturnedData.Error != nil {
		D.ControlAnswer.S.SetBadAnswer(RepoLoginRealizations.AnswerDetails{
			W:         D.Net.W,
			Operation: DomainLevel.Break,
			Error:     ReturnedData.Error.Error(),
		})
		return
	}

	D.ControlAnswer.S.SetGoodAnswer(RepoLoginRealizations.AnswerDetails{
		W:               D.Net.W,
		Operation:       DomainLevel.Success,
		UrlToRedistrict: "/main",
	})
	return
}
