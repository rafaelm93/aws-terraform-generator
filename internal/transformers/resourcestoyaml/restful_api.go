package resourcestoyaml

import (
	"github.com/diagram-code-generator/resources/pkg/resources"

	"github.com/joselitofilho/aws-terraform-generator/internal/generators/config"
	awsresources "github.com/joselitofilho/aws-terraform-generator/internal/resources"
)

func (t *Transformer) buildRestfulAPIRelationship(source, target resources.Resource) {
	if awsresources.ParseResourceType(source.ResourceType()) == awsresources.LambdaType {
		t.buildLambdaToRestfulAPI(source, target)
	}
}

func (t *Transformer) buildRestfulAPIs() []config.RestfulAPI {
	var restfulAPIs []config.RestfulAPI

	restfulAPINames := map[string]struct{}{}

	for _, restfulAPI := range t.resourcesByTypeMap[awsresources.RestfulAPIType] {
		name := restfulAPI.Value()
		if _, ok := restfulAPINames[name]; !ok {
			restfulAPIs = append(restfulAPIs, config.RestfulAPI{Name: name})
			restfulAPINames[name] = struct{}{}
		}
	}

	return restfulAPIs
}
