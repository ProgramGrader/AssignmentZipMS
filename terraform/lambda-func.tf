

// compile code into binary
resource "null_resource" "compile_binary" {
  triggers = {
    build_number = timestamp()
  }

  provisioner "local-exec" {
    command = "go build -ldflags '-w' -o handler handler.go"
    interpreter = ["PowerShell"]
    working_dir = "../src"
  }
}

// zipping code
data "archive_file" "lambda_zip"{
  source_file = "../src/handler"
  type        = "zip"
  output_path = "handler.zip"
  depends_on = [null_resource.compile_binary]
}


resource "aws_lambda_function" "redirect_lambda" {
  function_name = "redirect"
  filename      = data.archive_file.lambda_zip.output_path
  source_code_hash = data.archive_file.lambda_zip.output_base64sha256
  handler       = "handler"
  role          = "test"
  runtime       = "go1.x"
  timeout       = 5
  memory_size   = 128
  depends_on = [data.archive_file.lambda_zip]
}

resource "aws_lambda_permission" "allow_api" {
  statement_id  = "AllowApigatewayInvocation"
  function_name = aws_lambda_function.redirect_lambda.function_name
  principal     = "apigateway.amazonaws.com"
  action        = "lambda:InvokeFunction"
}

// AWS COMMANDS TO MAKE SURE FUNCTION EXISTS/WORKS:
// awslocal lambda list-functions
// awslocal lambda invoke --function-name redirect out.txt