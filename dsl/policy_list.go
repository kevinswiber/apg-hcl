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
	"assign_message":       assignmessage.DecodeAssignMessageHCL,
	"extract_variables":    extractvariables.DecodeExtractVariablesHCL,
	"javascript":           javascript.DecodeJavaScriptHCL,
	"quota":                quota.DecodeQuotaHCL,
	"raise_fault":          raisefault.DecodeRaiseFaultHCL,
	"response_cache":       responsecache.DecodeResponseCacheHCL,
	"script":               script.DecodeScriptHCL,
	"service_callout":      servicecallout.DecodeServiceCalloutHCL,
	"spike_arrest":         spikearrest.DecodeSpikeArrestHCL,
	"statistics_collector": statisticscollector.DecodeStatisticsCollectorHCL,
	"verify_api_key":       verifyapikey.DecodeVerifyAPIKeyHCL,
	"xml_to_json":          xmltojson.DecodeXMLToJSONHCL,
}
