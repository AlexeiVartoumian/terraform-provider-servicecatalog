

terraform {
    required_providers {
        servicecatalog = {
            source = "hashicorp.com/edu/servicecatalog"
        }
    }
}


provider "servicecatalog" {
    region = "eu-west-2"
}

data "servicecatalog_applications" "all" {}

output "applications" {
    value = data.servicecatalog_applications.all.applications
}

resource "servicecatalog_resource_association" "example" {
  application_arn = data.servicecatalog_applications.all.applications[0].arn
  resource_type   = "CFN_STACK"
  resource_name   = "instance-scheduler"
  options         = ["APPLY_APPLICATION_TAG"]  # Optional, defaults to ["APPLY_APPLICATION_TAG"] if not specified
}