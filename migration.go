package main

import (
	"fmt"
	"github.com/stobias123/terraform_migrator/types"
	"github.com/stobias123/terraform_migrator/util"
	"io/ioutil"
	"log"
	"os"
)

func MigrateDirectory(terraformDirectory *string, migrationConfigs []*types.Migration)  {
	if terraformDirectory == nil {
		log.Fatal("You must pass non nil terraformDirectory to this function")
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
				EditModules(hclFile.Body(), *migConfig)
				EditResource(hclFile.Body(), *migConfig)
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