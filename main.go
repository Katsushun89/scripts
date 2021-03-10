package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
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
	compDone := false
	wg.Done()
	compDone = true
	for {
		select {
		case <-ctx.Done():
			if !compDone {
				wg.Done()
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

func execScripts(scripts []string) error {
	ch := make(chan int)
	c := make(chan int)
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
		select {
		case c <- 1:
		default:
		}
		for {
			select {
			case <-ctx.Done():
				return
			}
		}
	}(ctx)

	timeout := time.Duration(5) * time.Second

	s := make(chan os.Signal)
	signal.Notify(s, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(s)

	select {
	case sig := <-s:
		fmt.Println("signal:", sig)
	case <-c:
		fmt.Printf("all scripts finished\n")
	case <-time.After(timeout):
		fmt.Printf("timed out\n")
	}

	for _, pgid := range pgids {
		//fmt.Println("kill process pgid:", pgid)
		syscall.Kill(-pgid, syscall.SIGKILL)
	}

	cancel()
	time.Sleep(1 * time.Second) //wait cancel
	//fmt.Println("after cancel and sleep")
	return nil
}

func getScritps(args []string) ([]string, error) {
	if len(args) < 2 {
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
	err = execScripts(scripts)
	if err != nil {
		fmt.Printf("exec2scripts error:%s\n", err)
		return
	}
	return

}
