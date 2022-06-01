package my_find

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
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

func Find(addr string, fl *Fl) {
	err := filepath.Walk(addr,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if (*fl.D && info.IsDir()) || (*fl.F && info.IsDir() == false) {
				if info.Mode().Type()&fs.ModeSymlink == 0 {
					fmt.Println(path)
				} else {
					fmt.Printf("%s -> %s", path, getRealPath(path))
				}
			}
			return nil
		})
	if err != nil {
		log.Println(err)
	}
}

func Find_ext(addr, ext string) {
	err := filepath.Walk(addr,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() == false && filepath.Ext(path)[1:] == ext {
				if info.Mode().Type()&fs.ModeSymlink == 0 {
					fmt.Println(path)
				} else {
					fmt.Printf("%s -> %s", path, getRealPath(path))
				}
			}
			return nil
		})
	if err != nil {
		log.Println(err)
	}
}
