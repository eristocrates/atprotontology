package ontology

import (
    "fmt"
    "path/filepath"
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
        Prefix:    prefix,
        TypeInfer: NewTypeInferAgent(),
        IRIMapper: NewIRIMapperAgent(prefix),
        LexSchema: NewLexSchemaAgent(),
        GoReflect: NewGoReflectAgent(),
        Godoc:     NewGodocAgent(),
        RdfWriter: NewRdfWriterAgent(prefix),
        Validator: NewValidationAgent(),
        Provenance: NewProvenanceAgent("atproto") ,
    }
}

// Run executes the ontology extraction pipeline.
func (o *OntologyPipeline) Run(srcDir string, outPath string) error {
    lexDir := filepath.Join(srcDir, "indigo", "lex")
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

    for name, def := range o.LexSchema.Defs {
        classIRI := fmt.Sprintf("<%s>", o.IRIMapper.ClassIRI(name))
        o.RdfWriter.WriteTriple(Triple{classIRI, "a", "owl:Class"})
        if def.Description != "" {
            o.RdfWriter.WriteTriple(Triple{classIRI, "rdfs:comment", fmt.Sprintf("\"%s\"", def.Description)})
        }
        for propName, prop := range def.Properties {
            fieldIRI := fmt.Sprintf("<%s>", o.IRIMapper.FieldIRI(name, propName))
            o.RdfWriter.WriteTriple(Triple{fieldIRI, "a", "owl:DatatypeProperty"})
            o.RdfWriter.WriteTriple(Triple{fieldIRI, "rdfs:domain", classIRI})
            dt := o.TypeInfer.InferDatatype(prop.Type)
            o.RdfWriter.WriteTriple(Triple{fieldIRI, "rdfs:range", dt})
            if prop.Description != "" {
                o.RdfWriter.WriteTriple(Triple{fieldIRI, "rdfs:comment", fmt.Sprintf("\"%s\"", prop.Description)})
            }
        }
    }
    if err := o.Validator.Validate(o.RdfWriter.Buffer.String()); err != nil {
        return err
    }
    return o.RdfWriter.Save(outPath)
}
