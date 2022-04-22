package model

type Role struct {
	ID   string `json:"id,omitempty" binding:"required"`
	Name string `json:"name" binding:"required"`
}

type User struct {
	ID       int    `json:"id,omitempty" binding:"required"`
	Role     *Role  `json:"role,omitempty"`
	Username string `json:"username"`
}
