package main

import (
	"testing"

	"github.com/go-playground/validator/v10"
)

type Peopole struct {
	Name      string   `validate:"min=3,max=2"`       // 长度范围
	Gender    string   `validate:"oneof=male female"` // 性别
	Email     string   `validate:"required,email"`    //邮件
	Password  string   `validate:"min=6"`
	Password2 string   `validate:"eqfield=Password"` //二次密码
	Hobbies   []string `validate:"unique"`           // 唯一性
}

func TestPeopleValidation(t *testing.T) {
	testCases := []struct {
		desc string
		p    Peopole
	}{
		{
			desc: "ok",
			p: Peopole{
				Name:      "smith",
				Gender:    "male",
				Email:     "88@mm.com",
				Password:  "123456",
				Password2: "123456",
				Hobbies:   []string{"pinpong", "hiking", "swimming"},
			},
		},
		{
			desc: "false",
			p: Peopole{
				Name:      "jack",
				Gender:    "prefer_not_to",
				Email:     "8888",
				Password:  "1234",
				Password2: "123",
				Hobbies:   []string{"pinpong", "hiking", "swimming", "swimming"},
			},
		},
	}
	validate := validator.New()
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			err := validate.Struct(tC.p)
			if err != nil {
				t.Error(err)
			}
		})
	}
}
