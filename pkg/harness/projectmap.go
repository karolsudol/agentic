package harness

import (
	"fmt"
	"strings"
)

// ProjectMap implements Progressive Disclosure (Mori Rule 3)
// Instead of dumping entire repos into context, the harness maintains a structured map.
type ProjectMap struct {
	Areas map[string]AreaEntry `json:"areas"`
}

type AreaEntry struct {
	Path        string   `json:"path"`
	Description string   `json:"description"`
	KeyFiles    []string `json:"key_files"`
}

func NewProjectMap() *ProjectMap {
	return &ProjectMap{
		Areas: make(map[string]AreaEntry),
	}
}

func (pm *ProjectMap) RegisterArea(name, path, desc string, keyFiles []string) {
	pm.Areas[name] = AreaEntry{
		Path:        path,
		Description: desc,
		KeyFiles:    keyFiles,
	}
}

func (pm *ProjectMap) Summary() string {
	var sb strings.Builder
	sb.WriteString("=== PROJECT MAP (Attention Index) ===\n")
	for name, area := range pm.Areas {
		sb.WriteString(fmt.Sprintf("• %s -> %s (%s)\n", name, area.Path, area.Description))
		if len(area.KeyFiles) > 0 {
			sb.WriteString(fmt.Sprintf("    Key files: %s\n", strings.Join(area.KeyFiles, ", ")))
		}
	}
	sb.WriteString("Use tool 'load_system_context' with specific area to inspect details.\n")
	return sb.String()
}
