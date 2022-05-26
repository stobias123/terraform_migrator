package util

import (
	"fmt"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/stobias123/terraform_migrator/types"
	"github.com/zclconf/go-cty/cty"
	"io/ioutil"
	"log"
	"os"
	"strings"
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
		if ! file.IsDir() && strings.Contains(file.Name(), "tf") {
			hclFile, err := LoadHCLFile(filePath)
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


func EditProviders(body *hclwrite.Body, editConfig types.Migration) {
	for _, provider := range editConfig.ProviderBlocks {
		blocks := FindBlocks(body, "provider" , []string{provider.Name})
		for _, block := range blocks {
			editAttributes(block, provider.Attributes)
			//todo - edit sub blocks etc.
		}
	}
}

func EditModules(body *hclwrite.Body, editConfig types.Migration) {
	for _, module := range editConfig.ModuleBlocks {
		blocks := FindBlocks(body, "module" , []string{module.Name})
		for _, block := range blocks {
			editAttributes(block, module.Attributes)
			//todo - edit sub blocks etc.
		}
	}
}

func EditResource(body *hclwrite.Body, editConfig types.Migration) {
	for _, resource := range editConfig.ResourceBlocks {
		blocks := FindBlocks(body, "resource" , []string{resource.Name})
		for _, block := range blocks {
			editAttributes(block, resource.Attributes)
			//todo - edit sub blocks etc.
		}
	}
}

func editAttributes(block *hclwrite.Block, attributes []types.AttributeConfig) {
	for _, attribute := range attributes {
		switch {
		case attribute.Action == types.Delete:
			block.Body().RemoveAttribute(attribute.Name)
		case attribute.Action == types.Update:
			block.Body().SetAttributeValue(attribute.Name, cty.StringVal(*attribute.Value))
		case attribute.Action == types.Replace:
			if attribute.OriginalValue == nil || attribute.Value == nil {
				log.Fatalf("[ERROR] Block: %s - Attribute: %s. You must set attribute OriginalValue and Value to replace values.", block.Labels(), attribute.Name)
			}
			attr := block.Body().GetAttribute(attribute.Name)
			if attr == nil {
				log.Fatalf("Could not find attribute %s in provider", attribute.Name)
			}
			value := attr.BuildTokens(nil)
			log.Println(fmt.Sprintf("Found attr %s - value %v ", attr.BuildTokens(nil).Bytes(), value))
		case attribute.Action == types.Add:
			block.Body().SetAttributeValue(attribute.Name, cty.StringVal(*attribute.Value))
		default:
			log.Fatalf("Could not find action %s", attribute.Action)
		}
	}
}

