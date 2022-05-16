// compile code into binary
resource "null_resource" "compile_binary" {
  triggers = {
    build_number = timestamp()

  }
  provisioner "local-exec" {
    command = "GOOS=linux GOARCH=amd64 go build -ldflags '-w' -o  ../src/sam_lambda_test/handler/main  ../src/sam_lambda_test/handler/main.go"
  }
}

// zipping code
data "archive_file" "lambda_zip" {
  source_file = "../src/sam_lambda_test/handler/main"
  type        = "zip"
  output_path = "handler.zip"
  depends_on  = [null_resource.compile_binary]
}


resource "aws_lambda_function" "redirect_lambda" {
  function_name = "assignment-url-redirect"
  filename      = data.archive_file.lambda_zip.output_path
  #  source_code_hash = data.archive_file.lambda_zip.output_base64sha256
  handler       = "main"
  role          = aws_iam_role.lambda-role.arn
  runtime       = "go1.x"
  timeout       = 5
  memory_size   = 128

  tracing_config {
    mode = "Active"
  }

}

resource "aws_lambda_permission" "allow_api" {
  statement_id  = "AllowApigatewayInvocation"
  function_name = aws_lambda_function.redirect_lambda.function_name
  principal     = "apigateway.amazonaws.com"
  action        = "lambda:InvokeFunction"
  source_arn    = "${aws_api_gateway_rest_api.url_shortener_proxy.execution_arn}/*/*"
}


resource "aws_iam_role" "lambda-role" {
  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "lambda_policy" {
  role       = aws_iam_role.lambda-role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

// allows s3:getobject access

resource "aws_iam_policy" "s3-policy" {
  policy = <<EOF
{"Version": "2012-10-17",
 "Statement": [
  {
  "Effect": "Allow",
  "Action": ["s3:GetObject"],
  "Resource": ["arn:aws:s3:::*"]
  }
]
}
EOF
}

resource "aws_iam_role_policy_attachment" "attach_s3_policy" {
  role       =  aws_iam_role.lambda-role.name
  policy_arn = aws_iam_policy.s3-policy.arn
}

// Allows dynamodb getObject
resource "aws_iam_policy" "dynamodb-policy" {
  policy = <<EOF
{"Version": "2012-10-17",
 "Statement": [
  {
   "Sid" : "ReadWriteTable",
  "Effect": "Allow",
  "Action": [
        "dynamodb:GetItem",
        "dynamodb:BatchGetItem"
],
  "Resource": "arn:aws:dynamodb:${local.region}:${aws_dynamodb_table.S3AssignmentFileSource.arn}:table/${aws_dynamodb_table.S3AssignmentFileSource.name}"
  }
]
}
EOF
}

resource "aws_iam_role_policy_attachment" "attach_dynamodb_policy" {
  role       =  aws_iam_role.lambda-role.name
  policy_arn = aws_iam_policy.dynamodb-policy.arn
}


// AWS COMMANDS TO MAKE SURE FUNCTION EXISTS/WORKS: // TODO create tests instead of manual commands
// awslocal lambda list-functions
// awslocal lambda invoke --function-name redirect out.txt