package main

import (
	"fmt"
	"strconv"

	"github.com/rehacktive/kvod/kvod"
)

type User struct {
	Name  string
	Email string
}

func main() {
	db := kvod.Init("./db/", "whatever")
	userContainer := kvod.CreateContainer[User](db, "users")

	for i := 0; i < 10; i++ {
		err := userContainer.Put(strconv.Itoa(i), User{"name" + strconv.Itoa(i), "email"})
		if err != nil {
			panic(err)
		}
	}

	user, err := userContainer.Get("3")
	if err != nil {
		panic(err)
	}
	fmt.Println(*user)

	keys, err := userContainer.GetKeys()
	if err != nil {
		panic(err)
	}
	fmt.Println(keys)

	userValues, err := userContainer.GetData()
	if err != nil {
		panic(err)
	}
	fmt.Println(userValues)

	users, err := userContainer.GetAll()
	if err != nil {
		panic(err)
	}
	fmt.Println(users)
}
