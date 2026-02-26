package utils

import (
	"fmt"
	"log"
	"os"
)

/*
* Prints the man page specified to the console
 */
func PrintManPage(page string) {
	file := fmt.Sprintf("Man/%s.txt", page)

	// Read man page
	data, err := os.ReadFile(file)
	if err != nil {
		log.Fatal("Error reading man page: ", err)
	}

	// Output man page
	fmt.Println(string(data))
}
