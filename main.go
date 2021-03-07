package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)

func execScript(script string, ctx context.Context, ch chan int) {
	cmd := exec.Command("sh", "-c", script)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Start()
	ch <- cmd.Process.Pid

	cmd.Wait()
	for {
		select {
		case <-ctx.Done():
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

func execScripts(scripts [2]string) error {
	ch := make(chan int)

	for _, script := range scripts {
		if !isExistsScript(script) {
			err := fmt.Errorf("not exists script:%s", script)
			return err
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
		defer cancel()
		go execScript(script, ctx, ch)
	}

	var pgids []int

	for _, _ = range scripts {
		select {
		case v := <-ch:
			pgids = append(pgids, v)
		}
	}

	fmt.Printf("wait %d sec ", 5)
	t := time.NewTimer(5 * time.Second)
	fmt.Println("-> end of waiting")

	s := make(chan os.Signal)
	signal.Notify(s, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	select {
	case sig := <-s:
		fmt.Println("signal:", sig)
	case <-t.C:
		fmt.Println("timeout")
	}

	for _, pgid := range pgids {
		fmt.Println("kill process pgid:", pgid)
		syscall.Kill(-pgid, syscall.SIGKILL)
	}

	return nil
}

func getScritps(args []string) ([2]string, error) {
	if len(args) < 3 {
		err := fmt.Errorf("You must set 1 script file")
		return [2]string{"", ""}, err
	}
	scripts := [2]string{args[1], args[2]}
	return scripts, nil
}

func showPS() {
	b, err := exec.Command("ps", "j").Output()
	fmt.Println(string(b), err)
}

func main() {
	scripts, err := getScritps(os.Args)

	fmt.Println("scripts[0]:", scripts[0])
	fmt.Println("scripts[1]:", scripts[1])
	fmt.Println(err)
	//return

	err = execScripts(scripts)
	if err != nil {
		fmt.Printf("exec2scripts error:%s\n", err)
		return
	}
	/*
		s := make(chan os.Signal)
		signal.Notify(s, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
		select {
		case sig := <-s:
			fmt.Println("signal:", sig)
		}
	*/
	return

}
