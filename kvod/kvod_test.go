package kvod

// import (
// 	"fmt"
// 	"os"
// 	"path/filepath"
// 	"reflect"
// 	"strconv"
// 	"testing"
// )

// type User struct {
// 	Name  string
// 	Email string
// }

// const (
// 	path     = "/tmp/kvod_test"
// 	password = "whatever"
// 	n        = 10
// )

// // TODO TEST FOR CONCURRENCY
// func TestConcurrency(t *testing.T) {
// 	kvod := Init(path, password)

// 	//c := make(chan int)

// 	var maxval = 100
// 	var incr = 0

// 	for incr = 0; incr <= maxval; incr++ {
// 		go func(v int) {
// 			err := kvod.Put("key", v)
// 			if err != nil {
// 				t.Error("error Put ", err)
// 			}
// 			fmt.Printf("send on channel: %v\n", v)
// 			//c <- v
// 		}(incr)
// 	}

// 	fmt.Println(incr)

// 	// for {
// 	// 	if incr == maxval {
// 	// 		//close(c)

// 	// 		var ret int
// 	// 		err := kvod.Get("key", &ret)
// 	// 		if err != nil {
// 	// 			t.Error("error Get ", err)
// 	// 		}
// 	// 		// if x != maxval {
// 	// 		// 	t.Errorf("error : last item arrived is %v", x)
// 	// 		// }
// 	// 		if ret != maxval {
// 	// 			t.Errorf("error put/get: expected %v found %v", maxval, ret)
// 	// 		}
// 	// 		return
// 	// 	}
// 	// }
// }

// func TestInit(t *testing.T) {
// 	os.RemoveAll(path)

// 	// repeat it twice, so second time will not create salt again
// 	for i := 0; i <= 1; i++ {
// 		kvod := Init(path, password)
// 		if kvod == nil {
// 			t.Error("error init kvod")
// 		}
// 		if kvod.path != path {
// 			t.Errorf("error kvod path: expected %x found %x", path, kvod.path)
// 		}
// 		// check that .salt file exists
// 		saltFile := filepath.Join(path, saltFilename)
// 		if _, err := os.Stat(saltFile); os.IsNotExist(err) {
// 			t.Error("error .salt file not found")
// 		}
// 	}
// }

// func TestPutAndGet(t *testing.T) {
// 	kvod := Init(path, password)

// 	sampleUser := User{"test user", "test@user.org"}

// 	err := kvod.Put("1", sampleUser)
// 	if err != nil {
// 		t.Error("error Put ", err)
// 	}

// 	var user User
// 	err = kvod.Get("1", &user)
// 	if err != nil {
// 		t.Error("error Get ", err)
// 	}
// 	if !reflect.DeepEqual(sampleUser, user) {
// 		t.Errorf("error put/get: expected %v found %v", sampleUser, user)
// 	}
// }

// func TestGetKeys(t *testing.T) {
// 	kvod := Init(path, password)

// 	sampleUser := User{"test user", "test@user.org"}

// 	for i := 0; i < n; i++ {
// 		kvod.Put(strconv.Itoa(i), sampleUser)
// 	}

// 	keys, err := kvod.GetKeys()
// 	if err != nil {
// 		t.Error("error GetKeys ", err)
// 	}

// 	if len(keys) != n {
// 		t.Errorf("error get keys: expected %v found %v", n, len(keys))
// 	}
// }

// func TestDelete(t *testing.T) {
// 	kvod := Init(path, password)

// 	sampleUser := User{"test user", "test@user.org"}

// 	err := kvod.Put("1", sampleUser)
// 	if err != nil {
// 		t.Error("error Put ", err)
// 	}

// 	kvod.Delete("1")

// 	var user User
// 	err = kvod.Get("1", &user)
// 	// this should raise an error
// 	if err == nil {
// 		t.Error("error delete ", err)
// 	}
// }
// func BenchmarkInit(b *testing.B) {
// 	os.RemoveAll(path)
// 	for i := 0; i < b.N; i++ {
// 		Init(path, password)
// 	}
// }

// func BenchmarkPut(b *testing.B) {
// 	kvod := Init(path, password)
// 	sampleUser := User{"test user", "test@user.org"}
// 	for i := 0; i < b.N; i++ {
// 		kvod.Put("1", sampleUser)
// 	}
// }

// func BenchmarkGet(b *testing.B) {
// 	kvod := Init(path, password)
// 	sampleUser := User{"test user", "test@user.org"}
// 	kvod.Put("1", sampleUser)
// 	var user User

// 	for i := 0; i < b.N; i++ {
// 		kvod.Get("1", &user)
// 	}
// }

// func BenchmarkGetKeys(b *testing.B) {
// 	kvod := Init(path, password)

// 	sampleUser := User{"test user", "test@user.org"}

// 	for i := 0; i < n; i++ {
// 		kvod.Put(strconv.Itoa(i), sampleUser)
// 	}
// 	for i := 0; i < b.N; i++ {
// 		kvod.GetKeys()
// 	}
// }
