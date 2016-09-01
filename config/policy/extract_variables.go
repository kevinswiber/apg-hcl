package policy

import (
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
)

type ExtractVariablesPolicy struct {
	XMLName                   string `xml:"ExtractVariables" hcl:"-"`
	Policy                    `hcl:",squash"`
	DisplayName               string          `xml:",omitempty" hcl:"display_name"`
	Source                    *evSource       `xml:",omitempty" hcl:"source"`
	VariablePrefix            string          `xml:",omitempty" hcl:"variable_prefix"`
	IgnoreUnresolvedVariables bool            `xml:",omitempty" hcl:"ignore_unresolved_variables"`
	URIPath                   *evURIPath      `xml:",omitempty" hcl:"uri_path"`
	QueryParams               []*evQueryParam `xml:"QueryParam,omitempty" hcl:"query_param"`
	Headers                   []*evHeader     `xml:"Header,omitempty" hcl:"header"`
	FormParams                []*evFormParam  `xml:"FormParam,omitempty" hcl:"form_param"`
	Variables                 []*evVariable   `xml:"Variable,omitempty" hcl:"variable"`
	JSONPayload               *evJSONPayload  `xml:",omitempty" hcl:"json_payload"`
	XMLPayload                *evXMLPayload   `xml:",omitempty" hcl:"xml_payload"`
}

type evSource struct {
	XMLName      string `xml:"Source" hcl:"-"`
	ClearPayload bool   `xml:"clearPayload,attr,omitempty" hcl:"clear_payload"`
	Value        string `xml:",chardata" hcl:"value"`
}

type evURIPath struct {
	XMLName  string       `xml:"URIPath" hcl:"-"`
	Patterns []*evPattern `xml:"Pattern" hcl:"pattern"`
}

type evQueryParam struct {
	XMLName  string       `xml:"QueryParam" hcl:"-"`
	Name     string       `xml:",attr" hcl:"name"`
	Patterns []*evPattern `xml:"Pattern" hcl:"pattern"`
}

type evHeader struct {
	XMLName  string       `xml:"Header" hcl:"-"`
	Name     string       `xml:",attr" hcl:"name"`
	Patterns []*evPattern `xml:"Pattern" hcl:"pattern"`
}

type evFormParam struct {
	XMLName  string       `xml:"FormParam" hcl:"-"`
	Name     string       `xml:",attr" hcl:"name"`
	Patterns []*evPattern `xml:"Pattern" hcl:"pattern"`
}

type evVariable struct {
	XMLName  string       `xml:"Variable" hcl:"-"`
	Name     string       `xml:",attr" hcl:"name"`
	Patterns []*evPattern `xml:"Pattern" hcl:"pattern"`
}

type evJSONPayload struct {
	XMLName   string                   `xml:"JSONPayload" hcl:"-"`
	Variables []*evJSONPayloadVariable `xml:"Variable" hcl:"variable"`
}

type evJSONPayloadVariable struct {
	XMLName  string `xml:"Variable" hcl:"-"`
	Name     string `xml:"name,attr" hcl:"name"`
	Type     string `xml:"type,attr,omitempty" hcl:"type"`
	JSONPath string `hcl:"json_path"`
}

type evXMLPayload struct {
	XMLName               string                  `xml:"XMLPayload" hcl:"-"`
	StopPayloadProcessing bool                    `xml:"stopPayloadProcessing,attr,omitempty" hcl:"stop_payload_processing"`
	Variables             []*evXMLPayloadVariable `xml:"Variable" hcl:"-"`
}

type evXMLNamespace struct {
	Prefix string `xml:"prefix,attr,omitempty" hcl:"prefix"`
	Value  string `xml:",chardata" hcl:"value"`
}

type evXMLPayloadVariable struct {
	XMLName string `xml:"Variable" hcl:"-"`
	Name    string `xml:"name,attr" hcl:"name"`
	Type    string `xml:"type,attr,omitempty" hcl:"type"`
	XPath   string `hcl:"xpath"`
}

type evPattern struct {
	XMLName    string `xml:"Pattern" hcl:"-"`
	IgnoreCase bool   `xml:"ignoreCase,attr,omitempty" hcl:"ignore_case"`
	Value      string `xml:",chardata" hcl:"value"`
}

func LoadExtractVariablesHCL(item *ast.ObjectItem) (interface{}, error) {
	var p ExtractVariablesPolicy

	if err := LoadCommonPolicyHCL(item, &p.Policy); err != nil {
		return nil, err
	}

	if err := hcl.DecodeObject(&p, item.Val.(*ast.ObjectType)); err != nil {
		return nil, err
	}

	return p, nil
}
