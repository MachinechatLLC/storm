package storm

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAllByIndex(t *testing.T) {
	dir, _ := ioutil.TempDir(os.TempDir(), "storm")
	defer os.RemoveAll(dir)
	db, _ := Open(filepath.Join(dir, "storm.db"))

	for i := 0; i < 100; i++ {
		w := User{Name: "John", ID: i + 1, Slug: fmt.Sprintf("John%d", i+1), DateOfBirth: time.Now().Add(-time.Duration(i*10) * time.Minute)}
		err := db.Save(&w)
		assert.NoError(t, err)
	}

	err := db.AllByIndex("", nil)
	assert.Error(t, err)
	assert.EqualError(t, err, "provided target must be a pointer to a slice")

	var users []User

	err = db.All(&users)
	assert.NoError(t, err)
	assert.Len(t, users, 100)
	assert.Equal(t, 1, users[0].ID)
	assert.Equal(t, 100, users[99].ID)

	err = db.AllByIndex("DateOfBirth", &users)
	assert.NoError(t, err)
	assert.Len(t, users, 100)
	assert.Equal(t, 100, users[0].ID)
	assert.Equal(t, 1, users[99].ID)

	err = db.AllByIndex("Name", &users)
	assert.NoError(t, err)
	assert.Len(t, users, 100)
	assert.Equal(t, 1, users[0].ID)
	assert.Equal(t, 100, users[99].ID)

	var unknowns []UserWithNoID
	err = db.All(&unknowns)
	assert.Error(t, err)
	assert.EqualError(t, err, "missing struct tag id or ID field")

	err = db.Save(&NestedID{
		ToEmbed: ToEmbed{ID: "id1"},
		Name:    "John",
	})
	assert.NoError(t, err)

	err = db.Save(&NestedID{
		ToEmbed: ToEmbed{ID: "id2"},
		Name:    "Mike",
	})
	assert.NoError(t, err)

	db.Save(&NestedID{
		ToEmbed: ToEmbed{ID: "id3"},
		Name:    "Steve",
	})
	assert.NoError(t, err)

	var nested []NestedID
	err = db.All(&nested)
	assert.NoError(t, err)
	assert.Len(t, nested, 3)
}
