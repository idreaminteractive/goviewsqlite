// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.747
package components

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

import (
	"fmt"
	"github.com/idreaminteractive/goviewsqlite/internal/common"
)

type BasePageData struct {
	Title   string
	Content templ.Component
}

// CSRF needs to be auto injected here
func BasePageComponent(props BasePageData) templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<!doctype html><html data-theme=\"dark\"><head>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if props.Title == "" {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<title>Ignite Home</title>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		} else {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<title>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var2 string
			templ_7745c5c3_Var2, templ_7745c5c3_Err = templ.JoinStringErrs(props.Title)
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `internal/views/components/basepage.templ`, Line: 21, Col: 24}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var2))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</title>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<meta charset=\"UTF-8\"><meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\"><link rel=\"stylesheet\" href=\"https://cdn.jsdelivr.net/npm/@picocss/pico@2/css/pico.classless.min.css\"><link rel=\"icon\" href=\"/static/favicon.ico\"><script src=\"https://unpkg.com/htmx.org@2.0.0\"></script><script src=\"https://unpkg.com/htmx.org/dist/ext/alpine-morph.js\"></script><script defer src=\"https://cdn.jsdelivr.net/npm/alpinejs@3.x.x/dist/cdn.min.js\"></script><script type=\"text/javascript\">\n\t\t\t\t// uncomment to see all the htmx logging\n\t\t\t\thtmx.logAll();\t\t\n\t\t\t</script><style>\n\t\t\t[x-cloak] { display: none !important; }\n\t\t\t</style></head><body hx-headers=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var3 string
		templ_7745c5c3_Var3, templ_7745c5c3_Err = templ.JoinStringErrs(fmt.Sprintf("{\"X-CSRF-Token\": \"%s\"}", common.GetCtxCSRF(ctx)))
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `internal/views/components/basepage.templ`, Line: 42, Col: 86}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var3))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\"><script type=\"text/javascript\">\n\t\t\t\t// capture unknown responses\n\t\t\t\tdocument.body.addEventListener('htmx:responseError', function (evt) {\n\t\t\t\t\t// handle non-200s for htmx\n\t\t\t\t\t// todo -  find a clean way to handle these items (rather than render internal server error in templ)\n\t\t\t\t\tswitch(evt.detail.xhr.status) {\n\t\t\t\t\t\tcase 401:\n\t\t\t\t\t\t\twindow.alert(\"Unauthorized\")\n\t\t\t\t\t\t\tbreak\n\t\t\t\t\t\tcase 403:\n\t\t\t\t\t\t\twindow.alert(\"Forbidden\")\n\t\t\t\t\t\t\tbreak\n\t\t\t\t\t\tcase 404:\n\t\t\t\t\t\t\twindow.alert(\"Not Found\")\n\t\t\t\t\t\t\tbreak\n\t\t\t\t\t\tcase 500:\t\n\t\t\t\t\t\t\twindow.alert(\"Internal server error\")\n\t\t\t\t\t\t\t// log out error in console.error here.\n\t\t\t\t\t\t\tbreak\n\t\t\t\t\t\tdefault:\n\t\t\t\t\t\t\tbreak\n\t\t\t\t\t}\n\t\t\t\t});\n\n\t\t\t</script><main>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = props.Content.Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</main></body></html>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}
