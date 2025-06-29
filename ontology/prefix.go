package ontology

// PrefixAgent maintains ontology prefixes.
type PrefixAgent struct {
    Prefixes map[string]string
}

// NewPrefixAgent returns a PrefixAgent with default prefixes.
func NewPrefixAgent() *PrefixAgent {
    return &PrefixAgent{Prefixes: map[string]string{
        "owl":  "http://www.w3.org/2002/07/owl#",
        "rdfs": "http://www.w3.org/2000/01/rdf-schema#",
        "xsd":  "http://www.w3.org/2001/XMLSchema#",
        "bsky": "https://atproto.social/ontology/",
    }}
}

// IRI returns the expanded IRI for a given prefix and local name.
func (p *PrefixAgent) IRI(prefix, local string) string {
    base, ok := p.Prefixes[prefix]
    if !ok {
        return prefix + ":" + local
    }
    return base + local
}
