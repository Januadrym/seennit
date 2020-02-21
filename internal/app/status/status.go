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
		Success    Status         `yaml:"success"`
		NotFound   Status         `yaml:"not_found"`
		Timeout    status.Timeout `yaml:"timeout"`
		BadRequest Status         `yaml:"bad_request"`
		Internal   Status         `yaml:"internal"`
	}

	UserStatus struct {
		DuplicatedEmail Status `yaml:"duplicated_email"`
		RegisterFail    Status `yaml:"register_fail"`
	}

	AuthStatus struct {
		InvalidUserPassword Status `yaml:"invalid_user_password"`
	}

	PolicyStatus struct {
		Unauthorized Status `yaml:"unauthorized"`
	}

	CommunityStatus struct {
		NameTaken    Status `yaml:"name_taken"`
		UserEnrolled Status `yaml:"user_enrolled"`
		NotFound     Status `yaml:"community_not_found"`
		CreateFail   Status `yaml:"create_failed"`
	}
	PostStatus struct {
		Archived Status `yaml:"archived"`
		NotFound Status `yaml:"post_not_found"`
	}
	statuses struct {
		Gen       GenStatus
		User      UserStatus
		Auth      AuthStatus
		Policy    PolicyStatus
		Community CommunityStatus
		Post      PostStatus
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

func Community() CommunityStatus {
	return load().Community
}
func Post() PostStatus {
	return load().Post
}
