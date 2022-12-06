# 0. Template Introduction

With `template`, Go can render the View layer for web.

template can implement the field mapping, custom function and nesting template.   

For the detail information of tempalte can refer to [template](https://github.com/astaxie/build-web-application-with-golang/blob/master/zh/07.4.md).


# 1. helm template
helm use the same template to render and customize the function of helm template.

Take an example of `helm lint`, let's see how helm render the template and how to customize the function.

with `helm lint` command, the [cobra](https://github.com/spf13/cobra) has been used as the entry of helm:
```
// helm/cmd/helm/lint.go
func newLintCmd(out io.Writer) *cobra.Command {
    ...
    result := client.Run([]string{path}, vals)
}

// helm/pkg/action/lint.go
func (l *Lint) Run(paths []string, vals map[string]interface{}) *LintResult {
    linter, err := lintChart(path, vals, l.Namespace, l.Strict)
}

// helm/pkg/action/lint.go
func lintChart(path string, vals map[string]interface{}, namespace string, strict bool) (support.Linter, error) {
    ...
    return lint.All(chartPath, vals, namespace, strict), nil
}

// helm/pkg/lint/lint.go
func All(basedir string, values map[string]interface{}, namespace string, strict bool) support.Linter {
    ...
    rules.Templates(&linter, values, namespace, strict)
}

// helm/pkg/lint/rules/template.go
func Templates(linter *support.Linter, values map[string]interface{}, namespace string, strict bool) {
    ...
    var e engine.Engine
    renderedContentMap, err := e.Render(chart, valuesToRender)
}
```

Then going to the engine part which is the real executer to render the template:
```
// helm/pkg/engine/engine.go
func (e Engine) Render(chrt *chart.Chart, values chartutil.Values) (map[string]string, error) {
    tmap := allTemplates(chrt, values)
    return e.render(tmap)
}

func (e Engine) render(tpls map[string]renderable) (map[string]string, error) {
    return e.renderWithReferences(tpls, tpls)
}

func (e Engine) renderWithReferences(tpls, referenceTpls map[string]renderable) (rendered map[string]string, err error) {
    t := template.New("gotpl")
    e.initFunMap(t, referenceTpls)
}

func (e Engine) initFunMap(t *template.Template, referenceTpls map[string]renderable) {
    funcMap := funcMap()

    // If we are not linting and have a cluster connection, provide a Kubernetes-backed
    // implementation.
    if !e.LintMode && e.config != nil {
        funcMap["lookup"] = NewLookupFunction(e.config)
    }

    t.Funcs(funcMap)
}

// helm/pkg/funcs.go
func funcMap() template.FuncMap {
    // Add some extra functionality
    extra := template.FuncMap{
        "toToml":        toTOML,
        "toYaml":        toYAML,
        ...
    }
}
```

we can see it almost the same implementation for the demo as the template example:
```
t := template.New("test.tpl")
t = t.Funcs(template.FuncMap{"tToLower": TranslateToLower})
t, _ = t.ParseFiles("test.tpl")
```

In the `funcs.go` the customize function has been defined for helm template, and the template function added in template with `t.Funcs`, after that the template can be used to render the template of helm.

Here to defined a function to customize a function for example, defined the function `toEmptyString` to return the empty string like:
```
func toEmptyString(s string) string {
    return ""
}

func funcMap() template.FuncMap {
    extra := template.FuncMap{
        "toEmptyString": toEmptyString,
        ...
    }
}
```

Use the function `toEmptyString` in helm template, and then `helm lint` to verify whether the function are effective:
```
# pwd
/root/go/src/helm/bin
# ./helm lint test
WARNING: Kubernetes configuration file is group-readable. This is insecure. Location: /root/.kube/config
WARNING: Kubernetes configuration file is world-readable. This is insecure. Location: /root/.kube/config
==> Linting test
[INFO] Chart.yaml: icon is recommended

1 chart(s) linted, 0 chart(s) failed
```

As expected, the customize function are effective and can be render by helm.
