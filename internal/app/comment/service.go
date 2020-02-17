package comment

type (
	RepoProvider interface {
	}
	Service struct {
		Repo RepoProvider
	}
)

func NewService(repo RepoProvider) *Service {
	return &Service{
		Repo: repo,
	}
}
