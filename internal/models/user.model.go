package models

type User struct {
	Id
	FullName string     `json:"full_name"`
	Email    string     `gorm:"uniqueIndex" json:"email"`                    // auth username
	Password string     `json:"-"`                                           // auth password
	Role     UserRole   `gorm:"type:user_role;default:unsetted" json:"role"` // "student", "teacher", "admin", "unsetted"
	Gender   UserGender `gorm:"type:user_gender" json:"gender"`              // "male", "female"
	Phone    string     `gorm:"uniqueIndex" json:"phone"`                    // phone number

	StudentProfile *Student `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"student_profile,omitempty" swaggerignore:"true"`
	TeacherProfile *Teacher `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"teacher_profile,omitempty" swaggerignore:"true"`
	AdminProfile   *Admin   `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"admin_profile,omitempty" swaggerignore:"true"`

	TimestampArchivable
}

type Student struct {
	Id
	ClassID string    `json:"-"`
	Class   *Class    `json:"class,omitempty"`
	Parents []*Parent `gorm:"many2many:student_parents" json:"parents"`
	NISN    string    `gorm:"unique;not null"`

	UserID string `gorm:"uniqueIndex;not null" json:"-"`

	Timestamp
}

type Teacher struct {
	Id
	Subjects   []*Subject `gorm:"many2many:teacher_subjects" json:"subjects"`
	NUPTK      string     `gorm:"unique;not null"`
	EmployeeID string     `gorm:"not null" json:"employee_id"`

	UserID string `gorm:"uniqueIndex;not null" json:"-"`

	TimestampJoinTime
}

type Admin struct {
	Id
	StaffRole   string        `json:"staff_role"`
	Permissions []*Permission `gorm:"many2many:admin_permissions" json:"permissions,omitempty" swaggerignore:"true"`
	EmployeeID  string        `gorm:"not null" json:"employee_id"`

	UserID string `gorm:"uniqueIndex;not null" json:"-"`

	TimestampJoinTime
}

type Parent struct {
	Id
	FullName string     `json:"full_name"`
	Phone    string     `gorm:"unique" json:"phone"` // phone number
	Email    string     `gorm:"unique" json:"email"`
	Gender   UserGender `gorm:"type:user_gender" json:"gender"` // "male", "female"

	Timestamp
}
