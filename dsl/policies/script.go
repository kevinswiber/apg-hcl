package policies

import (
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
)

// Script represents a <Script/> element.
//
// Documentation: http://docs.apigee.com/api-services/reference/python-script-policy
type Script struct {
	XMLName     string `xml:"Script" hcl:"-"`
	Policy      `hcl:",squash"`
	DisplayName string   `xml:",omitempty" hcl:"display_name"`
	ResourceURL string   `hcl:"resource_url"`
	IncludeURL  []string `xml:",omitempty" hcl:"include_url"`
	Content     string   `xml:"-" hcl:"content"`
}

// NewScriptFromHCL converts an HCL ast.ObjectItem into a Script object.
func NewScriptFromHCL(item *ast.ObjectItem) (interface{}, error) {
	var p Script

	if err := LoadCommonPolicyHCL(item, &p.Policy); err != nil {
		return nil, err
	}

	if err := hcl.DecodeObject(&p, item.Val.(*ast.ObjectType)); err != nil {
		return nil, err
	}

	return p, nil
}
