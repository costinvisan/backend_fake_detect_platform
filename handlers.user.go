// handlers.user.go

package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type authenticate struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func showLoginPage(c *gin.Context) {
	// Call the render function with the name of the template to render
	render(c, gin.H{
		"title": "Login",
	}, "login.html")
}

func performLogin(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	var body authenticate
	c.ShouldBindJSON(&body)

	if (authenticate{} == body) {

		// Check if the username/password combination is valid
		if u := isUserValid(username, password); (u != user{}) {
			// If the username/password is valid set the token in a cookie
			token := generateSessionToken()
			c.SetCookie("token", token, 3600, "", "", false, true)
			c.Set("is_logged_in", true)
			render(c, gin.H{
				"title": "Successful Login"}, "login-successful.html")

		} else {
			// If the username/password combination is invalid,
			// show the error message on the login page
			c.HTML(http.StatusBadRequest, "login.html", gin.H{
				"ErrorTitle":   "Login Failed",
				"ErrorMessage": "Invalid credentials provided"})
		}
	} else {
		username := body.Username
		password := body.Password

		if u := isUserValid(username, password); (u != user{}) {
			c.SecureJSON(http.StatusOK, u)
		} else {
			c.SecureJSON(http.StatusOK, "")
		}

	}
	// Obtain the POSTed username and password values

}

func generateSessionToken() string {
	// We're using a random 16 character string as the session token
	// This is NOT a secure way of generating session tokens
	// DO NOT USE THIS IN PRODUCTION
	return strconv.FormatInt(rand.Int63(), 16)
}

func logout(c *gin.Context) {
	// Clear the cookie
	c.SetCookie("token", "", -1, "", "", false, true)

	// Redirect to the home page
	c.Redirect(http.StatusTemporaryRedirect, "/")
}

func showRegistrationPage(c *gin.Context) {
	// Call the render function with the name of the template to render
	render(c, gin.H{
		"title": "Register"}, "register.html")
}

func register(c *gin.Context) {
	// Obtain the POSTed username and password values
	username := c.PostForm("username")
	password := c.PostForm("password")

	var body authenticate
	c.ShouldBindJSON(&body)

	if (authenticate{} == body) {

		if _, err := registerNewUser(username, password); err == nil {
			// If the user is created, set the token in a cookie and log the user in
			token := generateSessionToken()
			c.SetCookie("token", token, 3600, "", "", false, true)
			c.Set("is_logged_in", true)

			render(c, gin.H{
				"title": "Successful registration & Login"}, "login-successful.html")

		} else {
			// If the username/password combination is invalid,
			// show the error message on the login page
			c.HTML(http.StatusBadRequest, "register.html", gin.H{
				"ErrorTitle":   "Registration Failed",
				"ErrorMessage": err.Error()})

		}
	} else {
		username := body.Username
		password := body.Password
		if u, err := registerNewUser(username, password); err == nil {
			c.SecureJSON(http.StatusOK, u)
		} else {
			c.SecureJSON(http.StatusServiceUnavailable, "")
		}
	}
}

func manageUsers(c *gin.Context) {
	var users []user
	users = getAllUsers()
	fmt.Println(users)
	render(c, gin.H{
		"title":   "Manage Users",
		"payload": users}, "create-article.html")
}

func deleteUser(c *gin.Context) {
	if userID, err := strconv.Atoi(c.Param("user_id")); err == nil {
		deleteUserById(userID)
		render(c, gin.H{
			"title":   "Submission Successful",
			"payload": nil}, "submission-successful.html")
	}
}
