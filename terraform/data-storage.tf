
resource "aws_s3_bucket" "url_s3_b" {
  bucket = "url-bucket"

  tags = {
    Name = "url bucket"
    Environment = "dev"
  }
}

// Not sure if these are needed -
data "aws_iam_policy_document" "role-policy" {
  policy_id = "url-bucket-policy"

  statement {
    effect = "Allow"
    actions = ["s3:GetObject"]

    resources = ["arn::aws:s3:::url-bucket"]
  }
}

resource "aws_iam_role" "role" {
  assume_role_policy = data.aws_iam_policy_document.role-policy.json
}

// ^ Did not solve the pre signed calculation error //

resource "aws_dynamodb_table" "urls" {
  name         = "S3AssignmentFileSource"
  billing_mode = "PAY_PER_REQUEST"

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
}
