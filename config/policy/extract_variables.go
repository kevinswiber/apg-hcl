package policy

import (
	"fmt"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/kevinswiber/apigee-hcl/config/hclerror"
)

type ExtractVariablesPolicy struct {
	XMLName                   string `xml:"ExtractVariables" hcl:"-"`
	Policy                    `hcl:",squash"`
	DisplayName               string          `xml:",omitempty" hcl:"display_name"`
	Source                    *evSource       `xml:",omitempty" hcl:"source"`
	VariablePrefix            string          `xml:",omitempty" hcl:"variable_prefix"`
	IgnoreUnresolvedVariables bool            `xml:",omitempty" hcl:"ignore_unresolved_variables"`
	URIPaths                  []*evURIPath    `xml:"URIPath,omitempty" hcl:"uri_path"`
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
	Name     string       `xml:"name,attr" hcl:"-"`
	Patterns []*evPattern `xml:"Pattern" hcl:"pattern"`
}

type evHeader struct {
	XMLName  string       `xml:"Header" hcl:"-"`
	Name     string       `xml:"name,attr" hcl:"-"`
	Patterns []*evPattern `xml:"Pattern" hcl:"pattern"`
}

type evFormParam struct {
	XMLName  string       `xml:"FormParam" hcl:"-"`
	Name     string       `xml:"name,attr" hcl:"-"`
	Patterns []*evPattern `xml:"Pattern" hcl:"pattern"`
}

type evVariable struct {
	XMLName  string       `xml:"Variable" hcl:"-"`
	Name     string       `xml:"name,attr" hcl:"-"`
	Patterns []*evPattern `xml:"Pattern" hcl:"pattern"`
}

type evJSONPayload struct {
	XMLName   string                   `xml:"JSONPayload" hcl:"-"`
	Variables []*evJSONPayloadVariable `xml:"Variable" hcl:"variable"`
}

type evJSONPayloadVariable struct {
	XMLName  string `xml:"Variable" hcl:"-"`
	Name     string `xml:"name,attr" hcl:",key"`
	Type     string `xml:"type,attr,omitempty" hcl:"type"`
	JSONPath string `hcl:"json_path"`
}

type evXMLPayload struct {
	XMLName               string                  `xml:"XMLPayload" hcl:"-"`
	StopPayloadProcessing bool                    `xml:"stopPayloadProcessing,attr,omitempty" hcl:"stop_payload_processing"`
	Namespaces            []*evXMLNamespace       `xml:"Namespaces>Namespace,omitempty" hcl:"namespace"`
	Variables             []*evXMLPayloadVariable `xml:"Variable" hcl:"variable"`
}

type evXMLNamespace struct {
	Prefix string `xml:"prefix,attr,omitempty" hcl:",key"`
	Value  string `xml:",chardata" hcl:"value"`
}

type evXMLPayloadVariable struct {
	XMLName string `xml:"Variable" hcl:"-"`
	Name    string `xml:"name,attr" hcl:",key"`
	Type    string `xml:"type,attr,omitempty" hcl:"type"`
	XPath   string `hcl:"xpath"`
}

type evPattern struct {
	XMLName    string `xml:"Pattern" hcl:"-"`
	IgnoreCase bool   `xml:"ignoreCase,attr,omitempty" hcl:"ignore_case"`
	Value      string `xml:",chardata" hcl:"value"`
}

func LoadExtractVariablesHCL(item *ast.ObjectItem) (interface{}, error) {
	var errors *multierror.Error
	var p ExtractVariablesPolicy

	if err := LoadCommonPolicyHCL(item, &p.Policy); err != nil {
		return nil, err
	}

	var listVal *ast.ObjectList
	if ot, ok := item.Val.(*ast.ObjectType); ok {
		listVal = ot.List
	} else {
		pos := item.Val.Pos()
		newError := hclerror.PosError{
			Pos: pos,
			Err: fmt.Errorf("extract variables policy not an object"),
		}
		return nil, &newError
	}

	if err := hcl.DecodeObject(&p, item.Val.(*ast.ObjectType)); err != nil {
		return nil, err
	}

	if uriPathList := listVal.Filter("uri_path"); len(uriPathList.Items) > 0 {
		uriPaths, err := loadExtractVariablesURIPathsHCL(uriPathList.Items)
		if err != nil {
			errors = multierror.Append(errors, err)
		} else {
			p.URIPaths = uriPaths
		}
	}

	if queryParamList := listVal.Filter("query_param"); len(queryParamList.Items) > 0 {
		queryParams, err := loadExtractVariablesQueryParamsHCL(queryParamList.Items)
		if err != nil {
			errors = multierror.Append(errors, err)
		} else {
			p.QueryParams = queryParams
		}
	}

	if headerList := listVal.Filter("header"); len(headerList.Items) > 0 {
		headers, err := loadExtractVariablesHeadersHCL(headerList.Items)
		if err != nil {
			errors = multierror.Append(errors, err)
		} else {
			fmt.Printf("headers: %+v\n", headers)
			p.Headers = headers
		}
	}

	if formParamList := listVal.Filter("form_param"); len(formParamList.Items) > 0 {
		formParams, err := loadExtractVariablesFormParamsHCL(formParamList.Items)
		if err != nil {
			errors = multierror.Append(errors, err)
		} else {
			p.FormParams = formParams
		}
	}

	if variableList := listVal.Filter("variable"); len(variableList.Items) > 0 {
		variables, err := loadExtractVariablesVariablesHCL(variableList.Items)
		if err != nil {
			errors = multierror.Append(errors, err)
		} else {
			p.Variables = variables
		}
	}

	if errors != nil {
		return nil, errors
	}

	return p, nil
}

func loadExtractVariablesURIPathsHCL(items []*ast.ObjectItem) ([]*evURIPath, error) {
	var uriPaths []*evURIPath
	for _, item := range items {
		var up evURIPath

		var listVal *ast.ObjectList
		if ot, ok := item.Val.(*ast.ObjectType); ok {
			listVal = ot.List
		} else {
			pos := item.Val.Pos()
			newError := hclerror.PosError{
				Pos: pos,
				Err: fmt.Errorf("uri_path not an object"),
			}
			return nil, &newError
		}

		if err := hcl.DecodeObject(&up, item.Val.(*ast.ObjectType)); err != nil {
			return nil, err
		}

		if patternList := listVal.Filter("pattern"); len(patternList.Items) > 0 {
			patterns, err := loadExtractVariablesPatternsHCL(patternList.Items)
			if err != nil {
				return nil, err
			}

			up.Patterns = patterns
		}

		uriPaths = append(uriPaths, &up)
	}
	return uriPaths, nil
}

func loadExtractVariablesQueryParamsHCL(items []*ast.ObjectItem) ([]*evQueryParam, error) {
	var queryParams []*evQueryParam
	for _, item := range items {
		var qp evQueryParam

		var listVal *ast.ObjectList
		if ot, ok := item.Val.(*ast.ObjectType); ok {
			listVal = ot.List
		} else {
			pos := item.Val.Pos()
			newError := hclerror.PosError{
				Pos: pos,
				Err: fmt.Errorf("query_param not an object"),
			}
			return nil, &newError
		}

		if err := hcl.DecodeObject(&qp, item.Val.(*ast.ObjectType)); err != nil {
			return nil, err
		}

		if len(item.Keys) == 0 || item.Keys[0].Token.Value().(string) == "" {
			pos := item.Val.Pos()
			newError := hclerror.PosError{
				Pos: pos,
				Err: fmt.Errorf("query_param requires a name"),
			}
			return nil, &newError
		}

		qp.Name = item.Keys[0].Token.Value().(string)

		if patternList := listVal.Filter("pattern"); len(patternList.Items) > 0 {
			patterns, err := loadExtractVariablesPatternsHCL(patternList.Items)
			if err != nil {
				return nil, err
			}

			qp.Patterns = patterns
		}

		queryParams = append(queryParams, &qp)
	}
	return queryParams, nil
}

func loadExtractVariablesHeadersHCL(items []*ast.ObjectItem) ([]*evHeader, error) {
	var headers []*evHeader
	for _, item := range items {
		var hdr evHeader

		var listVal *ast.ObjectList
		if ot, ok := item.Val.(*ast.ObjectType); ok {
			listVal = ot.List
		} else {
			pos := item.Val.Pos()
			newError := hclerror.PosError{
				Pos: pos,
				Err: fmt.Errorf("header not an object"),
			}
			return nil, &newError
		}

		if err := hcl.DecodeObject(&hdr, item.Val.(*ast.ObjectType)); err != nil {
			return nil, err
		}

		if len(item.Keys) == 0 || item.Keys[0].Token.Value().(string) == "" {
			pos := item.Val.Pos()
			newError := hclerror.PosError{
				Pos: pos,
				Err: fmt.Errorf("header requires a name"),
			}
			return nil, &newError
		}

		hdr.Name = item.Keys[0].Token.Value().(string)

		if patternList := listVal.Filter("pattern"); len(patternList.Items) > 0 {
			fmt.Printf("header pattern count: %d", len(items))
			patterns, err := loadExtractVariablesPatternsHCL(patternList.Items)
			if err != nil {
				return nil, err
			}

			hdr.Patterns = patterns
		} else {
			fmt.Println("No header patterns")
		}

		headers = append(headers, &hdr)
	}
	return headers, nil
}

func loadExtractVariablesFormParamsHCL(items []*ast.ObjectItem) ([]*evFormParam, error) {
	var formParams []*evFormParam
	for _, item := range items {
		var fp evFormParam

		var listVal *ast.ObjectList
		if ot, ok := item.Val.(*ast.ObjectType); ok {
			listVal = ot.List
		} else {
			pos := item.Val.Pos()
			newError := hclerror.PosError{
				Pos: pos,
				Err: fmt.Errorf("form_param not an object"),
			}
			return nil, &newError
		}

		if err := hcl.DecodeObject(&fp, item.Val.(*ast.ObjectType)); err != nil {
			return nil, err
		}

		if len(item.Keys) == 0 || item.Keys[0].Token.Value().(string) == "" {
			pos := item.Val.Pos()
			newError := hclerror.PosError{
				Pos: pos,
				Err: fmt.Errorf("form_param requires a name"),
			}
			return nil, &newError
		}

		fp.Name = item.Keys[0].Token.Value().(string)

		if patternList := listVal.Filter("pattern"); len(patternList.Items) > 0 {
			patterns, err := loadExtractVariablesPatternsHCL(patternList.Items)
			if err != nil {
				return nil, err
			}

			fp.Patterns = patterns
		}

		formParams = append(formParams, &fp)
	}
	return formParams, nil
}

func loadExtractVariablesVariablesHCL(items []*ast.ObjectItem) ([]*evVariable, error) {
	var variables []*evVariable
	for _, item := range items {
		var v evVariable

		var listVal *ast.ObjectList
		if ot, ok := item.Val.(*ast.ObjectType); ok {
			listVal = ot.List
		} else {
			pos := item.Val.Pos()
			newError := hclerror.PosError{
				Pos: pos,
				Err: fmt.Errorf("variable not an object"),
			}
			return nil, &newError
		}

		if err := hcl.DecodeObject(&v, item.Val.(*ast.ObjectType)); err != nil {
			return nil, err
		}

		if len(item.Keys) == 0 || item.Keys[0].Token.Value().(string) == "" {
			pos := item.Val.Pos()
			newError := hclerror.PosError{
				Pos: pos,
				Err: fmt.Errorf("variable requires a name"),
			}
			return nil, &newError
		}

		v.Name = item.Keys[0].Token.Value().(string)

		if patternList := listVal.Filter("pattern"); len(patternList.Items) > 0 {
			patterns, err := loadExtractVariablesPatternsHCL(patternList.Items)
			if err != nil {
				return nil, err
			}

			v.Patterns = patterns
		}

		variables = append(variables, &v)
	}
	return variables, nil
}

func loadExtractVariablesPatternsHCL(items []*ast.ObjectItem) ([]*evPattern, error) {
	var patterns []*evPattern
	for _, item := range items {
		var pat evPattern

		if _, ok := item.Val.(*ast.ObjectType); !ok {
			pos := item.Val.Pos()
			newError := hclerror.PosError{
				Pos: pos,
				Err: fmt.Errorf("pattern not an object"),
			}
			return nil, &newError
		}

		if err := hcl.DecodeObject(&pat, item.Val.(*ast.ObjectType)); err != nil {
			return nil, err
		}

		patterns = append(patterns, &pat)
	}

	return patterns, nil
}
