package structure

import (
	_ "embed"
	"fmt"
	"os"

	"github.com/joselitofilho/aws-terraform-generator/internal/templates"
	"github.com/joselitofilho/aws-terraform-generator/internal/templates/config"
)

type Structure struct {
	input  string
	output string
}

func NewStructure(input, output string) *Structure {
	return &Structure{input: input, output: output}
}

func (s *Structure) Build() error {
	yamlParser := config.NewYAML(s.input)

	yamlConfig, err := yamlParser.Parse()
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	defaultTemplatesMap := map[string]string{}
	for i := range yamlConfig.Structure.DefaultTemplates {
		defaultTemplatesMap = yamlConfig.Structure.DefaultTemplates[i]
	}

	for i := range yamlConfig.Structure.Stacks {
		conf := yamlConfig.Structure.Stacks[i]

		data := Data{
			StackName: conf.StackName,
		}

		for _, folder := range conf.Folders {
			output := fmt.Sprintf("%s/%s/%s", s.output, conf.StackName, folder.Name)
			_ = os.MkdirAll(output, os.ModePerm)

			for _, file := range folder.Files {
				outputFile := fmt.Sprintf("%s/%s", output, file.Name)

				err = templates.GenerateFiles(defaultTemplatesMap, file.Name, file.Tmpl, data, outputFile)
				if err != nil {
					return fmt.Errorf("%w", err)
				}
			}
		}

		fmt.Printf("Structure '%s' has been generated successfully\n", conf.StackName)
	}

	return nil
}