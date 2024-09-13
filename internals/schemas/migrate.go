package schemas

import (
	"fmt"
	"io"
	"os"
	"strings"

	_ "ariga.io/atlas-go-sdk/recordriver"
	"ariga.io/atlas-provider-gorm/gormschema"
)

func LoadModels() {
	sb := &strings.Builder{}

	stmts, err := gormschema.New("postgres").Load(migratableApps...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load gorm schema: %v\n", err)
		os.Exit(1)
	}
	sb.WriteString(stmts)
	sb.WriteString(";\n")
	io.WriteString(os.Stdout, stmts)
	// fmt.Println(stmts)
}
