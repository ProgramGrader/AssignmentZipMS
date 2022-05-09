// Tools used to test this infrastructure locally: Localstacks, tflocal, and awslocal
// build localStacks: docker-compose up
// pip install terraform-local
// if the tflocal or awslocal commands aren't recognized try restarting your terminal

provider "aws" {
  region = "us-east-2" // looks like for some reason for awslocal to work you need to change the region to us-east-1
  profile = "dev"
}
