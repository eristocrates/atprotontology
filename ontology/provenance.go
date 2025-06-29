package ontology

// ProvenanceAgent attaches minimal provenance information.
type ProvenanceAgent struct {
    Info string
}

func NewProvenanceAgent(info string) *ProvenanceAgent {
    return &ProvenanceAgent{Info: info}
}
