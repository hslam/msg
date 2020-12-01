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
	IPC_CREAT = 01000

	// IPC_EXCL fails if key exists.
	IPC_EXCL = 02000

	// IPC_NOWAIT returns error no wait.
	IPC_NOWAIT = 04000

	// IPC_PRIVATE is private key
	IPC_PRIVATE = 00000

	// SEM_UNDO sets up adjust on exit entry
	SEM_UNDO = 010000

	// IPC_RMID removes identifier
	IPC_RMID = 0
	// IPC_SET sets ipc_perm options.
	IPC_SET = 1
	// IPC_STAT gets ipc_perm options.
	IPC_STAT = 2

	maxText = 8192
)

// ErrTooLong is returned when the Text length is bigger than maxText.
var ErrTooLong = errors.New("Text length is too long")

type message struct {
	Type uint
	Text [maxText]byte
}

// Get calls the msgget system call.
func Get(key int, msgflg int) (int, error) {
	r1, _, err := syscall.Syscall(syscall.SYS_MSGGET, uintptr(key), uintptr(msgflg), 0)
	msgid := int(r1)
	if msgid < 0 {
		return msgid, syscall.Errno(err)
	}
	return msgid, nil
}

// Snd calls the msgsnd system call.
//
// The msgsnd() and msgrcv() system calls are used to send messages to,
// and receive messages from, a System V message queue.  The calling
// process must have write permission on the message queue in order to
// send a message, and read permission to receive a message.
// The msgp argument is a pointer to a caller-defined structure of the
// following general form:
//
// struct msgbuf {
// 	long mtype;       /* message type, must be > 0 */
// 	char mtext[1];    /* message data */
// };
// The mtext field is an array (or other structure) whose size is speci‐
// fied by msgsz, a nonnegative integer value.  Messages of zero length
// (i.e., no mtext field) are permitted.  The mtype field must have a
// strictly positive integer value.  This value can be used by the re‐
// ceiving process for message selection (see the description of ms‐
// grcv() below).
//
func Snd(msgid int, msgp uintptr, msgsz int, msgflg int) error {
	_, _, err := syscall.Syscall6(syscall.SYS_MSGSND, uintptr(msgid), uintptr(msgp), uintptr(msgsz), uintptr(msgflg), 0, 0)
	if err != 0 {
		return err
	}
	return nil
}

// Rcv calls the msgrcv system call.
func Rcv(msgid int, msgp uintptr, msgsz int, msgtyp uint, msgflg int) (int, error) {
	r1, _, err := syscall.Syscall6(syscall.SYS_MSGRCV, uintptr(msgid), msgp, uintptr(msgsz), uintptr(msgtyp), uintptr(msgflg), 0)
	length := int(r1)
	if err != 0 {
		return length, err
	}
	return length, nil
}

// Send calls the msgsnd system call.
func Send(msgid int, msgType uint, msgText []byte, flags int) error {
	if len(msgText) > maxText {
		return ErrTooLong
	}
	m := message{Type: msgType}
	copy(m.Text[:], msgText)
	return Snd(msgid, uintptr(unsafe.Pointer(&m)), len(msgText), flags)
}

// Receive calls the msgrcv system call.
func Receive(msgid int, msgType uint, flags int) ([]byte, error) {
	m := message{Type: msgType}
	length, err := Rcv(msgid, uintptr(unsafe.Pointer(&m)), maxText, msgType, flags)
	if err != nil {
		return nil, err
	}
	return m.Text[:length], nil
}

// Remove removes the message queue with the given id.
func Remove(msgid int) error {
	r1, _, errno := syscall.Syscall(syscall.SYS_MSGCTL, uintptr(msgid), IPC_RMID, 0)
	if int(r1) < 0 {
		return syscall.Errno(errno)
	}
	return nil
}
