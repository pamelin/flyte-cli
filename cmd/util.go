package cmd

import (
	"path/filepath"
	"github.com/HotelsDotCom/flyte/httputil"
	"io/ioutil"
	"os"
	"strings"
	"bufio"
	"bytes"
	"log"
)

func getContentType(filename string) string {
	switch filepath.Ext(filename) {
	case ".json":
		return httputil.MediaTypeJson
	case ".yaml", ".yml":
		return httputil.MediaTypeYaml
	case ".sh":
		return "application/x-sh"
	default:
		return "application/octet-stream"
	}
}

func readFile(filename string) ([]byte, error) {
	if filename == "-" {
		return ioutil.ReadAll(os.Stdin)
	}
	return ioutil.ReadFile(filename)
}

// naive way to detect file extension
// we know it should be only JSON or YAML in specific format
func detectExt(filename string, data []byte) string {

	ext := filepath.Ext(filename)
	if ext != "" {
		return ext
	}

	if strings.HasPrefix(string(data), "{") {
		return ".json"
	}

	if sniffYaml(data) {
		return ".yaml"
	}

	return ""
}

// sniff naively for yaml content
// we know what we can expect
func sniffYaml(data []byte) bool {
	if strings.HasPrefix(string(data), "---") {
		return true
	}

	bufReader := bufio.NewReader(bytes.NewReader(data))
	fl, _, err := bufReader.ReadLine()
	if err != nil {
		log.Print(err)
		return false
	}

	return strings.HasSuffix(string(fl), ":")
}
