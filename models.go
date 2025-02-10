package main

type Device struct {
	UUID string `gorm:"primary_key"`
	Name string
}
