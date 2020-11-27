package msg

import (
	"github.com/hslam/ftok"
	"strings"
	"testing"
	"time"
)

func TestMsg(t *testing.T) {
	context := strings.Repeat("1", maxText+1)
	done := make(chan struct{})
	go func() {
		key, err := ftok.Ftok("/tmp", 0x22)
		if err != nil {
			panic(err)
		}
		msgid, err := Get(key, IPC_CREAT|0600)
		if err != nil {
			panic(err)
		}
		defer Remove(msgid)
		m := &Msg{Type: 1, Text: []byte(context)}
		err = Snd(msgid, m, 0600)
		if err != nil {
			panic(err)
		}
		time.Sleep(time.Second * 2)
		close(done)
	}()
	time.Sleep(time.Second)
	key, err := ftok.Ftok("/tmp", 0x22)
	if err != nil {
		panic(err)
	}
	msgid, err := Get(key, 0600)
	if err != nil {
		panic(err)
	}
	m := &Msg{Type: 1}
	err = Rcv(msgid, m, 0600)
	if err != nil {
		panic(err)
	}
	if context != string(m.Text) {
		t.Error(context, string(m.Text))
	}
	<-done
}
