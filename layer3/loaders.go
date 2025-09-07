package layer3

import (
	"fmt"
	"path"

	"github.com/ossf/gemara/internal/loaders"
)

// LoadFile loads data from a YAML or JSON file at the provided path.
// If run multiple times for the same data type, this method will override previous data.
func (c *PolicyDocument) LoadFile(sourcePath string) error {
	ext := path.Ext(sourcePath)
	switch ext {
	case ".yaml", ".yml":
		err := loaders.LoadYAML(sourcePath, c)
		if err != nil {
			return err
		}
	case ".json":
		err := loaders.LoadJSON(sourcePath, c)
		if err != nil {
			return fmt.Errorf("error loading json: %w", err)
		}
	default:
		return fmt.Errorf("unsupported file extension: %s", ext)
	}
	return nil
}
