## KVoD [Key/Value on Disk], a simple, encrypted key/value storage for Go (now using Generics)

Just a sample:
```golang
type User struct {
	Name  string
	Email string
}

func main() {
	store := kvod.Init("./db/", "secret")
	userContainer := kvod.CreateContainer[User](db, "users")

	userContainer.Put("1", User{"a username", "user@memorya.org"})

	user, _ := userContainer.Get("1")
	fmt.Println(user)

	userContainer.Delete("1")
}
```
The idea is to save each value to a file, under the specified path.
The value is serialized using gob and encrypted (with AES256 in GCM mode, with PBKDF2 as key derivation function with 10000 iterations).

The key is always a string.
The value can be any type gob supports (https://blog.golang.org/gobs-of-data), so structs, strings, slices and more, using Generics.

#### TESTS
```
ok  	github.com/rehacktive/kvod/kvod	0.299s	coverage: 82.1% of statements
Success: Tests passed.
```
#### BENCHMARKS
```
goos: linux
goarch: amd64
pkg: github.com/rehacktive/kvod/kvod
BenchmarkGenerateKey-4   	     100	  19887814 ns/op	     532 B/op	       8 allocs/op
BenchmarkEncrypt-4       	 1000000	      1945 ns/op	    1008 B/op	       7 allocs/op
BenchmarkDecrypt-4       	 2000000	      1079 ns/op	     976 B/op	       6 allocs/op
BenchmarkInit-4          	     100	  16506587 ns/op	    3483 B/op	      24 allocs/op
BenchmarkPut-4           	   20000	     84323 ns/op	    4744 B/op	      59 allocs/op
BenchmarkGet-4           	   20000	     70250 ns/op	   18346 B/op	     383 allocs/op
BenchmarkGetKeys-4       	    3000	    430877 ns/op	  121920 B/op	    2040 allocs/op
BenchmarkSerialize-4     	 1000000	      1461 ns/op	     912 B/op	      12 allocs/op
BenchmarkDeserialize-4   	 1000000	      2085 ns/op	    1088 B/op	      18 allocs/op
PASS
coverage: 79.5% of statements
ok  	github.com/rehacktive/kvod/kvod	19.396s
Success: Benchmarks passed.
```
