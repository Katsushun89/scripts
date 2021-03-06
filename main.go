package main

import (
	"fmt"
)
const (
	script_path1 string = "./time.sh"
	script_path2 string = "./test.sh"
)
func execCommand(cmd_path string) {
	cmd := exec.Command("sh", "-c", cmd_path)
	cmd.Run()
	cmd.Wait()
}

func execScriptBackground(script1 string, ctx context.Context) {
	//var home string = os.Getenv("HOME")
	cmdPath := /*home +*/ script1

	execCommand(cmdPath)
	//cmd := exec.Command("sh", "-c", cmdPath)
	//cmd.Run()
	//cmd.Wait()

	for {
		select {
		case <-ctx.Done():
			return
		}
	}
}

/*
func killProcess() {
	cmdstr := "ps aux | grep \"time.sh\" | grep -v grep | awk '{ print \"kill -9\", $2 }' | sh"

	cmd := exec.Command("sh", "-c", cmdstr)
	cmd.Run()
	cmd.Wait()
}
*/

func exec2Scripts(script1 string, script2 string) error {

	err := isExistsScript(script1)
	if err != nil {
		fmt.Println("not exists script:", script1)
		return err
	}
	err = isExistsScript(script2)
	if err != nil {
		fmt.Println("not exists script:", script2)
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Scond*15)
	defer cancel()
	go execScriptBackground(script1, ctx)

	//fmt.Print("wait 15sec ")
	//time.Sleep(time.Second * 15)
	//fmt.Println("-> end of waiting")

	execCommand(script2)

	cancel()
	return err
}

func main() {
	if len(args) < 2 {
		fmt.Println("You must set 2 script file")
		return -1
	}

	var script1 = args[0]
	var script2 = args[1]

	err := exec2Scripts(script1, script2)
	if err != nil {
		return -1
	}
	return 0

}