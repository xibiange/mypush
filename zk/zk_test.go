package zk

import (
	"fmt"
	log "github.com/xibiange/mypush/log"
	zk "github.com/xibiange/mypush/zk"
	"testing"
	"time"
)

func TestZK(t *testing.T) {
	addr := []string{"127.0.0.1:2181"}
	conn, err := zk.Connect(addr, time.Second*10)
	if err != nil {
		t.Error(err)
		log.Error("zk connect \"%s\" errors (%v)", addr, err)
		fmt.Println("haha")
	}
	zk.Create(conn, "/mytest/node1", "assd")

	defer conn.Close()

}
