proxy "ExtractVariablesFixture" {}

proxy_endpoint "default" {
  http_proxy_connection {
    base_path    = "/v0/variables"
    virtual_host = ["default", "secure"]
  }

  route_rule "default" {
    target_endpoint = "default"
  }
}

target_endpoint "default" {
  http_target_connection {
    url = "http://mocktarget.apigee.net"
  }

  pre_flow {
    response {
      step "extract-vars" {}
    }
  }
}

policy extract_variables "extract-vars" {
  source {
    value = "response"
  }

  json_payload {
    variable {
      name      = "method"
      type      = "string"
      json_path = "$.method"
    }
  }

  variable_prefix             = "resinfo"
  ignore_unresolved_variables = true
}
