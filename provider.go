package main

import (
	"fmt"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/stobias123/terraform_migrator/types"
	"github.com/stobias123/terraform_migrator/util"
	"github.com/zclconf/go-cty/cty"
	"log"
)

func EditProviders(body *hclwrite.Body, editConfig types.Migration) {
	for _, provider := range editConfig.ProviderBlocks {
		blocks := util.FindBlocks(body, "provider" , []string{provider.Name})
		for _, block := range blocks {
			editAttributes(block, provider.Attributes)
			//todo - edit sub blocks etc.
		}
	}
}

func EditModules(body *hclwrite.Body, editConfig types.Migration) {
	for _, module := range editConfig.ModuleBlocks {
		blocks := util.FindBlocks(body, "module" , []string{module.Name})
		for _, block := range blocks {
			editAttributes(block, module.Attributes)
			//todo - edit sub blocks etc.
		}
	}
}

func EditResource(body *hclwrite.Body, editConfig types.Migration) {
	for _, resource := range editConfig.ResourceBlocks {
		blocks := util.FindBlocks(body, "resource" , []string{resource.Name})
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
