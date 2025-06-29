package ontology

import (
    "go/ast"
    "go/parser"
    "go/token"
)

// GoStruct represents a Go struct with fields.
type GoStruct struct {
    Name   string
    Fields []GoField
}

type GoField struct {
    Name string
    Type string
    Tag  string
}

// GoReflectAgent parses Go structs.
type GoReflectAgent struct {
    Structs map[string]GoStruct
}

func NewGoReflectAgent() *GoReflectAgent { return &GoReflectAgent{Structs: map[string]GoStruct{}} }

func (g *GoReflectAgent) Load(dir string) error {
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
                    st, ok := ts.Type.(*ast.StructType)
                    if !ok { continue }
                    gs := GoStruct{Name: ts.Name.Name}
                    for _, f := range st.Fields.List {
                        t := ""
                        if se, ok := f.Type.(*ast.Ident); ok {
                            t = se.Name
                        } else if se, ok := f.Type.(*ast.ArrayType); ok {
                            if id, ok := se.Elt.(*ast.Ident); ok {
                                t = "[]" + id.Name
                            }
                        }
                        tag := ""
                        if f.Tag != nil {
                            tag = f.Tag.Value
                        }
                        for _, name := range f.Names {
                            gs.Fields = append(gs.Fields, GoField{Name: name.Name, Type: t, Tag: tag})
                        }
                    }
                    g.Structs[gs.Name] = gs
                }
            }
        }
    }
    return nil
}
