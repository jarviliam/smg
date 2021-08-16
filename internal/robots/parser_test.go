package robots_test

//TODO Redo
import (
	"io"
	"os"
)

func LoadAsset(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		return nil
	}
	data, err := io.ReadAll(file)
	if err != nil {
	}

}

func TestNewParser() {

}

func TestParseHeader() {

}
