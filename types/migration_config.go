package types

type Migration struct {
	MigrationName string `hcl:"migration_name"`
	SourceVersion string `hcl:"source_version"`
	DestVersion string `hcl:"destination_version"`
	ModuleBlocks []*Module `hcl:"module,block"`
	ResourceBlocks []*Resource `hcl:"resource,block"`
	ProviderBlocks []*Provider `hcl:"provider,block"`
}
