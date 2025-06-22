package tictok

import (
	"crypto/md5"
	"encoding/base64"
	"sort"
	"strings"
)

/*
Signature

		签名方式：开发者请务必校验数据签名，验证数据来源的合法性，否则存在被伪造数据攻击的危险，需自行担责对表格中header参数(x-signature，content-type除外)，
		按key字典序从小到大排序, 排序后，将key-value按顺序连接起来。如：key1=value1&key2=value2。再直接拼接(无需连接符)上body字符串和secret（推送配置中的字段),
		也就是前置工作中需要开发者提供给开平的签名秘钥。
	 	注意，字符串需要使用utf-8编码把拼接好的字符串进行md5计算(16bytes)，并对md5计算结果进行base64编码，编码结果便是signature
		example:
		header := map[string]string{
		             "x-nonce-str": "123456",
		             "x-timestamp": "456789",
		             "x-roomid":    "268",
		             "x-msg-type":  "live_gift",
		     }
		bodyStr := "abc123你好"
		secret := "123abc"

rawData为：x-msg-type=live_gift&x-nonce-str=123456&x-roomid=268&x-timestamp=456789abc123你好123abc
signature为：PDcKhdlsrKEJif6uMKD2dw==
*/
func Signature(header map[string]string, bodyStr, secret string) string {
	keyList := make([]string, 0, 4)
	for key, _ := range header {
		keyList = append(keyList, key)
	}
	sort.Slice(keyList, func(i, j int) bool {
		return keyList[i] < keyList[j]
	})
	kvList := make([]string, 0, 4)
	for _, key := range keyList {
		kvList = append(kvList, key+"="+header[key])
	}
	urlParams := strings.Join(kvList, "&")
	rawData := urlParams + bodyStr + secret
	md5Result := md5.Sum([]byte(rawData))
	return base64.StdEncoding.EncodeToString(md5Result[:])
}
