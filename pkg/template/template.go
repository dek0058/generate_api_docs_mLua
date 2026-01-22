package template

const DocumentTemplate = `
<table style="width:100%;border-collapse:collapse; border-color:#ccc;border-spacing:0;border-style:solid;border-width:1px">
    <thead>
        <tr>
            <th style="background-color:#f0f0f0;border-color:#ccc;border-style:solid;border-width:0px;color:#333;overflow:hidden;padding:10px 5px;text-align:left;vertical-align:top;word-break:normal">
                <span style="color:#3167ad">{{.ReturnType}}</span> <span style="font-weight:bold">{{.FunctionName}}</span>({{.FunctionParamsStr}}){{.BadgeHTML}}
            </th>
        </tr>
    </thead>
    <tbody>
        <tr>
            <td style="background-color:#fff;border-color:#ccc;border-style:solid;border-width:0px;color:#333;overflow:hidden;padding:10px 5px;text-align:left;vertical-align:top;word-break:normal">
                {{.Description}}
            </td>
        </tr>
        {{range .Params}}
        <tr>
            <td style="background-color:#fafafa; border-top: 1px solid #eee; padding: 10px 5px 10px 15px;">
                <code style="background-color: #e1e4e8; padding: 2px 5px; border-radius: 4px;">{{.Name}}</code>
                <span style="color: #57606a;"> &nbsp;|&nbsp; {{.Desc}}</span>
            </td>
        </tr>
        {{end}}
    </tbody>
</table>
<br>
`
