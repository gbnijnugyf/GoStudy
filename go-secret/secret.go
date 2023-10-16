package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

//var privateKey *rsa.PrivateKey
//var publicKeyBytes []byte

// RSA加密
func RSA_Encrypt(plainText []byte, path string) ([]byte, string) {
	//打开文件
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	//读取文件的内容
	info, _ := file.Stat()
	buf := make([]byte, info.Size())
	file.Read(buf)
	//pem解码
	block, _ := pem.Decode(buf)
	//x509解码
	publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		panic(err)
	}

	//类型断言
	publicKey := publicKeyInterface.(*rsa.PublicKey)
	publicKeyStr, _ := os.ReadFile(path)

	if string(plainText) == "" {
		return []byte(""), string(publicKeyStr)
	}
	//对明文进行加密
	cipherText, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, plainText)
	if err != nil {
		panic(err)
	}
	cpiherText_ := base64.StdEncoding.EncodeToString(cipherText)
	//返回密文
	return []byte(cpiherText_), string(publicKeyStr)
}

// RSA解密
func RSA_Decrypt(cipherText []byte, path string) ([]byte, *rsa.PrivateKey) {
	//打开文件
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	//获取文件内容
	info, _ := file.Stat()
	buf := make([]byte, info.Size())
	file.Read(buf)
	//pem解码
	block, _ := pem.Decode(buf)
	//X509解码
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		panic(err)
	}
	if string(cipherText) == "" {
		return []byte(""), privateKey
	}
	//对密文进行解密
	code, _ := base64.StdEncoding.DecodeString(string(cipherText)) //网上找的

	plainText, _ := rsa.DecryptPKCS1v15(rand.Reader, privateKey, code)
	//返回明文
	return plainText, privateKey
}

// 生成RSA私钥和公钥，保存到文件中
func GenerateRSAKey(bits int) *rsa.PrivateKey {
	//GenerateKey函数使用随机数据生成器random生成一对具有指定字位数的RSA密钥
	//Reader是一个全局、共享的密码用强随机数生成器
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		panic(err)
	}
	//保存私钥
	//通过x509标准将得到的ras私钥序列化为ASN.1 的 DER编码字符串
	X509PrivateKey := x509.MarshalPKCS1PrivateKey(privateKey)
	//使用pem格式对x509输出的内容进行编码
	//创建文件保存私钥
	privateFile, err := os.Create("private.pem")
	if err != nil {
		panic(err)
	}
	defer privateFile.Close()
	//构建一个pem.Block结构体对象
	privateBlock := pem.Block{Type: "RSA Private Key", Bytes: X509PrivateKey}
	//将数据保存到文件
	pem.Encode(privateFile, &privateBlock)
	//保存公钥
	//获取公钥的数据
	publicKey := privateKey.PublicKey
	//X509对公钥编码
	X509PublicKey, err := x509.MarshalPKIXPublicKey(&publicKey)
	if err != nil {
		panic(err)
	}
	//pem格式编码
	//创建用于保存公钥的文件
	publicFile, err := os.Create("public.pem")
	if err != nil {
		panic(err)
	}
	defer publicFile.Close()
	//创建一个pem.Block结构体对象
	publicBlock := pem.Block{Type: "RSA Public Key", Bytes: X509PublicKey}
	//保存到文件
	pem.Encode(publicFile, &publicBlock)
	return privateKey
}
func main() {
	// 生成公私钥对
	//var err error
	//privateKey, err = rsa.GenerateKey(rand.Reader, 1024)
	//if err != nil {
	//	fmt.Println("Failed to generate private key:", err)
	//	return
	//}
	//
	////保存私钥
	////通过x509标准将得到的ras私钥序列化为ASN.1 的 DER编码字符串
	//X509PrivateKey := x509.MarshalPKCS1PrivateKey(privateKey)
	////使用pem格式对x509输出的内容进行编码
	////创建文件保存私钥
	//privateFile, err := os.Create("private_key.pem")
	//if err != nil {
	//	panic(err)
	//}
	//defer privateFile.Close()
	////构建一个pem.Block结构体对象
	//privateBlock := &pem.Block{Type: "RSA Private Key", Bytes: X509PrivateKey}
	////将数据保存到文件
	//pem.Encode(privateFile, privateBlock)
	//
	//// 保存公钥
	//publicKeyBytes = x509.MarshalPKCS1PublicKey(&privateKey.PublicKey)
	//publicKeyFile, err := os.Create("public_key.pem")
	//if err != nil {
	//	fmt.Println("Failed to create public key file:", err)
	//	return
	//}
	//defer publicKeyFile.Close()
	//publicKeyBlock := &pem.Block{
	//	Type:  "RSA PUBLIC KEY",
	//	Bytes: publicKeyBytes,
	//}
	//err = pem.Encode(publicKeyFile, publicKeyBlock)
	//if err != nil {
	//	fmt.Println("Failed to save public key:", err)
	//	return
	//}

	GenerateRSAKey(2048)
	//message := []byte("1234567")
	////加密
	//cipherText, _ := RSA_Encrypt(message, "public.pem")
	//fmt.Println("加密后为：")
	//fmt.Println(string(cipherText))
	////解密
	//cipherText = []byte("vr1aJ2LCIvhERMV5HZgQMiT0PBZwNHlZk0+torIB9kzGAjcFIxw4DC2HAKWFDXMavKnL2efbw9OwIOdaRDzzSjp2NbsHndr5Hs27aiSWgFe55sOov6OHe4CoEy1wKRTwco3FpPi9SLufqLmA7OG9M+I63ZpIcIT7/2qf+uW7eqnxXi7PqrE1UbVoZ+pQuWSK+SihvGbgpAzLvx43PDMzWdxCx/vQSNj3SYAjdQljoBAALjvRHDYhHf1avrQndsRLFg5sa6lbMYqjEDbgE2K6jCprhGyS1tGPvU6dL/3qXljEU4kEjQJw8FjdlZtJGJJdft8aB+72VtFy2h5eF+kSgQ==")
	//plainText, _ := RSA_Decrypt(cipherText, "private.pem")
	//fmt.Println("解密后为：", string(plainText))
	// 使用gin框架创建HTTP服务器
	router := gin.Default()

	// 前端请求获取公钥的接口
	router.GET("/public-key", func(c *gin.Context) {
		_, publicKey := RSA_Encrypt([]byte(""), "public.pem")
		c.JSON(http.StatusOK, gin.H{
			"msg":    "successful",
			"status": 1,
			"data":   publicKey,
		})
	})
	type login struct {
		UserName string `json:"userName"`
		PassWord string `json:"passWord"`
	}
	// 后端解密请求的接口
	router.POST("/decrypt", func(c *gin.Context) {

		var jsons login
		err := c.ShouldBind(&jsons)
		if err != nil {
			fmt.Println(err)
		}

		// 使用私钥解密数据
		plainUserName, _ := RSA_Decrypt([]byte(jsons.UserName), "private.pem")
		plainPassWord, _ := RSA_Decrypt([]byte(jsons.PassWord), "private.pem")

		fmt.Println("userName:", string(plainUserName))
		fmt.Println("passWord:", string(plainPassWord))

		c.String(http.StatusOK, "ok")
	})

	// 启动HTTP服务器
	router.Run(":8080")
}
