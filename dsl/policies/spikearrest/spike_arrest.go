package spikearrest

import (
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/kevinswiber/apigee-hcl/dsl/policies"
)

// SpikeArrest represents a <SpikeArrest/> element.
//
// Documentation: http://docs.apigee.com/api-services/reference/spike-arrest-policy
type SpikeArrest struct {
	XMLName         string `xml:"SpikeArrest" hcl:"-"`
	policies.Policy `hcl:",squash"`
	DisplayName     string              `xml:",omitempty" hcl:"display_name"`
	Identifier      *spikeIdentifier    `hcl:"identifier"`
	MessageWeight   *spikeMessageWeight `hcl:"message_weight"`
	Rate            *spikeRate          `hcl:"rate"`
}

type spikeIdentifier struct {
	XMLName string `xml:"Identifier" hcl:"-"`
	Ref     string `xml:"ref,attr,omitempty" hcl:"ref"`
}

type spikeMessageWeight struct {
	XMLName string `xml:"MessageWeight" hcl:"-"`
	Ref     string `xml:"ref,attr,omitempty" hcl:"ref"`
}

type spikeRate struct {
	XMLName string `xml:"Rate" hcl:"-"`
	Ref     string `xml:"ref,attr,omitempty" hcl:"ref"`
	Value   string `xml:",chardata" hcl:"value"`
}

// DecodeSpikeArrestHCL converts an HCL ast.ObjectItem into a SpikeArrest object.
func DecodeSpikeArrestHCL(item *ast.ObjectItem) (interface{}, error) {
	var p SpikeArrest

	if err := policies.DecodePolicyHCL(item, &p.Policy); err != nil {
		return nil, err
	}

	if err := hcl.DecodeObject(&p, item.Val.(*ast.ObjectType)); err != nil {
		return nil, err
	}

	return &p, nil
}
