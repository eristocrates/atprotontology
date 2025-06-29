package ontology

import (
	"fmt"
	"path/filepath"
	"strings"
)

// OntologyPipeline wires all agents together.
type OntologyPipeline struct {
	Prefix     *PrefixAgent
	TypeInfer  *TypeInferAgent
	IRIMapper  *IRIMapperAgent
	LexSchema  *LexSchemaAgent
	GoReflect  *GoReflectAgent
	Godoc      *GodocAgent
	RdfWriter  *RdfWriterAgent
	Validator  *ValidationAgent
	Provenance *ProvenanceAgent
}

func NewOntologyPipeline() *OntologyPipeline {
	prefix := NewPrefixAgent()
	return &OntologyPipeline{
		Prefix:     prefix,
		TypeInfer:  NewTypeInferAgent(),
		IRIMapper:  NewIRIMapperAgent(prefix),
		LexSchema:  NewLexSchemaAgent(),
		GoReflect:  NewGoReflectAgent(),
		Godoc:      NewGodocAgent(),
		RdfWriter:  NewRdfWriterAgent(prefix),
		Validator:  NewValidationAgent(),
		Provenance: NewProvenanceAgent("atproto"),
	}
}

// Run executes the ontology extraction pipeline.
func (o *OntologyPipeline) Run(srcDir string, outPath string) error {
	// Actual lexicon JSON files live under the atproto/lexicons directory.
	lexDir := filepath.Join(srcDir, "atproto", "lexicons")
	goDir := filepath.Join(srcDir, "indigo", "api", "bsky")

	if err := o.LexSchema.Load(lexDir); err != nil {
		return err
	}
	if err := o.GoReflect.Load(goDir); err != nil {
		return err
	}
	if err := o.Godoc.Load(goDir); err != nil {
		return err
	}

	for _, def := range o.LexSchema.Defs {
		classIRI := fmt.Sprintf("<%s>", o.IRIMapper.ClassIRI(def.LexiconID, def.Name))
		o.RdfWriter.WriteTriple(Triple{classIRI, "a", "owl:Class"})
		desc := def.Description
		sname := structName(def.LexiconID, def.Name)
		if d, ok := o.Godoc.Docs[sname]; ok {
			desc = strings.TrimSpace(d)
		}
		if desc != "" {
			o.RdfWriter.WriteTriple(Triple{classIRI, "rdfs:comment", fmt.Sprintf("\"%s\"", desc)})
		}
		for propName, prop := range def.Properties {
			fieldIRI := fmt.Sprintf("<%s>", o.IRIMapper.FieldIRI(def.LexiconID, def.Name, propName))
			propType := "owl:DatatypeProperty"
			rng := o.TypeInfer.InferDatatype(prop.Type)
			if prop.Type == "ref" || (prop.Type == "array" && prop.Items != nil && prop.Items.Type == "ref") {
				propType = "owl:ObjectProperty"
				refName := strings.TrimPrefix(prop.Ref, "#")
				if prop.Items != nil && prop.Items.Ref != "" {
					refName = strings.TrimPrefix(prop.Items.Ref, "#")
				}
				rng = fmt.Sprintf("<%s>", o.IRIMapper.ClassIRI(def.LexiconID, refName))
			}
			o.RdfWriter.WriteTriple(Triple{fieldIRI, "a", propType})
			o.RdfWriter.WriteTriple(Triple{fieldIRI, "rdfs:domain", classIRI})
			o.RdfWriter.WriteTriple(Triple{fieldIRI, "rdfs:range", rng})
			cmt := prop.Description
			if d, ok := o.Godoc.Docs[sname+"."+toPascal(propName)]; ok {
				cmt = strings.TrimSpace(d)
			}
			if cmt != "" {
				o.RdfWriter.WriteTriple(Triple{fieldIRI, "rdfs:comment", fmt.Sprintf("\"%s\"", cmt)})
			}
		}
	}
	if err := o.Validator.Validate(o.RdfWriter.Buffer.String()); err != nil {
		return err
	}
	return o.RdfWriter.Save(outPath)
}

func structName(lexID, defName string) string {
	segs := strings.Split(lexID, ".")
	base := ""
	if len(segs) >= 2 {
		base = toPascal(segs[len(segs)-2]) + toPascal(segs[len(segs)-1])
	} else if len(segs) == 1 {
		base = toPascal(segs[0])
	}
	return base + "_" + toPascal(defName)
}
