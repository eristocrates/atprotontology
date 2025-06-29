package ontology

import (
    "go/ast"
    "go/parser"
    "go/token"
)

// GodocAgent extracts documentation from Go source.
type GodocAgent struct {
    Docs map[string]string
}

func NewGodocAgent() *GodocAgent { return &GodocAgent{Docs: map[string]string{}} }

func (g *GodocAgent) Load(dir string) error {
    fset := token.NewFileSet()
    pkgs, err := parser.ParseDir(fset, dir, nil, parser.ParseComments)
    if err != nil { return err }
    for _, pkg := range pkgs {
        for _, file := range pkg.Files {
            for _, decl := range file.Decls {
                gd, ok := decl.(*ast.GenDecl)
                if !ok || gd.Tok != token.TYPE { continue }
                for _, spec := range gd.Specs {
                    ts := spec.(*ast.TypeSpec)
                    if ts.Doc != nil {
                        g.Docs[ts.Name.Name] = ts.Doc.Text()
                    }
                    if st, ok := ts.Type.(*ast.StructType); ok {
                        for _, f := range st.Fields.List {
                            if f.Doc != nil && len(f.Names) > 0 {
                                key := ts.Name.Name + "." + f.Names[0].Name
                                g.Docs[key] = f.Doc.Text()
                            }
                        }
                    }
                }
            }
        }
    }
    return nil
}
