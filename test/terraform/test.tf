
resource "datadog_monitor" "replace" {
  message = "foo"
  type = "events"
}

resource "datadog_monitor" "dontreplace" {
  message = "foo"
  type = "metric"
}

module "k8s-api" {
  source = "./foo-bar"
}

provider "vault" {
  version = "= 100000"
}