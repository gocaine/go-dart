#!/bin/bash

ps -ef |grep go-dart|sed "s/  */ /g" | cut -d " " -f 2 | xargs sudo kill -9

for pin in $(ls /sys/class/gpio/ | xargs -0 basename | grep "gpio[0-9][0-9]*" | sed "s/gpio//")
do
	echo "clearing $pin"
	echo "$pin" > /sys/class/gpio/unexport
done
