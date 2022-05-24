# Terraform Source Migrator

When managing large amounts of terraform code, you might need to generate "migrations" to update terraform code in a stable
and defined way.

Terraform migrator helps with that.

## Usage

Generate a list of migration files in `.hcl` or `.go`

1. Specify top level blocks just as you would in terraform.
```hcl
module "aws_instance" "web" {
  ...
}
```
2. Add attributes that you'd like to modify with `attribute` blocks. Supported [AttributeActions](./types/attribute_config.go) are 
      * "updated" - update or add value
      * "replaced" - update attribute value, only if current value is `x`
      * "deleted" - removed from a block
   
```hcl
module "aws_instance" "web" {
  attribute {
    name = "source"
    action = "update"
    value = "./path-to-new-source"
  }
}
```

All HCL files will be run in order, and modifications applied in order, when given a dir of `.hcl` files.
The tool will apply migrations to all terraform files in the given directory (not recursive)

**Example:**
```shell
terraform_migrator -migrationDir test/migrations -terraformDir test/terraform
```

I don't need to do complex modifications yet, but might add them in the future if needed.