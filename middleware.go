package main

func CheckPassword(password string) bool {
	return password == cfg.Password
}