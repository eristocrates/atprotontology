package ontology

import (
	"bytes"
	"fmt"
	"os"
)

// Triple represents a simple RDF triple.
type Triple struct {
	S string
	P string
	O string
}

// RdfWriterAgent writes triples to Turtle.
type RdfWriterAgent struct {
	Prefix *PrefixAgent
	Buffer bytes.Buffer
}

func NewRdfWriterAgent(p *PrefixAgent) *RdfWriterAgent {
	return &RdfWriterAgent{Prefix: p}
}

func (r *RdfWriterAgent) WriteTriple(t Triple) {
	r.Buffer.WriteString(fmt.Sprintf("%s %s %s .\n", t.S, t.P, t.O))
}

func (r *RdfWriterAgent) Save(path string) error {
	var out bytes.Buffer
	for pref, iri := range r.Prefix.Prefixes {
		out.WriteString(fmt.Sprintf("@prefix %s: <%s> .\n", pref, iri))
	}
	out.Write(r.Buffer.Bytes())
	return os.WriteFile(path, out.Bytes(), 0644)
}
