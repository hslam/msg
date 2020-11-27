// Copyright (c) 2020 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

package msg

import (
	"syscall"
	"unsafe"
)

const (
	// IPC_CREAT creates if key is nonexistent
	IPC_CREAT = 00001000

	// IPC_RMID removes identifier
	IPC_RMID = 0

	maxText = 512
)

type message struct {
	Type uint
	flag uint8
	Text [maxText]byte
}

// Msg represents a message.
type Msg struct {
	Type uint
	Text []byte
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
func Snd(msgid uintptr, msg *Msg, flags uint) error {
	offset := 0
	m := message{Type: msg.Type}
	for len(msg.Text)-offset > 0 {
		if len(msg.Text)-offset > maxText {
			m.flag = 1
			copy(m.Text[:], msg.Text[offset:offset+maxText])
			_, _, err := syscall.Syscall6(syscall.SYS_MSGSND, msgid, uintptr(unsafe.Pointer(&m)), 1+maxText, uintptr(flags), 0, 0)
			if err != 0 {
				return err
			}
			offset += maxText
		} else {
			m.flag = 0
			copy(m.Text[:], msg.Text[offset:len(msg.Text)])
			_, _, err := syscall.Syscall6(syscall.SYS_MSGSND, msgid, uintptr(unsafe.Pointer(&m)), 1+uintptr(len(msg.Text)-offset), uintptr(flags), 0, 0)
			if err != 0 {
				return err
			}
			offset += len(msg.Text) - offset
		}
	}
	return nil
}

// Rcv calls the msgrcv system call.
func Rcv(msgid uintptr, msg *Msg, flags uint) error {
	m := message{Type: msg.Type}
	for {
		length, _, err := syscall.Syscall6(syscall.SYS_MSGRCV, msgid, uintptr(unsafe.Pointer(&m)), 1+maxText, uintptr(msg.Type), uintptr(flags), 0)
		if err != 0 {
			return err
		}
		msg.Type = m.Type
		msg.Text = append(msg.Text, m.Text[:length-1]...)
		if m.flag == 0 {
			break
		}
	}
	return nil
}

// Remove removes the message queue with the given id.
func Remove(msgid uintptr) error {
	r1, _, errno := syscall.Syscall(syscall.SYS_MSGCTL, msgid, IPC_RMID, 0)
	if int(r1) < 0 {
		return syscall.Errno(errno)
	}
	return nil
}
