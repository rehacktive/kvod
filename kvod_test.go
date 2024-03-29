package kvod

import (
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"testing"
)

type User struct {
	Name  string
	Email string
}

const (
	path     = "/tmp/kvod_test"
	password = "whatever"
	n        = 10
)

func TestInit(t *testing.T) {
	os.RemoveAll(path)

	// repeat it twice, so second time will not create salt again
	for i := 0; i <= 1; i++ {
		kvod := Init(path, password)
		if kvod == nil {
			t.Error("error init kvod")
		}
		if kvod.path != path {
			t.Errorf("error kvod path: expected %x found %x", path, kvod.path)
		}
		// check that .salt file exists
		dbFile := filepath.Join(path, dbFilename)
		if _, err := os.Stat(dbFile); os.IsNotExist(err) {
			t.Error("error dbFile not found")
		}
	}
}

func TestPutAndGet(t *testing.T) {
	os.RemoveAll(path)

	kvod := Init(path, password)
	container := CreateContainer[User](kvod, "users")

	sampleUser := User{"test user", "test@user.org"}

	err := container.Put("1", sampleUser)
	if err != nil {
		t.Error("error Put ", err)
	}

	user, err := container.Get("1")
	if err != nil {
		t.Error("error Get ", err)
	}
	if !reflect.DeepEqual(sampleUser, *user) {
		t.Errorf("error put/get: expected %v found %v", sampleUser, user)
	}
}

func TestGetKeys(t *testing.T) {
	kvod := Init(path, password)
	container := CreateContainer[User](kvod, "users")

	sampleUser := User{"test user", "test@user.org"}

	for i := 0; i < n; i++ {
		container.Put(strconv.Itoa(i), sampleUser)
	}

	keys, err := container.GetKeys()
	if err != nil {
		t.Error("error GetKeys ", err)
	}

	if len(keys) != n {
		t.Errorf("error get keys: expected %v found %v", n, len(keys))
	}
}

func TestGetData(t *testing.T) {
	kvod := Init(path, password)
	container := CreateContainer[User](kvod, "users")

	sampleUser := User{"test user", "test@user.org"}

	for i := 0; i < n; i++ {
		container.Put(strconv.Itoa(i), sampleUser)
	}

	data, err := container.GetData()
	if err != nil {
		t.Error("error GetKeys ", err)
	}

	if len(data) != n {
		t.Errorf("error get keys: expected %v found %v", n, len(data))
	}
}

func TestGetAll(t *testing.T) {
	kvod := Init(path, password)
	container := CreateContainer[User](kvod, "users")

	sampleUser := User{"test user", "test@user.org"}

	for i := 0; i < n; i++ {
		container.Put(strconv.Itoa(i), sampleUser)
	}

	data, err := container.GetAll()
	if err != nil {
		t.Error("error GetKeys ", err)
	}

	if len(data) != n {
		t.Errorf("error get keys: expected %v found %v", n, len(data))
	}
}

func TestDelete(t *testing.T) {
	kvod := Init(path, password)
	container := CreateContainer[User](kvod, "users")

	sampleUser := User{"test user", "test@user.org"}

	err := container.Put("1", sampleUser)
	if err != nil {
		t.Error("error Put ", err)
	}

	container.Delete("1")

	_, err = container.Get("1")
	// this should raise an error
	if err == nil {
		t.Error("error delete ", err)
	}
}
func BenchmarkInit(b *testing.B) {
	os.RemoveAll(path)
	for i := 0; i < b.N; i++ {
		Init(path, password)
	}
}

func BenchmarkPut(b *testing.B) {
	kvod := Init(path, password)
	container := CreateContainer[User](kvod, "users")

	sampleUser := User{"test user", "test@user.org"}
	for i := 0; i < b.N; i++ {
		container.Put("1", sampleUser)
	}
}

func BenchmarkGet(b *testing.B) {
	kvod := Init(path, password)
	container := CreateContainer[User](kvod, "users")

	sampleUser := User{"test user", "test@user.org"}
	container.Put("1", sampleUser)

	for i := 0; i < b.N; i++ {
		container.Get("1")
	}
}

func BenchmarkGetKeys(b *testing.B) {
	kvod := Init(path, password)
	container := CreateContainer[User](kvod, "users")

	sampleUser := User{"test user", "test@user.org"}

	for i := 0; i < n; i++ {
		container.Put(strconv.Itoa(i), sampleUser)
	}
	for i := 0; i < b.N; i++ {
		container.GetKeys()
	}
}
