package s3

import (
	_ "embed"
	"os"
	"path"
	"testing"

	generatorserrs "github.com/joselitofilho/aws-terraform-generator/internal/generators/errors"

	"github.com/stretchr/testify/require"
)

var (
	testdataFolder = "../testdata"
	testOutput     = "./testoutput"
)

func TestS3_Build(t *testing.T) {
	type fields struct {
		configFileName string
		output         string
	}

	tests := []struct {
		name             string
		fields           fields
		extraValidations func(testing.TB, string, error)
		targetErr        error
	}{
		{
			name: "default templates for multiple s3",
			fields: fields{
				configFileName: path.Join(testdataFolder, "s3.config.multiple.yaml"),
				output:         path.Join(testOutput, "multiple"),
			},
			extraValidations: func(tb testing.TB, output string, err error) {
				if err != nil {
					return
				}

				require.FileExists(tb, path.Join(output, "mod", "s3.tf"))
			},
		},
		{
			name: "override default template for multiple s3",
			fields: fields{
				configFileName: path.Join(testdataFolder, "s3.config.override.default.tmpls.yaml"),
				output:         path.Join(testOutput, "override"),
			},
			extraValidations: func(tb testing.TB, output string, err error) {
				if err != nil {
					return
				}

				require.FileExists(tb, path.Join(output, "mod", "s3.tf"))
			},
		},
		{
			name: "at least one s3 customising",
			fields: fields{
				configFileName: path.Join(testdataFolder, "s3.config.custom.yaml"),
				output:         path.Join(testOutput, "one"),
			},
			extraValidations: func(tb testing.TB, output string, err error) {
				if err != nil {
					return
				}

				modPath := path.Join(output, "mod")
				require.FileExists(tb, path.Join(modPath, "s3.tf"))
				require.FileExists(tb, path.Join(modPath, "my-second-bucket-s3.tf"))
			},
		},
		{
			name: "all custom s3",
			fields: fields{
				configFileName: path.Join(testdataFolder, "s3.config.allcustom.yaml"),
				output:         path.Join(testOutput, "all"),
			},
			extraValidations: func(tb testing.TB, output string, err error) {
				if err != nil {
					return
				}

				modPath := path.Join(output, "mod")
				require.NoFileExists(tb, path.Join(modPath, "s3.tf"))
				require.FileExists(tb, path.Join(modPath, "my-first-bucket-s3.tf"))
				require.FileExists(tb, path.Join(modPath, "my-second-bucket-s3.tf"))
				require.FileExists(tb, path.Join(modPath, "my-third-bucket-s3.tf"))
			},
		},
		{
			name: "when yaml parser fails should return an error",
			fields: fields{
				configFileName: "",
				output:         "",
			},
			targetErr: generatorserrs.ErrYAMLParser,
		},
	}

	defer func() {
		_ = os.RemoveAll(testOutput)
	}()

	for i := range tests {
		tc := tests[i]

		t.Run(tc.name, func(t *testing.T) {
			err := NewS3(tc.fields.configFileName, tc.fields.output).Build()

			require.ErrorIs(t, err, tc.targetErr)

			if tc.extraValidations != nil {
				tc.extraValidations(t, tc.fields.output, err)
			}
		})
	}
}
