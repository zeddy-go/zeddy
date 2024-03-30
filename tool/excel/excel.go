package excel

import (
	"fmt"
	"strings"
)

// Cell 根据行列索引计算坐标 从1开始
func Cell(col int, row int) string {
	return fmt.Sprintf("%s%d", ColIndexByNum(col), row)
}

func ColIndexByNum(num int) string {
	s := make([]string, 0, num/26)
	for num != 0 {
		tmp := num % 26
		num /= 26

		//此处略微关键，当为0时，其实是26，也就是Z，
		//而且当你将0调整为26后，需要从数字中去除26代表的这个数
		if tmp == 0 {
			tmp = 26
			num -= 1
		}
		s = append(s, fmt.Sprintf("%c", 'A'+tmp-1))
	}

	for i := 0; i < len(s)/2; i++ {
		temp := s[i]
		s[i] = s[len(s)-1-i]
		s[len(s)-1-i] = temp
	}

	return strings.Join(s, "")
}
