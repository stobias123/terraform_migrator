package main

/**
c.f.
- godoc:
  - https://godoc.org/github.com/hashicorp/hcl2/hclparse
  - https://godoc.org/github.com/hashicorp/hcl2/gohcl
- parsing: https://github.com/hashicorp/hcl2/blob/master/guide/go_parsing.rst
- decoding: https://github.com/hashicorp/hcl2/blob/master/guide/go_decoding_gohcl.rst
**/

import (
	"flag"
	"fmt"
	"github.com/stobias123/terraform_editor/types"
	"github.com/stobias123/terraform_editor/util"
	"io/ioutil"
	"log"
	"os"
)


func main() {

	//flag.String("file","","The Adhoc file you'd like to run migrations with")
	migrationDirectory := flag.String("migrationDir","","Directory containing migration files.")
	terraformDirectory := flag.String("terraformDir","","Directory containing migration files.")
	flag.Parse()
	if *migrationDirectory == "" || *terraformDirectory == "" {
		fmt.Println("call terraform-migrator -migrationDir foo -terraformDir foo")
		os.Exit(1)
	}
	var migrationConfigs []*types.Migration
	migrationFiles, err := ioutil.ReadDir(*migrationDirectory)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range migrationFiles {
		filePath := fmt.Sprintf("%s/%s",*migrationDirectory, file.Name())
		log.Println(fmt.Sprintf("[Info] trying to migrate file: %s", filePath))
		if ! file.IsDir() {
			config, err := util.LoadMigrationFile(filePath)
			if err != nil {
				log.Fatal(err)
			}
			migrationConfigs = append(migrationConfigs,config)
			fmt.Println(file.Name(), file.IsDir())
		}
	}

	terraformFiles, err := ioutil.ReadDir(*terraformDirectory)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range terraformFiles {
		filePath := fmt.Sprintf("%s/%s",*terraformDirectory, file.Name())
		log.Println(filePath)
		if ! file.IsDir() {
			hclFile, err := util.LoadHCLFile(filePath)
			if err != nil {
				log.Fatal(err)
			}
			for _, migConfig := range migrationConfigs {
				log.Println(fmt.Sprintf("[Info] trying to migrate %s", migConfig.MigrationName))
				EditProviders(hclFile.Body(), *migConfig)
			}
			f, err := os.Create(filePath)
			if err != nil{
				log.Fatalf("Problem writing file %s", filePath)
			}
			_, err = f.Write(hclFile.Bytes())
			//fmt.Printf("%s", hclFile.Bytes())
			if err != nil {
				log.Fatalf("Problem writing file %s", filePath)
			}
			err = f.Close()
			if err != nil {
				log.Fatalf("Problem closing file %s", filePath)
			}
		}
	}
}
