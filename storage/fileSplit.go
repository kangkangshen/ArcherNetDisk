package storage

import (
	"github.com/kangkangshen/ArcherNetDisk/config"
	"github.com/kangkangshen/ArcherNetDisk/utils/bufs"
	"hash"
	"io"
	"net"
	"os"
)

type FileSplit struct {
	SplitSize   int64 //默认是与全局配置一致，增加此字段是为了以后方便更改
	Addr        string
	fromOffset  int64
	writen      int64
	total       int64
	hasher      hash.Hash
	file        *os.File //使用本地文件系统
	hash        []byte
	fileMeta    *FileMeta
	done        bool
	conn        net.Conn //当前或者上一个传输的conn
	options     map[string]string
	contentChan chan *bufs.Buffer
	writenChan  chan int
	errChan     chan error
}
type FileSplitMeta struct {

}
func (f *FileSplit) Hash() []byte {
	if (f.done) {
		return f.hash
	}
	return nil
}

func (f *FileSplit) Start() error {
	go f.readToContentChan()
	go f.consumeContent()
	err:= <-f.errChan
	go f.updateMetaInfo()
	return err

}

func (f *FileSplit) readToContentChan() {
	for {
		buffer := config.BufPool.Get()
		if _, err := io.Copy(buffer, f.conn); err != nil {
			close(f.contentChan)
		} else {
			f.contentChan <- buffer
		}
	}
}

//hash and store to local file system ,use pipeline design
func (f *FileSplit) consumeContent() {
	hContentChan := make(chan *bufs.Buffer, len(f.contentChan))
	var buffer *bufs.Buffer
	for {
		select {
		case buffer = <-f.contentChan:
			go func() {
				if cwriten, err := f.file.Write(buffer.Bytes()); err != nil {
					//There is an inexplicable error,
					// maybe the disk is being read and written by other programs?
					// Cancel this file partition task
					f.errChan <- err
					f.done=true
					close(f.errChan)
				} else {
					f.writenChan <- cwriten
					hContentChan <- buffer
				}
			}()
		case buffer = <-hContentChan:
			f.hasher.Write(buffer.Bytes())
		case cwriten:=<-f.writenChan:
			f.writen+=int64(cwriten)
			if f.writen==f.total{
				f.done=true
				close(f.errChan)
				return
			}
		}
	}
}

func (f *FileSplit) handleErrOrSuccess(){

}

func (f *FileSplit) updateMetaInfo(){}

func (f *FileSplit) Resume() {

}

func (f *FileSplit)
