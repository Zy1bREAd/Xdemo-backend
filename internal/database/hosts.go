package database

import "time"

type Host struct {
	UUID          string
	Name          string
	Status        string
	CreateAt      time.Time
	CreateBy      string
	UpdateAt      time.Time
	UpdateBy      string
	OfflineAt     time.Time
	OfflineEvents string
	Arena         string
	InternalIP    []string
	PublicIP      []string
	Type          string
	Account       string
	Password      string
}
