package verifyapikey

import (
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/kevinswiber/apigee-hcl/dsl/policies"
)

// VerifyAPIKey represents a <VerifyAPIKey/> element.
//
// Documentation: http://docs.apigee.com/api-services/reference/verify-api-key-policy
type VerifyAPIKey struct {
	XMLName         string `xml:"VerifyAPIKey" hcl:"-"`
	policies.Policy `hcl:",squash"`
	DisplayName     string  `xml:",omitempty" hcl:"display_name"`
	APIKey          *apikey `hcl:"apikey"`
}

// GetName returns the policy name.
func (policy VerifyAPIKey) GetName() string {
	return policy.Name
}

type apikey struct {
	XMLName string `xml:"APIKey" hcl:"-"`
	Ref     string `xml:"ref,attr" hcl:"ref"`
	Value   string `xml:",chardata" hcl:"value"`
}

// NewVerifyAPIKeyFromHCL converts an HCL ast.ObjectItem into a VerifyAPIKey object.
func NewVerifyAPIKeyFromHCL(item *ast.ObjectItem) (interface{}, error) {
	var p VerifyAPIKey

	if err := policies.DecodePolicyHCL(item, &p.Policy); err != nil {
		return nil, err
	}

	if err := hcl.DecodeObject(&p, item.Val.(*ast.ObjectType)); err != nil {
		return nil, err
	}

	return p, nil
}
