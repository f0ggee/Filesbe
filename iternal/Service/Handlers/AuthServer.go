package Handlers

import (
	"Kaban/iternal/InfrastructureLayer"
)

func Auth(Rt string, JwtToken string) (string, error) {

	manageToken := *InfrastructureLayer.SetSittingsTokenInteraction()

	//ok := manageToken.Tokens.TokenDenyMapChecker(Rt)

	//if ok {
	//	return "", "", errors.New("token Deny")
	//}

	JwtTokenClaims, err := manageToken.Tokens.CheckLifeJwt(JwtToken)
	if err == nil || JwtTokenClaims == nil {
		return "", nil
	}
	RefreshToken, err := manageToken.Tokens.CheckLifeRt(Rt)
	if err != nil || RefreshToken == nil {
		return "", err
	}

	if !JwtTokenClaims.Valid {
		if RefreshToken.Valid {
			JwtToken, err = manageToken.Tokens.GenerateJWT(RefreshToken.Claims)
			if err != nil {
				return "", err
			}

			return JwtToken, nil
		}
	}
	return "", nil

}
