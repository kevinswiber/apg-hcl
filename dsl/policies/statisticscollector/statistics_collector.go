package statisticscollector

import (
	"fmt"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/kevinswiber/apigee-hcl/dsl/hclerror"
	"github.com/kevinswiber/apigee-hcl/dsl/policies"
	"strings"
)

// StatisticsCollector represents an <StatisticsCollector/> element.
//
// Documentation: http://docs.apigee.com/api-services/reference/statistics-collector-policy
type StatisticsCollector struct {
	XMLName         string `xml:"StatisticsCollector" hcl:"-"`
	policies.Policy `hcl:",squash"`
	DisplayName     string         `xml:",omitempty" hcl:"display_name"`
	Statistics      []*scStatistic `xml:"Statistics>Statistic" hcl:"statistic"`
}

// GetName returns the policy name.
func (policy StatisticsCollector) GetName() string {
	return policy.Name
}

type scStatistic struct {
	XMLName string `xml:"Statistic" hcl:"-"`
	Name    string `xml:"name,attr" hcl:"-"`
	Ref     string `xml:"ref,attr" hcl:"ref"`
	Type    string `xml:"type,attr,omitempty" hcl:"type"`
	Value   string `xml:",chardata" hcl:"value"`
}

// DecodeStatisticsCollectorHCL converts an HCL ast.ObjectItem into an StatisticsCollector object.
func DecodeStatisticsCollectorHCL(item *ast.ObjectItem) (interface{}, error) {
	var errors *multierror.Error
	var p StatisticsCollector

	if err := policies.DecodePolicyHCL(item, &p.Policy); err != nil {
		return nil, err
	}

	var listVal *ast.ObjectList
	if ot, ok := item.Val.(*ast.ObjectType); ok {
		listVal = ot.List
	} else {
		pos := item.Val.Pos()
		newError := hclerror.PosError{
			Pos: pos,
			Err: fmt.Errorf("statistics_collector policy not an object"),
		}
		return nil, &newError
	}

	if err := hcl.DecodeObject(&p, item.Val.(*ast.ObjectType)); err != nil {
		return nil, err
	}

	if statsList := listVal.Filter("statistic"); len(statsList.Items) > 0 {
		stats, err := loadStatisticsCollectorStatisticHCL(statsList.Items)
		if err != nil {
			errors = multierror.Append(errors, err)
		} else {
			p.Statistics = stats
		}
	}

	if errors != nil {
		return nil, errors
	}

	return p, nil
}

func loadStatisticsCollectorStatisticHCL(items []*ast.ObjectItem) ([]*scStatistic, error) {
	var errors *multierror.Error
	var stats []*scStatistic
	for _, item := range items {
		var stat scStatistic

		if _, ok := item.Val.(*ast.ObjectType); !ok {
			pos := item.Val.Pos()
			newError := hclerror.PosError{
				Pos: pos,
				Err: fmt.Errorf("statistic not an object"),
			}
			return nil, &newError
		}

		if err := hcl.DecodeObject(&stat, item.Val.(*ast.ObjectType)); err != nil {
			return nil, err
		}

		if len(item.Keys) == 0 || item.Keys[0].Token.Value().(string) == "" {
			pos := item.Val.Pos()
			newError := hclerror.PosError{
				Pos: pos,
				Err: fmt.Errorf("statistic requires a name"),
			}
			return nil, &newError
		}

		stat.Name = item.Keys[0].Token.Value().(string)

		if stat.Ref == "" {
			pos := item.Val.Pos()
			newError := hclerror.PosError{
				Pos: pos,
				Err: fmt.Errorf("statistic requires a ref value"),
			}
			errors = multierror.Append(errors, &newError)
		}

		validTypes := []string{"string", "integer", "float", "long", "double"}
		if stat.Type == "" {
			pos := item.Val.Pos()
			newError := hclerror.PosError{
				Pos: pos,
				Err: fmt.Errorf("statistic requires a type value"),
			}
			errors = multierror.Append(errors, &newError)
		} else {
			hasValidType := false
			for _, t := range validTypes {
				if stat.Type == t {
					hasValidType = true
					break
				}
			}

			if !hasValidType {
				pos := item.Val.Pos()
				newError := hclerror.PosError{
					Pos: pos,
					Err: fmt.Errorf("statistic requires a valid type value [%s]",
						strings.Join(validTypes, ", ")),
				}
				errors = multierror.Append(errors, &newError)
			}
		}

		stats = append(stats, &stat)
	}

	if errors != nil {
		return nil, errors
	}

	return stats, nil
}
