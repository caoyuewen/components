package util

import "strings"

// GetArrByString 根据逗号 分割字符串 返回数组
func GetArrByString(subs string) []string {
	var res []string
	if subs == "" {
		return res
	}
	// 去除首位 逗号 追加 下级到末尾
	if len(subs) >= 1 && subs[0] == ',' {
		subs = subs[1:]
	}
	if len(subs) >= 1 && subs[len(subs)-1] == ',' {
		subs = subs[:len(subs)-1]
	}
	if subs == "" {
		return res
	}
	res = strings.Split(subs, ",")
	return res
}

func FmtSubordinatesStr(subs string) string {
	if subs == "" {
		return subs
	}
	// 去除首位 逗号 追加 下级到末尾
	if len(subs) >= 1 && subs[0] == ',' {
		subs = subs[1:]
	}
	if len(subs) >= 1 && subs[len(subs)-1] == ',' {
		subs = subs[:len(subs)-1]
	}
	return subs
}

func SubordinatesStrToArr(subs string) []string {
	return GetArrByString(subs)
}

// MergeSubordinates 合并下级代理
func MergeSubordinates(subs, sub string) string {
	if subs == "" {
		return sub
	}
	if strings.Contains(subs, sub) {
		return subs
	}

	subsArr := StringToArrayStr(subs)
	realSubs := append(subsArr, sub)
	res := RemoveDuplicatesStr(realSubs)
	return ArrayToStringStr(res)

}

func RemoveDuplicatesStr(arr []string) []string {
	uniqueMap := make(map[string]bool)
	var uniqueArr []string

	for _, s := range arr {
		if !uniqueMap[s] {
			uniqueMap[s] = true
			uniqueArr = append(uniqueArr, s)
		}
	}

	return uniqueArr
}

func StringToArrayStr(str string) []string {
	strArr := strings.Split(str, ",")
	// 处理首位可能存在多余逗号的情况
	var nonEmptyStrArr []string
	for _, s := range strArr {
		if s != "" {
			nonEmptyStrArr = append(nonEmptyStrArr, s)
		}
	}
	return nonEmptyStrArr
}

func ArrayToStringStr(arr []string) string {
	return strings.Join(arr, ",")
}
