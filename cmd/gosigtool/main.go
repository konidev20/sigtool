package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/konidev20/sigtool"
)

func main() {
	inParam := flag.String("in", "", "This specifies the input Signed PE filename to read from")
	outParam := flag.String("out", "", "This specifies the output PKCS#7 filename to write to")
	isVerificationRequired := flag.Bool("validate", false, "This specifies if the PKCS#7 signature of the file should be verified")

	flag.Parse()
	if *inParam == "" {
		fmt.Fprintf(os.Stderr, "Error: input file (-in) is required\n\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	buf, err := sigtool.ExtractDigitalSignature(*inParam)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error extracting signature: %v\n", err)
		os.Exit(1)
	}

	if *isVerificationRequired {
		err = sigtool.IsValidDigitalSignature(*inParam)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error validating signature: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Signature is valid")
	}

	var outputPath string
	if *outParam != "" {
		outputPath = *outParam
	} else {
		fileName := filepath.Base(*inParam)
		if fileName == "." || fileName == "" {
			fmt.Fprintf(os.Stderr, "Error: unable to determine output filename from input path\n")
			os.Exit(1)
		}
		outputPath = fileName + ".pkcs7"
	}

	if err := os.WriteFile(outputPath, buf, 0600); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing output file %q: %v\n", outputPath, err)
		os.Exit(1)
	}

	fmt.Printf("Successfully extracted signature to %q\n", outputPath)
}
