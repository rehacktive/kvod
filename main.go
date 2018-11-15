package main

import (
	"fmt"
	"strconv"

	"github.com/rehacktive/kvod/kvod"
)

type user struct {
	Name  string
	Email string
}

func main() {
	db := kvod.Init("./db/", "whatever")

	for i := 0; i < 10; i++ {
		err := db.Put(strconv.Itoa(i), user{"name" + strconv.Itoa(i), "email"})
		if err != nil {
			panic(err)
		}
	}

	var u user
	err := db.Get("3", &u)
	if err != nil {
		panic(err)
	}
	fmt.Println(u)

	keys, err := db.GetKeys()
	fmt.Println(keys)

	keys, _ = db.GetKeys()
	for _, k := range keys {
		err := db.Get(k, &u)
		if err != nil {
			panic(err)
		}
		fmt.Println(k, "->", u)
	}
}
