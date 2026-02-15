package Handlers

import (
	"Kaban/iternal/InfrastructureLayer"
	"errors"
	"log/slog"
)

func Auth(Rt string, JwtToken string) (string, string, error) {

	manageToken := *InfrastructureLayer.SetSittingsTokenInteraction()

	ok := manageToken.Tokens.TokenDenyMapChecker(Rt)

	if ok {
		return "", "", errors.New("token Deny")
	}

	err := manageToken.Tokens.CheckLifeJwt(JwtToken)
	if err == nil {
		return "", "", nil
	}

	err = manageToken.Tokens.CheckLifeRt(Rt)
	if err != nil {
		return "", "", err
	}

	NewJwt, NewRet, err := manageToken.Tokens.GenerateNewTokens(Rt)
	if err != nil {
		slog.Error("error generate new tokens", "Value", err.Error())
		return "", "", err
	}

	manageToken.Tokens.DeleteAndSaveToken(Rt, NewRet)
	return NewJwt, NewRet, nil

}
