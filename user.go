package main

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"

	"github.com/go-acme/lego/lego"
	"github.com/go-acme/lego/registration"
)

func newUser(email string) (up *User, err error) {
	var u User
	// Set random key for user
	if u.key, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader); err != nil {
		return
	}

	// Set user email
	u.Email = email

	// Assign reference to User
	up = &u
	return
}

// User implements the acme.User interface
type User struct {
	key crypto.PrivateKey

	// Users email address
	Email string
	// Users Let's Encrypt registration
	Registration *registration.Resource
}

// GetEmail will get the email for a user
func (u *User) GetEmail() string {
	return u.Email
}

// GetRegistration will get the registration for a user
func (u User) GetRegistration() *registration.Resource {
	return u.Registration
}

// GetPrivateKey will get the private key for a user
func (u *User) GetPrivateKey() crypto.PrivateKey {
	return u.key
}

// Register will register with the provided Client
func (u *User) Register(client *lego.Client) (err error) {
	// Register utilizing the provided client
	if u.Registration, err = client.Registration.Register(registrationOpts); err != nil {
		return
	}

	return
}
