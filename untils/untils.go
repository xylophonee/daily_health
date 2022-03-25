package untils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"math/rand"

)

// BaseString 加密用的字符串
const BaseString = "ABCDEFGHJKMNPQRSTWXYZabcdefhijkmnprstwxyz2345678"

// RandStr 获取随机length长度的字符串
func RandStr(length int) string{
	b := make([]byte,length)
	for i := range b{
		b[i] = BaseString[rand.Intn(len(BaseString))]
	}
	return string(b)
}

//PKCS5Padding 用于AES padding 填充
func PKCS5Padding(cipherText []byte, blockSize int) []byte {
	padding := blockSize - len(cipherText)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(cipherText, padText...)
}

// AESCBCEncrypt  AES-128-CBC-PKCS5Padding 加密登录密码
func AESCBCEncrypt(plainText []byte,key []byte) []byte{
	//指定加密算法，返回一个AES算法的Block接口对象
	block,err:=aes.NewCipher(key)
	if err!=nil{
		panic(err)
	}
	//进行填充
	plainText=PKCS5Padding(plainText,16)
	//指定初始向量vi
	//指定分组模式，返回一个BlockMode接口对象
	iv := []byte(RandStr(16))
	blockMode:=cipher.NewCBCEncrypter(block,iv)
	cipherText:=make([]byte,len(plainText))
	blockMode.CryptBlocks(cipherText,plainText)
	//返回密文
	return cipherText
}
