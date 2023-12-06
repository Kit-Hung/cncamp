package module1

import (
	"fmt"
	"testing"
)

func TestSlice(t *testing.T) {
	/*
		给定一个字符串数组
		[“I”,“am”,“stupid”,“and”,“weak”]
		用 for 循环遍历该数组并修改为
		[“I”,“am”,“smart”,“and”,“strong”]
	*/

	strs := []string{"I", "am", "stupid", "and", "weak"}
	fmt.Printf("before change, strs: %v\n", strs)
	for index, str := range strs {
		switch str {
		case "stupid":
			strs[index] = "smart"
		case "weak":
			strs[index] = "strong"
		}
	}
	fmt.Printf("after change, strs: %v\n", strs)
}
