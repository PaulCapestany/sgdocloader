package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	couch "github.com/tleyden/go-couch"
)

func main() {
	url := flag.String("u", "http://127.0.0.1:4984", "Sync Gateway URL")
	bucketname := flag.String("b", "mybucket", "The name of the bucket")
	flag.Parse()
	args := flag.Args()

	if len(args) < 1 {
		fmt.Println("Usage: sgdocloader -u http://127.0.0.1:4984 -b mybucket [files and/or directories]")
		os.Exit(1)
	}

	db, err := couch.Connect(*url + "/" + *bucketname)
	if err != nil {
		fmt.Println(err)
	}

	for _, arg := range args {
		if fileInfo, err := os.Stat(arg); err == nil {
			if fileInfo.IsDir() {
				dir, err := os.Open(arg)
				if err != nil {
					fmt.Println(err)
				}
				defer dir.Close()

				filenames, err := dir.Readdirnames(0)
				if err != nil {
					fmt.Println(err)
				}
				for _, filename := range filenames {
					loadJSON(db, dir.Name()+"/"+filename)
				}
			} else {
				loadJSON(db, arg)
			}
		} else {
			fmt.Printf("Error: \"%v\" - no such file or directory exists!\n", arg)
		}
	}
}

func loadJSON(db couch.Database, filename string) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Error:", err)
	}

	var document interface{}
	err = json.Unmarshal(file, &document)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		db.Insert(document)
	}
}
