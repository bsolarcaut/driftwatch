package drift

import (
	"fmt"

	"github.com/your-org/driftwatch/internal/terraform"
)

// Diff represents a detected difference between state and live resource.
type Diff struct {
	ResourceType string
	ResourceName string
	Attribute    string
	StateValue   interface{}
	LiveValue    interface{}
}

func (d Diff) String() string {
	return fmt.Sprintf("%s.%s: attribute %q — state=%v live=%v",
		d.ResourceType, d.ResourceName, d.Attribute, d.StateValue, d.LiveValue)
}

// LiveFetcher retrieves live attribute values for a given resource.
type LiveFetcher interface {
	Fetch(resource terraform.Resource) (map[string]interface{}, error)
}

// Detector compares Terraform state resources against live cloud state.
type Detector struct {
	fetcher LiveFetcher
}

// NewDetector creates a new Detector with the provided LiveFetcher.
func NewDetector(fetcher LiveFetcher) *Detector {
	return &Detector{fetcher: fetcher}
}

// Detect compares each resource in state against live values and returns diffs.
func (d *Detector) Detect(state *terraform.State) ([]Diff, error) {
	var diffs []Diff

	for _, res := range state.Resources {
		liveAttrs, err := d.fetcher.Fetch(res)
		if err != nil {
			return nil, fmt.Errorf("fetching live state for %s.%s: %w", res.Type, res.Name, err)
		}

		for key, stateVal := range res.Attributes {
			liveVal, exists := liveAttrs[key]
			if !exists {
				diffs = append(diffs, Diff{
					ResourceType: res.Type,
					ResourceName: res.Name,
					Attribute:    key,
					StateValue:   stateVal,
					LiveValue:    nil,
				})
				continue
			}
			if fmt.Sprintf("%v", stateVal) != fmt.Sprintf("%v", liveVal) {
				diffs = append(diffs, Diff{
					ResourceType: res.Type,
					ResourceName: res.Name,
					Attribute:    key,
					StateValue:   stateVal,
					LiveValue:    liveVal,
				})
			}
		}
	}

	return diffs, nil
}
