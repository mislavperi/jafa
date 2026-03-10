package services

import "errors"

var ErrInvalidCredentials = errors.New("invalid username or password")
var ErrUsernameTaken = errors.New("username already taken")
