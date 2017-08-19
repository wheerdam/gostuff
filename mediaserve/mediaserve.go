package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"io"
	"time"
	"io/ioutil"
	"os/exec"
	"html/template"
	"net/http"
	"path/filepath"
	"bbi/netutil"
	"github.com/icza/session"
)

var path string
var users *netutil.Users

type MessagePage struct {
	Header 	string
	Message interface{}
}

type ViewPage struct {
	Header	string
	Up		string
	Options	interface{}
	DirInfo	string
	MPre	interface{}
	MPost	interface{}
	Dirs	[]interface{}
	Medias	[]interface{}
	Others	[]interface{}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
        t, _ := template.ParseFiles("templates/login.gtpl")
        t.Execute(w, nil)
    } else {
        r.ParseForm()
        // logic part of log in
		userName := r.FormValue("username")
		err := users.Login(userName, r.FormValue("password"))
		if err != nil {
			t, _ := template.ParseFiles("templates/message.gtpl")
			msg := MessagePage{
				Header: "Login Failed",
				Message: template.HTML(
					"<p>Incorrect name and/or password was provided</p>" +
					"<p><a href=\"./login\">Retry</a></p>",
				),
			}
			t.Execute(w, msg)
			s := session.Get(r)
			if s != nil {
				// logout of existing session if login failed
				session.Remove(s, w)
			}
		} else {
			s := session.NewSessionOptions(&session.SessOptions{
					CAttrs: map[string]interface{}{"UserName": userName},
			})
			session.Add(s, w)
			http.Redirect(w, r, "./view?path=.", 301)
		}
    }
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	s := session.Get(r)
	if s == nil {
		t, _ := template.ParseFiles("templates/message.gtpl")
		msg := MessagePage{
				Header:	"Logout Failed",
				Message: template.HTML(
						"<p>You were not logged in</p>" +
						"<p><a href=\"./login\">Login</a></p>",
				),
		}
		t.Execute(w, msg)
	} else {
		session.Remove(s, w)
		t, _ := template.ParseFiles("templates/message.gtpl")
		msg := MessagePage{
				Header: "Logged Out",
				Message: template.HTML(
						"<p><a href=\"./login\">Login</a></p>",
				),
		}
		t.Execute(w, msg)
	}
}

func RootHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "./view?path=.", 301)
}

func ViewHandler(w http.ResponseWriter, r *http.Request) {
	s := session.Get(r)
	if s == nil {
		t, _ := template.ParseFiles("templates/message.gtpl")
		msg := MessagePage{
				Header: "Not Logged In",
				Message: template.HTML(
						"<p>You were not logged in</p>" +
						"<p><a href=\"./login\">Login</a></p>",
				),
		}
		t.Execute(w, msg)
		return
	}
	userReqPath := r.URL.Query().Get("path")
	scaling := r.URL.Query().Get("scaling")
	showVid := r.URL.Query().Get("showvid")
	height := r.URL.Query().Get("height")
	if scaling == "" {
		scaling = "FillHorizontal"
	}
	if userReqPath == "" {
		t, _ := template.ParseFiles("templates/message.gtpl")
		msg := MessagePage{
				Header: "No Path Specified",
				Message: template.HTML(
					"<p>You must specify path to resource with 'path='</p>",
				),
		}
		t.Execute(w, msg)
		return
	}
	if showVid == "" {
		showVid = "No"
	}
	if height == "" {
		height = "25"
	}
	reqPath := path + "/" + userReqPath
	f, err := os.Stat(reqPath)
	if os.IsNotExist(err) {
		t, _ := template.ParseFiles("templates/message.gtpl")
		msg := MessagePage{
				Header: "Path not Found",
				Message: template.HTML(
					"<p>Can not find '" + reqPath + "'</p>",
				),
		}
		t.Execute(w, msg)
		return
	} else if err != nil {
		t, _ := template.ParseFiles("templates/message.gtpl")
		msg := MessagePage{
				Header: "File status check failed",
				Message: template.HTML(
					"<p>Failed to check status for  '" + reqPath + "'</p>",
				),
		}
		t.Execute(w, msg)
		return
	}
	if strings.HasSuffix(strings.ToLower(reqPath), ".jpg") ||
	   strings.HasSuffix(strings.ToLower(reqPath), ".jpeg") {
		img, err := os.Open(reqPath)
		if err != nil {
			fmt.Println("'" + reqPath + "' failed to open: " + err.Error())
			return // no response
		}
		defer img.Close()
		w.Header().Set("Content-Type", "image/jpeg")
		io.Copy(w, img)
	} else if strings.HasSuffix(strings.ToLower(reqPath), ".png") {
		img, err := os.Open(reqPath)
		if err != nil {
		fmt.Println("'" + reqPath + "' failed to open: " + err.Error())
			return // no response
		}
		defer img.Close()
		w.Header().Set("Content-Type", "image/png")
		io.Copy(w, img)
	} else if strings.HasSuffix(strings.ToLower(reqPath), ".gif") {
		img, err := os.Open(reqPath)
		if err != nil {
		fmt.Println("'" + reqPath + "' failed to open: " + err.Error())
			return // no response
		}
		defer img.Close()
		w.Header().Set("Content-Type", "image/gif")
		io.Copy(w, img)
	} else if strings.HasSuffix(strings.ToLower(reqPath), ".webm") || 
			strings.HasSuffix(strings.ToLower(reqPath), ".mkv") || 
			strings.HasSuffix(strings.ToLower(reqPath), ".mp4") {
		video, err := os.Open(reqPath)
		if err != nil {
			fmt.Println("'" + reqPath + "' failed to open: " + err.Error())
			return 
		}
		defer video.Close()
		http.ServeContent(w, r, reqPath, time.Now(), video)
	} else if strings.HasSuffix(strings.ToLower(reqPath), ".txt") {
		img, err := os.Open(reqPath)
		if err != nil {
			fmt.Println("'" + reqPath + "' failed to open: " + err.Error())
			return // no response
		}
		defer img.Close()
		w.Header().Set("Content-Type", "text/plain")
		io.Copy(w, img)
	} else if strings.HasSuffix(strings.ToLower(reqPath), ".html") {
		img, err := os.Open(reqPath)
		if err != nil {
			fmt.Println("'" + reqPath + "' failed to open: " + err.Error())
			return // no response
		}
		defer img.Close()
		w.Header().Set("Content-Type", "text/html")
		io.Copy(w, img)
	} else if f.IsDir() {
		files, _ := ioutil.ReadDir(reqPath)
		cur := "./view?path=" + userReqPath
		curVid := "showvid=" + showVid
		curScaling := "scaling=" + scaling
		curHeight := "height=" + height
		options := "<span style=\"margin-right: 10px\">"
		if showVid == "Yes" {
			options = options + "<a href=\"" + cur + "&showvid=No&" +
				curScaling + "&" + curHeight + "\">-Vid</a> " 
		} else {
			options = options + "<a href=\"" + cur + "&showvid=Yes&" +
				curScaling + "&" + curHeight + "\">+Vid</a> "
		}
		options = options + "</span><span style=\"margin-right: 10px\">"
		if scaling != "FillHorizontal" {
			options = options + "<a href=\"" + cur + "&scaling=FillHorizontal&" +
				curVid + "&" + curHeight + "\">H</a> "
		} else {
			options = options + "H "
		}
		options = options + "</span><span style=\"margin-right: 10px\">"
		if scaling != "FillVertical" {
			options = options + "<a href=\"" + cur + "&scaling=FillVertical&" +
				curVid + "&" + curHeight + "\">V</a> "
		} else {
			options = options + "V "
		}
		options = options + "</span><span style=\"margin-right: 5px\">"
		if scaling != "Thumbnail" {
			options = options + "<a href=\"" + cur + "&scaling=Thumbnail&" +
				curVid + "&" + curHeight + "\">T</a>:"
		} else {
			options = options + "T:"
		}
		options = options + "</span><span style=\"margin-right: 5px\">"
		if height != "25" {
			options = options + "<a href=\"" + cur + "&" + curScaling +
				"&" + curVid + "&height=25\">&#188;</a>"
		} else {
			options = options + "&#188;"
		}
		options = options + "</span><span style=\"margin-right: 10px\">"
		if height != "50" {
			options = options + "<a href=\"" + cur + "&" + curScaling +
				"&" + curVid + "&height=50\">&#189;</a> "
		} else {
			options = options + "&#189; "
		}
		options = options + "</span><span style=\"margin-right: 5px\">"
		if scaling != "List" {
			options = options + "<a href=\"" + cur + "&scaling=List&" +
				curVid + "&" + curHeight + "\">L</a> "
		} else {
			options = options + "L "
		}
		options = options + "</span>"
		options = options + "</span><span style=\"margin-right: 10px\">"
		if scaling != "ListPreview" {
			options = options + "<a href=\"" + cur + "&scaling=ListPreview&" +
				curVid + "&" + curHeight + "\">P</a> "
		} else {
			options = options + "P "
		}
		options = options + "</span><span style=\"margin-right: 10px\">"
		options = options + "<a href=\"thumbgen?path=" +
			userReqPath + "&done=" + userReqPath +
			"&" + curVid + "&" + curScaling + "&" + curHeight +
			"\">" +
			"TG </a>"
		options = options + "</span>"
		page := ViewPage{
			Header:	 userReqPath,
			Up:		 "./view?path=" + filepath.Dir("./" + userReqPath) +
						"&" + curVid + "&" + curScaling + "&" + curHeight,
			Options: template.HTML(options),
			MPre:	 "",
			MPost:	 "",
			Dirs:	 make([]interface{}, 0),
			Medias:	 make([]interface{}, 0),
			Others:	 make([]interface{}, 0)}
		fileCount := 0
		unknownCount := 0
		for _, file := range files {
			if file.IsDir() {
				page.Dirs = append(page.Dirs, 
					template.HTML("<p>&#128193; <a href=\"./view?" +
						"path=" + userReqPath + "/" + file.Name() +
						"&" + curVid + 
						"&" + curScaling +
						"&" + curHeight +
						"\">" +
						file.Name() + "</a></p>"))
			} else if scaling == "List" && isListable(file.Name()) {
				page.Others = append(page.Others, 
					template.HTML("<p><a href=\"" +
						cur + "/" + file.Name() + "\">" +					
						file.Name() + "</a></p>"))
				fileCount++
			} else if scaling == "ListPreview" && isListable(file.Name()) {
				page.MPre = template.HTML("<table>")
				page.MPost = template.HTML("</table>")
				if isImage(file.Name()) {
					prefix := "<tr><td style=\"vertical-align: middle;\"><a href=\"" + cur + "/" + 
						file.Name() + "\">"
					suffix := "</a></td><td><a href=\"" + cur + "/" + 
						file.Name() + "\"><span style=\"vertical-align: middle;\">" + file.Name() + "</span></a></td></tr>"
					imgAttr := "height=\"100px\" style=\"vertical-align: middle;\""
					page.Medias = append(page.Medias,
						template.HTML(prefix + "<img src=\"./view?path=" +
							userReqPath + "/" + file.Name() + "\" " +
							imgAttr + "> " + suffix))
				} else if isVideo(file.Name()) {
					prefix := "<tr><td style=\"vertical-align: middle;\"><a href=\"" + cur + "/" + 
						file.Name() + "\">"
					suffix := "</a></td><td><a href=\"" + cur + "/" + 
						file.Name() + "\"><span style=\"vertical-align: middle;\">" + file.Name() + "</span></a></td></tr>"
					imgAttr := "height=\"100px\" style=\"vertical-align: middle;\""
					page.Medias = append(page.Medias,
						template.HTML(prefix + "<img src=\"./view?path=" +
							userReqPath + "/" + file.Name() + ".thumb.jpg\" " +
							imgAttr + "> " + suffix))
				} else {
					page.Others = append(page.Others, 
						template.HTML("<p><a href=\"" +
							cur + "/" + file.Name() + "\">" +					
							file.Name() + "</a></p>"))
				fileCount++
				}
				fileCount++
			} else {
				imgAttr := ""
				switch scaling {
				case "FillHorizontal":
					imgAttr = "width=\"100%\""
				case "FillVertical":
					imgAttr = "height=\"100%\""
				case "Thumbnail":
					imgAttr = "height=\"" + height + "%\""
				}
				if isImage(file.Name()) && 
						!strings.HasSuffix(file.Name(), ".thumb.jpg") {
					prefix := ""
					suffix := ""
					if scaling == "Thumbnail" {
						prefix = "<a href=\"" + cur + "/" + 
							file.Name() + "\">"
						suffix = "</a>"
					}
					page.Medias = append(page.Medias,
						template.HTML(prefix + "<img src=\"./view?path=" +
							userReqPath + "/" + file.Name() + "\" " +
							imgAttr + ">" + suffix + " "))
					fileCount++
				} else if isVideo(file.Name()) {
					prefix := ""
					suffix := ""
					video := ""
					if showVid != "Yes" {
						prefix = "<a href=\"" + cur + "/" + 
							file.Name() + "\">"
						video = "<img " + imgAttr + 
							" src=\"" + cur + "/" +
							file.Name() + ".thumb.jpg\">"
						suffix = "</a>"
					} else {
						video = "<video " + imgAttr + 
							" src=\"" + cur + "/" +
							file.Name() + "\" autoplay loop muted></video>"
					}
					page.Medias = append(page.Medias,
						template.HTML(prefix + video + suffix))
					fileCount++
				} else if isListable(file.Name()) {
					page.Others = append(page.Others, 
						template.HTML("<p><a href=\"" +
							cur + "/" + file.Name() + "\">" +					
							file.Name() + "</a></p>"))
					fileCount++
				} else if !strings.HasSuffix(file.Name(), ".thumb.jpg") {
					unknownCount++
				}
			}
		}
		page.DirInfo = fmt.Sprintf("%s: %d%s%d%s%d%s", userReqPath,
			len(page.Dirs), " dirs, ", fileCount, " media files and ",
			unknownCount, " unknown files")
		t, _ := template.ParseFiles("templates/view.gtpl")
		t.Execute(w, page)
	} else {
		t, _ := template.ParseFiles("templates/message.gtpl")
		msg := MessagePage{
				Header: "Unable to Handle File",
				Message: template.HTML(
					"<p>Unable to handle '" + reqPath + "'</p>",
				),
		}
		t.Execute(w, msg)
		return
	}
}

func ThumbnailGenerator(w http.ResponseWriter, r *http.Request) {
	s := session.Get(r)
	if s == nil {
		t, _ := template.ParseFiles("templates/message.gtpl")
		msg := MessagePage{
				Header: "Not Logged In",
				Message: template.HTML(
						"<p>You were not logged in</p>" +
						"<p><a href=\"./login\">Login</a></p>",
				),
		}
		t.Execute(w, msg)
		return
	}
	userReqPath := r.URL.Query().Get("path")
	done := r.URL.Query().Get("done")
	curVid := "showvid=" + r.URL.Query().Get("showvid")
	curScaling := "scaling=" + r.URL.Query().Get("scaling")
	curHeight := "height=" + r.URL.Query().Get("height") 
	reqPath := path + "/" + userReqPath
	f, err := os.Stat(reqPath)
	if os.IsNotExist(err) {
		t, _ := template.ParseFiles("templates/message.gtpl")
		msg := MessagePage{
				Header: "Path not Found",
				Message: template.HTML(
					"<p>Can not find '" + reqPath + "'</p>",
				),
		}
		t.Execute(w, msg)
		return
	} else if err != nil {
		t, _ := template.ParseFiles("templates/message.gtpl")
		msg := MessagePage{
				Header: "File status check failed",
				Message: template.HTML(
					"<p>Failed to check status for  '" + reqPath + "'</p>",
				),
		}
		t.Execute(w, msg)
		return
	}
	if !f.IsDir() {
		t, _ := template.ParseFiles("templates/message.gtpl")
		msg := MessagePage{
				Header: "Path is not a directory",
				Message: template.HTML(
					"<p>'" + reqPath + "' needs to be a directory</p>",
				),
		}
		t.Execute(w, msg)
		return
	}
	files, _ := ioutil.ReadDir(reqPath)
	output := ""
	cmd := "ffmpeg"
	for _, file := range files {
		fileName := reqPath + "/" + file.Name()
		if isVideo(fileName) {
			args := []string{"-y", "-ss", "00:10:00", "-i",
				fileName, "-vframes", "1", fileName + ".thumb.jpg"}
			out, err := exec.Command(cmd, args...).CombinedOutput()
			output = output + fmt.Sprintf("%s", out)
			if err != nil {
				output = output + err.Error()
			}
			output = output + "\n"
			_, err = os.Stat(fileName + ".thumb.jpg")
			if os.IsNotExist(err) {
				args := []string{"-y", "-ss", "00:00:01", "-i",
				fileName, "-vframes", "1", fileName + ".thumb.jpg"}
				out, err := exec.Command(cmd, args...).CombinedOutput()
				output = output + fmt.Sprintf("%s", out)
				if err != nil {
					output = output + err.Error()
				}
				output = output + "\n"
			}
		}
	}
	if done == "" {
		t, _ := template.ParseFiles("templates/message.gtpl")
		msg := MessagePage{
				Header: "Thumbnail Generation Output",
				Message: template.HTML(
					"<div align=\"left\"><pre>" + output + "</pre></div>",
				),
		}
		t.Execute(w, msg)
	} else {
		http.Redirect(w, r, "./view?path=" + done +
			"&" + curVid + "&" + curScaling + "&" + curHeight, 301)
	}
}

func isImage(name string) (bool) {
	lName := strings.ToLower(name)
	return strings.HasSuffix(lName, ".png") ||
		   strings.HasSuffix(lName, ".gif") ||
		   strings.HasSuffix(lName, ".jpg") ||
		   strings.HasSuffix(lName, ".jpeg")
}

func isVideo(name string) (bool) {
	lName := strings.ToLower(name)
	return strings.HasSuffix(lName, ".mp4") ||
		   strings.HasSuffix(lName, ".mkv") ||
		   strings.HasSuffix(lName, ".webm")
}

func isText(name string) (bool) {
	lName := strings.ToLower(name)
	return strings.HasSuffix(lName, ".txt") ||
		   strings.HasSuffix(lName, ".text") ||
		   strings.HasSuffix(lName, ".html")
}

func isListable(name string) (bool) {
	lName := strings.ToLower(name)
	return (isImage(lName) || isVideo(lName) || isText(lName)) &&
		!strings.HasSuffix(lName, ".thumb.jpg")
}

func usage() {
	fmt.Println("usage: mediaserve [path] [users-file] [static-path] (cert) (key)\n")
}

func main() {
	if len(os.Args) != 4 && len(os.Args) != 6 {
		usage()
		return
	}
	path = os.Args[1]
	users = netutil.NewUsers()
	err := users.LoadFromFile(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}
	session.Global.Close()
	session.Global = session.NewCookieManagerOptions(session.NewInMemStore(), &session.CookieMngrOptions{AllowHTTP: true})
	fs := http.FileServer(http.Dir(os.Args[3]))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/view", ViewHandler)
	http.HandleFunc("/thumbgen", ThumbnailGenerator)
	http.HandleFunc("/login", LoginHandler)
	http.HandleFunc("/logout", LogoutHandler)
	http.HandleFunc("/", RootHandler)
	if len(os.Args) == 6 {
		err = http.ListenAndServeTLS(":18311", os.Args[4], os.Args[5], nil)
	} else {		
		err = http.ListenAndServe(":18310", nil)
	}
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}