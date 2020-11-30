package msg

import (
	"github.com/hslam/ftok"
	"strings"
	"testing"
	"time"
)

func TestMsg(t *testing.T) {
	context := strings.Repeat("1", 128)
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
		err = Snd(msgid, 1, []byte(context), 0600)
		if err != nil {
			panic(err)
		}
		time.Sleep(time.Millisecond * 200)
		close(done)
	}()
	time.Sleep(time.Millisecond * 100)
	key, err := ftok.Ftok("/tmp", 0x22)
	if err != nil {
		panic(err)
	}
	msgid, err := Get(key, 0600)
	if err != nil {
		panic(err)
	}
	text, err := Rcv(msgid, 1, 0600)
	if err != nil {
		panic(err)
	}
	if context != string(text) {
		t.Error(context, string(text))
	}
	<-done
}
