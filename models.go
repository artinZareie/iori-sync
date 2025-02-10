package main

import (
	"time"
)

type Device struct {
	UUID string `gorm:"primaryKey;type:varchar(255)"`
	Name string `gorm:"type:varchar(255)"`
}

type Group struct {
	ID       uint   `gorm:"primaryKey"`
	ParentID *uint  // Nullable parent ID
	Name     string `gorm:"type:varchar(255)"`

	// Relationships
	Parent      *Group       `gorm:"foreignKey:ParentID"`
	Files       []File       `gorm:"foreignKey:GroupID"`
	Subscribers []Subscriber `gorm:"foreignKey:GroupID"`
}

type File struct {
	ID             uint      `gorm:"primaryKey"`
	GroupID        uint      `gorm:"index"`
	LastUpdateDate time.Time `gorm:"type:datetime"`
	FileName       string    `gorm:"type:varchar(255)"`
	FileHash       string    `gorm:"type:varchar(255)"`

	// Relationships
	Group          *Group          `gorm:"foreignKey:GroupID"`
	UpdatedClients []UpdatedClient `gorm:"foreignKey:FileID"`
}

type Subscriber struct {
	ID         uint   `gorm:"primaryKey"`
	DeviceUUID string `gorm:"type:varchar(255);index"`
	GroupID    uint   `gorm:"index"`

	// Relationships
	Device *Device `gorm:"foreignKey:DeviceUUID;references:UUID"`
	Group  *Group  `gorm:"foreignKey:GroupID;references:ID"`
}

type UpdatedClient struct {
	ID         uint   `gorm:"primaryKey"`
	FileID     uint   `gorm:"index"`
	DeviceUUID string `gorm:"type:varchar(255);index"`

	// Relationships
	File   *File   `gorm:"foreignKey:FileID;references:ID"`
	Device *Device `gorm:"foreignKey:DeviceUUID;references:UUID"`
}
