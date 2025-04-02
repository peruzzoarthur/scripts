package file

import (
	"bufio"
	"fmt"
	"os"
	// "os/exec"
	// "path/filepath"
	"strings"
)

func GetFilename() string {
	// reader := bufio.NewReader(os.Stdin)
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Enter a filename: ")

		scanned := scanner.Scan()
		if scanned {
			text := scanner.Text()
			if len(text) > 0 {
				fmt.Print(text)
				fmt.Print('\n')
				fmt.Print(strings.TrimSpace(text))
				return strings.TrimSpace(text)
			}
		}

		fmt.Println("Error: Filename is empty")
	}
}

func main() {
	a := GetFilename()
	fmt.Printf("%s\n", a)
}
