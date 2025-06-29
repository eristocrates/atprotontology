package ontology

import "fmt"

// IRIMapperAgent generates IRIs and optional doc URLs.
type IRIMapperAgent struct{
    Prefix *PrefixAgent
}

func NewIRIMapperAgent(p *PrefixAgent) *IRIMapperAgent {
    return &IRIMapperAgent{Prefix: p}
}

func (m *IRIMapperAgent) ClassIRI(name string) string {
    return m.Prefix.IRI("bsky", name)
}

func (m *IRIMapperAgent) FieldIRI(class, field string) string {
    return fmt.Sprintf("%s/%s", m.ClassIRI(class), field)
}
