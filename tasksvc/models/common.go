package models

import "github.com/jmoiron/sqlx"

type StatusMessageResponse struct {
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Status  string      `json:"status,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

type TaskDB struct {
	Db *sqlx.DB
}
