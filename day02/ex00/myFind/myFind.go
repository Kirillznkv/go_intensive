package my_find

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"unicode/utf8"
)

type Fl struct {
	F, D, Sl *bool
	Ext      *string
}

func getRealPath(path string) string {
	cmd := exec.Command("readlink", path)
	out, _ := cmd.Output()
	realPath := filepath.Join(filepath.Dir(path), string(out))
	return realPath
}

func checkSymlink(path string) bool {
	file, err := os.Open(path)
	if err != nil {
		return false
	}
	file.Close()
	return true
}

func getExt(path string) string {
	ext := filepath.Ext(path)
	_, i := utf8.DecodeRuneInString(ext)
	return ext[i:]
}

func conditionDir(info os.FileInfo, fl *Fl) bool {
	if *fl.D && info.IsDir() {
		return true
	}
	return false
}

func conditionFile(info os.FileInfo, fl *Fl) bool {
	if *fl.F && *fl.Ext == "" && info.IsDir() == false {
		return true
	}
	return false
}

func conditionExt(path string, info os.FileInfo, fl *Fl) bool {
	if *fl.Ext != "" && info.IsDir() == false && getExt(path) == *fl.Ext {
		return true
	}
	return false
}

func conditionSl(info os.FileInfo, fl *Fl) bool {
	if *fl.Sl && info.IsDir() == false && info.Mode().Type()&fs.ModeSymlink != 0 {
		return true
	}
	return false
}

func conditionFlags(path string, info os.FileInfo, fl *Fl) bool {
	return conditionDir(info, fl) || conditionFile(info, fl) || conditionSl(info, fl) || conditionExt(path, info, fl)
}

func Find(addr string, fl *Fl) {
	err := filepath.Walk(addr,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if path != addr && conditionFlags(path, info, fl) {
				if info.Mode().Type()&fs.ModeSymlink == 0 {
					fmt.Println(path)
				} else if *fl.Sl {
					realPath := getRealPath(path)
					if ok := checkSymlink(path); ok == false {
						realPath = "[broken]\n"
					}
					fmt.Printf("%s -> %s", path, realPath)
				}
			}
			return nil
		})
	if err != nil {
		log.Fatal(err)
	}
}
