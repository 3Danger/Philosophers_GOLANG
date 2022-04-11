package main

import (
	"ProducerConsumer/src/Philosoph"
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Args struct {
	countPhi int
	eat      int
	sleep    int
	dead     int
}

func atoi(str string) int {
	num, ok := strconv.Atoi(str)
	if ok != nil {
		fmt.Print(str, ": is not valide!\n", ok)
		os.Exit(1)
	}
	return num
}

func getInfo() (a Args) {
	if len(os.Args) > 4 {
		a.countPhi = atoi(os.Args[1])
		a.eat = atoi(os.Args[2])
		a.sleep = atoi(os.Args[3])
		a.dead = atoi(os.Args[4])
		return a
	}
	read := func() func() int {
		reader := bufio.NewReader(os.Stdin)
		return func() int {
			str, _ := reader.ReadString('\n')
			str = strings.Replace(str, "\n", "", 1)
			num := atoi(str)
			return num
		}
	}()
	fmt.Print("how many philosophers? :")
	a.countPhi = read()
	fmt.Print("how long will it take to eat? :")
	a.eat = read()
	fmt.Print("how long will it take to sleep? :")
	a.sleep = read()
	fmt.Print("how many lives without food? :")
	a.dead = read()
	return a
}

func main() {

	args := getInfo()

	philo := Philosoph.NewPhilosopher(args.countPhi, args.eat, args.sleep, args.dead)
	philo.Run()
	for !philo.IsFinished() {
		time.Sleep(time.Second)
	}
}
