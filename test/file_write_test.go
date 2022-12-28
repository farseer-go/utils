package test

import (
	file2 "github.com/farseer-go/utils/file"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"os"
	"strconv"
	"testing"
)

func TestWriteString(t *testing.T) {
	file := "./Farseer.Go/" + strconv.Itoa(rand.Intn(999-100)) + ".txt"
	defer os.Remove(file)

	content := "aaa"
	file2.WriteString(file, content)
	assert.Equal(t, file2.ReadString(file), content)
}

func TestAppendString(t *testing.T) {
	file := "./Farseer.Go/" + strconv.Itoa(rand.Intn(999-100)) + ".txt"
	defer os.Remove(file)

	file2.WriteString(file, "aaa")
	file2.AppendString(file, "bbb")
	readString := file2.ReadString(file)
	assert.Equal(t, readString, "aaabbb")
}

func TestAppendLine(t *testing.T) {
	file := "./Farseer.Go/" + strconv.Itoa(rand.Intn(999-100)) + ".txt"
	defer os.Remove(file)

	file2.WriteString(file, "aaa")
	file2.AppendLine(file, "bbb")
	readString := file2.ReadString(file)
	assert.Equal(t, readString, "aaa\nbbb")
}

func TestAppendAllLine(t *testing.T) {
	file := "./Farseer.Go/" + strconv.Itoa(rand.Intn(999-100)) + ".txt"
	defer os.Remove(file)

	file2.WriteString(file, "aaa")
	str := []string{"bbb", "ccc"}
	file2.AppendAllLine(file, str)
	readString := file2.ReadString(file)
	assert.Equal(t, readString, "aaa\nbbb\nccc")
}
