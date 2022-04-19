package main

import (
	"flag"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"regexp"
)

var root = flag.String("root", "~", "the root folder from which to start listing the directories")

type set struct {
	values map[string]string
}

func (s *set) Put(value string) {
	if s.values == nil {
		s.values = map[string]string{}
	}

	_, ok := s.values[value]
	if !ok {
		s.values[value] = value
	}
}

func (s *set) Values() map[string]string {
	if s.values == nil {
		s.values = map[string]string{}
	}

	return s.values
}

var s = set{}

func main() {
	flag.Parse()

	if *root == "~" {
		*root = os.Getenv("HOME")
	}

	if err := listFiles(*root); err != nil {
		log.Fatal(err)
	}

	for _, directory := range s.Values() {
		fmt.Println(directory)
	}
}

func getFiles(dir string) ([]fs.FileInfo, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	return files, nil
}

func listFiles(dir string) error {
	files, err := getFiles(dir)
	if err != nil {
		return err
	}

	for _, file := range files {
		isDir := file.IsDir()
		name := file.Name()
		if isDir == false {
			continue
		}
		path := fmt.Sprintf("%s/%s", dir, name)
		s.Put(path)
		listGitProjects(path, 3)
	}

	return nil
}

var r *regexp.Regexp
var err error

func listGitProjects(root string, maxLevel int) error {
	if r == nil {
		r, err = regexp.Compile(`^\.git$|^HEAD$`)
		if err != nil {
			return err
		}
	}

	var recursive func(dir string, level int) error

	recursive = func(dir string, level int) error {
		if level >= maxLevel {
			return nil
		}
		files, err := getFiles(dir)
		if err != nil {
			return err
		}

		for _, file := range files {
			name := file.Name()
			path := fmt.Sprintf("%s/%s", dir, name)
			if r.MatchString(name) == true {
				s.Put(dir)
				continue
			}
			recursive(path, level+1)
		}

		return nil
	}

	return recursive(root, 0)
}
