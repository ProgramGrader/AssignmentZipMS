// Tools used to test this infrastructure locally: Localstacks, tflocal, and awslocal
// build localStacks: docker-compose up
// pip install terraform-local
// if the tflocal or awslocal commands aren't recognized try restarting your terminal

// TODO - Fix terraform vulnerabilities
// TODO - Test terraform using terragrunt

locals {
  region = "us-east-2"
}

provider "aws" {
  profile = "default"
  region = "us-east-2"
}
