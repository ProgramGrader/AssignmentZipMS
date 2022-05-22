<<<<<<< Updated upstream:terraform/assignment-lambda-func.tf
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
  function_name    = "assignment-url-redirect"
  filename         = data.archive_file.lambda_zip.output_path
  source_code_hash = data.archive_file.lambda_zip.output_base64sha256
  handler          = "main"
  role             = aws_iam_role.lambda-role.arn
  runtime          = "go1.x"
  timeout          = 5
  memory_size      = 128
  tracing_config {
    mode = "Active"
  }

}


resource "aws_iam_role" "lambda-role" {
  // name = "csgl-assignmentzipms-iam-role-lambda"
  assume_role_policy = jsonencode(
    {
      Version = "2012-10-17"
      Statement = [{
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Sid    = ""
        Principal = {
          Service = "lambda.amazonaws.com"
        }
        }
      ]
  })
}

resource "aws_iam_role_policy_attachment" "lambda_policy" {
  role       = aws_iam_role.lambda-role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

// allows s3:getobject access

resource "aws_iam_policy" "s3-policy" {
  policy = jsonencode(
    {
      "Version" : "2012-10-17",
      "Statement" : [
        {
          "Effect" : "Allow",
          "Action" : ["s3:GetObject"],
          "Resource" : ["arn:aws:s3:::*"]
        }
      ]
    })
}

resource "aws_iam_role_policy_attachment" "attach_s3_policy" {
  role       = aws_iam_role.lambda-role.name
  policy_arn = aws_iam_policy.s3-policy.arn
}

// Adds dynamodb getObject permission
resource "aws_iam_policy" "dynamodb-policy" {
  policy = jsonencode(
{
	"Version": "2012-10-17",
	"Statement": [{
		"Sid": "ReadWriteTable",
		"Effect": "Allow",
		"Action": ["dynamodb:GetItem"],
		"Resource": "arn:aws:dynamodb:${var.primary_aws_region}:${aws_dynamodb_table.S3AssignmentFileSource.arn}:table/${aws_dynamodb_table.S3AssignmentFileSource.name}"
	}]
})
}

resource "aws_iam_role_policy_attachment" "attach_dynamodb_policy" {
  role       = aws_iam_role.lambda-role.name
  policy_arn = aws_iam_policy.dynamodb-policy.arn
}


// Gives readOnly permission for dynamo
resource "aws_iam_policy" "readonly-policy" {
  policy = jsonencode({
    "Version": "2012-10-17",
    "Statement": [{
      "Action": [
        "dynamodb:BatchGetItem",
        "dynamodb:Describe*",
        "dynamodb:List*",
        "dynamodb:GetItem",
        "dynamodb:Query",
        "dynamodb:Scan",
        "dynamodb:PartiQLSelect"

      ],
      "Effect": "Allow",
      "Resource": "*"
    },
      {
        "Action": "cloudwatch:GetInsightRuleReport",
        "Effect": "Allow",
        "Resource": "arn:aws:cloudwatch:*:*:insight-rule/DynamoDBContributorInsights*"
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "readonly-attach-policy-lambda" {
  policy_arn = aws_iam_policy.readonly-policy.arn
  role       = aws_iam_role.lambda-role.name
}


// AWS COMMANDS TO MAKE SURE FUNCTION EXISTS/WORKS: // TODO create tests instead of manual commands
// awslocal lambda list-functions
// awslocal lambda invoke --function-name redirect out.txt

=======
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
  function_name    = "assignment-url-redirect"
  filename         = data.archive_file.lambda_zip.output_path
  source_code_hash = data.archive_file.lambda_zip.output_base64sha256
  handler          = "main"
  role             = aws_iam_role.lambda-role.arn
  runtime          = "go1.x"
  timeout          = 5
  memory_size      = 128
  tracing_config {
    mode = "Active"
  }

}


resource "aws_iam_role" "lambda-role" {
  // name = "csgl-assignmentzipms-iam-role-lambda"
  assume_role_policy = jsonencode(
    {
      Version = "2012-10-17"
      Statement = [{
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Sid    = ""
        Principal = {
          Service = "lambda.amazonaws.com"
        }
        }
      ]
  })
}

resource "aws_iam_role_policy_attachment" "lambda_policy" {
  role       = aws_iam_role.lambda-role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}


#resource "aws_iam_role_policy_attachment" "lambda_policy_SQS" {
#  role       = aws_iam_role.lambda-role.name
#  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonSQSFullAccess"
#}
#
#resource "aws_iam_role_policy_attachment" "lambda_policy_cloudwatch" {
#  role       = aws_iam_role.lambda-role.name
#  policy_arn = "arn:aws:iam::aws:policy/service-role/CloudWatchFullAccess"
#}
#
#resource "aws_iam_role_policy_attachment" "lambda_policy_sns" {
#  role       = aws_iam_role.lambda-role.name
#  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonSnsFullAccess"
#}

// Creating sns topic to push to given that the lambda runs into error


resource "aws_sns_topic" "lambda-sqs-topic" {
  name = "lambda-dlq-sqs-topic"

}

resource "aws_sqs_queue" "dlq-lambda" {
  name = "lambda-dlq"
  message_retention_seconds = 60
  visibility_timeout_seconds = 34560 // 4days
}

resource "aws_sqs_queue" "lambda-sqs-queue" {
  name = "lambda-sqs-queue"
  delay_seconds = 90
  max_message_size = 2048
  receive_wait_time_seconds = 5
  redrive_policy = jsonencode({
    "deadLetterTargetArn" = aws_sqs_queue.dlq-lambda.arn,
    "maxReceiveCount" = 3
  })

}

// allows our sqs queue to receive notifications form sns topic

resource "aws_sns_topic_subscription" "lambda-sqs-target" {
#  endpoint  = aws_sqs_queue.lambda-sqs-queue.arn
  endpoint = aws_lambda_function.redirect_lambda.arn
  protocol  = "lambda"
  topic_arn = aws_sns_topic.lambda-sqs-topic.arn
}

// policies for aws
resource "aws_sqs_queue_policy" "lambda-sqs-policy" {
  policy    = jsonencode(
    {
      "Version": "2012-10-17",
      "Id": "sqspolicy",
      "Statement": [
        {
          "Sid": "First",
          "Effect": "Allow",
          "Principal": "*",
          "Action": "sqs:SendMessage",
          "Resource": "${aws_sqs_queue.lambda-sqs-queue.arn}",
          "Condition": {
            "ArnEquals": {
              "aws:SourceArn": "${aws_sns_topic.lambda-sqs-topic.arn}"
            }
          }
        }
      ]
    }
  )
  queue_url = aws_sqs_queue.lambda-sqs-queue.id
}
// allows s3:getobject access

resource "aws_iam_policy" "s3-policy" {
  policy = jsonencode(
    {
      "Version" : "2012-10-17",
      "Statement" : [
        {
          "Effect" : "Allow",
          "Action" : ["s3:GetObject"],
          "Resource" : ["arn:aws:s3:::*"]
        }
      ]
    })
}

resource "aws_iam_role_policy_attachment" "attach_s3_policy" {
  role       = aws_iam_role.lambda-role.name
  policy_arn = aws_iam_policy.s3-policy.arn
}

// Adds dynamodb getObject permission
resource "aws_iam_policy" "dynamodb-policy" {
  policy = jsonencode(
{
	"Version": "2012-10-17",
	"Statement": [{
		"Sid": "ReadWriteTable",
		"Effect": "Allow",
		"Action": ["dynamodb:GetItem"],
		"Resource": "arn:aws:dynamodb:${var.aws_region}:${aws_dynamodb_table.S3AssignmentFileSource.arn}:table/${aws_dynamodb_table.S3AssignmentFileSource.name}"
	}]
})
}

resource "aws_iam_role_policy_attachment" "attach_dynamodb_policy" {
  role       = aws_iam_role.lambda-role.name
  policy_arn = aws_iam_policy.dynamodb-policy.arn
}


// Gives readOnly permission for dynamo
resource "aws_iam_policy" "readonly-policy" {
  policy = jsonencode({
    "Version": "2012-10-17",
    "Statement": [{
      "Action": [
        "dynamodb:BatchGetItem",
        "dynamodb:Describe*",
        "dynamodb:List*",
        "dynamodb:GetItem",
        "dynamodb:Query",
        "dynamodb:Scan",
        "dynamodb:PartiQLSelect"

      ],
      "Effect": "Allow",
      "Resource": "*"
    },
      {
        "Action": "cloudwatch:GetInsightRuleReport",
        "Effect": "Allow",
        "Resource": "arn:aws:cloudwatch:*:*:insight-rule/DynamoDBContributorInsights*"
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "readonly-attach-policy-lambda" {
  policy_arn = aws_iam_policy.readonly-policy.arn
  role       = aws_iam_role.lambda-role.name
}


// AWS COMMANDS TO MAKE SURE FUNCTION EXISTS/WORKS: // TODO create tests instead of manual commands
// awslocal lambda list-functions
// awslocal lambda invoke --function-name redirect out.txt

>>>>>>> Stashed changes:terraform/lambda-func.tf
