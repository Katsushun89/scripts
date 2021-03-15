# scripts
Go cli tool for parallel execution shell script

# Install
You can download binary file of cli tool from release.
Use it by deploying it where the path goes.

https://github.com/Katsushun89/scripts/releases

# Support environment
- os : linux
- arch : amd64, arm64

# How to use
If a shell script is specified after the command, it will be executed.

```
$ scripts test/echo_time_1sec.sh 
timeout duration: 60 [sec]
run scripts: [test/echo_time_1sec.sh]
all scripts finished
```

```bash:echo_time_1sec.sh
#!/bin/sh

sleep 1
date >> echo_time.log
```

If you want to run two shells, specify two.

```
$ scripts test/echo_time_1sec.sh test/echo_cnt_2sec.sh 
timeout duration: 60 [sec]
run scripts: [test/echo_time_1sec.sh test/echo_cnt_2sec.sh]
all scripts finished
```

```bash:echo_cnt_2sec.sh
#!/bin/sh

cnt=0
sleep 2 
cnt=$((cnt+1))
echo $cnt " : hoge" >> echo_cnt.log

```

The above example is a script that outputs a file, so we can confirm that it works by ls command.

```
$ ls -la
-rw-rw-r--  1 username username   10  3月 14 01:35 echo_cnt.log
-rw-rw-r--  1 username username   43  3月 14 01:35 echo_time.log
```

Also, the timeout time can be set with the option (-t, --timeout). (Default: 60 seconds)
Scripts that do not terminate, such as infinite loops, will terminate at the specified timeout time.

```
$ scripts -t 10 test/echo_time.sh test/echo_cnt.sh 
timeout duration: 10 [sec]
run scripts: [test/echo_time.sh test/echo_cnt.sh]
timed out
```

```bash:echo_time.sh
#!/bin/sh

while true
do
  sleep 1
  date >> echo_time.log
done

```

```bash:echo_cnt.sh
#!/bin/sh

cnt=0
while true
do
  sleep 1
  cnt=$((cnt+1))
  echo $cnt " : hoge" >> echo_cnt.log
done
```

You can find out how to set the options in the help section.

```
$ scripts -h
NAME:scripts - cli tool to run multiple shell scripts in parallel

USAGE:
   scripts [options] [arguments...]

OPTIONS:
   --timeout value, -t value  timeout duration [sec] (default: "60")
   --help, -h                 show help (default: false)
```

If you want to quit in the middle of the process, you can do so by pressing Ctrl+C.
```
$ scripts -t 10 test/echo_time.sh test/echo_cnt.sh 
timeout duration: 10 [sec]
run scripts: [test/echo_time.sh test/echo_cnt.sh]
^Csignal: interrupt
```
