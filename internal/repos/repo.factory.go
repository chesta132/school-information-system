package repos

import "gorm.io/gorm"

type Repos struct {
	db         *gorm.DB
	user       *User
	student    *Student
	teacher    *Teacher
	admin      *Admin
	parent     *Parent
	class      *Class
	permission *Permission
	revoked    *Revoked
	subject    *Subject
}

func New(db *gorm.DB) *Repos {
	return &Repos{db: db}
}

func (r *Repos) User() *User {
	if r.user == nil {
		r.user = NewUser(r.db)
	}
	return r.user
}

func (r *Repos) Student() *Student {
	if r.student == nil {
		r.student = NewStudent(r.db)
	}
	return r.student
}

func (r *Repos) Teacher() *Teacher {
	if r.teacher == nil {
		r.teacher = NewTeacher(r.db)
	}
	return r.teacher
}

func (r *Repos) Admin() *Admin {
	if r.admin == nil {
		r.admin = NewAdmin(r.db)
	}
	return r.admin
}

func (r *Repos) Parent() *Parent {
	if r.parent == nil {
		r.parent = NewParent(r.db)
	}
	return r.parent
}

func (r *Repos) Class() *Class {
	if r.class == nil {
		r.class = NewClass(r.db)
	}
	return r.class
}

func (r *Repos) Permission() *Permission {
	if r.permission == nil {
		r.permission = NewPermission(r.db)
	}
	return r.permission
}

func (r *Repos) Revoked() *Revoked {
	if r.revoked == nil {
		r.revoked = NewRevoked(r.db)
	}
	return r.revoked
}

func (r *Repos) Subject() *Subject {
	if r.subject == nil {
		r.subject = NewSubject(r.db)
	}
	return r.subject
}
