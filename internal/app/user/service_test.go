package user_test

import (
	"context"
	"testing"

	"github.com/Januadrym/seennit/internal/app/types"
	"github.com/Januadrym/seennit/internal/app/user"
	"github.com/golang/mock/gomock"
)

func TestRegister(t *testing.T) {
	mockedRepo := NewMockrepoProvider(gomock.NewController(t))
	service := user.NewService(mockedRepo)

	tUser := &types.RegisterRequest{
		// UserID:    "112",
		FirstName: "Quan",
		LastName:  "Vo",
		Email:     "this@gmail.com",
		Password:  "kinkou",
		// Locked:    false,
		// CreatedAt :time.Time,
		// UpdatedAt: time.Time,
	}

	outUser := &types.User{
		UserID:    "112",
		FirstName: "Quan",
		LastName:  "Vo",
		Email:     "this@gmail.com",
		Password:  "kinkou",
		Locked:    false,
		// CreatedAt :time.Time,
		// UpdatedAt: time.Time,
	}

	//dbErr := errors.New("cannot connect to db")

	testCases := []struct {
		name     string
		thisFunc func()
		input    *types.RegisterRequest
		output   *types.User
		tErr     error
	}{
		{
			name:  "user not exist, register successful",
			input: tUser,
			thisFunc: func() {
				mockedRepo.EXPECT().FindUserByMail(gomock.Any(), tUser.Email).Times(1).Return(nil, nil)
				mockedRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Times(1).Return(nil)
			},
			output: outUser,
			tErr:   nil,
		},
	}
	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			test.thisFunc()
			usr, err := service.Register(context.Background(), test.input)
			if err != test.tErr {
				t.Errorf("got err = %v, want err = %v", err, test.output)
			}
			if usr == nil {
				t.Fatal("no user return ")
			}
		})
	}
}
