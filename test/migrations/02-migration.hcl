migration_name = "test-migration2"
source_version = "1.2.0"
destination_version = "1.3.0"

module "k8s-api" {
  attribute {
    name = "source"
    action  = "add"
    value = "./PATH_TO_MODULE"
  }
}

## this will replace the value events in with events-v2 in all datadog monitors.
resource "datadog_monitor" {
  attribute {
    name = "type"
    action = "replace" // TODO: Replace should fail if not found.
    original_value = "events"
    value = "events-v2"
  }
}

## this will remove the the vault provider version
provider "vault" {
  attribute {
    name = "version"
    action  = "delete"
  }
}