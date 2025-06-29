package ontology

import (
	"fmt"
	"strings"
)

// IRIMapperAgent generates IRIs and optional doc URLs.
type IRIMapperAgent struct {
	Prefix *PrefixAgent
}

func NewIRIMapperAgent(p *PrefixAgent) *IRIMapperAgent {
	return &IRIMapperAgent{Prefix: p}
}

func toPascal(s string) string {
	if s == "" {
		return s
	}
	parts := strings.Split(s, "_")
	if len(parts) == 1 {
		parts = strings.Split(s, "-")
	}
	for i, p := range parts {
		if p == "" {
			continue
		}
		parts[i] = strings.ToUpper(p[:1]) + p[1:]
	}
	return strings.Join(parts, "")
}

func (m *IRIMapperAgent) ClassIRI(lexID, name string) string {
	return m.Prefix.IRI("bsky", fmt.Sprintf("%s#%s", lexID, toPascal(name)))
}

func (m *IRIMapperAgent) FieldIRI(lexID, class, field string) string {
	return fmt.Sprintf("%s/%s", m.ClassIRI(lexID, class), field)
}
