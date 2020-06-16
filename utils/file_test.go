package utils

import (
	"fmt"
	"os"
	"testing"
)

func Test_ReadFile(t *testing.T) {
	pwd, _:= os.Getwd()
	strs, err := ReadFileForLine(pwd + "/aes.go")
	if err != nil{
		t.Fatal(err)
	}
	fmt.Println(strs)
}
