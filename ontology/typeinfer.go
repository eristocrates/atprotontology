package ontology

import "strings"

// TypeInferAgent infers XSD datatypes and OWL property types.
type TypeInferAgent struct{}

func NewTypeInferAgent() *TypeInferAgent { return &TypeInferAgent{} }

// InferDatatype returns an XSD type based on input strings.
func (t *TypeInferAgent) InferDatatype(goType string) string {
	switch strings.TrimPrefix(goType, "[]") {
	case "string":
		return "xsd:string"
	case "boolean", "bool":
		return "xsd:boolean"
	case "int", "int64", "uint", "uint64":
		return "xsd:integer"
	case "float32", "float64":
		return "xsd:decimal"
	case "bytes", "blob":
		return "xsd:base64Binary"
	case "record", "object", "array", "ref":
		return "owl:Thing"
	default:
		return "owl:Thing"
	}
}
