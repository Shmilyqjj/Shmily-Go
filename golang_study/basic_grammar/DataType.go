package main

import (
	"fmt"
	"strconv"
)

/*
1	布尔型 布尔型的值只可以是常量 true 或者 false。一个简单的例子：var b bool = true。
2	数字类型 整型 int 和浮点型 float32、float64，Go 语言支持整型和浮点型数字，并且支持复数，其中位的运算采用补码。
3	字符串类型: 字符串就是一串固定长度的字符连接起来的字符序列。Go 的字符串是由单个字节连接起来的。Go 语言的字符串的字节使用 UTF-8 编码标识 Unicode 文本。
4	派生类型:包括：
(a) 指针类型（Pointer）
(b) 数组类型
(c) 结构化类型(struct)
(d) Channel 类型
(e) 函数类型
(f) 切片类型
(g) 接口类型（interface）
(h) Map 类型

数字类型
1	uint8 无符号 8 位整型 (0 到 255)
2	uint16 无符号 16 位整型 (0 到 65535)
3	uint32 无符号 32 位整型 (0 到 4294967295)
4	uint64 无符号 64 位整型 (0 到 18446744073709551615)
5	int8 有符号 8 位整型 (-128 到 127)
6	int16 有符号 16 位整型 (-32768 到 32767)
7	int32 有符号 32 位整型 (-2147483648 到 2147483647)
8	int64 有符号 64 位整型 (-9223372036854775808 到 9223372036854775807)

浮点型
1	float32 IEEE-754 32位浮点型数
2	float64 IEEE-754 64位浮点型数
3	complex64 32位实数和虚数
4	complex128 64 位实数和虚数

其他数字类型
1	byte 类似 uint8
2	rune 类似 int32
3	uint 32 或 64 位
4	int 与 uint 一样大小
5	uintptr 无符号整型，用于存放一个指针
*/

func main() {
	goBoolean()
	defaultValue()
	variableDeclaration()
	compute()
	str2int()
}

func defaultValue() {
	var i int     // 不初始化默认0
	var f float64 // 不初始化默认0
	var b bool    // 不初始化默认false
	var s string  // 不初始化默认空串
	fmt.Printf("%v %v %v %q\n", i, f, b, s)
}

func goBoolean() {
	var b bool = true
	fmt.Println(b)
}

func variableDeclaration() {
	intVal := 1 // := 为声明变量的符号  声明intVal值为1  如果变量已经用var声明了则不能用:=声明，而是直接用=赋值
	fmt.Println(intVal)

	//多变量声明
	var v1, v2, v3 int
	v1, v2, v3 = 1, 2, 3
	fmt.Printf("%d %d %d \n", v1, v2, v3)
	v4, v5, v6 := 4, 5, 6
	fmt.Printf("%d %d %d \n", v4, v5, v6)

	//全局变量声明
	var (
		vg1 int
		vg2 bool
	)
	fmt.Printf("%d %v \n", vg1, vg2)
}

func compute()  {
	int1 := 1
	int2 := 3
	fmt.Println(float32(int1)/float32(int2))

	// 计算百分比 保留三位小数
	f, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", float32(int1)/float32(int2)*100), 64)
	fmt.Println(f)

}

func str2int() {
	v := "1111"
	atoi, err := strconv.Atoi(v)
	if err != nil {
		return
	}
	fmt.Println(atoi)
}
