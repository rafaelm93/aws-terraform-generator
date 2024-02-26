package lambda

import (
	_ "embed"
	"fmt"
	"os"
	"strings"

	"github.com/joselitofilho/aws-terraform-generator/internal/generators"
	"github.com/joselitofilho/aws-terraform-generator/internal/generators/config"
	"github.com/joselitofilho/aws-terraform-generator/internal/utils"
)

type Lambda struct {
	configFileName string
	output         string
}

func NewLambda(configFileName, output string) *Lambda {
	return &Lambda{configFileName: configFileName, output: output}
}

func (l *Lambda) Build() error {
	yamlParser := config.NewYAML(l.configFileName)

	yamlConfig, err := yamlParser.Parse()
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	for i := range yamlConfig.Lambdas {
		lambdaConf := yamlConfig.Lambdas[i]

		envars := map[string]string{}
		for i := range lambdaConf.Envars {
			for key, value := range lambdaConf.Envars[i] {
				envars[key] = value
			}
		}

		sqsTriggers := make([]SQSTrigger, len(lambdaConf.SQSTriggers))
		for i := range lambdaConf.SQSTriggers {
			sqsTriggers[i] = SQSTrigger{
				SourceARN: lambdaConf.SQSTriggers[i].SourceARN,
			}
		}

		crons := make([]Cron, len(lambdaConf.Crons))
		for i := range lambdaConf.Crons {
			crons[i] = Cron{
				ScheduleExpression: lambdaConf.Crons[i].ScheduleExpression,
				IsEnabled:          lambdaConf.Crons[i].IsEnabled,
			}
		}

		filesConf := generators.CreateFilesMap(lambdaConf.Files)

		asModule := strings.Contains(lambdaConf.Source, "git@")

		roleName := lambdaConf.RoleName
		if roleName == "" {
			roleName = "iam_for_lambda"
		}

		data := Data{
			Name:        lambdaConf.Name,
			AsModule:    asModule,
			Source:      lambdaConf.Source,
			RoleName:    roleName,
			Runtime:     lambdaConf.Runtime,
			Description: lambdaConf.Description,
			Envars:      envars,
			SQSTriggers: sqsTriggers,
			Crons:       crons,
			Files:       filesConf,
		}

		output := fmt.Sprintf("%s/mod", l.output)
		_ = os.MkdirAll(output, os.ModePerm)

		outputFile := fmt.Sprintf("%s/%s.tf", output, lambdaConf.Name)

		tmplName := "lambda-tf-template"
		tmpl := string(lambdaTFTmpl)

		err = generators.BuildFile(data, tmplName, tmpl, outputFile)
		if err != nil {
			return fmt.Errorf("%w", err)
		}

		err = utils.TerraformFormat(output)
		if err != nil {
			fmt.Printf("%s: %s\n", lambdaConf.Name, err)
		}

		fmt.Printf("Terraform '%s' has been generated successfully\n", lambdaConf.Name)

		output = fmt.Sprintf("%s/lambda/%s", l.output, lambdaConf.Name)
		_ = os.MkdirAll(output, os.ModePerm)

		err = generators.GenerateFiles(defaultTemplatesMap, filesConf, output, data)
		if err != nil {
			return fmt.Errorf("%w", err)
		}

		fmt.Printf("Lambda '%s' has been generated successfully\n", lambdaConf.Name)
	}

	return nil
}