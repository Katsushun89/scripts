#!/bin/sh

cnt=0
while true
do
  sleep 1
  cnt=$((cnt+1))
  echo $cnt " : hoge" >> echo_cnt.log
done
