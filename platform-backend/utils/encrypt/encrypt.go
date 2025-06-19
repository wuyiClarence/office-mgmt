package encrypt

// XOR 加密/解密函数，使用字符串作为密钥
func XorEncryptDecrypt(input, key string) []byte {
	output := make([]byte, len(input))
	keyLen := len(key)
	for i := 0; i < len(input); i++ {
		output[i] = input[i] ^ key[i%keyLen] // 循环使用密钥的每个字符
	}
	return output
}
