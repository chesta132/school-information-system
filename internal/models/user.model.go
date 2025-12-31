package models

type User struct {
	Id
	FullName string     `json:"full_name" example:"Chesta Ardiona"`
	Email    string     `gorm:"uniqueIndex" json:"email" example:"chestaardi4@gmail.com"` // auth username
	Password string     `json:"-"`                                                        // auth password
	Role     UserRole   `gorm:"type:user_role;default:unsetted" json:"role"`              // "student", "teacher", "admin", "unsetted"
	Gender   UserGender `gorm:"type:user_gender" json:"gender"`                           // "male", "female"
	Phone    string     `gorm:"uniqueIndex" json:"phone" example:"+6281234567890"`        // phone number

	// empty if user's role not student
	StudentProfile *Student `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"student_profile,omitempty" swaggerignore:"true"`
	// empty if user's role not teacher
	TeacherProfile *Teacher `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"teacher_profile,omitempty" swaggerignore:"true"`
	// empty if user's role not admin
	AdminProfile *Admin `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"admin_profile,omitempty" swaggerignore:"true"`

	TimestampArchivable
}

type Student struct {
	Id
	ClassID string    `json:"-"`
	Class   *Class    `json:"class,omitempty" swaggerignore:"true"`
	Parents []*Parent `gorm:"many2many:student_parents" json:"parents,omitempty" swaggerignore:"true"`
	NISN    string    `gorm:"unique;not null" example:"0091913711"`

	UserID string `gorm:"uniqueIndex;not null" json:"-"`

	Timestamp
}

type Teacher struct {
	Id
	Subjects   []*Subject `gorm:"many2many:teacher_subjects" json:"subjects,omitempty" swaggerignore:"true"`
	NUPTK      string     `gorm:"unique;not null" json:"NUPTK" example:"1234567890123456"`
	EmployeeID string     `gorm:"not null" json:"employee_id" example:"TEA001"`

	UserID string `gorm:"uniqueIndex;not null" json:"-"`

	TimestampJoinTime
}

type Admin struct {
	Id
	StaffRole   string        `json:"staff_role" example:"developer"`
	Permissions []*Permission `gorm:"many2many:admin_permissions" json:"permissions,omitempty" swaggerignore:"true"`
	EmployeeID  string        `gorm:"not null" json:"employee_id" example:"DEV001"`

	UserID string `gorm:"uniqueIndex;not null" json:"-"`

	TimestampJoinTime
}

type Parent struct {
	Id
	FullName string     `json:"full_name" example:"Chesta Ardiona"`
	Phone    string     `gorm:"unique" json:"phone" example:"+6281234567890"` // phone number
	Email    string     `gorm:"unique" json:"email" example:"chestaardi4@gmail.com"`
	Gender   UserGender `gorm:"type:user_gender" json:"gender"` // "male", "female"

	Timestamp
}
