package status

import (
	"os"
	"sync"

	"github.com/Januadrym/seennit/internal/pkg/status"

	"gopkg.in/yaml.v2"
)

type (
	Status    = status.Status
	GenStatus struct {
		Success    Status
		NotFound   Status
		Timeout    status.Timeout
		BadRequest Status
		Internal   Status
	}

	UserStatus struct {
		DuplicatedEmail Status `yaml:"duplicated_email"`
	}

	AuthStatus struct {
		InvalidUserPassword Status `yaml:"invalid_user_password"`
	}
	PolicyStatus struct {
		Unauthorized Status
	}
	statuses struct {
		Gen    GenStatus
		User   UserStatus
		Auth   AuthStatus
		Policy PolicyStatus
	}
)

var (
	all  *statuses
	once sync.Once
)

func Init(conf string) {
	once.Do(func() {
		file, err := os.Open(conf)
		if err != nil {
			panic(err)
		}
		all = &statuses{}
		if err := yaml.NewDecoder(file).Decode(all); err != nil {
			panic(err)
		}
	})
}

func load() *statuses {
	conf := os.Getenv("STATUS_PATH")
	if conf == "" {
		conf = "configs/status.yml"
	}
	Init(conf)
	return all
}

func Gen() GenStatus {
	return load().Gen
}

func User() UserStatus {
	return load().User
}

func Success() Status {
	return Gen().Success
}

func Auth() AuthStatus {
	return load().Auth
}
func Policy() PolicyStatus {
	return load().Policy
}
