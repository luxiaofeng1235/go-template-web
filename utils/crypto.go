package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
)

// HashPassword 密码哈希（使用MD5，参考go-novel）
func HashPassword(password string) (string, error) {
	hash := md5.New()
	hash.Write([]byte(password))
	return hex.EncodeToString(hash.Sum(nil)), nil
}


// AesEncryptByCFB AES CFB模式加密（参考go-novel的实现）
func AesEncryptByCFB(key, plaintext string) (string, error) {
	// 使用MD5处理密钥，确保长度为32字节
	keyHash := md5.Sum([]byte(key))
	keyBytes := keyHash[:]
	
	// 创建AES cipher
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}
	
	// 转换明文为字节
	plaintextBytes := []byte(plaintext)
	
	// 创建加密后的字节切片，长度为块大小+明文长度
	ciphertext := make([]byte, aes.BlockSize+len(plaintextBytes))
	
	// 生成随机IV
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}
	
	// CFB模式加密
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintextBytes)
	
	// 返回base64编码的结果
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// AesDecryptByCFB AES CFB模式解密
func AesDecryptByCFB(key, ciphertext string) (string, error) {
	// 使用MD5处理密钥
	keyHash := md5.Sum([]byte(key))
	keyBytes := keyHash[:]
	
	// Base64解码
	ciphertextBytes, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}
	
	// 创建AES cipher
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}
	
	// 检查密文长度
	if len(ciphertextBytes) < aes.BlockSize {
		return "", errors.New("密文长度不足")
	}
	
	// 提取IV和密文
	iv := ciphertextBytes[:aes.BlockSize]
	ciphertextBytes = ciphertextBytes[aes.BlockSize:]
	
	// CFB模式解密
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertextBytes, ciphertextBytes)
	
	return string(ciphertextBytes), nil
}

// GenerateToken 生成简单的token（实际项目中应使用JWT）
func GenerateToken(data map[string]interface{}) (string, error) {
	// 简单示例，实际项目中应该使用JWT
	tokenData := fmt.Sprintf("%v_%d", data["username"], data["user_id"])
	
	// 使用AES加密token
	encrypted, err := AesEncryptByCFB("token-secret-key", tokenData)
	if err != nil {
		return "", err
	}
	
	return encrypted, nil
}

// ValidateToken 验证token
func ValidateToken(token string) (map[string]interface{}, error) {
	if token == "" {
		return nil, errors.New("token为空")
	}
	
	// 解密token
	decrypted, err := AesDecryptByCFB("token-secret-key", token)
	if err != nil {
		return nil, errors.New("token无效")
	}
	
	// 这里应该解析token内容
	// 简单示例
	return map[string]interface{}{
		"token_data": decrypted,
	}, nil
}