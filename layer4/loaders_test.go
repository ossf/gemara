package layer4

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var tests = []struct {
	name       string
	sourcePath string
	wantErr    bool
}{
	{
		name:       "Bad path",
		sourcePath: "file://bad-path.yaml",
		wantErr:    true,
	},
	{
		name:       "Bad YAML",
		sourcePath: "file://test-data/bad.yaml",
		wantErr:    true,
	},
	{
		name:       "Good YAML — Multi Tool Plan",
		sourcePath: "file://test-data/multi-tool-plan.yaml",
		wantErr:    false,
	},
}

func Test_LoadFile(t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &EvaluationPlan{}
			err := e.LoadFile(tt.sourcePath)
			if (err == nil) == tt.wantErr {
				t.Errorf("EvaluationPlan.LoadFile() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				// Validate that the evaluation plan was loaded successfully
				if len(e.Plans) > 0 {
					assert.NotEmpty(t, e.Plans[0].Control.ReferenceId, "Control reference ID should not be empty")
					assert.NotEmpty(t, e.Plans[0].Control.EntryId, "Control entry ID should not be empty")
				}
			}
		})
	}
}
