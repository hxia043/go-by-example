# interface

The interface is a pretty graceful design in Go. With interface the object which implement the action of interface can implement the polymorphic.

For example:
```
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

func main() {
    var p person

    s := &student{"三年二班"}
    p = s

    p.sleep()

    e := &employee{"Nokia"}
    p = e

    p.sleep()
}
```

output:
```
sleep at 三年二班
sleep at Nokia
```

The interface p can overloading by struct object s and e, it kindly an implemention of polymorphic.

Let's use the function parameter to make the polymorphic more clearly, as:
```
func sleep(p person) {
	p.sleep()
}

sleep(s)
sleep(e)
```

output:
```
sleep at 三年二班
sleep at Nokia
```

The sleep function has make the polymorphic more clear which to output different result with different object.
