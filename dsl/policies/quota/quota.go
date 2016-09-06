package quota

import (
	"fmt"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/kevinswiber/apigee-hcl/dsl/policies/policy"
)

// Quota represents a <Quota/> element.
//
// Documentation: http://docs.apigee.com/api-services/reference/quota-policy
type Quota struct {
	XMLName                   string `xml:"Quota" hcl:"-"`
	policy.Policy             `hcl:",squash"`
	Type                      string         `xml:"type,attr,omitempty" hcl:"type"`
	DisplayName               string         `xml:",omitempty" hcl:"display_name"`
	Allows                    []*allow       `xml:"Allow" hcl:"allow"`
	Interval                  *interval      `hcl:"interval"`
	TimeUnit                  *timeUnit      `hcl:"time_unit"`
	StartTime                 string         `xml:",omitempty" hcl:"start_time"`
	Distributed               bool           `xml:",omitempty" hcl:"distributed"`
	Synchronous               bool           `xml:",omitempty" hcl:"synchronous"`
	AsynchronousConfiguration *asyncConfig   `xml:",omitempty" hcl:"asynchronous_configuration"`
	Identifier                *identifier    `xml:",omitempty" hcl:"identifier"`
	MessageWeight             *messageWeight `xml:",omitempty" hcl:"message_weight"`
}

type allow struct {
	XMLName  string   `xml:"Allow" hcl:"-"`
	Count    int      `xml:"count,attr,omitempty" hcl:"count"`
	CountRef string   `xml:"countRef,attr,omitempty" hcl:"count_ref"`
	Classes  []*class `xml:",omitempty" hcl:"class"`
}

type class struct {
	XMLName string        `xml:"Class" hcl:"-"`
	Ref     string        `xml:"ref,attr,omitempty" hcl:"ref"`
	Allows  []*classAllow `xml:"Allow,omitempty" hcl:"allow"`
}

type classAllow struct {
	XMLName string `xml:"Allow" hcl:"-"`
	Class   string `xml:"class,attr,omitempty" hcl:"class"`
	Count   int    `xml:"count,attr,omitempty" hcl:"count"`
}

type interval struct {
	XMLName string `xml:"Interval" hcl:"-"`
	Ref     string `xml:"ref,attr" hcl:"ref"`
	Value   string `xml:",chardata" hcl:"value"`
}

type timeUnit struct {
	XMLName string `xml:"TimeUnit" hcl:"-"`
	Ref     string `xml:"ref,attr" hcl:"ref"`
	Value   string `xml:",chardata" hcl:"value"`
}

type asyncConfig struct {
	XMLName               string `xml:"AsynchronousConfiguration" hcl:"-"`
	SyncIntervalInSeconds int    `xml:",omitempty" hcl:"sync_interval_in_seconds"`
	SyncMessageCount      int    `xml:",omitempty" hcl:"sync_message_count"`
}

type identifier struct {
	XMLName string `xml:"Identifier" hcl:"-"`
	Ref     string `xml:"ref,attr,omitempty" hcl:"ref"`
}

type messageWeight struct {
	XMLName string `xml:"MessageWeight" hcl:"-"`
	Ref     string `xml:"ref,attr,omitempty" hcl:"ref"`
}

// DecodeHCL converts an HCL ast.ObjectItem into a Quota object.
func DecodeHCL(item *ast.ObjectItem) (interface{}, error) {
	var errors *multierror.Error
	var p Quota

	if err := policy.DecodeHCL(item, &p.Policy); err != nil {
		errors = multierror.Append(errors, err)
		return nil, errors
	}

	if err := hcl.DecodeObject(&p, item.Val.(*ast.ObjectType)); err != nil {
		errors = multierror.Append(errors, err)
		return nil, errors
	}

	var listVal *ast.ObjectList
	if ot, ok := item.Val.(*ast.ObjectType); ok {
		listVal = ot.List
	} else {
		errors = multierror.Append(errors, fmt.Errorf("quota policy not an object"))
		return nil, errors
	}

	if allowList := listVal.Filter("allow"); len(allowList.Items) > 0 {
		a, err := decodeQuotaAllowsHCL(allowList.Items)
		if err != nil {
			errors = multierror.Append(errors, err)
		} else {
			p.Allows = a
		}
	}

	if errors != nil {
		return nil, errors
	}

	return &p, nil
}

func decodeQuotaAllowsHCL(items []*ast.ObjectItem) ([]*allow, error) {
	var result []*allow

	for _, item := range items {
		var a allow
		if err := hcl.DecodeObject(&a, item.Val.(*ast.ObjectType)); err != nil {
			return nil, err
		}

		var listVal *ast.ObjectList
		if ot, ok := item.Val.(*ast.ObjectType); ok {
			listVal = ot.List
		} else {
			return nil, fmt.Errorf("allow not an object")
		}

		if cs := listVal.Filter("class"); len(cs.Items) > 0 {
			classes, err := decodeQuotaAllowClassHCL(cs.Items)
			if err != nil {
				return nil, err
			}
			a.Classes = classes
		}
		result = append(result, &a)
	}

	return result, nil
}

func decodeQuotaAllowClassHCL(items []*ast.ObjectItem) ([]*class, error) {
	var result []*class

	for _, item := range items {
		var c class
		if err := hcl.DecodeObject(&c, item.Val.(*ast.ObjectType)); err != nil {
			return nil, err
		}

		var listVal *ast.ObjectList
		if ot, ok := item.Val.(*ast.ObjectType); ok {
			listVal = ot.List
		} else {
			return nil, fmt.Errorf("class not an object")
		}

		if as := listVal.Filter("allow"); len(as.Items) > 0 {

			classAllows, err := decodeQuotaAllowClassAllowsHCL(as.Items)
			if err != nil {
				return nil, err
			}

			c.Allows = classAllows
		}

		result = append(result, &c)
	}

	return result, nil
}

func decodeQuotaAllowClassAllowsHCL(items []*ast.ObjectItem) ([]*classAllow, error) {
	var result []*classAllow

	for _, item := range items {
		var c classAllow
		if err := hcl.DecodeObject(&c, item.Val.(*ast.ObjectType)); err != nil {
			return nil, err
		}
		result = append(result, &c)
	}

	return result, nil
}
