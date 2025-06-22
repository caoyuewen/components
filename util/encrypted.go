package util

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
)

var iv = []byte{64, 64, 64, 64, 38, 38, 38, 38, 35, 35, 35, 35, 36, 36, 36, 36}

func GetSha256String(str string) string {
	hash := sha256.New()
	hash.Write([]byte(str))
	bytes := hash.Sum(nil)
	return hex.EncodeToString(bytes)
}

func GetSha512String(str string) string {
	hash := sha512.New()
	hash.Write([]byte(str))
	bytes := hash.Sum(nil)
	return hex.EncodeToString(bytes)
}

func GetHmacSha256Base64(str, key string) string {
	hash := hmac.New(sha256.New, []byte(key))
	hash.Write([]byte(str))
	bytes := hash.Sum(nil)
	return base64.StdEncoding.EncodeToString(bytes)
}

// 加密数据
func AESEncrypt(plaintext, k string) (string, error) {
	data := []byte(plaintext)
	key := []byte(k)
	aesBlockEncrypter, err := aes.NewCipher(key)
	content := PKCS5Padding(data, aesBlockEncrypter.BlockSize())
	encrypted := make([]byte, len(content))
	if err != nil {
		println(err.Error())
		return "", err
	}
	aesEncrypter := cipher.NewCBCEncrypter(aesBlockEncrypter, iv)
	aesEncrypter.CryptBlocks(encrypted, content)
	return base64.StdEncoding.EncodeToString(encrypted), nil
}

// 解密数据
func AESDecrypt(chiper, k string) (string, error) {
	src, err := base64.StdEncoding.DecodeString(chiper)
	if err != nil {
		return "", err
	}
	key := []byte(k)
	decrypted := make([]byte, len(src))
	aesBlockDecrypter, err := aes.NewCipher(key)
	if err != nil {
		println(err.Error())
		return "", err
	}
	aesDecrypter := cipher.NewCBCDecrypter(aesBlockDecrypter, iv)
	aesDecrypter.CryptBlocks(decrypted, src)
	trimming := PKCS5Trimming(decrypted)
	return string(trimming), nil
}

/*
*
PKCS5包装
*/
func PKCS5Padding(cipherText []byte, blockSize int) []byte {
	padding := blockSize - len(cipherText)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(cipherText, padText...)
}

/*
解包装
*/
func PKCS5Trimming(encrypt []byte) []byte {
	if len(encrypt) == 0 {
		return encrypt
	}
	padding := encrypt[len(encrypt)-1]
	l := len(encrypt) - int(padding)
	if l >= len(encrypt) || l < 0 {
		return []byte{}
	}
	return encrypt[:l]
}

// ------------------ AES-zeropadding-ECB ------------------
// 加密
func AesEncryptECB(origData []byte, key []byte) (encrypted []byte) {
	cipher, _ := aes.NewCipher(generateKey(key))
	length := (len(origData) + aes.BlockSize) / aes.BlockSize
	plain := make([]byte, length*aes.BlockSize)
	copy(plain, origData)
	pad := byte(len(plain) - len(origData))
	for i := len(origData); i < len(plain); i++ {
		plain[i] = pad
	}
	encrypted = make([]byte, len(plain))
	// 分组分块加密
	for bs, be := 0, cipher.BlockSize(); bs <= len(origData); bs, be = bs+cipher.BlockSize(), be+cipher.BlockSize() {
		cipher.Encrypt(encrypted[bs:be], plain[bs:be])
	}

	return encrypted
}

// 解密
func AesDecryptECB(encrypted []byte, key []byte) (decrypted []byte) {
	cipher, _ := aes.NewCipher(generateKey(key))
	decrypted = make([]byte, len(encrypted))

	for bs, be := 0, cipher.BlockSize(); bs < len(encrypted); bs, be = bs+cipher.BlockSize(), be+cipher.BlockSize() {
		cipher.Decrypt(decrypted[bs:be], encrypted[bs:be])
	}

	trim := 0
	if len(decrypted) > 0 {
		trim = len(decrypted) - int(decrypted[len(decrypted)-1])
	}

	return decrypted[:trim]
}

func generateKey(key []byte) (genKey []byte) {
	genKey = make([]byte, 16)
	copy(genKey, key)
	for i := 16; i < len(key); {
		for j := 0; j < 16 && i < len(key); j, i = j+1, i+1 {
			genKey[j] ^= key[i]
		}
	}
	return genKey
}

// ------------------ AES-zeropadding-ECB ------------------
