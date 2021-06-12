// models.user.go

package main

import (
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type user struct {
	gorm.Model
	Username string `json:"username"`
	Password string `json:"-"`
}

// Check if the username and password combination is valid
func isUserValid(username, password string) bool {
	var userList []user
	DB.Find(&userList)
	for _, u := range userList {
		fmt.Println("User:" + username + "\nPassword:" + password)
		if u.Username == username && u.Password == password {
			return true
		}
	}
	return false
}

// Register a new user with the given username and password
// NOTE: For this demo, we
func registerNewUser(username, password string) (*user, error) {
	if strings.TrimSpace(password) == "" {
		return nil, errors.New("The password can't be empty")
	} else if !isUsernameAvailable(username) {
		return nil, errors.New("The username isn't available")
	}

	u := user{Username: username, Password: password}

	DB.Create(&u)

	return &u, nil
}

// Check if the supplied username is available
func isUsernameAvailable(username string) bool {
	var userList []user
	DB.Find(&userList)
	for _, u := range userList {
		if u.Username == username {
			return false
		}
	}
	return true
}
