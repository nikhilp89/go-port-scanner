package main

import (
	"context"
	"fmt"
	"net"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/sync/semaphore"
	"github.com/pieterclaerhout/go-waitgroup"
)

type PortScanner struct {
	host string
	lock *semaphore.Weighted
}

func ulimit() int64 {
    out, err := exec.Command("ulimit", "-n").Output()
    if err != nil {
        panic(err)
    }
    s := strings.TrimSpace(string(out))

    i, err := strconv.ParseInt(s, 10, 64)

    if err != nil {
        panic(err)
    }
    return i
}

func scanPort(host string, port int, timeout time.Duration) {
	target := fmt.Sprintf("%s:%d", host, port)

	conn, err := net.DialTimeout("tcp", target, timeout)

	if err != nil {
		if strings.Contains(err.Error(), "too many open files") {
			time.Sleep(timeout)
            scanPort(host, port, timeout)
		} /*else {
			fmt.Println(port, "closed")
		}*/

		return
	}

	conn.Close()
	fmt.Printf("%d/tcp open\n", port)
}

func (ps *PortScanner) StartUnlimited(firstPort int, lastPort int, timeout time.Duration) {
	wg := sync.WaitGroup{}
	defer wg.Wait()

	for currentPort := firstPort; currentPort <= lastPort; currentPort++ {
		wg.Add(1)
		ps.lock.Acquire(context.TODO(), 1)

		go func(currentPort int) {
			defer ps.lock.Release(1)
			defer wg.Done()
			scanPort(ps.host, currentPort, timeout)
		}(currentPort)
	}
}

func (ps *PortScanner) StartLimited(firstPort int, lastPort int, timeout time.Duration) {
	wg := waitgroup.NewWaitGroup(1000)

	for currentPort := firstPort; currentPort <= lastPort; currentPort++ {
		wg.BlockAdd()

		go func(currentPort int) {
			defer wg.Done()
			scanPort(ps.host, currentPort, timeout)
		}(currentPort)
	}

	wg.Wait()
}


func main() {

	start_time := time.Now()

	ps := &PortScanner {
		//host: "17.57.0.206",
		host: "60.205.174.40",
		lock: semaphore.NewWeighted(ulimit()),
	}

	//ps.StartUnlimited(0, 65535, 500*time.Millisecond)
	ps.StartLimited(0, 65535, 500*time.Millisecond)

	time_elapsed := time.Since(start_time)
	fmt.Printf("\nScan took %s\n", time_elapsed)		
}