package providers

import "fmt"

type IHost interface {
	Deploy() bool
}

func (hp HostProvider) Add() bool {
	fmt.Println("host add")
	return true
}

func (hp HostProvider) List() bool {
	fmt.Println("host list")
	return true
}
