#// grab this url to interact with the api_client
output "api_url"{ // url of the api gateway
  value = aws_apigatewayv2_stage.api-gw_stage.invoke_url
}

output "permanent_url" {
  value = aws_apigatewayv2_domain_name.csgrader.domain_name
}



