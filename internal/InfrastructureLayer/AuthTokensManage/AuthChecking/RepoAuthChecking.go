package AuthChecking

import (
	"Kaban/internal/DomainLevel"
	"Kaban/internal/InfrastructureLayer/AuthTokensManage/Creating"
	"Kaban/internal/InfrastructureLayer/AuthTokensManage/ValidatingTokens"
)

type NewAuthChecker struct {
	Validate     *ValidatingTokens.Checking
	CreateTokens *Creating.CreatingTokens
}

func (n NewAuthChecker) CheckAuthTokens(data DomainLevel.AuthCheckIncomingData) DomainLevel.OutComingAuthData {
	err := n.Validate.CheckJwt(data.Jwt)
	if err == nil {
		return DomainLevel.OutComingAuthData{
			IsNewJwtCreated: true,
		}
	}
	Claims, err := n.Validate.CheckRt(data.Rft)
	if err != nil {
		return DomainLevel.OutComingAuthData{
			Err: err,
		}
	}

	JwtToken, err := n.CreateTokens.GenerateJWT(Claims)
	if err != nil {
		return DomainLevel.OutComingAuthData{Err: err}
	}
	return DomainLevel.OutComingAuthData{
		NewJwt: JwtToken,
	}
}
