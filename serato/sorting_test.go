package serato

import (
	"bufio"
	"github.com/kaubry/serato_tools/encoding"
	"io/ioutil"
	"os"
	"testing"
)

func createTemporaryTestFile(t *testing.T) (f *os.File) {
	tempFile, err := ioutil.TempFile("", "prefixTest")
	if err != nil {
		t.Fatalf("Failed to create temp file: %s", err)
	}

	// Sample file content
	content := []string{
		"[crate]DJ Yuma%%Latino%%Salsa%%New",
		"[crate]DJ Yuma%%Latino%%Salsa%%On2",
		"[crate]DJ Yuma%%Latino%%Salsa%%Other",
		"[crate]DJ Yuma%%Latino%%Salsa%%Pop",
		"[crate]DJ Yuma%%Latino%%Bachata",
		"[crate]DJ Yuma%%Latino%%Albums",
		"[crate]DJ Yuma%%Latino%%Albums%%A Man And His Music _ The Player _ Willie Colón",
		"[crate]DJ Yuma%%Latino%%Albums%%Alexander Abreu, Mayito Rivera y Alain Pérez _ A Romper El Coco (2019)",
		"[crate]DJ Yuma%%Latino%%Albums%%Azúcar Negra _ Bailando Sin Parar (2022)",
		"[crate]DJ Yuma%%Latino%%Albums%%Ricky Campanelli",
		"[crate]DJ Yuma%%Latino%%Albums%%Ricky Campanelli%%Nos Curamos Con la Rumba _ 2023",
		"[crate]DJ Yuma%%Latino%%Albums%%Ricky Campanelli%%Alma de Rumbero _ 2018",
		"[crate]DJ Yuma%%Latino%%Albums%%Leslie Grace",
		"[crate]DJ Yuma%%Latino%%Albums%%El Noro Y Primera Clase",
		"[crate]DJ Yuma%%Latino%%Albums%%El Noro Y Primera Clase%%Sin Escala _ 2015",
		"[crate]DJ Yuma%%Latino%%Albums%%El Noro Y Primera Clase%%El Espejo _ 2020",
		"[crate]DJ Yuma%%Latino%%Albums%%Sonora Poncena`"}

	w := bufio.NewWriter(tempFile)
	for _, s := range convertStringArrayToUTF16(content) {
		w.WriteString(s + "\n")
	}
	w.Flush()

	if err != nil {
		t.Fatalf("Failed to write to temp file: %s", err)
	}
	tempFile.Close()

	return tempFile
}

func TestReadLinesWithPrefix(t *testing.T) {
	// Create a temporary file
	var tempFile = createTemporaryTestFile(t)
	defer os.Remove(tempFile.Name())

	// Test readLinesWithPrefix function
	prefix := "[crate]DJ Yuma%%Latino%%Albums%%"
	prefix = string(encoding.EncodeUTF16(prefix, false))
	allLines, filteredLines, err := readLinesWithPrefix(tempFile.Name(), prefix)
	if err != nil {
		t.Fatalf("readLinesWithPrefix returned an error: %s", err)
	}

	// Verify all lines
	expectedAllLines := []string{
		"[crate]DJ Yuma%%Latino%%Salsa%%New",
		"[crate]DJ Yuma%%Latino%%Salsa%%On2",
		"[crate]DJ Yuma%%Latino%%Salsa%%Other",
		"[crate]DJ Yuma%%Latino%%Salsa%%Pop",
		"[crate]DJ Yuma%%Latino%%Bachata",
		"[crate]DJ Yuma%%Latino%%Albums",
		"[crate]DJ Yuma%%Latino%%Albums%%A Man And His Music _ The Player _ Willie Colón",
		"[crate]DJ Yuma%%Latino%%Albums%%Alexander Abreu, Mayito Rivera y Alain Pérez _ A Romper El Coco (2019)",
		"[crate]DJ Yuma%%Latino%%Albums%%Azúcar Negra _ Bailando Sin Parar (2022)",
		"[crate]DJ Yuma%%Latino%%Albums%%Ricky Campanelli",
		"[crate]DJ Yuma%%Latino%%Albums%%Ricky Campanelli%%Nos Curamos Con la Rumba _ 2023",
		"[crate]DJ Yuma%%Latino%%Albums%%Ricky Campanelli%%Alma de Rumbero _ 2018",
		"[crate]DJ Yuma%%Latino%%Albums%%Leslie Grace",
		"[crate]DJ Yuma%%Latino%%Albums%%El Noro Y Primera Clase",
		"[crate]DJ Yuma%%Latino%%Albums%%El Noro Y Primera Clase%%Sin Escala _ 2015",
		"[crate]DJ Yuma%%Latino%%Albums%%El Noro Y Primera Clase%%El Espejo _ 2020",
		"[crate]DJ Yuma%%Latino%%Albums%%Sonora Poncena",
	}
	expectedAllLines = convertStringArrayToUTF16(expectedAllLines)
	for i, line := range allLines {
		expected := expectedAllLines[i] + "\n"
		if expected != line {
			e := []byte(expectedAllLines[i])
			l := []byte(line)
			t.Errorf("Hex dump - got:\n%x, want:\n%x", l, e)
			//for y, el := range e {
			//	if el != l[y] {
			//		t.Errorf("Unexpected allLines at index %d: got %v, want %v\n", i, line, expectedAllLines[i])
			//	}
			//}
			//print(e, l)
			t.Errorf("Unexpected allLines at index %d: got %v, want %v\n", i, line, expectedAllLines[i])
		}
	}

	// Verify filtered lines
	expectedFilteredLines := []string{
		"[crate]DJ Yuma%%Latino%%Albums%%A Man And His Music _ The Player _ Willie Colón",
		"[crate]DJ Yuma%%Latino%%Albums%%Alexander Abreu, Mayito Rivera y Alain Pérez _ A Romper El Coco (2019)",
		"[crate]DJ Yuma%%Latino%%Albums%%Azúcar Negra _ Bailando Sin Parar (2022)",
		"[crate]DJ Yuma%%Latino%%Albums%%Ricky Campanelli",
		"[crate]DJ Yuma%%Latino%%Albums%%Ricky Campanelli%%Nos Curamos Con la Rumba _ 2023",
		"[crate]DJ Yuma%%Latino%%Albums%%Ricky Campanelli%%Alma de Rumbero _ 2018",
		"[crate]DJ Yuma%%Latino%%Albums%%Leslie Grace",
		"[crate]DJ Yuma%%Latino%%Albums%%El Noro Y Primera Clase",
		"[crate]DJ Yuma%%Latino%%Albums%%El Noro Y Primera Clase%%Sin Escala _ 2015",
		"[crate]DJ Yuma%%Latino%%Albums%%El Noro Y Primera Clase%%El Espejo _ 2020",
		"[crate]DJ Yuma%%Latino%%Albums%%Sonora Poncena",
	}
	for i, line := range filteredLines {
		if expectedFilteredLines[i]+"\n" != line {
			t.Errorf("Unexpected filteredLines at index %d: got %v, want %v\n", i, line, expectedFilteredLines[i])
		}
	}
}

func TestSortPrefixLineInFile(t *testing.T) {
	var tempFile = createTemporaryTestFile(t)
	defer os.Remove(tempFile.Name())
	prefix := "[crate]DJ Yuma%%Latino%%Albums%%"
	SortPrefixLineInFile(tempFile.Name(), prefix)

	sortedLines := []string{
		"[crate]DJ Yuma%%Latino%%Salsa%%New",
		"[crate]DJ Yuma%%Latino%%Salsa%%On2",
		"[crate]DJ Yuma%%Latino%%Salsa%%Other",
		"[crate]DJ Yuma%%Latino%%Salsa%%Pop",
		"[crate]DJ Yuma%%Latino%%Bachata",
		"[crate]DJ Yuma%%Latino%%Albums",
		"[crate]DJ Yuma%%Latino%%Albums%%A Man And His Music _ The Player _ Willie Colón",
		"[crate]DJ Yuma%%Latino%%Albums%%Alexander Abreu, Mayito Rivera y Alain Pérez _ A Romper El Coco (2019)",
		"[crate]DJ Yuma%%Latino%%Albums%%Azúcar Negra _ Bailando Sin Parar (2022)",
		"[crate]DJ Yuma%%Latino%%Albums%%El Noro Y Primera Clase",
		"[crate]DJ Yuma%%Latino%%Albums%%El Noro Y Primera Clase%%El Espejo _ 2020",
		"[crate]DJ Yuma%%Latino%%Albums%%El Noro Y Primera Clase%%Sin Escala _ 2015",
		"[crate]DJ Yuma%%Latino%%Albums%%Leslie Grace",
		"[crate]DJ Yuma%%Latino%%Albums%%Ricky Campanelli",
		"[crate]DJ Yuma%%Latino%%Albums%%Ricky Campanelli%%Alma de Rumbero _ 2018",
		"[crate]DJ Yuma%%Latino%%Albums%%Ricky Campanelli%%Nos Curamos Con la Rumba _ 2023",
		"[crate]DJ Yuma%%Latino%%Albums%%Sonora Poncena",
	}
	sortedLines = convertStringArrayToUTF16(sortedLines)

	allLines, _, _ := readLinesWithPrefix(tempFile.Name(), prefix)

	for i, line := range allLines {
		if sortedLines[i] != line {
			t.Errorf("Unexpected allLines at index %d: got %v, want %v\n", i, line, sortedLines[i])
		}
	}
}

func TestConvertPathToPrefix(t *testing.T) {
	path := "DJ Yuma/Latino/Albums/"
	expected := "[crate]DJ Yuma%%Latino%%Albums%%"
	actual := ConvertPathToPrefix(path)
	if actual != expected {
		t.Errorf("Wrong path convertion: got %v, want %v\n", actual, expected)
	}
}

func TestConvertPathToPrefixWithVolume(t *testing.T) {
	path := "/Volumes/TestVolume1/DJ Yuma/Latino/Albums/"
	expected := "[crate]DJ Yuma%%Latino%%Albums%%"
	actual := ConvertPathToPrefix(path)
	if actual != expected {
		t.Errorf("Wrong path convertion: got %v, want %v\n", actual, expected)
	}
}

func convertStringArrayToUTF16(content []string) []string {
	var converted []string
	for _, line := range content {
		converted = append(converted, string(encoding.EncodeUTF16(line, false)))
	}
	return converted
}
