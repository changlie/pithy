package pkg

import (
    "bufio"
    "bytes"
    "changlie/pithy/common"
    "fmt"
    "io/ioutil"
    "os"
    "regexp"
    "strings"
    "text/template"
)

const (
    PkgName = "internel_pithy_gen"
    InitHandlerFile = "init_handler.go"
)

func init() {
    pkgPath := "src/"+PkgName
    if _, err := os.Stat(pkgPath); os.IsNotExist(err) {
        os.MkdirAll(pkgPath, 0666)
    }

    importPkgInMain()
}

func importPkgInMain() {
    srcDir, _ := ioutil.ReadDir("src")
    for _, f := range srcDir {
        if !f.IsDir() {
            continue
        }
        subDirPath := "src/"+f.Name()
        subDir, _ := ioutil.ReadDir(subDirPath)
        for _, f1 := range subDir {
            fname := f1.Name()
            if f1.IsDir() || !strings.HasSuffix(fname, ".go") {
                continue
            }
            fpath := fmt.Sprintf("%v/%v", subDirPath, fname)
            if isNotMainPkg(fpath) {
                break
            }
            insertImportPkgExpr2main(fpath)
        }
    }
}

func insertImportPkgExpr2main(fpath string) {
    if existPkgImport(fpath) {
        return
    }

    f, _ := os.Open(fpath)
    scanner := bufio.NewScanner(f)
    var buf bytes.Buffer
    var overPkg bool
    var insert bool
    //var preLine string
    for scanner.Scan() {
        line := scanner.Text()
        rawLine := strings.TrimSpace(line)
        buf.WriteString(line)
        buf.WriteString("\r\n")
        if strings.HasPrefix(rawLine, "package") {
            overPkg = true
        }else if !insert && overPkg {
            insert = true
            buf.WriteString(fmt.Sprintf(`import _ "%v"%v`, PkgName, "\r\n"))
        }
        //preLine = line
    }
    ioutil.WriteFile(fpath, buf.Bytes(), 066)
}

func existPkgImport(fpath string) bool {
    keyword := fmt.Sprintf(`_ "%v"`, PkgName)
    f, _ := os.Open(fpath)
    scanner := bufio.NewScanner(f)
    for scanner.Scan() {
        line := scanner.Text()
        if strings.Contains(line, keyword) {
            return true
        }
    }
    return false
}

func isNotMainPkg(fpath string) bool {
    f, _ := os.Open(fpath)
    scanner := bufio.NewScanner(f)
    for scanner.Scan() {
        line := strings.TrimSpace(scanner.Text())
        if strings.HasPrefix(line, "package") {
            if strings.Contains(line, "main"){
                return false
            }else{
                return true
            }
        }
    }
    return true
}

const initHandlerFileTempl = `package {{ .PkgName }}

import pt "{{ .PithyPkg }}"
import s "{{ .ServicePkg }}"

func init() {
{{range .HandlerInfos }}    pt.SetHandler("{{ .Url }}", s.{{ .MethodName }})
{{end}}
{{range .ServiceInfos -}}
{{$varname := .VarName}}
    var {{ $varname }} *s.{{ .ServiceName }}
{{range .Methods }}
{{- if not .IsFunc }}    pt.SetHandler("{{ .Url }}", {{ $varname }}.{{ .MethodName }})
{{ else }}    pt.SetHandler("{{ .Url }}", s.{{ .MethodName }})
{{end}}
{{- end}}
{{end}}
}
`

type InitHandlerInfo struct {
    PkgName string
    PithyPkg string
    ServicePkg string
    ServiceInfos []ServiceInfo
    HandlerInfos []HandlerInfo
}

type ServiceInfo struct {
    ServiceName string
    Methods []HandlerInfo
}

func (s *ServiceInfo) VarName() string {
    return strings.ToLower(s.ServiceName)
}

type HandlerInfo struct {
    IsFunc bool
    HttpMethod string
    Url string
    MethodName string
}

// generate go source file that initialize Handlers Information
func GenerateInitHandlerFile() {
    f, err := os.Create(fmt.Sprintf("src/%v/%v", PkgName, InitHandlerFile))
    common.ExitWhenError(err)

    sp := common.GetConfig("service.pkg", "service")

    data := InitHandlerInfo{
        PkgName:      PkgName,
        PithyPkg:     "changlie/pithy",
        ServicePkg:   sp,
        ServiceInfos: nil,
        HandlerInfos: nil,
    }

    scanForHandlerInfos(&data, sp)
    //fmt.Println("InitHandlerInfo:", data)

    t := template.Must(template.New("GenerateInitHandlerFile").Parse(initHandlerFileTempl))

    err = t.Execute(f, data)
    common.ExitWhenError(err)
}

func scanForHandlerInfos(data *InitHandlerInfo, servicePkg string) {
    var hs []HandlerInfo
    var ss []ServiceInfo
    path := "src/"+servicePkg
    fs, err := ioutil.ReadDir(path)
    common.ExitWhenError(err)
    for _, info := range fs {
        fname := info.Name()
        if info.IsDir() || !strings.HasSuffix(fname, ".go") {
            continue
        }
        fpath := fmt.Sprintf("%v/%v", path, fname)
        subHs, serviceName := getFileInfo(fpath)
        if serviceName != "" {
            sinfo := ServiceInfo{
                ServiceName: serviceName,
                Methods:     subHs,
            }
            ss = append(ss, sinfo)
        } else {
            hs = append(hs, subHs...)
        }
    }
    data.HandlerInfos = hs
    data.ServiceInfos = ss
}

func getFileInfo(fpath string) ([]HandlerInfo, string) {
    var hInfos []HandlerInfo
    var serviceName string
    f, _ := os.Open(fpath)
    scanner := bufio.NewScanner(f)
    for scanner.Scan() {
        var methodName string
        line := strings.TrimSpace(scanner.Text())
        if !strings.HasPrefix(line, "func") {
            continue
        }
        serviceName, methodName = getfuncInfo(line)
        if isPrivateFunc(methodName) {
            continue
        }

        isfunc := serviceName == ""
        fmt.Println("gen#getFileInfo:", serviceName, methodName, isfunc)
        url := getUrl(methodName)
        h := HandlerInfo{
            IsFunc:     isfunc,
            HttpMethod: "",
            Url:        url,
            MethodName: methodName,
        }
        hInfos = append(hInfos, h)
    }
    return hInfos, serviceName
}

func isPrivateFunc(funcName string) bool {
    b := funcName[0]
    if b >= 'A' && b <= 'Z' {
        return false
    }
    return true
}

func getUrl(s string) string{
    var paths []string
    var buf bytes.Buffer
    for i, b := range []byte(s) {
        if (b >= 'A') && (b <= 'Z') {
            if i>0 {
                paths = append(paths, buf.String())
                buf.Reset()
            }
            buf.Write(bytes.ToLower([]byte{b}))
        }else {
            buf.WriteByte(b)
        }
    }
    if buf.Len() > 0 {
        paths = append(paths, buf.String())
    }
    url := "/" + strings.Join(paths, "/")
    return url
}

var re = regexp.MustCompile(`func\s+(\(\s*[^\)\s]+\s+\*?([^\)\s]+)\s*\)\s+)?(\w+)`)
func getfuncInfo(line string) (serviceName string, methodName string) {
    arr := re.FindAllStringSubmatch(line, -1)[0]
    for i, item := range arr {
        if i == 2 {
            serviceName = item
        }else if i == 3 {
            methodName = item
        }
    }
    return
}


func GenerateInitStaticResourcesFile() {

}


