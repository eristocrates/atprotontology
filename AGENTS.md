
## ✅ `AGENTS.md` — Project: ATProto Semantic Mining Pipeline

This system semantically mines the ATProto source code, documentation, and lexicons to generate a unified OWL ontology. The pipeline follows an **Extract–Transform–Lift** (ETL) model:

* **Extract** code structure, comments, and type metadata
* **Transform** into intermediate semantic representations
* **Lift** into OWL individuals, classes, and properties with resolvable IRIs and proper prefixing

---

### 🧠 Core Philosophy

* Code *is* documentation: Go structs, JSON schemas, and godoc are all semantic input
* Output must be **semantically fused**, not layered or parallel
* IRI structure must be **shortened via prefixes** and **resolvable** when possible
* Knowledge extraction is **ongoing and extensible**, enabling later enrichment from other sources

---

### 📦 AGENT ROSTER (ETL-Based)

| Agent Name            | Responsibility                                                                 |
| --------------------- | ------------------------------------------------------------------------------ |
| `ExtractionAgent`     | Parses lexicon JSON and Go structs; extracts type structure and doc comments   |
| `SemanticIndexAgent`  | Builds a central semantic map of all known terms, fields, types, and comments  |
| `TypeMapperAgent`     | Infers OWL types and XSD datatypes for each field or property                  |
| `PrefixAgent`         | Generates CURIEs from long IRIs; manages consistent prefix registration        |
| `IRIAgent`            | Constructs resolvable or canonical IRIs for every term; aligns with godoc URLs |
| `EnrichmentAgent`     | Attaches `rdfs:label`, `rdfs:comment`, `skos:definition` from godoc and schema |
| `OntologyWriterAgent` | Emits fully prefixed TTL/RDF/XML files with clear ontology partitioning        |
| `ProvenanceAgent`     | (Optional) Adds provenance from file paths, repo commits, or comment origins   |
| `ValidationAgent`     | Validates output against OWL 2 DL and Turtle syntax using riot or Jena         |

---

### 🧱 Prefix and IRI Strategy

* Base ontology IRI: `https://atproto.social/ontology/`
* CURIEs like:

  * `actor:labelersPref/labelers`
  * `actor:VerificationState`
  * `bsky:ViewerState`
* Prefixes:

  * `xsd:` → `http://www.w3.org/2001/XMLSchema#`
  * `owl:` → `http://www.w3.org/2002/07/owl#`
  * `rdfs:` → `http://www.w3.org/2000/01/rdf-schema#`
  * `actor:` → `https://atproto.social/ontology/app.bsky.actor.defs#`

Prefix mappings emitted to `prefixes.ttl`.

---

### 🔁 Example Transformation

**Go Field + Lexicon:**

```go
// The user's status as a verified account.
VerifiedStatus *string `json:"verifiedStatus,omitempty"`
```

Yields:

```ttl
actor:VerificationState/verifiedStatus
    a             owl:DatatypeProperty ;
    rdfs:label    "verifiedStatus" ;
    rdfs:comment  "The user's status as a verified account." ;
    rdfs:domain   actor:VerificationState ;
    rdfs:range    xsd:string .
```
