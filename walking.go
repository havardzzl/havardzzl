package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func WalkDir(f string, level int) ([]string, error) {
	prefix := "|"
	for i := 0; i < level; i++ {
		prefix += "------"
	}
	files, err := ioutil.ReadDir(f) // files为当前目录下的所有文件名称【包括文件夹】
	if err != nil {
		return nil, err
	}

	var allfile []string
	for _, v := range files {
		fullPath := f + "\\" + v.Name() // 全路径 + 文件名称
		if v.IsDir() {                  // 如果是目录
			allfile = append(allfile, prefix+v.Name())
			a, _ := WalkDir(fullPath, level+1) // 遍历改路径下的所有文件
			allfile = append(allfile, a...)
		} else {
			allfile = append(allfile, prefix+v.Name()) // 如果不是文件夹，就直接追加到路径下
		}
	}

	return allfile, nil
}

func Wk(f string) {
	err := filepath.Walk(f,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			relp, _ := filepath.Rel(f, path)
			fmt.Println(relp, info.Size())
			return nil
		})
	if err != nil {
		log.Println(err)
	}
}
