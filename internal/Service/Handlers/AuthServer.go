package Handlers

func (sa *HandlerPackCollect) Auth(Rt string, JwtToken string) (string, error) {

	JwtTokenClaims, err := sa.AuthTokens.Checking.CheckJwt(JwtToken)
	if err != nil {
		return "", nil
	}
	RefreshToken, err := sa.AuthTokens.Checking.CheckRt(Rt)
	if err != nil || RefreshToken == nil {
		return "", err
	}

	if !JwtTokenClaims.Valid && RefreshToken.Valid {
		JwtToken, err = sa.AuthTokens.GeneratingToken.GenerateJWT(RefreshToken.Claims)
		if err != nil {
			return "", err
		}

		return JwtToken, nil
	}
	return "", nil

}
