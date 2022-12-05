package main

import (
	"fmt"
	"html/template"
	"os"
	"strings"
)

type Login struct {
	Username string
	Token    string
	Week     []*Data
}

type Data struct {
	Data string
}

func TranslateToLower(args ...interface{}) string {
	for _, arg := range args {
		s := fmt.Sprint(arg)
		return strings.ToLower(s)
	}

	return " "
}

func main() {
	d1, d2 := Data{Data: "Saturday"}, Data{Data: "Sunday"}
	p := Login{Username: "Huyun", Token: "1234567890", Week: []*Data{&d1, &d2}}

	t := template.New("test.tpl")
	t = t.Funcs(template.FuncMap{"tToLower": TranslateToLower})
	t, _ = t.ParseFiles("test.tpl")

	if err := t.ExecuteTemplate(os.Stdout, "test", p); err != nil {
		panic(err)
	}

	fmt.Println("=================================================")

	t2, _ := template.ParseFiles("header.tpl", "content.tpl", "footer.tpl")
	if err := t2.ExecuteTemplate(os.Stdout, "content", nil); err != nil {
		panic(err)
	}
}
