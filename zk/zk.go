package zk

import (
	//log "code.google.com/p/log4go"
	"errors"
	"github.com/samuel/go-zookeeper/zk"
	log "github.com/xibiange/mypush/log"
	"os"
	"path"
	"strings"
	"syscall"
	"time"
)

var (
	ErrNoChild      = errors.New("zk:child is nil")
	ErrNodeNotExist = errors.New("zk: node not exist")
)

func Connect(addr []string, timeout time.Duration) (*zk.Conn, error) {
	conn, session, err := zk.Connect(addr, timeout)
	if err != nil {
		log.Error("zk connec (\"%v\",\"%d\") error (\"%v\")", addr, timeout, err)
	}
	go func() {
		for {
			event := <-session
			log.Debug("zookeeper get a event :%s\n", event.State.String())
		}
	}()
	return conn, nil
}

func Create(conn *zk.Conn, fpath, data string) error {
	tpath := ""
	for _, mypath := range strings.Split(fpath, "/")[1:] {
		tpath = path.Join(tpath, "/", mypath)
		//节点, 数据,节点类型|永久,权限|所有
		_, err := conn.Create(tpath, []byte(""), 0, zk.WorldACL(zk.PermAll))
		if err != nil {
			if err == zk.ErrNodeExists {
				log.Warn("zk create node exists (\"%s\")\n", tpath)
			} else {
				log.Error("zk create (\"%s\") error (%v)\n", tpath, err)
				return err
			}

		}
	}
	if data != "" {
		//节点 数据,版本|跳过版本判断
		_, err := conn.Set(tpath, []byte(data), -1)
		if err != nil {
			log.Error("zk set Data (\"%s\",\"%s\",%v)", tpath, data, err)
		}
	}
	return nil
}

func RegisterTemp(conn *zk.Conn, fpath, data string) error {
	tpath, err := conn.Create(path.Join(fpath)+"/", []byte(data), zk.FlagEphemeral|zk.FlagSequence, zk.WorldACL(zk.PermAll))
	if err != nil {
		log.Error("conn Create(\"%s\",\"%s\") error (%v)", fpath, data, err)
	}
	log.Debug("create a zookeeper node :%s", tpath)
	go func() {
		for {
			log.Info("zk path \"%s\" set a watch ", tpath)
			exist, _, watch, err := conn.ExistsW(tpath)
			if err != nil {
				log.Error("zk path \"%s\" watch err,errors (%v)", tpath, err)
				KillSelf()
				return
			}
			if !exist {
				log.Warn("zk path :\"%s\" not exists ,kill itself", tpath)
				KillSelf()
				return
			}
			event := <-watch
			log.Info("zk path \"%s\" receive a event %v", tpath, event)
		}
	}()
	return nil
}

func KillSelf() {
	if err := syscall.Kill(os.Getpid(), syscall.SIGQUIT); err != nil {
		log.Error("syscall.kill (%d,SIQUIT) error (%v)", os.Getpid(), err)
	}
}
