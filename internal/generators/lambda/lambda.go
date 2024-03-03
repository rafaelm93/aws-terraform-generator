package lambda

import (
	_ "embed"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/joselitofilho/aws-terraform-generator/internal/generators"
	"github.com/joselitofilho/aws-terraform-generator/internal/generators/config"
	generatorserrs "github.com/joselitofilho/aws-terraform-generator/internal/generators/errors"
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
		return fmt.Errorf("%w: %w", generatorserrs.ErrYAMLParse, err)
	}

	tfTemplates := utils.MergeStringMap(
		generators.FilterTemplatesMap(".tf", generators.CreateTemplatesMap(yamlConfig.OverrideDefaultTemplates.Lambda)),
		defaultTfTemplatesMap)

	goTemplates := utils.MergeStringMap(
		generators.FilterTemplatesMap(".go", generators.CreateTemplatesMap(yamlConfig.OverrideDefaultTemplates.Lambda)),
		defaultGoTemplatesMap)

	for i := range yamlConfig.Lambdas {
		lambdaConf := yamlConfig.Lambdas[i]

		envars := map[string]string{}

		for i := range lambdaConf.Envars {
			for key, value := range lambdaConf.Envars[i] {
				envars[key] = value
			}
		}

		crons := buildCrons(&lambdaConf)
		kinesisTriggers := buildKinesisTriggers(&lambdaConf)
		sqsTriggers := buildSQSTriggers(&lambdaConf)

		filesConf := generators.CreateFilesMap(lambdaConf.Files)

		asModule := strings.Contains(lambdaConf.Source, "git@")

		roleName := lambdaConf.RoleName
		if roleName == "" {
			roleName = "iam_for_lambda"
		}

		data := Data{
			Name:            lambdaConf.Name,
			AsModule:        asModule,
			Source:          lambdaConf.Source,
			RoleName:        roleName,
			Runtime:         lambdaConf.Runtime,
			Description:     lambdaConf.Description,
			Envars:          envars,
			KinesisTriggers: kinesisTriggers,
			SQSTriggers:     sqsTriggers,
			Crons:           crons,
			Files:           filesConf,
		}

		output := path.Join(l.output, "mod")
		_ = os.MkdirAll(output, os.ModePerm)

		outputFile := path.Join(output, lambdaConf.Name+".tf")

		err = generators.GenerateFile(tfTemplates, filenameTfLambda, "", outputFile, data)
		if err != nil {
			return fmt.Errorf("%w", err)
		}

		fmt.Printf("Terraform '%s' has been generated successfully\n", lambdaConf.Name)

		output = fmt.Sprintf("%s/lambda/%s", l.output, lambdaConf.Name)
		_ = os.MkdirAll(output, os.ModePerm)

		err = generators.GenerateFiles(goTemplates, filesConf, data, output)
		if err != nil {
			return fmt.Errorf("%w", err)
		}

		fmt.Printf("Lambda '%s' has been generated successfully\n", lambdaConf.Name)
	}

	return nil
}

func buildCrons(lambdaConf *config.Lambda) []Cron {
	crons := make([]Cron, len(lambdaConf.Crons))
	for i := range lambdaConf.Crons {
		crons[i] = Cron{
			ScheduleExpression: lambdaConf.Crons[i].ScheduleExpression,
			IsEnabled:          lambdaConf.Crons[i].IsEnabled,
		}
	}

	return crons
}

func buildKinesisTriggers(lambdaConf *config.Lambda) []KinesisTrigger {
	kinesisTriggers := make([]KinesisTrigger, len(lambdaConf.KinesisTriggers))
	for i := range lambdaConf.KinesisTriggers {
		kinesisTriggers[i] = KinesisTrigger{
			SourceARN: lambdaConf.KinesisTriggers[i].SourceARN,
		}
	}

	return kinesisTriggers
}

func buildSQSTriggers(lambdaConf *config.Lambda) []SQSTrigger {
	sqsTriggers := make([]SQSTrigger, len(lambdaConf.SQSTriggers))
	for i := range lambdaConf.SQSTriggers {
		sqsTriggers[i] = SQSTrigger{
			SourceARN: lambdaConf.SQSTriggers[i].SourceARN,
		}
	}

	return sqsTriggers
}
