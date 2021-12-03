package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"reflect"
)

//type UserX struct {
//	Name  string `mytag:"MyName"`
//	Email string `mytag:"MyEmail"`
//}

type User struct {
	Name  string `json:"name" xml:"name" tag:"name"`
	Email string `json:"email" xml:"email"`
	Info  string `json:"-"`
}

func (u *User) String() string {
	return fmt.Sprintf("Hi! My name is %s, My Email is %s", u.Name, u.Email)
}
func main() {
	u := &User{
		Name:  "Sammy",
		Email: "xxx@163.com",
	}
	v := User{
		Info: "hackerman",
	}

	t := reflect.TypeOf(*u)

	for _, fieldName := range []string{"Name", "Email"} {
		field, found := t.FieldByName(fieldName)
		if !found {
			continue
		}
		fmt.Printf("\nField: User.%s\n", fieldName)
		fmt.Printf("\tWhole tag value : %q\n", field.Tag)
		fmt.Printf("\tValue of 'json': %q\n", field.Tag.Get("json"))
		fmt.Printf("\tValue of 'xml': %q\n", field.Tag.Get("xml"))

		if _, ok := field.Tag.Lookup("example"); !ok {
			fmt.Printf("\texample is not set\n\n\n")
		}
	}

	fmt.Println(u)
	fmt.Println(v)

	jsonString := `
    {
        "name": "John",
        "email": "Smith.john@hacker.com",
"info":"Professor"
    }`

	person := new(User)
	json.Unmarshal([]byte(jsonString), person)
	fmt.Println(person)

	newJson, _ := json.Marshal(person)
	fmt.Printf("%s\n", newJson)

	xmlString := `<xml>
<name>Mike</name>
<email>Mike.Brown@hacker.com</email>
<info>this is it!</info>
</xml>`

	anotherPerson := new(User)
	xml.Unmarshal([]byte(xmlString), anotherPerson)
	fmt.Println(anotherPerson)

	newXml, _ := xml.Marshal(anotherPerson)
	fmt.Println((newXml))

}
