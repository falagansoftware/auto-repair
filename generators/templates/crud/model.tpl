type {{.Model.Name}} struct {
    {{range  .Model.Properties}}
	{{.Name}} {{.Type}}
    {{end}}
}
