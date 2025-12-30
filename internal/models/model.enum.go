package models

type UserRole string // "student", "teacher", "admin", "unsetted"
const (
	RoleStudent  UserRole = "student"
	RoleTeacher  UserRole = "teacher"
	RoleAdmin    UserRole = "admin"
	RoleUnsetted UserRole = "unsetted"
)

var UserRoles = []UserRole{
	RoleStudent,
	RoleTeacher,
	RoleAdmin,
	RoleUnsetted,
}

type UserGender string // "male", "female"
const (
	GenderMale   UserGender = "male"
	GenderFemale UserGender = "female"
)

var UserGenders = []UserGender{GenderMale, GenderFemale}

type PermissionAction string // "create", "read", "update", "delete"
const (
	ActionCreate PermissionAction = "create"
	ActionRead   PermissionAction = "read"
	ActionUpdate PermissionAction = "update"
	ActionDelete PermissionAction = "delete"
)

var PermissionActions = []PermissionAction{
	ActionCreate,
	ActionRead,
	ActionUpdate,
	ActionDelete,
}

type PermissionResource string // "role", "permission", "admin", "teacher", "student", "subject", "class"
const (
	ResourceRole       PermissionResource = "role"
	ResourcePermission PermissionResource = "permission"

	ResourceAdmin   PermissionResource = "admin"
	ResourceTeacher PermissionResource = "teacher"
	ResourceStudent PermissionResource = "student"

	ResourceSubject PermissionResource = "subject"
	ResourceClass   PermissionResource = "class"
)

var PermissionResources = []PermissionResource{
	ResourceRole,
	ResourcePermission,

	ResourceAdmin,
	ResourceTeacher,
	ResourceStudent,

	ResourceSubject,
	ResourceClass,
}
