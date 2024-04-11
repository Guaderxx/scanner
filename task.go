package main

import (
	"fmt"
	"log/slog"
	"runtime"
	"strings"
	"sync"
)

func RunTask(tasks []IpPort) {
	wg := &sync.WaitGroup{}

	taskChan := make(chan IpPort, concurrency)

	for i := 0; i < int(concurrency); i++ {
		go Scan(taskChan, wg)
	}

	for _, task := range tasks {
		wg.Add(1)
		taskChan <- task
	}

	close(taskChan)
	wg.Wait()
}

func Scan(taskChan chan IpPort, wg *sync.WaitGroup) {
	for task := range taskChan {
		if strings.ToLower(mode) == "syn" {
			SaveResult(SynScan(task.ip, task.port))
		} else {
			SaveResult(Connect(task.ip, task.port))
		}
		wg.Done()
	}
}

func SaveResult(ip string, port int, err error) {
	slog.Info("task", slog.Group("info",
		slog.String("ip", ip),
		slog.Int("port", port),
		slog.Int("goroutine-num", runtime.NumGoroutine()),
	))
	if err != nil {
		return
	}

	if port > 0 {
		v, ok := result.Load(ip)
		if ok {
			ports, ok1 := v.([]int)
			if ok1 {
				ports = append(ports, port)
				result.Store(ip, ports)
			}
		} else {
			ports := make([]int, 0)
			ports = append(ports, port)
			result.Store(ip, ports)
		}
	}
}

func PrintResult() {
	result.Range(func(key, value interface{}) bool {
		fmt.Printf("ip: %v\n", key)
		fmt.Printf("ports: %v\n", value)
		fmt.Println(strings.Repeat("-", 79))
		return true
	})
}
