package utils

import (
	"fmt"
	"os"
)

func CreateFile(name string) (*os.File, error) {
	file, err := os.Create("hello.txt")

	if err != nil {
		fmt.Println("Unable to create file:", err)
		return file, err
	}

	return file, nil
}

func WriteFile(text string, file *os.File) error {
	_, err := file.WriteString(text)
	if err != nil {
		fmt.Println("Not Write in file", text, err)
		return err
	}

	return nil
}
