package main

import (
	"html/template"
	"log"
	"os"

	"gitee.com/deep-spark/go-ixdcgm/pkg/ixdcgm"
)

const diagOutput = `Software:
  {{range $t := .Software}}
  {{printf "%-50s" $t.TestName}} {{$t.Status}}	{{$t.TestOutput}}
  {{- end}}
{{range $g := .PerGpu}}
  
GPU: {{$g.GPU}}
  {{range $t := $g.DiagResults}}
  {{printf "%-20s" $t.TestName}} {{$t.Status}}	{{$t.TestOutput}}
  {{- end}}
{{- end}}
`

func main() {
	cleanup, err := ixdcgm.Init(ixdcgm.Embedded)
	if err != nil {
		log.Panicln(err)
	}
	defer cleanup()

	// Choose ixdcgm diag type and input group handle
	// - ixdcgm.DiagQuick    -> run a very basic health check on the system
	// - ixdcgm.DiagMedium   -> run a medium-length diagnostic (a few minutes)
	// - ixdcgm.DiagLong     -> run a extensive diagnostic (several minutes)
	// - ixdcgm.DiagExtended -> run a very extensive diagnostic (many minutes)
	// Tip: to run diag within a time limit, please use ixdcgm.RunDiagWithTimeout
	result, err := ixdcgm.RunDiag(ixdcgm.DiagQuick, ixdcgm.GroupAllGPUs())
	if err != nil {
		log.Panicln(err)
	}

	t := template.Must(template.New("Diag").Parse(diagOutput))
	if err = t.Execute(os.Stdout, result); err != nil {
		log.Panicln("Template error:", err)
	}
}
