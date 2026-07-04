package newAuthCheck

import (
	"Kaban/internal/DomainLevel"
	"Kaban/internal/InfrastructureLayer/AuthTokensManage/Creating"
	"Kaban/internal/InfrastructureLayer/AuthTokensManage/ValidatingTokens"
	"errors"
)

type NewAuthChecker struct {
	TokenCreate *Creating.CreatingTokens
	TokenCheck  *ValidatingTokens.Checking
}

func (sa NewAuthChecker) CheckUserAuth(data DomainLevel.UserAuthCheck) (string, error) {
	err := sa.TokenCheck.CheckJwt(data.Jwt)
	if err == nil {
		return "", nil
	}
	Claims, err := sa.TokenCheck.CheckRt(data.Rft)
	if err != nil {
		return "", err
	}

	newJwtToken, err := sa.TokenCreate.GenerateJWT(Claims)
	if err != nil {
		return "", errors.New(DomainLevel.ErrorUserToken)
	}
	return newJwtToken, nil

}
