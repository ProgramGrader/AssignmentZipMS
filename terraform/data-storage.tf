
resource "aws_s3_bucket" "url_s3_b" {
  bucket = "assignment-url-bucket"

  force_destroy = true

  tags = {
    Name = "url bucket"
    Environment = "dev"
  }
#  logging {
#    target_bucket = "target-bucket"
#  }
}


resource "aws_dynamodb_table" "S3AssignmentFileSource" {
  name         = "S3AssignmentFileSource"
  billing_mode = "PAY_PER_REQUEST"

#  point_in_time_recovery {
#    enabled = true
#  }
  attribute {
    name = "UUID"
    type = "S"
  }

  hash_key = "UUID"

}

