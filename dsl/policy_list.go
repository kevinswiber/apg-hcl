package dsl

import (
	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/kevinswiber/apigee-hcl/dsl/policies"
)

// PolicyList is a map of HCL policy types to policy factory functions.
var PolicyList = map[string]func(*ast.ObjectItem) (interface{}, error){
	"assign_message":       policies.NewAssignMessageFromHCL,
	"extract_variables":    policies.NewExtractVariablesFromHCL,
	"javascript":           policies.NewJavaScriptFromHCL,
	"quota":                policies.NewQuotaFromHCL,
	"raise_fault":          policies.NewRaiseFaultFromHCL,
	"response_cache":       policies.NewResponseCacheFromHCL,
	"script":               policies.NewScriptFromHCL,
	"service_callout":      policies.NewServiceCalloutFromHCL,
	"spike_arrest":         policies.NewSpikeArrestFromHCL,
	"statistics_collector": policies.NewStatisticsCollectorFromHCL,
	"verify_api_key":       policies.NewVerifyAPIKeyFromHCL,
	"xml_to_json":          policies.NewXMLToJSONFromHCL,
}
