package DomainLevel

import "context"

type RepositorysTest interface {
	GetIdPassowrdTest(string) (int, string, error)

	GeTScryptTest(context.Context, int) (string, error)
}
