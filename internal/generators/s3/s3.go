package s3

import (
	_ "embed"
	"fmt"

	"github.com/joselitofilho/aws-terraform-generator/internal/generators"
	"github.com/joselitofilho/aws-terraform-generator/internal/generators/config"
	"github.com/joselitofilho/aws-terraform-generator/internal/utils"
)

type Data struct {
	Name           string
	ExpirationDays int
}

type S3 struct {
	configFileName string
	output         string
}

func NewS3(configFileName, output string) *S3 {
	return &S3{configFileName: configFileName, output: output}
}

func (s *S3) Build() error {
	yamlParser := config.NewYAML(s.configFileName)

	yamlConfig, err := yamlParser.Parse()
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	tmplName := "s3-tf-template"
	result := ""

	for i := range yamlConfig.Buckets {
		conf := yamlConfig.Buckets[i]

		data := Data{
			Name:           conf.Name,
			ExpirationDays: conf.ExpirationDays,
		}

		if len(conf.Files) > 0 {
			filesConf := generators.CreateFilesMap(conf.Files)

			err = generators.GenerateFiles(defaultTemplatesMap, filesConf, fmt.Sprintf("%s/mod", s.output), data)
			if err != nil {
				return fmt.Errorf("%w", err)
			}

			fmt.Printf("S3 '%s' has been generated successfully\n", conf.Name)

			continue
		}

		outputData, err := generators.Build(data, tmplName, string(s3TFTmpl))
		if err != nil {
			return fmt.Errorf("%w", err)
		}

		result = fmt.Sprintf("%s\n%s", result, outputData)
	}

	if result != "" {
		outputFile := fmt.Sprintf("%s/mod/s3.tf", s.output)

		if err := generators.BuildFile(Data{}, tmplName, result, outputFile); err != nil {
			return fmt.Errorf("%w", err)
		}

		if err := utils.TerraformFormat(outputFile); err != nil {
			fmt.Println(err)
		}

		fmt.Println("S3 has been generated successfully")
	}

	return nil
}
