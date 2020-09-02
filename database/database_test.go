package database

import (
	"reflect"
	"testing"
)

//Note Tests to be run in full sequence

var mockDbPath = "../mocks"

func MocktDb() Database {
	db, err := Open(mockDbPath)
	if err != nil {
		return db
	}

	return db
}

func TestSet(t *testing.T) {
	d := MocktDb()
	defer d.Close()

	err := d.Set("testKey", []byte("testData"))
	if err != nil {
		t.Errorf("error in set testKey %s", err)
	}

	err = d.Set("testKey2", []byte("testData2"))
	if err != nil {
		t.Errorf("error in set testKey2 %s", err)
	}

	err = d.Set("testKey3", []byte("testData3"))
	if err != nil {
		t.Errorf("error in set testKey3 %s", err)
	}
}

func TestGet(t *testing.T) {
	d := MocktDb()
	defer d.Close()
	expect := "testData"

	actual, err := d.Get("testKey")
	if err != nil {
		t.Errorf("error in get %s", err)
	}

	if expect != string(actual) {
		t.Errorf("expected %s got %s", expect, string(actual))
	}
}

func TestDelete(t *testing.T) {
	d := MocktDb()
	defer d.Close()

	err := d.Delete("testKey")
	if err != nil {
		t.Errorf("error in delete %s", err)
	}

	_, err = d.Get("testKey")
	if err == nil {
		t.Errorf("get should have failed as the key should not exist")
	}

}

func TestGetAllKeys(t *testing.T) {
	d := MocktDb()
	defer d.Close()
	expect := []string{"testKey2", "testKey3"}

	actual, err := d.GetAllKeys()
	if err != nil {
		t.Errorf("error in get all %s", err)
	}

	if reflect.DeepEqual(expect, &actual) {
		t.Errorf("expected %s got %s", expect, actual)
	}
}

func TestDeleteAll(t *testing.T) {
	d := MocktDb()
	defer d.Close()
	expected := 0

	err := d.DeleteAll()
	if err != nil {
		t.Errorf("error in delete all %s", err)
	}

	actual, err := d.GetAllKeys()
	if err != nil {
		t.Errorf("error in get all %s", err)
	}

	if len(actual) != expected {
		t.Errorf("expected %d keys got %d keys", expected, len(actual))
	}
}
