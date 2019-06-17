package conc

import (
	"github.com/etcd-io/etcd/clientv3"
	"testing"
	"time"
)

func TestNewDlocker(t *testing.T) {
	config := clientv3.Config{Endpoints: []string{"127.0.0.1:2381"}}
	if client, err := clientv3.New(config); err != nil {
		t.Fatal(err)
	} else {
		dLocker := NewDlocker(client, "wukang")
		if dLocker != nil {
			t.Log("dLocker created successfully")
		}
	}
}

func TestEtcdDLocker_Lock(t *testing.T) {
	config := clientv3.Config{Endpoints: []string{"127.0.0.1:2381"}}
	if client, err := clientv3.New(config); err != nil {
		t.Fatal(err)
	} else {
		dlocker := NewDlocker(client, "wukang3")
		dlocker.Lock()
		time.Sleep(time.Second * 10)
		//dlocker.Unlock()
		return
	}

}
