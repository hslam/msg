// Copyright (c) 2020 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

package msg

import (
	"errors"
	"syscall"
	"unsafe"
)

const (
	// IPC_CREAT creates if key is nonexistent
	IPC_CREAT = 00001000

	// IPC_RMID removes identifier
	IPC_RMID = 0

	maxText = 8192
)

// ErrTooLong is returned when the Text length is bigger than maxText.
var ErrTooLong = errors.New("Text length is too long")

type message struct {
	Type uint
	Text [maxText]byte
}

// Get calls the msgget system call.
func Get(key int, msgflg int) (uintptr, error) {
	msgid, _, err := syscall.Syscall(syscall.SYS_MSGGET, uintptr(key), uintptr(msgflg), 0)
	if int(msgid) < 0 {
		return 0, syscall.Errno(err)
	}
	return msgid, nil
}

// Snd calls the msgsnd system call.
func Snd(msgid uintptr, msgType uint, msgText []byte, flags uint) error {
	if len(msgText) > maxText {
		return ErrTooLong
	}
	m := message{Type: msgType}
	copy(m.Text[:], msgText)
	_, _, err := syscall.Syscall6(syscall.SYS_MSGSND, msgid, uintptr(unsafe.Pointer(&m)), uintptr(len(msgText)), uintptr(flags), 0, 0)
	if err != 0 {
		return err
	}
	return nil
}

// Rcv calls the msgrcv system call.
func Rcv(msgid uintptr, msgType uint, flags uint) ([]byte, error) {
	m := message{Type: msgType}
	length, _, err := syscall.Syscall6(syscall.SYS_MSGRCV, msgid, uintptr(unsafe.Pointer(&m)), maxText, uintptr(msgType), uintptr(flags), 0)
	if err != 0 {
		return nil, err
	}
	return m.Text[:length], nil
}

// Remove removes the message queue with the given id.
func Remove(msgid uintptr) error {
	r1, _, errno := syscall.Syscall(syscall.SYS_MSGCTL, msgid, IPC_RMID, 0)
	if int(r1) < 0 {
		return syscall.Errno(errno)
	}
	return nil
}
