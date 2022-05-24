package types

type Provider struct {
	Name string `hcl:"name,label"`
	Attributes []AttributeConfig `hcl:"attribute,block"`
}
