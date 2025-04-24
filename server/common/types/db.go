package types

import "time"

type User struct {
	Id       int       `json:"id" gorm:"primaryKey,autoIncrement"`
	Username string    `json:"username" gorm:"size:50"`
	Password string    `json:"password,omitempty" gorm:"size:100"`
	Role     string    `json:"role" gorm:"size:50"`
	Profile  StringMap `json:"profile,omitempty" gorm:"type:json"`
	UpdateAt time.Time `json:"update_at"`
}

type S3Repository struct {
	Id          int       `json:"id" gorm:"primaryKey,autoIncrement"`
	Name        string    `json:"name" gorm:"size:100,unique"`
	S3Endpoint  string    `json:"s3_endpoint" gorm:"size:300"`
	S3Region    string    `json:"s3_region" gorm:"size:50"`
	S3AccessKey string    `json:"s3_access_key" gorm:"size:100"`
	S3SecretKey string    `json:"s3_secret_key" gorm:"size:100"`
	UseSecure   bool      `json:"use_secure"`
	UpdateAt    time.Time `json:"update_at"`
}

type Settings struct {
	Id           int    `json:"id" gorm:"primaryKey;autoIncrement"`
	SettingKey   string `json:"setting_key" gorm:"unique"`
	SettingValue string `json:"setting_value"`
}
