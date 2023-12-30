package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/konidev20/sigtool"
)

func main() {
	inParam := flag.String("in", "filename", "This specifies the input Signed PE filename to read from")
	outParam := flag.String("out", "filename", "This specifies the output PKCS#7 filename to write to")
	isVerificationRequired := flag.Bool("validate", false, "This specifies if the PKCS#7 signature of the file should be verfied")

	flag.Parse()
	if *inParam == "filename" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	buf, err := sigtool.ExtractDigitalSignature(*inParam)
	if err != nil {
		log.Fatal(err)
	}

	if *isVerificationRequired {
		err = sigtool.IsValidDigitalSignature(*inParam)
		if err != nil {
			log.Fatal(err)
		}
	}

	if *outParam != "filename" {
		os.WriteFile(*outParam, buf, 0644)
	} else {
		fileName := filepath.Base(*inParam)
		if fileName == "." {
			fmt.Println("Input file path is not correct.")
			os.Exit(1)
		}
		os.WriteFile(fileName+".pkcs7", buf, 0644)
	}
}
