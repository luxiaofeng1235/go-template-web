package utils

// CheckPassword 验证密码
func CheckPassword(password, hashed string) bool {
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return false
	}
	return hashedPassword == hashed
}