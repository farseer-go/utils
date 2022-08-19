package file

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"os"
	"strconv"
	"testing"
)

func TestWriteString(t *testing.T) {
	file := "/Users/steden/Desktop/code/project/Farseer.Go/" + strconv.Itoa(rand.Intn(999-100)) + ".txt"
	defer os.Remove(file)

	content := "aaa"
	WriteString(file, content)
	assert.Equal(t, ReadString(file), content)
}

func TestAppendString(t *testing.T) {
	file := "/Users/steden/Desktop/code/project/Farseer.Go/" + strconv.Itoa(rand.Intn(999-100)) + ".txt"
	defer os.Remove(file)

	WriteString(file, "aaa")
	AppendString(file, "bbb")
	readString := ReadString(file)
	assert.Equal(t, readString, "aaabbb")
}

func TestAppendLine(t *testing.T) {
	file := "/Users/steden/Desktop/code/project/Farseer.Go/" + strconv.Itoa(rand.Intn(999-100)) + ".txt"
	defer os.Remove(file)

	WriteString(file, "aaa")
	AppendLine(file, "bbb")
	readString := ReadString(file)
	assert.Equal(t, readString, "aaa\nbbb")
}

func TestAppendAllLine(t *testing.T) {
	file := "/Users/steden/Desktop/code/project/Farseer.Go/" + strconv.Itoa(rand.Intn(999-100)) + ".txt"
	defer os.Remove(file)

	WriteString(file, "aaa")
	str := []string{"bbb", "ccc"}
	AppendAllLine(file, str)
	readString := ReadString(file)
	assert.Equal(t, readString, "aaa\nbbb\nccc")
}
