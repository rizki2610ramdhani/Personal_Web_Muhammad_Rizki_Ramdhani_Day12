package routes

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"personal-web/connection"
	"personal-web/utilities"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

// meta data declaration
type MetaData struct {
	Title     string
	IsLogin   bool
	Username  string
	FlashData string
	Id        int
}

// data declaration
var Data = MetaData{
	Title: "Personal Web",
}

// project declaration
type Project struct {
	Id           int
	ProjectName  string
	StartDate    time.Time //Star Date Type Time
	EndDate      time.Time //End Date Type Time
	StrSD        string    //Star Date Type String
	StrED        string    //End Date Type String
	Duration     string
	Description  string
	Technologies []string
	Image        string
	UserId       int
	IsLogin      bool
}

// function handling index.html
func Home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")

	// parsing template html
	var tmpl, err = template.ParseFiles("views/index.html")
	// error handling
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")

	if session.Values["IsLogin"] != true {
		Data.IsLogin = false
	} else {
		Data.IsLogin = session.Values["IsLogin"].(bool)
		Data.Username = session.Values["Name"].(string)
		Data.Id = session.Values["Id"].(int)
	}

	fm := session.Flashes("Message")

	var flashes []string

	if len(fm) > 0 {
		session.Save(r, w)

		for _, fl := range fm {
			flashes = append(flashes, fl.(string))
		}
	}
	Data.FlashData = strings.Join(flashes, "")

	//query select projects from tb project (get data)
	rows, _ := connection.Conn.Query(context.Background(), "SELECT tb_projects.id, tb_projects.name, tb_projects.start_date, tb_projects.end_date, tb_projects.description, tb_projects.technologies, tb_projects.image, tb_projects.user_id FROM tb_projects  JOIN tb_users ON tb_projects.user_id = tb_users.id ORDER BY user_id DESC")

	var result []Project
	for rows.Next() {
		var each = Project{}

		var err = rows.Scan(&each.Id, &each.ProjectName, &each.StartDate, &each.EndDate, &each.Description, &each.Technologies, &each.Image, &each.UserId)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		var store = sessions.NewCookieStore([]byte("SESSION_ID"))
		session, _ := store.Get(r, "SESSION_ID")

		if session.Values["IsLogin"] != true {
			each.IsLogin = false
		} else {
			each.IsLogin = session.Values["IsLogin"].(bool)
		}

		each.Duration = utilities.GetDuration(each.StartDate, each.EndDate)

		result = append(result, each)
	}

	resp := map[string]interface{}{
		"Data":     Data,
		"Projects": result,
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, resp)
}

// function delete projects\
func DeleteProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	_, err := connection.Conn.Exec(context.Background(), "DELETE FROM tb_projects WHERE id=$1", id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message : " + err.Error()))
		return
	}

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

// function handling contact me
func ContactMe(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")

	// parsing template html
	var tmpl, err = template.ParseFiles("views/contact-me.html")
	// error handling
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, Data)
}

// function handling add-project.html
func AddProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")

	// parsing template html
	var tmpl, err = template.ParseFiles("views/add-project.html")
	// error handling
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")

	if session.Values["IsLogin"] != true {
		Data.IsLogin = false
	} else {
		Data.IsLogin = session.Values["IsLogin"].(bool)
		Data.Username = session.Values["Name"].(string)
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, Data)
}

func StoreProject(w http.ResponseWriter, r *http.Request) {
	//check parse form error
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	//get value from form
	ProjectName := r.PostForm.Get("projectName")
	StartDate := r.PostForm.Get("startDate")
	EndDate := r.PostForm.Get("endDate")
	Description := r.PostForm.Get("description")
	Technologies := r.Form["technology"]
	dataContext := r.Context().Value("dataFile")
	Image := dataContext.(string)

	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")

	User_Id := session.Values["ID"].(int)

	_, err = connection.Conn.Exec(context.Background(), "INSERT INTO public.tb_projects (name, start_date, end_date, description, technologies, image, user_id) VALUES ($1, $2, $3, $4, $5, $6, $7)", ProjectName, StartDate, EndDate, Description, Technologies, Image, User_Id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

// function handling project-detail.html with query string
func ProjectDetail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	// parsing template html
	var tmpl, err = template.ParseFiles("views/project-detail.html")
	// error handling
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	ProjectDetail := Project{}
	err = connection.Conn.QueryRow(context.Background(), "SELECT * FROM public.tb_projects WHERE id=$1", id).Scan(&ProjectDetail.Id, &ProjectDetail.ProjectName, &ProjectDetail.StartDate, &ProjectDetail.EndDate, &ProjectDetail.Description, &ProjectDetail.Technologies, &ProjectDetail.Image, &ProjectDetail.UserId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message : " + err.Error()))
		return
	}

	ProjectDetail.StrSD = ProjectDetail.StartDate.Format("2006-01-02")
	ProjectDetail.StrED = ProjectDetail.EndDate.Format("2006-01-02")
	ProjectDetail.Duration = utilities.GetDuration(ProjectDetail.StartDate, ProjectDetail.EndDate)

	resp := map[string]interface{}{
		"Data":          Data,
		"DetailProject": ProjectDetail,
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, resp)
}

// route to edit project
func FormEditProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	//parsing template
	templ, err := template.ParseFiles("views/edit-project.html")

	//error handling of parsing template
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	//create type struct edit
	type Edit struct {
		Id           int
		ProjectName  string
		StartDate    time.Time //Star Date Type Time
		EndDate      time.Time //End Date Type Time
		StrSD        string    //Star Date Type String
		StrED        string    //End Date Type String
		Description  string
		Technologies []string
		Image        string
		UserId       int
		IsUsePhp     bool
		IsUseLaravel bool
		IsUseJava    bool
		IsUseMysql   bool
	}

	//create object EditData with type struct edit
	EditData := Edit{}

	//get id from url
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	err = connection.Conn.QueryRow(context.Background(), "SELECT * FROM public.tb_projects WHERE id=$1", id).Scan(&EditData.Id, &EditData.ProjectName, &EditData.StartDate, &EditData.EndDate, &EditData.Description, &EditData.Technologies, &EditData.Image, &EditData.UserId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	EditData.StrSD = EditData.StartDate.Format("2006-01-02")
	EditData.StrED = EditData.EndDate.Format("2006-01-02")

	//get technologies detail
	for _, tech := range EditData.Technologies {
		if tech == "fa-brands fa-php" {
			EditData.IsUsePhp = true
		} else if tech == "fa-solid fa-database" {
			EditData.IsUseMysql = true
		} else if tech == "fa-brands fa-laravel" {
			EditData.IsUseLaravel = true
		} else if tech == "fa-brands fa-java" {
			EditData.IsUseJava = true
		}
	}

	//parsing data to template
	Resp := map[string]interface{}{
		"Data":    Data,
		"Project": EditData,
	}

	w.WriteHeader(http.StatusOK)
	templ.Execute(w, Resp)
}

func StoreEdit(w http.ResponseWriter, r *http.Request) {
	//check parse form error
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	//get value from form
	ProjectName := r.PostForm.Get("projectName")
	StartDate := r.PostForm.Get("startDate")
	EndDate := r.PostForm.Get("endDate")
	Description := r.PostForm.Get("description")
	Technologies := r.Form["technology"]
	dataContext := r.Context().Value("dataFile")
	Image := dataContext.(string)

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	_, err = connection.Conn.Exec(context.Background(), "UPDATE public.tb_projects SET name=$1, start_date=$2, end_date=$3, description=$4, technologies=$5, image=$6 WHERE id=$7;", ProjectName, StartDate, EndDate, Description, Technologies, Image, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}
