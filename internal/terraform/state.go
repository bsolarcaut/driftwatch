package terraform

import (
	"encoding/json"
	"fmt"
	"os"
)

// Resource represents a single resource parsed from Terraform state.
type Resource struct {
	Type       string                 `json:"type"`
	Name       string                 `json:"name"`
	Provider   string                 `json:"provider"`
	Attributes map[string]interface{} `json:"attributes"`
}

// State holds the parsed contents of a Terraform state file.
type State struct {
	Version   int        `json:"version"`
	Resources []Resource `json:"resources"`
}

// tfStateFile mirrors the top-level structure of a terraform.tfstate file.
type tfStateFile struct {
	Version   int              `json:"version"`
	Resources []tfStateResource `json:"resources"`
}

type tfStateResource struct {
	Type      string            `json:"type"`
	Name      string            `json:"name"`
	Provider  string            `json:"provider"`
	Instances []tfStateInstance `json:"instances"`
}

type tfStateInstance struct {
	Attributes map[string]interface{} `json:"attributes"`
}

// LoadState reads and parses a Terraform state file from the given path.
func LoadState(path string) (*State, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading state file %q: %w", path, err)
	}

	var raw tfStateFile
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("parsing state file %q: %w", path, err)
	}

	state := &State{Version: raw.Version}
	for _, r := range raw.Resources {
		attrs := map[string]interface{}{}
		if len(r.Instances) > 0 {
			attrs = r.Instances[0].Attributes
		}
		state.Resources = append(state.Resources, Resource{
			Type:       r.Type,
			Name:       r.Name,
			Provider:   r.Provider,
			Attributes: attrs,
		})
	}

	return state, nil
}
