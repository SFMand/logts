package cmd

import "fmt"

func targzComp(fromPath *string, toPath *string) error {
	fmt.Println(*fromPath, *toPath)
	return nil
}
