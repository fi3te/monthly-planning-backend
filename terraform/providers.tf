terraform {
  required_version = ">= 1.4"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~>6.0"
    }
    archive = {
      source  = "hashicorp/archive"
      version = "~>2.0"
    }
  }
}

provider "aws" {
  default_tags {
    tags = {
      Tool = "Terraform"
    }
  }
}
