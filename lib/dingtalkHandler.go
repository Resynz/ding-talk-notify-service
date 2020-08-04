/**
 * @Author: Resynz
 * @Date: 2020/8/4 10:11
 */
package lib

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"
	"ding-talk-notify-service/structs"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"sort"
	"time"
)

var _rule = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

const _block_size = 32

func _PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func _PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func GetRandomBytes(n int) []byte {
	b := make([]byte, n)
	rc := len(_rule)
	for i := 0; i < n; i++ {
		b[i] = _rule[rand.Intn(rc)]
	}
	return b
}

func _msgToPad(msgLen int) []byte {
	paddingLen := _block_size - (msgLen % _block_size)
	if paddingLen == 0 {
		paddingLen = _block_size
	}
	padByte := byte(paddingLen)
	pad := make([]byte, paddingLen)
	for i := 0; i < paddingLen; i++ {
		pad[i] = padByte
	}
	return pad
}

type DingTalkHandler struct {
	aesKey      []byte
	config      *structs.DingTalkConfig
	registerMap structs.RegisterMap
}

func NewDingTalkHandler(conf *structs.DingTalkConfig) (*DingTalkHandler, error) {
	key, err := base64.StdEncoding.DecodeString(fmt.Sprintf("%s=", conf.AesKey))
	if err != nil {
		return nil, err
	}
	if len(key) != 32 {
		return nil, errors.New("invalid aes_key")
	}
	return &DingTalkHandler{
		aesKey:      key,
		config:      conf,
		registerMap: make(structs.RegisterMap, 0),
	}, nil
}

func (dh *DingTalkHandler) GenSignature(timestamp, nonce, encrypt string) string {
	sl := []string{dh.config.Token, timestamp, nonce, encrypt}
	sort.Strings(sl)
	h := sha1.New()
	for _, s := range sl {
		_, _ = io.WriteString(h, s)
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}

// 解密回调内容 返回解密后的明文string，和解密校验结果bool
func (dh *DingTalkHandler) DecryptMsg(encrypted string) ([]byte, bool) {
	// 1.base解码
	cryptedMsg, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		log.Printf("base64 decode failed! error:%v\n", err)
		return nil, false
	}

	key := dh.aesKey
	aesBlk, err := aes.NewCipher(key)
	if err != nil {
		log.Printf("ase NewCipher failed! error:%v\n", err)
		return nil, false
	}
	blockSize := aesBlk.BlockSize()
	iv := key[:blockSize]
	blockMode := cipher.NewCBCDecrypter(aesBlk, iv)
	plainText := make([]byte, len(cryptedMsg))
	blockMode.CryptBlocks(plainText, cryptedMsg)
	plainText = _PKCS5UnPadding(plainText)
	msgLen := binary.BigEndian.Uint32(plainText[16:20])
	return plainText[20 : msgLen+20], true
}

func (dh *DingTalkHandler) EncryptMsg(content []byte) (*structs.DingTalkResponse, error) {
	timestamp := fmt.Sprintf("%d", time.Now().UnixNano()/1e6)
	nonce := string(GetRandomBytes(8))
	randBytes := GetRandomBytes(16)
	msgLenInNet := make([]byte, 4)
	binary.BigEndian.PutUint32(msgLenInNet, uint32(len(content)))
	originData := bytes.Buffer{}
	originData.Write(randBytes)
	originData.Write(msgLenInNet)
	originData.Write(content)
	originData.WriteString(dh.config.Key)
	pad := _msgToPad(originData.Len())
	originData.Write(pad)
	key := dh.aesKey
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	iv := key[:block.BlockSize()]
	blockMode := cipher.NewCBCEncrypter(block, iv)
	crypted := make([]byte, originData.Len())
	blockMode.CryptBlocks(crypted, originData.Bytes())
	cryptedText := base64.StdEncoding.EncodeToString(crypted)
	return &structs.DingTalkResponse{
		MsgSignature: dh.GenSignature(timestamp, nonce, cryptedText),
		Encrypt:      cryptedText,
		TimeStamp:    timestamp,
		Nonce:        nonce,
	}, nil
}
