package models

import "time"

type TokenDetail struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	AccessUuid   string    `json:"-"`
	RefreshUuid  string    `json:"-"`
	AtExpires    time.Time `json:"-"`
	RtExpires    time.Time `json:"-"`
}
