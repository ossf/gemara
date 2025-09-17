package layer1

import (
	"os"

	"github.com/goccy/go-yaml"
)

func goodAIGFExample() (GuidanceDocument, error) {
	testdataPath := "./test-data/good-aigf.yaml"
	data, err := os.ReadFile(testdataPath)
	if err != nil {
		return GuidanceDocument{}, err
	}
	var guidance GuidanceDocument
	if err := yaml.Unmarshal(data, &guidance); err != nil {
		return GuidanceDocument{}, err
	}
	return guidance, nil
}
