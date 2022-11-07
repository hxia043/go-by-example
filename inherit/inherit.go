package main

import "fmt"

type person struct {
	name string
	age  int
}

type student struct {
	young person
	class string
}

type empolyee struct {
	person
	company string
	name    string
}

func (p *person) sayHello() {
	fmt.Println(p.name)
}

func (e *empolyee) sayHello() {
	fmt.Println("hello, world")
}

func main() {
	s := student{
		young: person{
			name: "hxia",
			age:  21,
		},
		class: "三年二班",
	}

	fmt.Printf("%p, %p, %p, %p, %p\n", &s, &s.young, &s.young.name, &s.young.age, &s.class)

	e := empolyee{
		person:  s.young,
		company: "Nokia",
	}

	fmt.Println(e.name, e.age)

	eo := empolyee{
		person:  s.young,
		company: "Nokia",
		name:    "Troy",
	}

	fmt.Println(eo.name, eo.person.name)

	eo.sayHello()
}
