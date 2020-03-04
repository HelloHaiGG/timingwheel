package main

import (
	"HelloMyWorld/common/icache/L1"
	"fmt"
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	if v, err := L1.Get("key1");err != nil{
		//t.Fatal(err)
		fmt.Println(err)
	}else{
		fmt.Println(v)
	}

	if err := L1.Put("key2", "haha");err!= nil{
		t.Fatal(err)
	}

	if v, err := L1.Get("key2");err != nil{
		fmt.Println(err)
	}else{
		fmt.Println(v)
	}

	time.Sleep(time.Second * 2)

	if v, err := L1.Get("key1");err != nil{
		fmt.Println(err)
	}else{
		fmt.Printf(v)
	}

}
