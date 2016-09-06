package script

import (
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/kevinswiber/apigee-hcl/dsl/policies"
)

// Script represents a <Script/> element.
//
// Documentation: http://docs.apigee.com/api-services/reference/python-script-policy
type Script struct {
	XMLName         string `xml:"Script" hcl:"-"`
	policies.Policy `hcl:",squash"`
	DisplayName     string   `xml:",omitempty" hcl:"display_name"`
	ResourceURL     string   `hcl:"resource_url"`
	IncludeURL      []string `xml:",omitempty" hcl:"include_url"`
	Content         string   `xml:"-" hcl:"content"`
}

// Resource represents an included file in a proxy bundle
func (s *Script) Resource() *policies.Resource {
	return &policies.Resource{
		URL:     s.ResourceURL,
		Content: s.Content,
	}
}

// DecodeScriptHCL converts an HCL ast.ObjectItem into a Script object.
func DecodeScriptHCL(item *ast.ObjectItem) (interface{}, error) {
	var p Script

	if err := policies.DecodePolicyHCL(item, &p.Policy); err != nil {
		return nil, err
	}

	if err := hcl.DecodeObject(&p, item.Val.(*ast.ObjectType)); err != nil {
		return nil, err
	}

	return &p, nil
}
