terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.0.0"
    }
    random = {
      source  = "hashicorp/random"
      version = "~> 3.1.0"
    }
    archive = {
      source  = "hashicorp/archive"
      version = "~> 2.2.0"
    }
  }

  required_version = "~> 1.0"
}

resource "aws_apigatewayv2_api" "url-shortener-proxy" {
  name          = "url-shortener-proxy"
  protocol_type = "HTTP"
}

resource "aws_cloudwatch_log_group" "api-gw" {
  retention_in_days = 30
}

resource "aws_apigatewayv2_stage" "dev_stage" {
  api_id = aws_apigatewayv2_api.url-shortener-proxy.id
  name   = "dev_stage"
  auto_deploy = true

  access_log_settings {
    destination_arn = aws_cloudwatch_log_group.api-gw.arn

    format = jsonencode({
      requestId               = "$context.requestId"
      sourceIp                = "$context.identity.sourceIp"
      requestTime             = "$context.requestTime"
      protocol                = "$context.protocol"
      httpMethod              = "$context.httpMethod"
      resourcePath            = "$context.resourcePath"
      routeKey                = "$context.routeKey"
      status                  = "$context.status"
      responseLength          = "$context.responseLength"
      integrationErrorMessage = "$context.integrationErrorMessage"
    }
    )
  }
}

resource "aws_apigatewayv2_integration" "redirect_lambda" {
  api_id           = aws_apigatewayv2_api.url-shortener-proxy.id
  integration_type = "AWS_PROXY"
  integration_uri = aws_lambda_function.redirect_lambda.invoke_arn
  integration_method = "POST"
}

resource "aws_apigatewayv2_route" "redirect_lambda" {
  api_id    = aws_apigatewayv2_api.url-shortener-proxy.id
  route_key = "GET /" // matching any get request matching the / path
  target = "integrations/${aws_apigatewayv2_integration.redirect_lambda.id}"
}

resource "aws_cloudwatch_log_group" "api_gw" {
  //name = "url-shortener-proxy/${aws_apigatewayv2_api.url-shortener-proxy.name}"

  retention_in_days = 30
}

resource "aws_lambda_permission" "api_gw" {
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.redirect_lambda.function_name
  principal     = "apigateway.amazonaws.com"

  source_arn = "${aws_apigatewayv2_api.url-shortener-proxy.execution_arn}/*/*"
}


