package dsl

import (
	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/kevinswiber/apigee-hcl/dsl/policies/assignmessage"
	"github.com/kevinswiber/apigee-hcl/dsl/policies/extractvariables"
	"github.com/kevinswiber/apigee-hcl/dsl/policies/javascript"
	"github.com/kevinswiber/apigee-hcl/dsl/policies/quota"
	"github.com/kevinswiber/apigee-hcl/dsl/policies/raisefault"
	"github.com/kevinswiber/apigee-hcl/dsl/policies/responsecache"
	"github.com/kevinswiber/apigee-hcl/dsl/policies/script"
	"github.com/kevinswiber/apigee-hcl/dsl/policies/servicecallout"
	"github.com/kevinswiber/apigee-hcl/dsl/policies/spikearrest"
	"github.com/kevinswiber/apigee-hcl/dsl/policies/statisticscollector"
	"github.com/kevinswiber/apigee-hcl/dsl/policies/verifyapikey"
	"github.com/kevinswiber/apigee-hcl/dsl/policies/xmltojson"
)

// PolicyList is a map of HCL policy types to policy factory functions.
var PolicyList = map[string]func(*ast.ObjectItem) (interface{}, error){
	"assign_message":       assignmessage.NewAssignMessageFromHCL,
	"extract_variables":    extractvariables.NewExtractVariablesFromHCL,
	"javascript":           javascript.NewJavaScriptFromHCL,
	"quota":                quota.NewQuotaFromHCL,
	"raise_fault":          raisefault.NewRaiseFaultFromHCL,
	"response_cache":       responsecache.NewResponseCacheFromHCL,
	"script":               script.NewScriptFromHCL,
	"service_callout":      servicecallout.NewServiceCalloutFromHCL,
	"spike_arrest":         spikearrest.NewSpikeArrestFromHCL,
	"statistics_collector": statisticscollector.NewStatisticsCollectorFromHCL,
	"verify_api_key":       verifyapikey.NewVerifyAPIKeyFromHCL,
	"xml_to_json":          xmltojson.NewXMLToJSONFromHCL,
}
