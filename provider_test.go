package main

import (
	hcl2 "github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/stobias123/terraform_editor/types"
	"github.com/stobias123/terraform_editor/util"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

// TODO: Add some tests to load migration files in order
// TODO: validate migration file order based on source/dest version

func setup() *hclwrite.File {
	b, err := os.ReadFile("test.tf")
	if err != nil {
		log.Fatalf("Couldn't read file %v", err)
	}
	file, diags := hclwrite.ParseConfig(b, "test.tf", hcl2.Pos{Line: 1, Column: 1})
	if diags.HasErrors() {
		log.Fatal(diags.Error())
	}
	return file
}

func TestProviderAttributeDelete(t *testing.T) {
	t.Parallel()
	file := setup()
	rootBody := file.Body()
	editConfig := types.Migration{
		SourceVersion: "1.1.0",
		DestVersion:   "1.2.0",
		ProviderBlocks: []*types.Provider{
			{
				Name: "vault",
				Attributes: []types.AttributeConfig{
					{
						Name:          "version",
						Action:        "delete",
					},
				},
			},
		},
	}
	EditProviders(rootBody, editConfig)
	blocks := util.FindBlocks(rootBody,"provider", []string{"vault"})
	assert.Nil(t,blocks[0].Body().GetAttribute("foo"))
}

// Proceed when we want to delete a non existant file.
func TestProviderAttributeDeleteNonExistant(t *testing.T) {
	t.Parallel()
	file := setup()
	rootBody := file.Body()
	editConfig := types.Migration{
		SourceVersion: "1.1.0",
		DestVersion:   "1.2.0",
		ProviderBlocks: []*types.Provider{
			{
				Name: "foo2",
				Attributes: []types.AttributeConfig{
					{
						Name:          "version",
						Action:        "delete",
					},
				},
			},
		},
	}
	EditProviders(rootBody, editConfig)
	blocks := util.FindBlocks(rootBody,"provider", []string{"vault"})
	assert.Nil(t,blocks[0].Body().GetAttribute("foo"))
}

// Update the resources value.
func TestProviderAttributeUpdate(t *testing.T) {
	t.Parallel()
	file := setup()
	rootBody := file.Body()
	versionString := "2.3.4"
	editConfig := types.Migration{
		SourceVersion: "1.1.0",
		DestVersion:   "1.2.0",
		ProviderBlocks: []*types.Provider{
			{
				Name: "vault",
				Attributes: []types.AttributeConfig{
					{
						Name:          "version",
						Action:        "update",
						Value: &versionString,
					},
				},
			},
		},
	}
	EditProviders(rootBody,editConfig)
	blocks := util.FindBlocks(rootBody,"provider", []string{"vault"})
	assert.Equal(t,"  version =\"2.3.4\"\n", string(blocks[0].Body().GetAttribute("version").BuildTokens(nil).Bytes()))
	assert.Contains(t,string(file.Bytes()),"2.3.4")
}

//Add when the attribute doens't exist
func TestProviderAttributeAddOrUpdate(t *testing.T) {
	t.Parallel()
	file := setup()
	rootBody := file.Body()
	versionString := "2.3.4"
	editConfig := types.Migration{
		SourceVersion: "1.1.0",
		DestVersion:   "1.2.0",
		ProviderBlocks: []*types.Provider{
			{
				Name: "vault",
				Attributes: []types.AttributeConfig{
					{
						Name:          "fooidontexist",
						Action:        "update",
						Value: &versionString,
					},
				},
			},
		},
	}
	EditProviders(rootBody,editConfig)
	assert.Contains(t,string(file.Bytes()),"fooidontexist")
}

// Assert nothing happens when the provider label doesn't exist.
func TestProviderProviderNonExistent(t *testing.T) {
	t.Parallel()
	file := setup()
	rootBody := file.Body()
	versionString := "2.3.4"
	editConfig := types.Migration{
		SourceVersion: "1.1.0",
		DestVersion:   "1.2.0",
		ProviderBlocks: []*types.Provider{
			{
				Name: "fooidontexist",
				Attributes: []types.AttributeConfig{
					{
						Name:          "version",
						Action:        "update",
						Value: &versionString,
					},
				},
			},
		},
	}
	EditProviders(rootBody,editConfig)
	//fmt.Printf("%s", file.Bytes())
	assert.NotContains(t, "fooidontexist",string(file.Bytes()))
}

