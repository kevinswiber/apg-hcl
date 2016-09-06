package javascript

import (
	"fmt"
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/kevinswiber/apigee-hcl/dsl/common"
	"github.com/kevinswiber/apigee-hcl/dsl/policies"
)

// JavaScript represents a <Javascript/> element.
//
// Documentation: http://docs.apigee.com/api-services/reference/javascript-policy
type JavaScript struct {
	XMLName         string `xml:"Javascript" hcl:"-"`
	policies.Policy `hcl:",squash"`
	TimeLimit       int                `xml:"timeLimit,attr" hcl:"time_limit"`
	DisplayName     string             `xml:",omitempty" hcl:"display_name"`
	ResourceURL     string             `hcl:"resource_url"`
	IncludeURL      []string           `xml:",omitempty" hcl:"include_url"`
	Properties      []*common.Property `xml:"Properties>Property" hcl:"properties"`
	Content         string             `xml:"-" hcl:"content"`
}

// GetName returns the policy name.
func (policy JavaScript) GetName() string {
	return policy.Name
}

// GetResourceURL returns the resource URL for the policy
func (policy JavaScript) GetResourceURL() string {
	return policy.ResourceURL
}

// GetResourceContent returns the reousrce content
func (policy JavaScript) GetResourceContent() string {
	return policy.Content
}

// DecodeJavaScriptHCL converts HCL into an JavaScript object.
func DecodeJavaScriptHCL(item *ast.ObjectItem) (interface{}, error) {
	var p JavaScript

	if err := policies.DecodePolicyHCL(item, &p.Policy); err != nil {
		return nil, err
	}

	if err := hcl.DecodeObject(&p, item.Val.(*ast.ObjectType)); err != nil {
		return nil, err
	}

	var listVal *ast.ObjectList
	if ot, ok := item.Val.(*ast.ObjectType); ok {
		listVal = ot.List
	} else {
		return nil, fmt.Errorf("javascript policy not an object")
	}

	if propsList := listVal.Filter("properties"); len(propsList.Items) > 0 {
		props, err := common.DecodePropertiesHCL(propsList.Items[0])
		if err != nil {
			return nil, err
		}

		p.Properties = props
	}

	return p, nil
}
