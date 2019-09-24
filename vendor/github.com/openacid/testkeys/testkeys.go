package testkeys

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
)

var Keys = []string{
	"10ll16k",
	"10vl5",
	"11vl5",
	"1mvl5_10",
	"20kl10",
	"20kvl10",
	"300vl50",
	"50kl10",
	"50kvl10",
	"empty",
}

func Load(dir, fn string) []string {

	p := filepath.Join(dir, fn)
	fmt.Println("load from:", p)

	f, err := os.Open(p)
	if err != nil {
		panic("fail to open: " + p + ":" + err.Error())
	}
	defer f.Close()

	rst := []string{}

	s := bufio.NewScanner(f)
	for s.Scan() {
		str := s.Text()
		rst = append(rst, str)
	}

	err = s.Err()
	if err != nil {
		panic("scan: " + p + ":" + err.Error())
	}

	return rst
}
