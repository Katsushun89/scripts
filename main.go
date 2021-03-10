package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"
)

func execScript(script string, ctx context.Context, ch chan int, wg *sync.WaitGroup) {
	cmd := exec.Command("sh", "-c", script)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Start()
	ch <- cmd.Process.Pid

	cmd.Wait()
	fmt.Println("end exec ", script)
	compDone := false
	wg.Done()
	fmt.Println("wg.Done()")
	compDone = true
	for {
		select {
		case <-ctx.Done():
			fmt.Println("execScript Done")
			if !compDone {
				wg.Done()
				fmt.Println("wg.Done()")
			}
			return
		}
	}
}

func isExistsScript(file_name string) bool {
	_, err := os.Stat(file_name)
	if err == nil {
		return true
	}
	return false
}

func waitTimeout(wg *sync.WaitGroup, timeout time.Duration) bool {
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()
	select {
	case <-c:
		return false // completed normally
	case <-time.After(timeout):
		return true // timed out
	}
}

func execScripts(scripts []string) error {
	ch := make(chan int)
	c := make(chan struct{})
	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for _, script := range scripts {
		if !isExistsScript(script) {
			err := fmt.Errorf("not exists script:%s", script)
			return err
		}

		wg.Add(1)
		go execScript(script, ctx, ch, &wg)
	}

	var pgids []int

	for _, _ = range scripts {
		select {
		case v := <-ch:
			pgids = append(pgids, v)
		}
	}

	go func(ctx context.Context) {
		wg.Wait()
		//c <- struct{}{}
		for {
			select {
			case <-ctx.Done():
				fmt.Println("go func Done")
				return
			}
		}
	}(ctx)

	timeout := time.Duration(5) * time.Second
	fmt.Printf("Wait for waitgroup (up to %s)\n", timeout)

	s := make(chan os.Signal)
	signal.Notify(s, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(s)

	select {
	case sig := <-s:
		fmt.Println("signal:", sig)
	case <-c:
		fmt.Printf("Wait group finished\n")
	case <-time.After(timeout):
		fmt.Printf("Timed out waiting for wait group\n")
	}

	for _, pgid := range pgids {
		fmt.Println("kill process pgid:", pgid)
		syscall.Kill(-pgid, syscall.SIGKILL)
	}

	fmt.Println("before cancel")
	cancel()
	time.Sleep(2 * time.Second) //wait cancel
	fmt.Println("after cancel and sleep")
	return nil
}

func getScritps(args []string) ([]string, error) {
	if len(args) < 3 {
		err := fmt.Errorf("Set one or more scripts")
		return []string{""}, err
	}
	scripts := []string{}
	for i, arg := range args {
		if i != 0 {
			scripts = append(scripts, arg)
		}
	}
	fmt.Println(scripts)
	return scripts, nil
}

func main() {
	scripts, err := getScritps(os.Args)
	if err != nil {
		fmt.Printf("arg err :%s", err)
		return
	}
	fmt.Println("A num goroutine: ", runtime.NumGoroutine())
	err = execScripts(scripts)
	fmt.Println("B num goroutine: ", runtime.NumGoroutine())
	if err != nil {
		fmt.Printf("exec2scripts error:%s\n", err)
		return
	}
	return

}
