package utils

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

// 结论：当大量字符串拼接的时候，使用 buffer 做拼接性能能较优，当只有一次字符串拼接的时候，join 性能较优

//////////////////多次相加字符串的性能比较

//fmt
func Benchmark_FmtStr(b *testing.B) {
	str := ""
	for i := 0; i < b.N; i++ {
		str += fmt.Sprintf("%s%s", "Hello", "HaiGG")
	}
}

//join
func Benchmark_JoinStr(b *testing.B) {
	str := ""
	for i := 0; i < b.N; i++ {
		str += strings.Join([]string{"Hello", "HaiGG"}, "")
	}
}

// +

func Benchmark_AddStr(b *testing.B) {
	str := ""
	for i := 0; i < b.N; i++ {
		str += "Hello" + "HaiGG"
	}
}

// buffer

func Benchmark_BufferStr(b *testing.B) {
	buf := bytes.Buffer{}
	for i := 0; i < b.N; i++ {
		buf.WriteString("Hello")
		buf.WriteString("HaiGG")
	}
}

//////////////////////////////单次相加的字符串性能比较

//fmt
func Benchmark_FmtStr_2(b *testing.B) {
	str := ""
	for i := 0; i < b.N; i++ {
		str = fmt.Sprintf("%s%s", "Hello", "HaiGG")
	}
	fmt.Println(str)
}

//join
func Benchmark_JoinStr_2(b *testing.B) {
	str := ""
	for i := 0; i < b.N; i++ {
		str = strings.Join([]string{"Hello", "HaiGG", "xx", "yy", "ww", "xx", "zz", "ww"}, "")
	}
	fmt.Println(str)
}

// buffer
func Benchmark_BufferStr_2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		buf := bytes.Buffer{}
		buf.WriteString("Hello")
		buf.WriteString("HaiGG")
		buf.WriteString("xx")
		buf.WriteString("yy")
		buf.WriteString("ww")
		buf.WriteString("xx")
		buf.WriteString("zz")
		buf.WriteString("ww")
	}
}

// +

func Benchmark_AddStr_2(b *testing.B) {
	str := ""
	for i := 0; i < b.N; i++ {
		str = "Hello" + "HaiGG"
	}
	fmt.Println(str)
}
