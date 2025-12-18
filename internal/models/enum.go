package models

type UserRole string // "student", "teacher", "admin"
const (
	RoleStudent UserRole = "student"
	RoleTeacher UserRole = "teacher"
	RoleAdmin   UserRole = "admin"
)

type UserGender string // "male", "female"
const (
	GenderMale   UserGender = "male"
	GenderFemale UserGender = "female"
)

type PermissionAction string // "create", "read", "update", "delete"
const (
	ActionCreate PermissionAction = "create"
	ActionRead   PermissionAction = "read"
	ActionUpdate PermissionAction = "update"
	ActionDelete PermissionAction = "delete"
)
