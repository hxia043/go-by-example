Q: Why lambda struct needed?
A: With the structure like:
```
type person struct {
	name string
	age  int
}

type student struct {
	young person
	class  string
}

s := student{
    young: person{
        name: "hxia",
        age:  21,
    },
    class: "三年二班",
}

fmt.Printf("%p, %p, %p, %p, %p\n", &s, &s.young, &s.young.name, &s.young.age, &s.class)
```

output:
```
0xc0000a0150, 0xc0000a0150, 0xc0000a0150, 0xc0000a0160, 0xc0000a0168
```

With the named struct, the student access to the filed of name/age, it should access to the identifier young first.

With the lambda structure nesting like:
```
type empolyee struct {
	person
	company string
}

e := empolyee{
    person:  s.young,
    company: "Nokia",
}

fmt.Printf("%p, %p, %p, %p, %p\n", &e, &e.person, &e.person.name, &e.person.age, &e.company)
fmt.Println(e.name, e.age)
```

output:
```
0xc000074150, 0xc000074150, 0xc000074150, 0xc000074160, 0xc000074168
hxia 21
```

The memory allocated is same as specific struct nesting, but can directly access to the filed of name/age by object e.
It can implement the effect of inherit, which to make a relationship with A has B, like e <template empolyee> has a field name.
With the specific struct nesting, the relationship is e has a object young which has a filed name.

Q: how to overloading the property?
A: Yes, with the lambda struct nesting, we can implement the effect like inherit. So according can implement the overloading which to over write the property. such as:
```
type empolyee struct {
	person
	company string
	name    string
}

eo := empolyee{
    person:  s.young,
    company: "Nokia",
    name:    "Troy",
}

fmt.Println(eo.name, eo.person.name)
```

output:
```
Troy hxia
```

Q: how can implement the inherit and overloading of method for object?
A: Not only can implement the overloading of property, but also can implement the overloading of method for object, like inherit:
```
func (p *person) sayHello() {
	fmt.Println(p.name)
}

eo.sayHello()
```

output:
```
hxia
```

overloading:
```
func (e *empolyee) sayHello() {
	fmt.Println("hello, world")
}

eo.sayHello()
```

output:
```
hello, world
```
