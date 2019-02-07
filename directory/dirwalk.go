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
// value : `^D:/HOGEMOGE/FUGAFUGA/csharp/src/.*/bin/.+/netcoreapp2.2$`
// result : `D:/HOGEMOGE/FUGAFUGA/csharp/src/`
//
// example2.
// value : `^D:/PIYOPIYO/ghoasd.*`
// result : `^D:/PIYOPIYO/`
func GetBasePath(path string) (string, error) {
	return getBasePath(path)
}

// IsExists checks directory is exists or not
func IsExists(path string) bool {
	return isExists(path)
}

// Dirwalk walk path recursively and return array of directory full path
func Dirwalk(path string) (fullPaths []string, err error) {
	return dirwalk(path, true)
}

func getBasePath(path string) (string, error) {
	asts, err := regexast.ParseRegex(path, syntax.Perl)
	if err != nil {
		//fmt.Println(err)
		return "", err
	}

	var b strings.Builder
	begin := false

	var s string
	for _, a := range asts {
		if begin && !a.IsRune {
			// check path is valid and fix
			s = b.String()
			if isExists(s) {
				break
			}

			if getLastRune(s, 1) != "/" {
				i := strings.LastIndex(s, "/") + 1
				s = string(s[:i])
			}
			break
		}
		if begin && a.IsRune {
			b.WriteString(a.Value)
		}
		if a.IsStart {
			begin = true
		}
	}

	return s, nil
}

func isExists(path string) bool {
	if f, err := os.Stat(path); !os.IsNotExist(err) {
		return f.IsDir()
	}
	return false
}

func dirwalk(path string, toSlash bool) (fullPaths []string, err error) {
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

			walkPaths, _err := dirwalk(filepath.Join(path, file.Name()), toSlash)
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
