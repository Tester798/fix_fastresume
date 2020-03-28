package main

import (
	"fmt"
	"github.com/marksamman/bencode"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	files, err := filepath.Glob("*.fastresume")
	if err != nil {
		log.Fatal(err)
	}

	if len(files) == 0 {
		fmt.Println("No fastresume files found in current directory")
		os.Exit(2)
	}

	for _, file := range files {
		fmt.Println("\nProcessing file", file)

		f, err := os.Open(file)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		dict, err := bencode.Decode(f)
		if err != nil {
			log.Fatal(err)
		}

		if dict["save_path"] == nil {
			fmt.Println("  No save_path found")
			continue
		}
		save_path := dict["save_path"].(string)
		fmt.Println("  save_path =", save_path)

		re := regexp.MustCompile(`^.+(//.+)$`)
		save_path_new := re.ReplaceAllString(save_path, `$1`)
		save_path_new = strings.ReplaceAll(save_path_new, "/", "\\")

		if save_path == save_path_new {
			continue
		}

		fmt.Println("  changing to")
		fmt.Println("  save_path =", save_path_new)

		dict["save_path"] = save_path_new
		file_out_str := bencode.Encode(dict)
		ioutil.WriteFile(file, file_out_str, 0644)
	}

	fmt.Println("\nDone")
}
