terraform {
  required_version = ">= 1.14"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "= 6.29.0"
    }
  }
}

provider "aws" {
  region = var.aws_region
}
