package main

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"text/template"
)

const (
	opensslConfTemplate = `[req]
distinguished_name = req_distinguished_name
req_extensions = v3_req
prompt = no

[v3_req]
keyUsage = digitalSignature,keyEncipherment
extendedKeyUsage = serverAuth
subjectAltName = @alt_names

[req_distinguished_name]
{{ range $value := .DN -}}
{{ $value }}
{{ end -}}

[alt_names]
{{ range $index, $value := .DNS -}}
DNS.{{ $index }} = {{ $value }}
{{ end -}}
{{ range $index, $value := .IP -}}
IP.{{ $index }} = {{ $value }}
{{ end -}}`
)

func main() {
	DN := []string{"CN=jenting", "OU=io"}
	DNS := []string{"jenting.io"}
	IP := []string{"8.8.8.8", "8.8.4.4"}

	template, err := template.New("").Parse(opensslConfTemplate)
	if err != nil {
		fmt.Printf("Err: %v\n", err)
		return
	}

	var rendered bytes.Buffer
	err = template.Execute(&rendered, struct {
		DN  []string
		DNS []string
		IP  []string
	}{
		DN:  DN,
		DNS: DNS,
		IP:  IP,
	})

	cmd1 := exec.Command("echo", rendered.String())
	cmd2 := exec.Command("tee", "openssl.conf")

	r, w := io.Pipe()
	cmd1.Stdout = w
	cmd2.Stdin = r

	if err := cmd1.Start(); err != nil {
		fmt.Printf("Err: %v\n", err)
		return
	}
	if err := cmd2.Start(); err != nil {
		fmt.Printf("Err: %v\n", err)
		return
	}

	if err := cmd1.Wait(); err != nil {
		fmt.Printf("Err: %v\n", err)
		return
	}
	if err := w.Close(); err != nil {
		fmt.Printf("Err: %v\n", err)
		return
	}
	if err := cmd2.Wait(); err != nil {
		fmt.Printf("Err: %v\n", err)
		return
	}
}
