package guides

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGuideDiagram(t *testing.T) {
	type args struct {
		workdir string
		fileMap map[string][]string
	}

	tests := []struct {
		name      string
		args      args
		want      *DiagramAnswers
		targetErr error
	}{
		// TODO: Add test cases.
	}

	for i := range tests {
		tc := tests[i]

		t.Run(tc.name, func(t *testing.T) {
			got, err := GuideDiagram(tc.args.workdir, tc.args.fileMap)

			if tc.targetErr == nil {
				require.NoError(t, err)
				require.Equal(t, tc.want, got)
			} else {
				require.ErrorIs(t, err, tc.targetErr)
			}
		})
	}
}
