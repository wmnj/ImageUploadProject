package main

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"runtime/debug"
)
const (
	UPLOAD_DIR = "./uploads"
	TEMPLATE_DIR = "./views"
	ListDir = 0x0001
)
// var templates map[string]*template.Template
var templates = make(map[string]*template.Template)
func check(err error) {
	if err != nil {
		panic(err)
	}
}
// 在服务启动的时候初始化模板 防止多次渲染比较慢
func init() {
	fileInfoArr, err := ioutil.ReadDir(TEMPLATE_DIR)
	check(err)
	var templateName, templatePath string
	for _, fileInfo := range fileInfoArr {
		templateName = fileInfo.Name()
		if ext := path.Ext(templateName); ext != ".html" {
			continue
		}
		tmpl := templateName[:len(templateName)-5]
		// fmt.Println("tmpl:",templateName,tmpl)
		templatePath = TEMPLATE_DIR + "/" + templateName
		log.Println("Loading template:", templatePath)
		t := template.Must(template.ParseFiles(templatePath))
		templates[tmpl] = t
	}
}
// 防止函数出错的时候就闪退 ,在业务逻辑函数中上包装一层函数，如果发生错误的时候执行匿名函数捕捉错误
 func safeHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if e, ok := recover().(error);ok {
				http.Error(w, e.Error(), http.StatusInternalServerError)
				log.Println("WARN: panic in %v - %v",fn,e)
				log.Println(string(debug.Stack()))
			}
		}()
		fn(w,r)
	}
 }

func renderHtml(w http.ResponseWriter, tmpl string, locals map[string]interface{}) {
	err := templates[tmpl].Execute(w, locals)
	check(err)
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w,"Hello, world!")
}
// 上传图片 
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" { // 如果是get方法则返回上穿图片的html表单，
		// 使用模版渲染后的写法
		renderHtml(w, "upload", nil)
	    // 最原始的直接输出html代码的写法
		// io.WriteString(w, "<form method=\"POST\" action=\"/upload\" "+ 
		// " enctype=\"multipart/form-data\">"+
		// "Choose an image to upload: <input name=\"image\" type=\"file\" />" + "<input type=\"submit\" value=\"Upload\" />"+
		// "</form>")
		// return
	}
	if r.Method == "POST" { // 如果是post则接受文件内容，新建临时文件进行保存
		f,h,err := r.FormFile("image")
		check(err)
		filename := h.Filename
		defer f.Close()
		fmt.Println("上传文件: ",filename)
		// t, err := ioutil.TempFile(UPLOAD_DIR, filename)
		t, err := os.Create(UPLOAD_DIR + "/" + filename)
		check(err)
		defer t.Close()
		_, err1 := io.Copy(t,f)
		check(err1)
		// 跳转到查看图片界面
		http.Redirect(w,r, "/view?id=" + filename, http.StatusFound)
	}
}
// 检测文件是否存在
func isExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return os.IsExist(err)
}
// 查看图片 根据请求中访问的图片id找到文件的位置并传回客户端
func viewHandler(w http.ResponseWriter, r *http.Request) {
	imageId := r.FormValue("id")
	imagePath := UPLOAD_DIR + "/" + imageId
	if exists := isExists(imagePath);!exists {
		http.NotFound(w,r)
		return
	}
	w.Header().Set("Content-Type","image")
	// 将服务端的一个文件内容读写到http.res并返回给请求来源的客户端
	http.ServeFile(w, r, imagePath)
} 
// 列出所有的图片 从uploads目录中读取所有的图片 并将列表展示的代码发给客户端
func listHandler(w http.ResponseWriter, r *http.Request) {
	fileInfoArr, err := ioutil.ReadDir("./uploads")
	check(err)
	locals := make(map[string]interface{})
	images := []string{}
	// var listHtml string
	for _, fileInfo := range fileInfoArr {
		images = append(images, fileInfo.Name())
		// imgid := fileInfo.Name
		// // listHtml += "<li><a href=\"/view?id=" + imgid + "\">imgid</a></li>"
		// listHtml += "<li><a href=\"/view?id="+imgid+"\">imgid</a></li>"
	}
	fmt.Println("文件列表：",images)
	locals["images"] = images
	renderHtml(w,"list",locals)
	// io.WriteString(w,"<ol>" + listHtml + "</ol>")
}
// 访问静态资源 
func staticDirHandler(mux *http.ServeMux, prefix string, staticDir string, flags int) {
	mux.HandleFunc(prefix,func(w http.ResponseWriter, r *http.Request){
		file := staticDir + r.URL.Path[len(prefix)-1:]
		if (flags & ListDir) == 0 {
			if exists := isExists(file); !exists {
				http.NotFound(w,r)
				return
			}
		}
		http.ServeFile(w,r,file)
	})
}
func main() {
	mux := http.NewServeMux()
	staticDirHandler(mux,"/assets/","./public",0)
	http.HandleFunc("/hello",safeHandler(helloHandler))
	http.HandleFunc("/upload",safeHandler(uploadHandler))
	http.HandleFunc("/view",safeHandler(viewHandler))
	http.HandleFunc("/",safeHandler(listHandler))
	err := http.ListenAndServe(":8080",nil)
	if err != nil {
		log.Fatal("ListenAndServe: ",err.Error())
	}
}