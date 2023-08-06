package dao

import (
	MySQL "github.com/akgarg0472/urlshortener-auth-service/database"
	AuthModel "github.com/akgarg0472/urlshortener-auth-service/model"
)

var db = MySQL.GetInstance()

func AddUser(signupRequest AuthModel.SignupRequest) {

}
