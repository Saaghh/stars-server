// Code generated by ogen, DO NOT EDIT.

package api

type CreateUserBadRequest Error

func (*CreateUserBadRequest) createUserRes() {}

type CreateUserInternalServerError Error

func (*CreateUserInternalServerError) createUserRes() {}

// CreateUserOK is response for CreateUser operation.
type CreateUserOK struct{}

func (*CreateUserOK) createUserRes() {}

// Ref: #/components/schemas/error
type Error struct {
	Code        int    `json:"code"`
	Description string `json:"description"`
}

// GetCode returns the value of Code.
func (s *Error) GetCode() int {
	return s.Code
}

// GetDescription returns the value of Description.
func (s *Error) GetDescription() string {
	return s.Description
}

// SetCode sets the value of Code.
func (s *Error) SetCode(val int) {
	s.Code = val
}

// SetDescription sets the value of Description.
func (s *Error) SetDescription(val string) {
	s.Description = val
}

type GameTickBadRequest Error

func (*GameTickBadRequest) gameTickRes() {}

type GameTickInternalServerError Error

func (*GameTickInternalServerError) gameTickRes() {}

// GameTickOK is response for GameTick operation.
type GameTickOK struct{}

func (*GameTickOK) gameTickRes() {}

// Ref: #/components/schemas/userCreate
type UserCreate struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// GetUsername returns the value of Username.
func (s *UserCreate) GetUsername() string {
	return s.Username
}

// GetPassword returns the value of Password.
func (s *UserCreate) GetPassword() string {
	return s.Password
}

// SetUsername sets the value of Username.
func (s *UserCreate) SetUsername(val string) {
	s.Username = val
}

// SetPassword sets the value of Password.
func (s *UserCreate) SetPassword(val string) {
	s.Password = val
}