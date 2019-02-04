package directory

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp/syntax"
	"strings"
	"unicode/utf8"

	"github.com/guitarrapc/watchdog-symlinker/regexast"
)

// GetBasePath extract basepath from regexp path string
//
// example1.
// value : `^D:/GitHub/guitarrapc/MixedContentChecker/csharp/src/.*/bin/.+/netcoreapp2.2$`
// result : `D:/GitHub/guitarrapc/MixedContentChecker/csharp/src/`
//
// example2.
// value : `^D:/GitHub/ghoasd.*`
// result : `^D:/GitHub/`
func GetBasePath(path string) (string, error) {
	asts, err := regexast.ParseRegex(path, syntax.Perl)
	if err != nil {
		//fmt.Println(err)
		return "", err
	}

	var b strings.Builder
	begin := false
	// for _, a := range asts {
	// 	fmt.Println(a)
	// }
	for _, a := range asts {
		if begin && !a.IsRune {
			break
		}
		if begin && a.IsRune {
			b.WriteString(a.Value)
		}
		if a.IsStart {
			begin = true
		}
	}

	// check path is valid and fix
	r := b.String()
	if getLastRune(r, 1) != "/" {
		l := strings.LastIndex(r, "/")
		r = substring(r, 0, l+1)
	}
	return r, nil
}

// IsExists checks directory is exists or not
func IsExists(path string) bool {
	if f, err := os.Stat(path); !os.IsNotExist(err) {
		return f.IsDir()
	}
	return false
}

// Dirwalk walk path recursively and return array of directory full path
func Dirwalk(path string, toSlash bool) (fullPaths []string, err error) {
	fullPaths = nil
	files, _err := ioutil.ReadDir(path)
	if _err != nil {
		err = _err
		return
	}

	for _, file := range files {
		if file.IsDir() {
			if toSlash {
				fullPaths = append(fullPaths, filepath.ToSlash(filepath.Join(path, file.Name())))
			} else {
				fullPaths = append(fullPaths, filepath.Join(path, file.Name()))
			}

			walkPaths, _err := Dirwalk(filepath.Join(path, file.Name()), toSlash)
			fullPaths = append(fullPaths, walkPaths...)
			err = _err
			continue
		}
	}
	return
}

func getLastRune(s string, c int) string {
	j := len(s)
	for i := 0; i < c && j > 0; i++ {
		_, size := utf8.DecodeLastRuneInString(s[:j])
		j -= size
	}
	return s[j:]
}

func substring(str string, start int, length int) string {
	if start < 0 || length <= 0 {
		return str
	}
	r := []rune(str)
	if start+length > len(r) {
		return string(r[start:])
	} else {
		return string(r[start : start+length])
	}
}
