package main

import (
	"fmt"

	"github.com/zhsso/funcMap"
)

type test struct {
}

func (a test) Test(aa string, bb int) (string, int) {
	fmt.Println(aa)
	fmt.Println(bb)
	panic(1111)
	return aa, bb
}

func main() {
	a := test{}
	f := funcMap.NewFuncMap()
	f.Register(a)
	ret, err := f.Invoke("Test", "a", 12)
	fmt.Println(ret)
	fmt.Println(err)
}
