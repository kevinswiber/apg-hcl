package config

import (
	"fmt"

	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
)

type ProxyEndpoint struct {
	XMLName             string               `xml:"ProxyEndpoint" hcl:"-"`
	Name                string               `xml:"name,attr" hcl:"-"`
	PreFlow             *PreFlow             `hcl:"pre_flow"`
	Flows               []*Flow              `xml:"Flows>Flow", hcl:"flows"`
	PostFlow            *PostFlow            `hcl:"post_flow"`
	PostClientFlow      *PostClientFlow      `hcl:"post_client_flow"`
	HTTPProxyConnection *HTTPProxyConnection `hcl:"http_proxy_connection"`
	RouteRules          []*RouteRule         `xml:"RouteRule" hcl:"route_rule"`
}

type HTTPProxyConnection struct {
	XMLName      string      `xml:"HTTPProxyConnection", hcl:"-"`
	BasePath     string      `hcl:"base_path"`
	VirtualHosts []string    `xml:"VirtualHost" hcl:"virtual_host"`
	Properties   []*Property `xml:"Properties>Property" hcl:"properties"`
}

type Property struct {
	XMLName string      `xml:"Property"`
	Name    string      `xml:"name,attr" hcl:",key"`
	Value   interface{} `xml:",chardata"`
}

type RouteRule struct {
	XMLName        string `xml:"RouteRule"`
	Name           string `xml:"name,attr", hcl:"-"`
	Condition      string `xml:",omitempty" hcl:"condition"`
	TargetEndpoint string `xml:",omitempty" hcl:"target_endpoint"`
	URL            string `xml:",omitempty" hcl:"url"`
}

func loadProxyEndpointsHCL(list *ast.ObjectList) ([]*ProxyEndpoint, error) {
	var result []*ProxyEndpoint
	for _, item := range list.Items {
		n := item.Keys[0].Token.Value().(string)

		var listVal *ast.ObjectList
		if ot, ok := item.Val.(*ast.ObjectType); ok {
			listVal = ot.List
		} else {
			return nil, fmt.Errorf("proxy endpoint not an object")
		}

		var proxyEndpoint ProxyEndpoint

		if err := hcl.DecodeObject(&proxyEndpoint, item.Val); err != nil {
			return nil, err
		}

		proxyEndpoint.Name = n

		if preFlow := listVal.Filter("pre_flow"); len(preFlow.Items) > 0 {
			preFlow, err := loadPreFlowHCL(preFlow)
			if err != nil {
				return nil, err
			}

			proxyEndpoint.PreFlow = preFlow
		}

		if flows := listVal.Filter("flow"); len(flows.Items) > 0 {
			flows, err := loadFlowsHCL(flows)
			if err != nil {
				return nil, err
			}

			proxyEndpoint.Flows = flows
		}

		if postFlow := listVal.Filter("post_flow"); len(postFlow.Items) > 0 {
			postFlow, err := loadPostFlowHCL(postFlow)
			if err != nil {
				return nil, err
			}

			proxyEndpoint.PostFlow = postFlow
		}

		if postClientFlow := listVal.Filter("post_client_flow"); len(postClientFlow.Items) > 0 {
			postClientFlow, err := loadPostClientFlowHCL(postClientFlow)
			if err != nil {
				return nil, err
			}

			proxyEndpoint.PostClientFlow = postClientFlow
		}

		if hpcList := listVal.Filter("http_proxy_connection"); len(hpcList.Items) > 0 {
			fmt.Printf("hpcList len: %d\n", len(hpcList.Items[0].Keys))
			hpc, err := loadProxyEndpointHTTPProxyConnectionHCL(hpcList.Items[0])
			if err != nil {
				return nil, err
			}

			proxyEndpoint.HTTPProxyConnection = hpc
		}

		if routeRulesList := listVal.Filter("route_rule"); len(routeRulesList.Items) > 0 {
			routeRules, err := loadProxyEndpointRouteRulesHCL(routeRulesList)
			if err != nil {
				return nil, err
			}

			proxyEndpoint.RouteRules = routeRules
		}

		result = append(result, &proxyEndpoint)
	}

	return result, nil
}

func loadProxyEndpointHTTPProxyConnectionHCL(item *ast.ObjectItem) (*HTTPProxyConnection, error) {
	var hpc HTTPProxyConnection

	if err := hcl.DecodeObject(&hpc, item.Val); err != nil {
		return nil, fmt.Errorf("error decoding http proxy connection")
	}

	var listVal *ast.ObjectList
	if ot, ok := item.Val.(*ast.ObjectType); ok {
		listVal = ot.List
	} else {
		return nil, fmt.Errorf("http proxy connection not an object")
	}

	if props := listVal.Filter("properties"); len(props.Items) > 0 {
		var propsVal *ast.ObjectList
		if ot, ok := props.Items[0].Val.(*ast.ObjectType); ok {
			propsVal = ot.List
		}

		var newProps []*Property
		for _, p := range propsVal.Items {
			var val interface{}
			if err := hcl.DecodeObject(&val, p.Val); err != nil {
				fmt.Println("can't decode property object")
			}

			newProp := Property{Name: p.Keys[0].Token.Value().(string), Value: val}
			newProps = append(newProps, &newProp)
		}

		hpc.Properties = newProps
	}

	return &hpc, nil
}

func loadProxyEndpointRouteRulesHCL(list *ast.ObjectList) ([]*RouteRule, error) {
	var result []*RouteRule

	for _, item := range list.Items {
		var rule RouteRule
		if err := hcl.DecodeObject(&rule, item.Val); err != nil {
			return nil, fmt.Errorf("error decoding route rule object")
		}
		rule.Name = item.Keys[0].Token.Value().(string)

		result = append(result, &rule)
	}

	return result, nil
}
