package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	couch "github.com/tleyden/go-couch"
)

var (
	nodeAddress     = flag.String("n", "http://127.0.0.1:4984", "Sync Gateway node address")
	bucket          = flag.String("b", "mybucket", "The name of the bucket")
	useBulkDocs     = flag.Bool("k", false, "Use _bulk_docs")
	recursiveSearch = flag.Bool("r", false, "Recursively search subdirectories for documents to add")
)

// TODO: use same flags and offer same capabilities as `cbdocloader` - http://docs.couchbase.com/admin/admin/CLI/cbdocloader_tool.html
// --------------------------------------------------------------------------------------
// Usage: cbdocloader [options] <directory>|zipfile
//
// Example: cbdocloader -u Administrator -p password -n 127.0.0.1:8091 -b mybucket -s 100 gamesim-sample.zip
//
// Options:
//   -h, --help         show this help message and exit
//   -u Administrator   Username
//   -p password        Password
//   -b mybucket        Bucket
//   -n 127.0.0.1:8091  Node address
//   -s 100             RAM quota in MB

// TODO: add tests
func main() {
	flag.Parse()
	filesOrDirs := flag.Args()

	if len(filesOrDirs) < 1 {
		log.Println("Usage: sgdocloader -n http://127.0.0.1:4984 -b mybucket [files and/or directories]")
		os.Exit(1)
	}

	db, err := couch.Connect(*nodeAddress + "/" + *bucket)
	if err != nil {
		log.Println(err)
	}

	var recursedFiles []string

	for _, arg := range filesOrDirs {
		if *recursiveSearch {
			err := filepath.Walk(arg, func(path string, f os.FileInfo, err error) error {
				if !f.IsDir() {
					recursedFiles = append(recursedFiles, path)
				}
				return nil
			})
			if err != nil {
				log.Printf("Error (%v): no such file or directory exists!\n", arg)
			}
		}
	}

	for _, thing := range recursedFiles {
		filesOrDirs = append(filesOrDirs, thing)
	}

	for _, arg := range filesOrDirs {

		if fileInfo, err := os.Stat(arg); err == nil {
			// TODO: use goroutines to load data faster (limited by maxfiles?)
			if fileInfo.IsDir() && !*recursiveSearch {
				dir, err := os.Open(arg)
				if err != nil {
					log.Printf("Error: %v\n", err)
				}
				defer dir.Close()

				filenames, err := dir.Readdirnames(0)
				if err != nil {
					log.Println(err)
				}
				for _, filename := range filenames {
					loadJSON(db, filepath.Join(dir.Name(), filename))
				}
			} else {
				loadJSON(db, arg)
			}
		} else {
			log.Printf("Error (%v): no such file or directory exists!\n", arg)
		}
	}
}

func loadJSON(db couch.Database, filename string) {
	baseName := filepath.Base(filename)
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Printf("Error (%v): %v\n", baseName, err)
	}

	var document interface{}
	err = json.Unmarshal(file, &document)
	if err != nil {
		log.Printf("Error (%v): %v\n", baseName, err)
	} else {
		// TODO: need to flesh out type determination/handling…
		switch document.(type) {
		case map[string]interface{}:
			log.Println("map[string]interface{}")
			if docs, ok := document.(map[string]interface{})["docs"].([]interface{}); ok && *useBulkDocs {
				_, err := db.Bulk(docs)
				if err != nil {
					log.Printf("Error (%v): %v\n", baseName, err)
				}
			} else {
				_, _, err := db.Insert(document)
				if err != nil {
					log.Printf("Error (%v): %v\n", baseName, err)
				}
			}
		default:
			log.Printf("Error (%v): %v\n", baseName, err)
		}
	}
}
