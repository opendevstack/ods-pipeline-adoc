package main

import (
	"fmt"
	"strings"
)

type multiFlag []string

func (i *multiFlag) String() string {
	return fmt.Sprintf("%v", *i)
}

func (i *multiFlag) Set(value string) error {
	*i = append(*i, strings.Split(value, ",")...)
	return nil
}
