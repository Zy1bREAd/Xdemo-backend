package app

import "fmt"

type Resourcer interface {
	List() error
	Add() error
	Del() error
	Update() error
}

func ListResources(rs Resourcer) {
	fmt.Println("打印出我的资源列表")
	fmt.Printf("这是一个%v资源\n", rs)
}
