
resource "aws_s3_bucket" "url_s3_b" {
  bucket = "assignment-url-bucket"

  tags = {
    Name = "url bucket"
    Environment = "dev"
  }

#  logging {
#    target_bucket = "target-bucket"
#  }

}

#
#resource "aws_s3_bucket_public_access_block" "block-public-access" {
#  bucket = aws_s3_bucket.url_s3_b.bucket
#  block_public_acls = true
#  block_public_policy = true
#  ignore_public_acls = true
#  restrict_public_buckets = true
#}

#resource "aws_s3_bucket_cors_configuration" "allow_access" {
#  bucket = aws_s3_bucket.url_s3_b.bucket
#
#  cors_rule {
#    allowed_headers = ["*"]
#    max_age_seconds = 3000
#    allowed_methods = ["GET"]
#    allowed_origins = ["*"]
#  }
#
#}

#


// ^ Did not solve the pre signed calculation error //


#resource "aws_kms_key" "dynamo_db_kms" {
#  enable_key_rotation = true
#}

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

  attribute {
    name = "bucket"
    type = "S"
  }

  attribute {
    name = "region"
    type = "S"
  }

  attribute {
    name = "filename"
    type = "S"
  }

  hash_key = "UUID"

  global_secondary_index {
    hash_key           = "region"
    name               = "region"
    projection_type    = "INCLUDE"
    non_key_attributes = [ "bucket", "filename"]
  }

  global_secondary_index {
    hash_key           = "bucket"
    name               = "bucket"
    projection_type    = "INCLUDE"
    non_key_attributes = ["region", "filename"]

  }

  global_secondary_index {
    hash_key           = "filename"
    name               = "filename"
    projection_type    = "INCLUDE"
    non_key_attributes = ["region", "bucket"]
  }

#  server_side_encryption {
#    enabled = true
#    kms_key_arn = aws_kms_key.dynamo_db_kms.key_id
#  }
}

