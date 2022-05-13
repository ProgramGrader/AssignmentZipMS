// Tools used to test this infrastructure locally: Localstacks, tflocal, and awslocal
// build localStacks: docker-compose up
// pip install terraform-local
// if the tflocal or awslocal commands aren't recognized try restarting your terminal

terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.0"
    }
  }

  required_version = ">= 1.1.0"
}

locals {
  ProjectName = "AssignmentZipMS"
}

provider "aws" {

  access_key = "test"
  secret_key = "test"

  region = "us-east-2" // looks like for some reason for awslocal to work you need to change the region to us-east-1
  #profile = "dev"

  # only required for non virtual hosted-style endpoint use case.
  # https://registry.terraform.io/providers/hashicorp/aws/latest/docs#s3_force_path_style
  skip_credentials_validation = true
  skip_metadata_api_check     = true
  skip_requesting_account_id  = true

  default_tags {
    tags = {
      Terraform   = "true"
      Project     = local.ProjectName
    }
  }

}
// Current pain point : receiving bucket not found error when trying to connect to the restapi using links similar to this http://localhost:4566/restapis/sl539f7o55/dev