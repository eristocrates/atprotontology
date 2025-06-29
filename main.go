package main

import (
    "flag"
    "fmt"
    "log"

    "atprotontology/ontology"
)

func main() {
    emit := flag.Bool("emit-ontology", false, "generate ontology")
    doc := flag.Bool("doc-reflection", false, "print doc reflection")
    flag.Parse()

    if *emit {
        pipe := ontology.NewOntologyPipeline()
        out := "build/lexicon.ttl"
        if err := pipe.Run(".", out); err != nil {
            log.Fatalf("pipeline error: %v", err)
        }
        fmt.Println("ontology written to", out)
    }

    if *doc {
        pipe := ontology.NewOntologyPipeline()
        if err := pipe.Godoc.Load("indigo/api/bsky"); err != nil {
            log.Fatal(err)
        }
        for k, v := range pipe.Godoc.Docs {
            fmt.Printf("%s: %s\n", k, v)
        }
    }
}
