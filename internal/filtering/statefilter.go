package filtering

import (
	"github.com/snyk/snyk-iac-capture/internal/terraform"
)

var globalWhitelistedAttributes = []string{"id"}

func FilterState(state *terraform.State) (*terraform.State, error) {
	artifact := terraform.State{
		Version:          state.Version,
		TerraformVersion: state.TerraformVersion,
		Lineage:          state.Lineage,
	}

	for _, resource := range state.Resources {
		if resource.Mode != "managed" {
			continue
		}
		var instances []terraform.ResourceInstance
		for _, instance := range resource.Instances {
			attributes := map[string]interface{}{}
			for _, attr := range globalWhitelistedAttributes {
				if _, exists := instance.Attributes[attr]; exists {
					attributes[attr] = instance.Attributes[attr]
				}
			}
			instance.Attributes = attributes
			instances = append(instances, instance)
		}
		resource.Instances = instances
		artifact.Resources = append(artifact.Resources, resource)
	}

	return &artifact, nil
}
