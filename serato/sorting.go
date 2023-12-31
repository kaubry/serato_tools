package serato

import (
	"bufio"
	"fmt"
	"github.com/kaubry/serato_tools/encoding"
	"log"
	"os"
	"sort"
	"strings"
)

// SortPrefixLineInFile Read and sort lines starting with a prefix from a files with string lines
func SortPrefixLineInFile(filePath string, prefix string) error {

	prefix = ConvertPathToPrefix(prefix)
	prefix = string(encoding.EncodeUTF16(prefix, false))
	var allLines, targetLines, err = readLinesWithPrefix(filePath, prefix)
	if err != nil {
		return err
	}

	// Sort target lines alphabetically
	sort.Strings(targetLines)

	// Replace the original lines with the sorted lines
	var sortedLines = replaceLines(allLines, targetLines, prefix)

	// Write the sorted lines back to the file
	return writeLinesToFile(filePath, sortedLines)
}

func ConvertPathToPrefix(path string) string {
	pathWithoutVolume, err := RemoveVolumeFromPath(path)
	if err == nil {
		path = pathWithoutVolume
	}
	path = strings.ReplaceAll(path, "/", "%%")
	path = "[crate]" + path
	log.Printf("Prefix to be sorted: %s\n", path)
	return path
}

func readLinesWithPrefix(filePath string, prefix string) (allLines []string, filteredLines []string, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, nil, fmt.Errorf("could not open file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		allLines = append(allLines, line)

		// Check if line starts with the specific prefix
		if strings.HasPrefix(line, prefix) {
			filteredLines = append(filteredLines, line)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, nil, fmt.Errorf("error while scanning file: %w", err)
	}
	return allLines, filteredLines, nil
}

func replaceLines(allLines []string, filteredLines []string, prefix string) (sortedLines []string) {
	// Replace the original lines with the sorted lines
	var j int
	sortedLines = make([]string, len(allLines))
	for i, line := range allLines {
		if strings.HasPrefix(line, prefix) {
			sortedLines[i] = filteredLines[j]
			j++
		} else {
			sortedLines[i] = line
		}
	}
	return sortedLines
}

func writeLinesToFile(filePath string, lines []string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(writer, line)
	}
	err = writer.Flush()
	if err != nil {
		return err
	}
	return nil
}
