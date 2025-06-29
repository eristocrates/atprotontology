package ontology

import (
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
)

// LexDef represents a simplified lexicon definition.
type LexDef struct {
	LexiconID   string             `json:"-"`
	Name        string             `json:"-"`
	Type        string             `json:"type"`
	Description string             `json:"description"`
	Required    []string           `json:"required"`
	Properties  map[string]LexProp `json:"properties"`
}

type LexProp struct {
	Type        string   `json:"type"`
	Ref         string   `json:"ref"`
	Items       *LexProp `json:"items"`
	Description string   `json:"description"`
}

type Lexicon struct {
	ID   string            `json:"id"`
	Defs map[string]LexDef `json:"defs"`
}

// LexSchemaAgent parses lexicon JSON into definitions.
type LexSchemaAgent struct {
	Defs map[string]LexDef
}

func NewLexSchemaAgent() *LexSchemaAgent {
	return &LexSchemaAgent{Defs: map[string]LexDef{}}
}

func (l *LexSchemaAgent) Load(dir string) error {
	return filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || filepath.Ext(path) != ".json" {
			return nil
		}
		b, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		var lx Lexicon
		if err := json.Unmarshal(b, &lx); err != nil {
			return err
		}
		for name, def := range lx.Defs {
			def.LexiconID = lx.ID
			def.Name = name
			l.Defs[lx.ID+"#"+name] = def
		}
		return nil
	})
}
