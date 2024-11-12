package models

import "database/sql"

// ユーザー名とパスワードで共通の構造体を使用
type UserAuth struct {
	UserName string `binding:"required,min=3,max=50" json:"user_name"`
	Password string `binding:"required,min=8,max=255" json:"password"`
}

// 表示時に使用する構造体
type UserResponse struct {
	UserName string         `binding:"required" json:"user_name"`
	Bio      sql.NullString `json:"bio"`
}

type UpdateUserRequest struct {
	NewUserName string `binding:"required,min=3,max=50" json:"new_user_name"`
	NewBio      string `binding:"required,max=255" json:"new_bio"`
}
