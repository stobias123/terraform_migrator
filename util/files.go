package util

import (
	"errors"
	hcl2 "github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/stobias123/terraform_editor/types"
	"os"
)

func LoadMigrationFile(migrationFilePath string) (*types.Migration, error) {
	var migrationConfig types.Migration
	parser := hclparse.NewParser()
	f, parseDiags := parser.ParseHCLFile(migrationFilePath)
	if parseDiags.HasErrors() {
		//log(parseDiags.Error())
		return nil, errors.New(parseDiags.Error())
	}
	decodeDiags := gohcl.DecodeBody(f.Body, nil, &migrationConfig)
	if decodeDiags.HasErrors() {
		//log.Fatal(decodeDiags.Error())
		return nil, errors.New(parseDiags.Error())
	}
	return &migrationConfig, nil
}

func LoadHCLFile(hclFilePath string) (*hclwrite.File, error) {
	//loadMigrationFile
	b, err := os.ReadFile(hclFilePath)
	if err != nil {
		return nil, err
	}
	file, diags := hclwrite.ParseConfig(b, hclFilePath, hcl2.Pos{Line: 1, Column: 1})
	if diags.HasErrors() {
		return nil, errors.New(diags.Error())
	}
	return file, nil
}