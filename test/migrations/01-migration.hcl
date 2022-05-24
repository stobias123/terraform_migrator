migration_name = "test-migration1"
source_version = "1.1.0"
destination_version = "1.2.0"

module "k8s-api" {
  attribute {
    name = "foo"
    action  = "delete"
  }

  attribute {
    name = "source"
    action  = "modify"
    value = "./PATH_TO_MODULE"
  }

  attribute {
    name = "bar"
    action  = "modify"
    value = "10000"
  }
}

## this will replace the value events in with events-v2 in all datadog monitors.
resource "datadog_monitor" {
  attribute {
    name = "type"
    action = "replace"
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