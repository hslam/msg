# msg
[![PkgGoDev](https://pkg.go.dev/badge/github.com/hslam/msg)](https://pkg.go.dev/github.com/hslam/msg)
[![Build Status](https://travis-ci.org/hslam/msg.svg?branch=master)](https://travis-ci.org/hslam/msg)
[![Go Report Card](https://goreportcard.com/badge/github.com/hslam/msg)](https://goreportcard.com/report/github.com/hslam/msg)
[![LICENSE](https://img.shields.io/github/license/hslam/msg.svg?style=flat-square)](https://github.com/hslam/msg/blob/master/LICENSE)

Package msg provides a way to use System V message queue.

## Get started

### Install
```
go get github.com/hslam/msg
```
### Import
```
import "github.com/hslam/msg"
```
### Usage
#### Example
msgsnd
```go
package main

import (
	"github.com/hslam/ftok"
	"github.com/hslam/msg"
	"time"
)

func main() {
	key, err := ftok.Ftok("/tmp", 0x22)
	if err != nil {
		panic(err)
	}
	msgid, err := msg.Get(key, msg.IPC_CREAT|0600)
	if err != nil {
		panic(err)
	}
	defer msg.Remove(msgid)
	m := &msg.Msg{Type: 1, Text: []byte("Hello World")}
	err = msg.Snd(msgid, m, 0600)
	if err != nil {
		panic(err)
	}
	time.Sleep(time.Second * 10)
}
```
msgrcv
```go
package main

import (
	"fmt"
	"github.com/hslam/ftok"
	"github.com/hslam/msg"
)

func main() {
	key, err := ftok.Ftok("/tmp", 0x22)
	if err != nil {
		panic(err)
	}
	msgid, err := msg.Get(key, 0600)
	if err != nil {
		panic(err)
	}
	m := &msg.Msg{Type: 1}
	err = msg.Rcv(msgid, m, 0600)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(m.Text))
}
```

#### Output
```
Hello World
```

### License
This package is licensed under a MIT license (Copyright (c) 2020 Meng Huang)


### Author
msg was written by Meng Huang.


