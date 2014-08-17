package main

type User interface {
	Name() string
	IsAdmin() bool
}

type userT string

func NewUser(name string) User {
	u := userT(name)
	return &u
}

func (u userT) Name() string {
	return string(u)
}

func (u userT) IsAdmin() bool {
	return false
}

type UsersT struct {
	m map[string]User
}

var Users UsersT = UsersT{make(map[string]User)}

func (us *UsersT) GetUser(name string) User {
	return us.m[name]
}

func init() {
	Users.m["admin"] = NewUser("admin")
	Users.m[""] = NewUser("anonymous")
}
