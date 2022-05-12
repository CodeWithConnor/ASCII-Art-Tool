package asciiart

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

// here we're defining a custom io.Writer by implementing a write function
type logWriter struct{}

// we're removing the date/time before the log contents so it has a clean output when we write to output.txt
func (writer logWriter) Write(bytes []byte) (int, error) {
	return fmt.Print(time.Now().UTC().Format(" [DEBUG] " + string(bytes)))
}

func GenerateArt(input string, banner string) {
	// continue only if the text and banner is passed to function
	mainArg := input
	map1 := make(map[rune]string)
	tempString := ""

	/* standard.txt contains text art corresponding to ascii
	characters (in ascending order) from ascii 32 - 126*/

	// using 'banner' POST data to determine which banner to use
	file, err := os.Open(banner + ".txt")
	// if there's an error, log it
	if err != nil {
		log.Fatalf("Failed to open banner")
	}

	// scan the file and split the contents line by line
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	// scan the contents and append it to variable 'text'
	var text []string
	for scanner.Scan() {
		text = append(text, scanner.Text())
	}
	file.Close()

	/* here we create a map record for ascii characters
	in range of 32 - 126. ascii chars in order */
	for i := 1; i < 96; i++ {
		// for each str from text slice
		// j = index
		for j, line := range text {
			// get first line of string and then returns each value
			if j >= GetFirstStr(i+31) && j < GetFirstStr(i+31)+8 {
				// concats strings together to create a complete string
				tempString += string(line)
				// start at 31 | +8 because each char is 8x8 lines
				// 31 is first value used in ascii table/art
				// choosing each char in the text file
				if j == GetFirstStr(i+31)+8 {
					continue
				} else {
					// if it reaches end of line, do \n
					tempString += string(rune(10))
				}
			}
		}

		// map creates a key of current char, then adds string of line 1 - 8
		// map is func that creates arrays inside of arrays
		// creates table and assigns ascii val accordingly
		// tempString = \n because standard.txt starts with new line
		map1[rune(i+31)] = tempString
		// clears tempString
		tempString = ""
	}

	// split main argument
	PrintAscii(SplitMainArg(mainArg), map1)
}

// skips empty lines between chars
func GetFirstStr(ascii int) int {
	return (1 + (ascii-32)*8) + (ascii - 32)
}

// splits str at each new line ('\n') then appends to slice of str
func SplitMainArg(s string) []string {
	slicedStr := make([]string, 0)

	// arg will be split into two strings in a slice
	newLineIndex := strings.Index(s, "\\n")

	// -1 to reset between each character
	if newLineIndex != -1 {
		slicedStr = append(slicedStr, s[0:newLineIndex])
		slicedStr = append(slicedStr, s[newLineIndex+2:])
	} else {
		slicedStr = append(slicedStr, s)
	}

	return slicedStr
}

// range over a slice containing words which need to be printed
func PrintAscii(text []string, map1 map[rune]string) {
	mergedLines := ""
	// word = line
	for _, word := range text {
		// prints only 8 lines because each char is always 8 lines
		for i := 0; i <= 7; i++ {
			// if it's not a new line
			if word != string(rune(10)) {
				// if j is < len(line),
				for j := 0; j < len(word); j++ {
					// concat mer
					mergedLines += ConvMapElem2Slice(map1[rune(word[j])], i)
				}

				// here we write, line by line, the generated ASCII art to the text file

				// try and create/open output.txt with read/write/create/append rights
				f, err := os.OpenFile("output.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
				if err != nil {
					log.Fatalf("error opening file: %v", err)
				}
				defer f.Close()

				// SetFlags(0) to take full control
				log.SetFlags(0)
				log.SetOutput(new(logWriter))

				// write ASCII art to file, line by line
				log.SetOutput(f)
				log.Printf("%s\n", mergedLines)
				mergedLines = ""
			}
		}
	}
}

// converts str val from map to slice of str
func ConvMapElem2Slice(s string, index int) string {
	start := 0
	sliceFromMap := make([]string, 0)

	// increase until length of string
	for i := 0; i < len(s); i++ {
		if s[i] == 10 {
			sliceFromMap = append(sliceFromMap, s[start:i])
			start = i + 1
		}
	}
	// returns slice used for printing line
	return sliceFromMap[index]
}
