package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID       string     `gorm:"default:gen_random_uuid()" json:"id"`
	FullName string     `json:"full_name"`
	Email    string     `gorm:"uniqueIndex" json:"email"`                    // auth username
	Password string     `json:"-"`                                           // auth password
	Role     UserRole   `gorm:"type:user_role;default:unsetted" json:"role"` // "student", "teacher", "admin", "unsetted"
	Gender   UserGender `gorm:"type:user_gender" json:"gender"`              // "male", "female"
	Phone    string     `gorm:"uniqueIndex" json:"phone"`                    // phone number

	StudentProfile *Student `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"student_profile,omitempty" swaggerignore:"true"`
	TeacherProfile *Teacher `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"teacher_profile,omitempty" swaggerignore:"true"`
	AdminProfile   *Admin   `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"admin_profile,omitempty" swaggerignore:"true"`

	DeletedAt gorm.DeletedAt `gorm:"index" json:"archived_at,omitzero" swaggertype:"string"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at,omitzero"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at,omitzero"`
}

type Student struct {
	ID      string    `gorm:"default:gen_random_uuid()" json:"id"`
	ClassID string    `json:"-"`
	Class   *Class    `json:"class,omitempty"`
	Parents []*Parent `gorm:"many2many:student_parents" json:"parents"`
	NISN    string    `gorm:"unique;not null"`

	UserID string `gorm:"uniqueIndex;not null" json:"-"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at,omitzero"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at,omitzero"`
}

type Teacher struct {
	ID         string     `gorm:"default:gen_random_uuid()" json:"id"`
	Subjects   []*Subject `gorm:"many2many:teacher_subjects" json:"subjects"`
	NUPTK      string     `gorm:"unique;not null"`
	EmployeeID string     `gorm:"not null" json:"employee_id"`

	UserID string `gorm:"uniqueIndex;not null" json:"-"`

	JoinedAt  time.Time `gorm:"not null" json:"joined_at"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at,omitzero"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at,omitzero"`
}

type Admin struct {
	ID          string        `gorm:"default:gen_random_uuid()" json:"id"`
	StaffRole   string        `json:"staff_role"`
	Permissions []*Permission `gorm:"many2many:admin_permissions" json:"permissions,omitempty" swaggerignore:"true"`
	EmployeeID  string        `gorm:"not null" json:"employee_id"`

	UserID string `gorm:"uniqueIndex;not null" json:"-"`

	JoinedAt  time.Time `gorm:"not null" json:"joined_at"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at,omitzero"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at,omitzero"`
}

type Parent struct {
	ID       string     `gorm:"default:gen_random_uuid()" json:"id"`
	FullName string     `json:"full_name"`
	Phone    string     `gorm:"unique" json:"phone"` // phone number
	Email    string     `gorm:"unique" json:"email"`
	Gender   UserGender `gorm:"type:user_gender" json:"gender"` // "male", "female"

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at,omitzero"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at,omitzero"`
}
