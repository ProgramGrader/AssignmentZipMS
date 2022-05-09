
resource "aws_s3_bucket" "url_s3_b" {
  bucket = "url-s3-bucket"
}

resource "aws_dynamodb_table" "urls" {
  name = "S3AssignmentFileSource"
  billing_mode = "PAY_PER_REQUEST"
  attribute {
    name = "UUID"
    type = "S"
  }

  attribute {
    name = "key"
    type = "S"
  }

  attribute {
    name = "region"
    type = "S"
  }

  attribute {
    name = "bucket"
    type = "S"
  }

  hash_key = "UUID"
}
