

resource "aws_apigatewayv2_api" "serverless_lambda_gw" {
  name          = "serverless_lambda_gw"
  protocol_type = "HTTP"
}

resource "aws_cloudwatch_log_group" "api-gw" {
  retention_in_days = 30
}

resource "aws_apigatewayv2_stage" "api-gw_stage" {
  api_id = aws_apigatewayv2_api.serverless_lambda_gw.id
  name   = var.environment
  auto_deploy = true
  //PascalCase

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
  api_id           = aws_apigatewayv2_api.serverless_lambda_gw.id
  integration_type = "AWS_PROXY"
  integration_uri = aws_lambda_function.document_lambda.invoke_arn
  integration_method = "POST"
}

resource "aws_apigatewayv2_route" "redirect_lambda" {
  api_id    = aws_apigatewayv2_api.serverless_lambda_gw.id
  route_key = "GET /" // matching any get request matching the / path
  target = "integrations/${aws_apigatewayv2_integration.redirect_lambda.id}"
}

resource "aws_lambda_permission" "api_gw" {
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.document_lambda
  principal     = "apigateway.amazonaws.com"

  source_arn = "${aws_apigatewayv2_api.serverless_lambda_gw.execution_arn}/*/*"
}
# -----
resource "aws_apigatewayv2_integration" "document_lambda" {
  api_id           = aws_apigatewayv2_api.serverless_lambda_gw.id
  integration_type = "AWS_PROXY"
  integration_uri = aws_lambda_function.document_lambda.invoke_arn
  integration_method = "POST"
}

resource "aws_apigatewayv2_route" "redirect_lambda" {
  api_id    = aws_apigatewayv2_api.serverless_lambda_gw.id
  route_key = "GET /install_doc/" // matching any get request matching the / path
  target = "integrations/${aws_apigatewayv2_integration.document_lambda.id}"
}

resource "aws_lambda_permission" "api_gw" {
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.document_lambda.function_name
  principal     = "apigateway.amazonaws.com"

  source_arn = "${aws_apigatewayv2_api.serverless_lambda_gw.execution_arn}/*/*"
}


resource "aws_cloudwatch_log_group" "api_gw" {
  //name = "url-shortener-proxy/${aws_apigatewayv2_api.url-shortener-proxy.name}"

  retention_in_days = 30
}

