## ✅ `AGENTS.md` — Project: ATProto Ontology Reflection System

This project constructs a semantically rich OWL ontology from the ATProto codebase, using source code, inline documentation (e.g., godoc), and lexicon definitions. The ontology will serve both as a machine-readable API model and as a basis for automated reasoning, alignment, and semantic services.

---

### 🧠 Core Goals

* Parse Lexicon JSON and Go struct definitions
* Extract semantic descriptions from Go comments (godoc format)
* Emit OWL ontologies using meaningful IRIs, labels, and descriptions
* Infer correct datatypes (XSD) and structural relationships (OWL object/data properties)
* Emit  PROV metadata
* Reconcile IRIs with resolvable links when possible (e.g., Go pkg links or godoc URLs)

---

### 📦 AGENT ROSTER

| Agent Name        | Responsibility                                                                              |
| ----------------- | ------------------------------------------------------------------------------------------- |
| `MainAgent`       | Entry in `main.go` to route control via flags (e.g., `--emit-ontology`, `--doc-reflection`) |
| `LexSchemaAgent`  | Parses Lexicon JSON files and emits `owl:Class`/`owl:Property` triples                      |
| `GoReflectAgent`  | Parses Go structs and field comments; aligns with lexicon types using naming conventions    |
| `GodocAgent`      | Extracts structured documentation (e.g., via AST or comment parsing)                        |
| `IRIMapperAgent`  | Creates canonical IRIs, handles prefixing, and attempts resolvable documentation links      |
| `RdfWriterAgent`  | Emits TTL/RDF/XML from intermediate model using OWL-compliant formatting                    |
| `TypeInferAgent`  | Maps Go and Lex types to XSD/OWL types, supports optionality and multiplicity               |
| `PrefixAgent`     | Registers and emits consistent namespace prefixes (e.g., `bsky`, `actor`, `xsd`, `owl`)     |
| `ProvenanceAgent` | Attaches provenance metadata (e.g., extracted from comments or repo commits)                |
| `ValidationAgent` | Ensures output TTL complies with OWL 2 DL and passes RDF validators (e.g., Protégé, riot)   |
| `DocServeAgent`   | Optionally generates human-facing HTML documentation from the ontology                      |

---

### 🗂️ Input Corpus

* `indigo/api/bsky/*.go`: Go source code with structs and field comments
* `indigo/lex/*.json`: ATProto Lexicon JSON schemas
* `references/owl`: W3C OWL and XSD specs
* `references/rdf`: RDF and RDFS specs
* `references/prov`: W3C PROV specifications
* `references/atproto`: Static mirror of the ATProto site (optional linking target)

---

### 🧾 IRI and Prefix Conventions

* Ontology base: `https://atproto.social/ontology/`
* Default IRI pattern:
  `:VerificationState` → `https://atproto.social/ontology/app.bsky.actor.defs#VerificationState`
  `:verifiedStatus` → `...#VerificationState/verifiedStatus`
* Use godoc URL when resolvable:
  `https://pkg.go.dev/github.com/bluesky-social/indigo/api/bsky#ActorDefs_VerificationState`

---

### 🧠 Ontological Semantics

| Code Construct     | Ontology Construct                             |
| ------------------ | ---------------------------------------------- |
| Go struct          | `owl:Class`                                    |
| Go field           | `owl:ObjectProperty` or `owl:DatatypeProperty` |
| JSON type `ref`    | `owl:ObjectProperty`                           |
| JSON type `string` | `xsd:string`                                   |
| Field doc comment  | `rdfs:comment`                                 |
| Field requiredness | OWL cardinality restrictions                   |
| Naming path        | Class local name, IRI suffix                   |

---

### 🧪 Deliverables

* `build/lexicon.ttl`: complete ontology with full lexical and structural semantics
* Optional SHACL or JSON-LD context files
* `prefixes.ttl`: prefix map (auto-generated or curated)
* `iri-map.tsv`: crosswalk of Go doc links ↔ ontology IRIs
* `doc/index.html`: human-readable OWL doc (optional)

---

