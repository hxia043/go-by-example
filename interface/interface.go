package main

import "fmt"

type person interface {
	sleep()
}

type student struct {
	class string
}

type employee struct {
	company string
}

func (s *student) sleep() {
	fmt.Println("sleep at", s.class)
}

func (e *employee) sleep() {
	fmt.Println("sleep at", e.company)
}

func sleep(p person) {
	p.sleep()
}

func main() {
	var p person

	s := &student{"三年二班"}
	p = s

	p.sleep()

	e := &employee{"Nokia"}
	p = e

	p.sleep()

	sleep(s)
	sleep(e)
}
