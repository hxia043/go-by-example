package main

import (
	"reflect"
)

func handleMethod(i interface{}) {
	switch f := reflect.ValueOf(i); f.Kind() {
	case reflect.Struct:
		print("{")
		for i := 0; i < f.NumField(); i++ {
			if i > 0 {
				print(" ")
			}
			switch ff := f.Field(i); ff.Kind() {
			case reflect.String:
				//println(ff.Interface().(string))
				print(ff.String())
			case reflect.Int:
				print(ff.Int())
			default:
				print("unknown")
			}
		}
		print("}\n")
	}
}

func printer(i interface{}) {
	switch v := i.(type) {
	case string:
		println("string: ", v)
	case reflect.Type:
		println("reflect type: ", v)
	case reflect.Value:
		switch v.Kind() {
		case reflect.String:
			println(v.Interface().(string))
			println(v.String())
		}
	default:
		handleMethod(i)
	}
}

func main() {
	var s string = "hxia"
	printer(s)

	printer(reflect.TypeOf(s))
	printer(reflect.ValueOf(s))

	var dummy struct {
		name string
		age  int
	}
	dummy.name = "hxia"
	dummy.age = 28

	printer(dummy)
}
