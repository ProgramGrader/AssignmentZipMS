// Tools used to test this infrastructure locally: Localstacks, tflocal, and awslocal
// build localStacks: docker-compose up
// pip install terraform-local
// if the tflocal or awslocal commands aren't recognized try restarting your terminal

provider "aws" {
  region = "us-east-2" // looks like for some reason for awslocal to work you need to change the region to us-east-1
  profile = "dev"
}
// Current pain point : receiving bucket not found error when trying to connect to the restapi using links similar to this http://localhost:4566/restapis/sl539f7o55/dev