// HongYu project HongYu.go
package main

import (
	"HongYu/conf"
	"HongYu/models"
	"fmt"
)

func init() {
	conf.LoadConfig()
	models.Show()
}

func main() {
	fmt.Println("hello ")
}
