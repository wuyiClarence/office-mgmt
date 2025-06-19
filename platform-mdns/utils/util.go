package utils

import (
	"encoding/json"
	"platform-mdns/dto"
)

// ToJSONString 将结构体转换为 JSON 字符串
func ToJSONString(info dto.MdnsInfo) (string, error) {
	jsonData, err := json.Marshal(info)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

// XorEncryptDecrypt XOR 加密/解密函数，使用字符串作为密钥
func XorEncryptDecrypt(input, key string) []byte {
	output := make([]byte, len(input))
	keyLen := len(key)
	for i := 0; i < len(input); i++ {
		output[i] = input[i] ^ key[i%keyLen] // 循环使用密钥的每个字符
	}
	return output
}
