package Handlers

import (
	"errors"
)

func (sa *HandlerPackCollect) Auth(Rt string, JwtToken string) (string, error) {

	JwtTokenClaims, err := sa.AuthTokens.Checking.CheckJwt(JwtToken)
	if err == nil && JwtTokenClaims.Valid {
		return "", nil
	}
	RefreshToken, err := sa.AuthTokens.Checking.CheckRt(Rt)
	if err != nil || RefreshToken == nil {
		return "", err
	}
	if RefreshToken.Valid {
		JwtToken, err = sa.AuthTokens.GeneratingToken.GenerateJWT(RefreshToken.Claims)
		if err != nil {
			return "", err
		}

		return JwtToken, nil
	}
	return "", nil

}

func (sa *HandlerPackAuthTokens) AuthTest(Rt string, JwtToken string) (string, error) {

	JwtTokenClaims, err := sa.Checking.CheckJwt(JwtToken)
	if err == nil && JwtTokenClaims.Valid {
		return "", nil
	}
	RefreshToken, err := sa.Checking.CheckRt(Rt)
	if err != nil || RefreshToken == nil {
		return "", err
	}
	if RefreshToken.Valid {
		JwtToken, err = sa.GeneratingToken.GenerateJWT(RefreshToken.Claims)
		if err != nil {
			return "", err
		}

		return JwtToken, nil
	}
	return "", nil
}
func funcName(Rt string, JwtToken string, sa *HandlerPackAuthTokens) error {
	JwtTokenClaims, err := sa.Checking.CheckJwt(JwtToken)
	if err == nil && JwtTokenClaims.Valid {
		return nil
	}
	RefreshToken, err := sa.Checking.CheckRt(Rt)
	if err != nil || RefreshToken == nil {
		return errors.New("errror")
	}
	return nil
}
