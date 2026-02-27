package Handlers

import "log/slog"

func (sa *HandlerPackCollect) Auth(Rt string, JwtToken string) (string, error) {

	//ok := manageToken.Tokens.TokenDenyMapChecker(Rt)

	//if ok {
	//	return "", "", errors.New("token Deny")
	//}
	slog.Info("Auth system", "Hy", sa.S.FileInfo.SayHi())

	JwtTokenClaims, err := sa.S.Tokens.CheckLifeJwt(JwtToken)
	if err == nil || JwtTokenClaims == nil {
		return "", nil
	}
	RefreshToken, err := sa.S.Tokens.CheckLifeRt(Rt)
	if err != nil || RefreshToken == nil {
		return "", err
	}

	if !JwtTokenClaims.Valid {
		if RefreshToken.Valid {
			JwtToken, err = sa.S.Tokens.GenerateJWT(RefreshToken.Claims)
			if err != nil {
				return "", err
			}

			return JwtToken, nil
		}
	}
	return "", nil

}
