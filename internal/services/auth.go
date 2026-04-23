package services

import "crypto/subtle"

// AuthService выполняет минимальную проверку учетных данных администратора.
type AuthService struct {
	adminUsername string
	adminPassword string
}

func NewAuthService(adminUsername, adminPassword string) *AuthService {
	return &AuthService{adminUsername: adminUsername, adminPassword: adminPassword}
}

func (s *AuthService) Authenticate(username, password string) bool {
	userMatch := subtle.ConstantTimeCompare([]byte(username), []byte(s.adminUsername)) == 1
	passMatch := subtle.ConstantTimeCompare([]byte(password), []byte(s.adminPassword)) == 1
	return userMatch && passMatch
}
