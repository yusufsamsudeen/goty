package main

import (
	"fmt"
	"github.com/yusufsamsudeen/goty/goty"
	"gorm.io/gorm"
)

type Person struct {
	gorm.Model
	FirstName string
	LastName string
}
func main() {
	person := Person{
		FirstName: "Yusuf",
		LastName: "Samsudeen",
	}
	save := goty.Save(&person)
	fmt.Println(save.Error)

	fmt.Println(person.Model.CreatedAt)

}
