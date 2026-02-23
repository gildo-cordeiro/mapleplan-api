package profile

import (
	"errors"
	"time"
)

type MemberRole string

const (
	RolePrimary MemberRole = "Primary"
	RoleSpouse  MemberRole = "Spouse"
)

func isValidRole(role MemberRole) bool {
	return role == RolePrimary || role == RoleSpouse
}

func (m *ProfileMember) ChangeRole(newRole MemberRole) error {
	now := time.Now()

	if !isValidRole(newRole) {
		return errors.New("invalid role")
	}
	m.Role = newRole
	m.UpdatedAt = now
	return nil
}
