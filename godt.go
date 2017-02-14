// godt is a program that operates on .odt documents in the current directory.
// To see all available options run 'godt help'
package main

import (
	"archive/zip"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func unzipDocument(file string) *zip.ReadCloser {
	unzipped, err := zip.OpenReader(file)
	check(err)
	return unzipped
}

func removeExtension(f os.FileInfo) {
	name := f.Name()
	lastChars := ""
	if len(name) > 4 {
		lastChars = name[(len(name) - 4):]
	}
	if lastChars == ".odt" {
		nameTruncated := name[0:(len(name) - 4)]
		err := os.Rename(name, nameTruncated)
		check(err)
	}
}

func addExtension(f os.FileInfo) {
	name := f.Name()
	lastChars := ""
	if len(name) >= 4 {
		lastChars = name[(len(name) - 4):]
	}
	r, err := os.Open(name)
	check(err)
	defer r.Close()
	var header [2]byte
	io.ReadFull(r, header[:])
	if lastChars != ".odt" && string(header[:]) == "PK" {
		nameAppended := name + ".odt"
		document := unzipDocument(name)
		defer document.Close()
		for _, file := range document.File {
			fileName := file.FileHeader.Name
			if fileName == "mimetype" {
				r, err := file.Open()
				mimeType, err := ioutil.ReadAll(r)
				check(err)
				if string(mimeType) == "application/vnd.oasis.opendocument.text" {
					err := os.Rename(name, nameAppended)
					check(err)
				}
			}
		}
	}
}

type Stats struct {
	XMLName          xml.Name
	PageCount        int `xml:"page-count,attr"`
	WordCount        int `xml:"word-count,attr"`
	CharacterCount   int `xml:"character-count,attr"`
	ParagraphCount   int `xml:"paragraph-count,attr"`
	ImageCount       int `xml:"image-count,attr"`
	TableCount       int `xml:"table-count,attr"`
	NWCharacterCount int `xml:"non-whitespace-character-count,attr"`
	ObjectCount      int `xml:"object-count,attr"`
}

type Meta struct {
	XMLName xml.Name
	Date    string `xml:"creation-date"`
	Title   string `xml:"title"`
	Stats   Stats  `xml:"document-statistic"`
}

type DocumentMeta struct {
	XMLName xml.Name
	Meta    Meta `xml:"meta"`
}

// Extracts metadata from any .odt file that has meta.xml inside it.
func extractMetadata(f os.FileInfo) (DocumentMeta, error) {
	var dm DocumentMeta
	document := unzipDocument(f.Name())
	defer document.Close()
	for _, file := range document.File {
		fileName := file.FileHeader.Name
		if fileName == "meta.xml" {
			r, err := file.Open()
			data, err := ioutil.ReadAll(r)
			check(err)
			xml.Unmarshal(data, &dm)
			return dm, nil
		}
	}
	err := errors.New("This zip file does not contain meta.xml.")
	return dm, err
}

// Sorts a map of document statistics and prints out each element
func sortAndPrint(docValues map[string]int) {
	var sortedMap = map[int][]string{}
	var a []int
	for k, v := range docValues {
		sortedMap[v] = append(sortedMap[v], k)
	}
	for k := range sortedMap {
		a = append(a, k)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(a)))
	for _, k := range a {
		for _, s := range sortedMap[k] {
			spaces := strings.Repeat(" ", (15 - len(strconv.Itoa(k))))
			fmt.Printf("%d"+spaces+"%s\n", k, s)
		}
	}
}

func findInside(name string, fileList []os.FileInfo) bool {
	for _, f := range fileList {
		if name == f.Name() {
			return true
		}
	}
	return false
}

func main() {
	arg := os.Args
	if len(arg) <= 1 {
		fmt.Println("Error: godt expects an argument.\n" + "Run 'godt help' to see all possible arguments.")
		os.Exit(1)
	}
	switch arg[1] {
	case "rmex":
		files, err := ioutil.ReadDir("./")
		check(err)
		for _, f := range files {
			name := f.Name()
			if len(name) > 4 {
				truncatedName := name[0:(len(name) - 4)]
				if findInside(truncatedName, files) {
					fmt.Println("Error: "+truncatedName+" already exists in the current directory.")
					os.Exit(1)
				}
			}
			if !f.IsDir() {
				removeExtension(f)
			}
		}
	case "adex":
		files, err := ioutil.ReadDir("./")
		check(err)
		for _, f := range files {
			if findInside(f.Name()+".odt", files) {
				fmt.Println("Error: "+f.Name()+".odt"+" already exists in the current directory.")
				os.Exit(1)
			}
			if !f.IsDir() {
				addExtension(f)
			}
		}
	case "lsdc":
		var docValues = map[string]int{}
		files, err := ioutil.ReadDir("./")
		check(err)
		for _, f := range files {
			name := f.Name()
			r, err := os.Open(name)
			check(err)
			var header [2]byte
			io.ReadFull(r, header[:])
			// Makes sure that the file is not a directory and is a valid .zip file
			if !f.IsDir() && string(header[:]) == "PK" {
				metaData, err := extractMetadata(f)
				if err != nil {
					continue
				}
				// Special case formatting for the date string
				date := metaData.Meta.Date[:10]
				compactDate := date[:4] + date[5:7] + date[8:]
				dateAsInteger, err := strconv.Atoi(compactDate)
				check(err)
				docValues[name] = dateAsInteger
			}
			r.Close()
		}
		// Sorts the documents by date created
		// Special case sort for the date format
		if len(docValues) != 0 {
			var sortedMap = map[int][]string{}
			var a []int
			for k, v := range docValues {
				sortedMap[v] = append(sortedMap[v], k)
			}
			for k := range sortedMap {
				a = append(a, k)
			}
			sort.Sort(sort.Reverse(sort.IntSlice(a)))
			for _, k := range a {
				for _, s := range sortedMap[k] {
					spaces := strings.Repeat(" ", 10)
					date := strconv.Itoa(k)
					formattedDate := date[6:] + "/" + date[4:6] + "/" + date[:4]
					fmt.Printf("%s"+spaces+"%s\n", formattedDate, s)
				}
			}
		}
	case "lswd":
		var docValues = map[string]int{}
		files, err := ioutil.ReadDir("./")
		check(err)
		for _, f := range files {
			name := f.Name()
			r, err := os.Open(name)
			check(err)
			var header [2]byte
			io.ReadFull(r, header[:])
			// Makes sure that the file is not a directory and is a valid .zip file
			if !f.IsDir() && string(header[:]) == "PK" {
				metaData, err := extractMetadata(f)
				if err != nil {
					continue
				}
				docValues[name] = metaData.Meta.Stats.WordCount
			}
			r.Close()
		}
		// Sorts the documents by word count
		if len(docValues) != 0 {
			sortAndPrint(docValues)
		}
	case "lsch":
		var docValues = map[string]int{}
		files, err := ioutil.ReadDir("./")
		check(err)
		for _, f := range files {
			name := f.Name()
			r, err := os.Open(name)
			check(err)
			var header [2]byte
			io.ReadFull(r, header[:])
			// Makes sure that the file is not a directory and is a valid .zip file
			if !f.IsDir() && string(header[:]) == "PK" {
				metaData, err := extractMetadata(f)
				if err != nil {
					continue
				}
				docValues[name] = metaData.Meta.Stats.CharacterCount
			}
			r.Close()
		}
		// Sorts the documents by character count
		if len(docValues) != 0 {
			sortAndPrint(docValues)

		}
	case "lspg":
		var docValues = map[string]int{}
		files, err := ioutil.ReadDir("./")
		check(err)
		for _, f := range files {
			name := f.Name()
			r, err := os.Open(name)
			check(err)
			var header [2]byte
			io.ReadFull(r, header[:])
			// Makes sure that the file is not a directory and is a valid .zip file
			if !f.IsDir() && string(header[:]) == "PK" {
				metaData, err := extractMetadata(f)
				if err != nil {
					continue
				}
				docValues[name] = metaData.Meta.Stats.PageCount
			}
			r.Close()
		}
		// Sorts the documents by page count
		if len(docValues) != 0 {
			sortAndPrint(docValues)
		}
	case "lspa":
		var docValues = map[string]int{}
		files, err := ioutil.ReadDir("./")
		check(err)
		for _, f := range files {
			name := f.Name()
			r, err := os.Open(name)
			check(err)
			var header [2]byte
			io.ReadFull(r, header[:])
			// Makes sure that the file is not a directory and is a valid .zip file
			if !f.IsDir() && string(header[:]) == "PK" {
				metaData, err := extractMetadata(f)
				if err != nil {
					continue
				}
				docValues[name] = metaData.Meta.Stats.ParagraphCount
			}
			r.Close()
		}
		// Sorts the documents by paragraph count
		if len(docValues) != 0 {
			sortAndPrint(docValues)
		}
	case "lsim":
		var docValues = map[string]int{}
		files, err := ioutil.ReadDir("./")
		check(err)
		for _, f := range files {
			name := f.Name()
			r, err := os.Open(name)
			check(err)
			var header [2]byte
			io.ReadFull(r, header[:])
			// Makes sure that the file is not a directory and is a valid .zip file
			if !f.IsDir() && string(header[:]) == "PK" {
				metaData, err := extractMetadata(f)
				if err != nil {
					continue
				}
				docValues[name] = metaData.Meta.Stats.ImageCount
			}
			r.Close()
		}
		// Sorts the documents by image count
		if len(docValues) != 0 {
			sortAndPrint(docValues)
		}
	case "lstb":
		var docValues = map[string]int{}
		files, err := ioutil.ReadDir("./")
		check(err)
		for _, f := range files {
			name := f.Name()
			r, err := os.Open(name)
			check(err)
			var header [2]byte
			io.ReadFull(r, header[:])
			// Makes sure that the file is not a directory and is a valid .zip file
			if !f.IsDir() && string(header[:]) == "PK" {
				metaData, err := extractMetadata(f)
				if err != nil {
					continue
				}
				docValues[name] = metaData.Meta.Stats.TableCount
			}
			r.Close()
		}
		// Sorts the documents by table count
		if len(docValues) != 0 {
			sortAndPrint(docValues)
		}
	case "lsnw":
		var docValues = map[string]int{}
		files, err := ioutil.ReadDir("./")
		check(err)
		for _, f := range files {
			name := f.Name()
			r, err := os.Open(name)
			check(err)
			var header [2]byte
			io.ReadFull(r, header[:])
			// Makes sure that the file is not a directory and is a valid .zip file
			if !f.IsDir() && string(header[:]) == "PK" {
				metaData, err := extractMetadata(f)
				if err != nil {
					continue
				}
				docValues[name] = metaData.Meta.Stats.NWCharacterCount
			}
			r.Close()
		}
		// Sorts the documents by non-whitespace character count
		if len(docValues) != 0 {
			sortAndPrint(docValues)
		}
	case "lsob":
		var docValues = map[string]int{}
		files, err := ioutil.ReadDir("./")
		check(err)
		for _, f := range files {
			name := f.Name()
			r, err := os.Open(name)
			check(err)
			var header [2]byte
			io.ReadFull(r, header[:])
			// Makes sure that the file is not a directory and is a valid .zip file
			if !f.IsDir() && string(header[:]) == "PK" {
				metaData, err := extractMetadata(f)
				if err != nil {
					continue
				}
				docValues[name] = metaData.Meta.Stats.ObjectCount
			}
			r.Close()
		}
		// Sorts the documents by object count
		if len(docValues) != 0 {
			sortAndPrint(docValues)
		}
	case "help":
		fmt.Println("godt operates on .odt documents in the current directory\n\n" +
			"godt [argument]\n\n" +
			"Possible arguments:\n\n" +
			"rmex - Remove .odt extensions\n" +
			"adex - Add .odt extensions\n" +
			"lsdc - List by date created\n" +
			"lswd - List by word count\n" +
			"lsch - List by character count\n" +
			"lspg - List by page count\n" +
			"lspa - List by paragraph count\n" +
			"lsim - List by image count\n" +
			"lstb - List by table count\n" +
			"lsnw - List by non-white-space character count\n" +
			"lsob - List by object count\n" +
			"help - Show help")
	default:
		fmt.Println("Error: argument not recognized.\n" + "Run 'godt help' to see all possible arguments.")
		os.Exit(1)
	}

}
