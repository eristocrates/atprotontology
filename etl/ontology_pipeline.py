import os
import json
import re
from dataclasses import dataclass, field
from typing import Dict, List, Optional
from rdflib import Graph, Namespace, URIRef, BNode, Literal
from rdflib.namespace import RDF, RDFS, OWL, XSD

@dataclass
class FieldInfo:
    name: str
    typ: str
    comment: str = ""

@dataclass
class StructInfo:
    name: str
    schema_id: str
    def_name: str
    comment: str = ""
    fields: List[FieldInfo] = field(default_factory=list)

class ExtractionAgent:
    def __init__(self, repo_root: str):
        self.repo_root = repo_root
        self.structs: Dict[str, StructInfo] = {}
        self.lex_defs: Dict[str, dict] = {}

    def run(self):
        self.extract_go_structs()
        self.extract_lexicons()

    def extract_go_structs(self):
        go_dir = os.path.join(self.repo_root, "indigo")
        for root, _, files in os.walk(go_dir):
            for f in files:
                if f.endswith(".go"):
                    self._parse_go_file(os.path.join(root, f))

    def _parse_go_file(self, path: str):
        with open(path, "r") as fh:
            lines = fh.readlines()
        i = 0
        while i < len(lines):
            line = lines[i]
            m = re.match(r"type\s+(\w+)\s+struct\s*{", line)
            if m:
                name = m.group(1)
                # gather preceding comments
                j = i - 1
                comments = []
                while j >= 0 and lines[j].strip().startswith("//"):
                    comments.insert(0, lines[j].strip().lstrip("// "))
                    j -= 1
                comment = " ".join(comments)
                # attempt to parse schema id from comment
                schema_id = ""
                def_name = name
                m2 = re.search(r'"([^"]+)"\s+in\s+the\s+([\w.]+)\s+schema', comment)
                if m2:
                    def_name = m2.group(1)
                    schema_id = m2.group(2)
                struct = StructInfo(name=name, schema_id=schema_id, def_name=def_name, comment=comment)
                # parse fields
                i += 1
                while i < len(lines) and not lines[i].startswith("}"):
                    fl = lines[i].rstrip()
                    if fl.strip().startswith("//"):
                        # field comment
                        field_comment = fl.strip().lstrip("// ")
                        # next line should be field definition
                        i += 1
                        fl = lines[i].rstrip()
                        fm = re.match(r"(\w+)\s+([\w*\[\]]+)", fl.strip())
                        if fm:
                            field_name = fm.group(1)
                            field_type = fm.group(2)
                            struct.fields.append(FieldInfo(name=field_name, typ=field_type, comment=field_comment))
                    else:
                        fm = re.match(r"(\w+)\s+([\w*\[\]]+)", fl.strip())
                        if fm:
                            struct.fields.append(FieldInfo(name=fm.group(1), typ=fm.group(2)))
                    i += 1
                self.structs[name] = struct
            else:
                i += 1

    def extract_lexicons(self):
        lex_root = os.path.join(self.repo_root, "atproto/lexicons")
        for root, _, files in os.walk(lex_root):
            for f in files:
                if f.endswith(".json"):
                    p = os.path.join(root, f)
                    with open(p, "r") as fh:
                        try:
                            data = json.load(fh)
                        except Exception:
                            continue
                    schema_id = data.get("id")
                    if not schema_id:
                        continue
                    defs = data.get("defs", {})
                    for name, val in defs.items():
                        key = f"{schema_id}#{name}"
                        self.lex_defs[key] = val

class SemanticIndexAgent:
    def __init__(self, extractor: ExtractionAgent):
        self.extractor = extractor
        self.entities: Dict[str, StructInfo] = {}

    def run(self):
        for struct in self.extractor.structs.values():
            if struct.schema_id:
                key = f"{struct.schema_id}#{struct.def_name}"
                self.entities[key] = struct

class TypeMapperAgent:
    PRIMITIVES = {
        "string": XSD.string,
        "int": XSD.integer,
        "int64": XSD.integer,
        "bool": XSD.boolean,
        "float64": XSD.double,
    }

    def map_field(self, field: FieldInfo):
        ft = field.typ.strip("*[]")
        if ft in self.PRIMITIVES:
            return "datatype", self.PRIMITIVES[ft]
        return "object", None

class PrefixAgent:
    def __init__(self):
        self.prefixes = {
            "owl": OWL,
            "rdfs": RDFS,
            "xsd": XSD,
            "prov": Namespace("http://www.w3.org/ns/prov#"),
        }

    def add(self, prefix: str, iri: str):
        self.prefixes[prefix] = Namespace(iri)

    def bind(self, graph: Graph):
        for p, ns in self.prefixes.items():
            graph.bind(p, ns)

    def write_prefixes(self, path: str):
        with open(path, "w") as fh:
            for p, ns in self.prefixes.items():
                fh.write(f"@prefix {p}: <{ns}> .\n")

class IRIAgent:
    BASE = "https://atproto.social/ontology/"
    PKG = "https://pkg.go.dev/github.com/bluesky-social/indigo/api/bsky#"

    def __init__(self, prefix_agent: PrefixAgent):
        self.prefix_agent = prefix_agent
        self.iri_map: List[str] = []

    def class_iri(self, schema_id: str, name: str) -> URIRef:
        pref_name = schema_id.replace(".", "/")
        ns = f"{self.BASE}{schema_id}#"
        prefix = schema_id.split(".")[-2]
        pfx = prefix if prefix not in self.prefix_agent.prefixes else prefix + "_"
        self.prefix_agent.add(pfx, ns)
        return URIRef(f"{ns}{name}")

    def prop_iri(self, class_iri: URIRef, field: str) -> URIRef:
        return URIRef(f"{class_iri}/{field}")

    def add_mapping(self, class_iri: URIRef, struct_name: str):
        self.iri_map.append(f"{self.PKG}{struct_name}\t{class_iri}")

    def write_map(self, path: str):
        with open(path, "w") as fh:
            for row in self.iri_map:
                fh.write(row + "\n")

class EnrichmentAgent:
    def __init__(self, iri_agent: IRIAgent, type_mapper: TypeMapperAgent, prefix_agent: PrefixAgent):
        self.iri_agent = iri_agent
        self.tm = type_mapper
        self.prefix_agent = prefix_agent

    def enrich(self, struct: StructInfo, g: Graph):
        class_iri = self.iri_agent.class_iri(struct.schema_id, struct.def_name)
        g.add((class_iri, RDF.type, OWL.Class))
        if struct.comment:
            g.add((class_iri, RDFS.comment, Literal(struct.comment)))
        self.iri_agent.add_mapping(class_iri, struct.name)
        for field in struct.fields:
            piri = self.iri_agent.prop_iri(class_iri, field.name)
            ftype, rng = self.tm.map_field(field)
            if ftype == "datatype":
                g.add((piri, RDF.type, OWL.DatatypeProperty))
                g.add((piri, RDFS.range, rng))
            else:
                g.add((piri, RDF.type, OWL.ObjectProperty))
            g.add((piri, RDFS.domain, class_iri))
            if field.comment:
                g.add((piri, RDFS.comment, Literal(field.comment)))

class OntologyWriterAgent:
    def __init__(self, prefix_agent: PrefixAgent):
        self.prefix_agent = prefix_agent
        self.graph = Graph()

    def emit(self, path: str):
        self.prefix_agent.bind(self.graph)
        self.graph.serialize(destination=path, format="turtle")

class ValidationAgent:
    def validate(self, path: str):
        g = Graph()
        g.parse(path, format="turtle")
        return True


def main():
    repo = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
    build_dir = os.path.join(repo, "build")
    os.makedirs(build_dir, exist_ok=True)

    extractor = ExtractionAgent(repo)
    extractor.run()

    indexer = SemanticIndexAgent(extractor)
    indexer.run()

    prefix_agent = PrefixAgent()
    iri_agent = IRIAgent(prefix_agent)
    tm = TypeMapperAgent()
    enrich = EnrichmentAgent(iri_agent, tm, prefix_agent)
    writer = OntologyWriterAgent(prefix_agent)

    for ent in indexer.entities.values():
        enrich.enrich(ent, writer.graph)

    ttl_path = os.path.join(build_dir, "lexicon.ttl")
    writer.emit(ttl_path)

    prefix_agent.write_prefixes(os.path.join(build_dir, "prefixes.ttl"))
    iri_agent.write_map(os.path.join(build_dir, "iri-map.tsv"))

    ValidationAgent().validate(ttl_path)

if __name__ == "__main__":
    main()
