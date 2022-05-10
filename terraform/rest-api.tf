
resource "aws_api_gateway_rest_api" "url_shortener_proxy" {
  name = "url_shortener_proxy"
  description = "proxy used to handle the requests to lambda function"
  endpoint_configuration {
    types = ["REGIONAL"]
  }
}

resource "aws_api_gateway_resource" "url" { // since we are accessing the hash from the url
  rest_api_id = aws_api_gateway_rest_api.url_shortener_proxy.id
  parent_id = aws_api_gateway_rest_api.url_shortener_proxy.root_resource_id
  path_part = "{proxy+}" // this string represents the endpoint path, for this resource

}

resource "aws_api_gateway_method" "get" {
  authorization = "NONE"
  api_key_required = false
  http_method   = "POST"
  resource_id   = aws_api_gateway_resource.url.id
  rest_api_id   = aws_api_gateway_rest_api.url_shortener_proxy.id
}

// this resource describes how we are going to react to the request received
// in our case the api receives a get request, we extract the hash from the url in the backend which will be passed
// into the event handler for our lambda


// we're not actually getting here anything here this, sends the request received by the proxy to the lambda function
resource "aws_api_gateway_integration" "integration-get" {
  resource_id         = aws_api_gateway_resource.url.id
  rest_api_id         = aws_api_gateway_rest_api.url_shortener_proxy.id
  integration_http_method  = aws_api_gateway_method.get.http_method // represents the HTTP method that will be done from the integration to the backend
  http_method         = "POST"
  type                = "AWS_PROXY" // https://docs.aws.amazon.com/apigateway/latest/developerguide/api-gateway-api-integration-types.html
  uri                 = aws_lambda_function.redirect_lambda.invoke_arn // contains the endpoint to which we are proxying too. In our case its a lambda function
}


#resource "aws_api_gateway_resource" "redirect" {
#  rest_api_id = aws_api_gateway_rest_api.url_shortener_proxy.id
#  parent_id = aws_api_gateway_rest_api.url_shortener_proxy.root_resource_id
#  path_part = "{url+}"
#}
#
#resource "aws_api_gateway_method" "redirect-handler" {
#  authorization = "NONE"
#  http_method   = "POST"
#  resource_id   = aws_api_gateway_resource.url.id
#  rest_api_id   = aws_api_gateway_rest_api.url_shortener_proxy.id
#}


// proxy cannot match a empty path
resource "aws_api_gateway_method" "url_root" {

  authorization = "NONE"
  http_method   = "POST"
  resource_id   = aws_api_gateway_rest_api.url_shortener_proxy.root_resource_id
  rest_api_id   = aws_api_gateway_rest_api.url_shortener_proxy.id
}

resource "aws_api_gateway_integration" "lambda_root" {
  http_method = "POST"
  resource_id = aws_api_gateway_method.url_root.resource_id
  rest_api_id = aws_api_gateway_rest_api.url_shortener_proxy.id

  type        = "AWS_PROXY"
  uri = aws_lambda_function.redirect_lambda.invoke_arn
}


resource "aws_api_gateway_deployment" "deploy-1" {
  rest_api_id = aws_api_gateway_rest_api.url_shortener_proxy.id

  depends_on = [aws_api_gateway_integration.integration-get]

  lifecycle {
    create_before_destroy = true // creates a new deployment before trashing old one
  }
  description = "Deployed endpoint at ${timestamp()}"
}

resource "aws_api_gateway_stage" "dev"{
  stage_name = "dev"
  deployment_id = aws_api_gateway_deployment.deploy-1.id
  rest_api_id = aws_api_gateway_rest_api.url_shortener_proxy.id
  xray_tracing_enabled = true
}