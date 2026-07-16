package Deliver

import (
	"Kaban/internal/DomainLevel"
	"Kaban/internal/Dto"
	"Kaban/internal/InfrastructureLayer/DeliverPackages/RepoLoginRealizations"
	"Kaban/internal/InfrastructureLayer/RepoParsers"
	"Kaban/internal/InfrastructureLayer/RepoSessionHandle"
	"Kaban/internal/Service/Application"
	"net/http"
)

type LoginNet struct {
	W http.ResponseWriter
	R *http.Request
}

type LoginDepends struct {
	S    *RepoLoginRealizations.LoginAnswers
	Sess RepoSessionHandle.Session
}
type ParseLogin struct {
	Parses RepoParsers.Decode
}
type LoginApplication struct {
	Application.NewLogin
}
type NewLogin struct {
	LoginNet
	LoginDepends
	ParseLogin
	LoginApplication
}

func GetNewLogin(networkLogin LoginNet, loginDepends LoginDepends, parseLogin ParseLogin, loginApplication LoginApplication) *NewLogin {
	return &NewLogin{LoginNet: networkLogin, LoginDepends: loginDepends, ParseLogin: parseLogin, LoginApplication: loginApplication}
}

func (D *NewLogin) Login() {
	if D.R.Method != http.MethodPost {
		D.S.SetBadAnswer(RepoLoginRealizations.AnswerDetails{
			W:         D.W,
			Operation: DomainLevel.NotStart,
			Error:     DomainLevel.MethodNotAllowed,
		})
		return
	}

	DataUserLogin := &Dto.UserLoginData{}
	err := D.Parses.JsonDecode(DataUserLogin, D.R.Body)
	if err != nil {
		D.S.SetBadAnswer(RepoLoginRealizations.AnswerDetails{
			W:         D.W,
			Operation: DomainLevel.NotStart,
			Error:     err.Error(),
		})
		return
	}

	err = DataUserLogin.ValidateData()
	if err != nil {
		D.S.SetBadAnswer(RepoLoginRealizations.AnswerDetails{
			W:         D.W,
			Operation: DomainLevel.NotStart,
			Error:     err.Error(),
		})
		return
	}

	loginDataOutput := D.LoginService(*DataUserLogin, D.R.Context())
	if loginDataOutput.Err != nil {
		D.S.SetGoodAnswer(RepoLoginRealizations.AnswerDetails{
			W:         D.W,
			Operation: DomainLevel.Break,
			Error:     loginDataOutput.Err.Error(),
		})
		return
	}
	ReturnedData := D.Sess.SetNewSession(RepoSessionHandle.IncomingSessionData{
		Writer:  D.W,
		Request: D.R,
		Jwt:     loginDataOutput.Jwt,
		Rt:      loginDataOutput.Rft,
	})
	if ReturnedData.Error != nil {
		D.S.SetBadAnswer(RepoLoginRealizations.AnswerDetails{
			W:         D.W,
			Operation: DomainLevel.Break,
			Error:     ReturnedData.Error.Error(),
		})
		return
	}
	D.S.SetGoodAnswer(RepoLoginRealizations.AnswerDetails{
		W:               D.W,
		Operation:       DomainLevel.Success,
		UrlToRedistrict: "/main",
	})
	return
}
