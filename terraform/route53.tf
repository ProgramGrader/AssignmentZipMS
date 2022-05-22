
data "aws_route53_zone" "csgrader" {
  zone_id      = "Z01322341QHBI9TU4903A"
  private_zone = false
}

resource "aws_acm_certificate" "dns_csgrader_acm" {
  domain_name       = data.aws_route53_zone.csgrader.name
  validation_method = "DNS"
}

resource "aws_route53_record" "dns_csgrader_record" {
  for_each = {
  for dvo in aws_acm_certificate.dns_csgrader_acm.domain_validation_options : dvo.domain_name => {
    name   = dvo.resource_record_name
    record = dvo.resource_record_value
    type   = dvo.resource_record_type
  }
  }

  allow_overwrite = true
  name            = each.value.name
  records         = [each.value.record]
  ttl             = 60
  type            = each.value.type
  zone_id         = data.aws_route53_zone.csgrader.zone_id
}

resource "aws_acm_certificate_validation" "dns_csgrader_validation" {
  certificate_arn         = aws_acm_certificate.dns_csgrader_acm.arn
  validation_record_fqdns = [for record in aws_route53_record.dns_csgrader_record : record.fqdn]
}

resource "aws_apigatewayv2_domain_name" "csgrader" {
  depends_on  = [aws_apigatewayv2_stage.api-gw_stage]
  domain_name = "assignmentFile.${data.aws_route53_zone.csgrader.name}"
  domain_name_configuration {
    certificate_arn = aws_acm_certificate.dns_csgrader_acm.arn
    endpoint_type   = "REGIONAL"
    security_policy = "TLS_1_2"
  }
}

resource "aws_route53_record" "microserviceUrl" {
  zone_id = data.aws_route53_zone.csgrader.zone_id
  name    = aws_apigatewayv2_domain_name.csgrader.domain_name
  type    = "A"
  alias {
    evaluate_target_health = false
    name                   = aws_apigatewayv2_domain_name.csgrader.domain_name_configuration.0.target_domain_name
    zone_id                = aws_apigatewayv2_domain_name.csgrader.domain_name_configuration.0.hosted_zone_id
  }
}

resource "aws_apigatewayv2_api_mapping" "mapping-invoke-url" {
  api_id = aws_apigatewayv2_api.url-shortener-proxy.id
  domain_name = aws_apigatewayv2_domain_name.csgrader.id
  stage = aws_apigatewayv2_stage.api-gw_stage.id
}
##// clean up source
##// package, then lambda folder with functions.