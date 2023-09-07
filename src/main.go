package main

import (
	"fmt"
	"log"
	"net"
	"sync"
	"syscall"
)

type Pool struct {
	mu   sync.RWMutex
	work chan func()
	cls  chan struct{}
}

func NewPool(poolsize int) *Pool {
	return &Pool{
		work: make(chan func()),
		cls:  make(chan struct{}, poolsize),
	}
}

func (p *Pool) Schedule(task func()) {
	select {
	case p.work <- task:
	case p.cls <- struct{}{}:
		go p.worker(task)
	}

}

func (p *Pool) worker(task func()) {
	defer func() { <-p.cls }()
	for {
		task()
		task = <-p.work
	}
}

func main() {

	ls, err := net.Listen("tcp", ":5050")
	if err != nil {
		log.Fatalf("Listen err: %+v", err)
	}
	defer ls.Close()

	listenerFile, err := ls.(*net.TCPListener).File()
	if err != nil {
		log.Fatalf("File returns a copy of the underlying os.File. err: %+v", err)
	}
	listenerFd := int(listenerFile.Fd())

	fd, err := syscall.EpollCreate(syscall.EPOLL_NONBLOCK)
	if err != nil {
		log.Fatalf("File description creator err: %+v", err)
	}
	defer syscall.Close(fd)

	event := syscall.EpollEvent{
		Events: syscall.EPOLLIN,
		Fd:     int32(listenerFd),
	}

	err = syscall.EpollCtl(fd, syscall.EPOLL_CTL_ADD, listenerFd, &event)
	if err != nil {
		log.Fatalf("Epoll ctl err: %+v", err)
	}

	events := make([]syscall.EpollEvent, 1)

	// Worker pool.
	size := 100
	WorkPool := NewPool(size)
	for i := 0; i < size; i++ {
		WorkPool.Schedule(func() {})
	}

	for {
		n, err := syscall.EpollWait(fd, events, -1)
		if err != nil {
			log.Fatalf("EpollWait err: %+v", err)
		}
		for i := 0; i < n; i++ {
			if events[i].Events&syscall.EPOLLIN != 0 {
				fmt.Println("Read event occurred on listenerFd")

				sock, err := ls.Accept()
				if err != nil {
					log.Printf("Accept error: %v", err)
					continue
				}
				log.Printf("New socket accepted: %v", sock.RemoteAddr().String())

				WorkPool.Schedule(func() {
					datas := make([]byte, 1024)
					for {
						n, err := sock.Read(datas)
						if err != nil {
							log.Printf("Socket read err: %+v. Socket is closed.", err)
							break
						}
						log.Printf("Byte size: %v, Value: %v", n, string(datas))
					}

					sock.Close()
				})

			}
		}
	}

}
