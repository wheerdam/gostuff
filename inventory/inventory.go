package inventory

import (
	"fmt"
	"os"
	"strings"
	"strconv"
	"time"
	"errors"
	"net/http"
	"html/template"
	"database/sql"
	"bbi/netutil"	
	"github.com/icza/session"
)

var	invTemplatePath string
var invPrefix string
var invUsers *netutil.Users
var invDB *sql.DB

type ViewPageFields struct {
	Prefix			string
	UserName 		string
	ViewTitle		string
	Types			[]string
	Manufacturers	[]string
	ViewOps			interface{}
	Data 			[]Item
}

type ItemPageFields struct {
	Prefix		string
	UserName	string
	Info		*Item
	InvEntries	[]InventoryEntry
}

type AddEditPageFields struct {
	Prefix		string
	UserName	string
	Header		string
	Info		*Item
	InvEntries	[]InventoryEntry
	Footer		interface{}
}

type MessagePage struct {
	Prefix		string
	Header 		string
	Message	 	interface{}
}

type SearchPage struct {
	Prefix			string
	UserName		string
}

type BrowsePageFields struct {
	Prefix			string
	UserName		string
	Types			[]Type
	Manufacturers	[]string
}

type Type struct {
	Name			string
	Subtypes		[]string
	Manufacturers	[]string
}

func LoginPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
        t, _ := template.ParseFiles(invTemplatePath + "/login.gtpl")
        t.Execute(w, SearchPage{invPrefix, ""})
    } else {
        r.ParseForm()
        // logic part of log in
		userName := r.FormValue("username")
		err := invUsers.Login(userName, r.FormValue("password"))
		if err != nil {
			t, _ := template.ParseFiles(invTemplatePath + "/message.gtpl")
			msg := MessagePage{
					Prefix: invPrefix,
					Header: "Login Failed",
					Message: template.HTML(
						"<p>Incorrect name and/or password was provided</p>" +
						"<p><a href=\"" + invPrefix + "/login\">Retry</a></p>",
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
					CAttrs: map[string]interface{}{
						"UserName": userName,
						"Inventory": "Yes",
					},
			})
			session.Add(s, w)
			http.Redirect(w, r, invPrefix + "/browse", 301)
		}
    }
}

func LogoutPage(w http.ResponseWriter, r *http.Request) {
	if !checkSession(w, r) {
		t, _ := template.ParseFiles(invTemplatePath + "/message.gtpl")
		msg := MessagePage{
				Prefix: invPrefix,
				Header:	"Logout Failed",
				Message: template.HTML(
						"<p>You were not logged in</p>" +
						"<p><a href=\"" + invPrefix + "/login\">Login</a></p>",
				),
		}
		t.Execute(w, msg)
	} else {
		s := session.Get(r)
		session.Remove(s, w)
		t, _ := template.ParseFiles(invTemplatePath + "/message.gtpl")
		msg := MessagePage{
				Prefix: invPrefix,
				Header: "Logged Out",
				Message: template.HTML(
						"<p><a href=\"" + invPrefix + "/login\">Login</a></p>",
				),
		}
		t.Execute(w, msg)
	}
}

func ListingPage(w http.ResponseWriter, r *http.Request) {
	if !checkSession(w, r) {
		return
	}
	s := session.Get(r)
	t, _ := template.ParseFiles(invTemplatePath + "/list.gtpl")
	userName := s.CAttr("UserName")
	viewOps := ""
	logicOr := false
	var items []Item		
	if len(r.URL.Query()) > 0 {
		conditions := make([]Condition, 0)
		for k, v := range r.URL.Query() {
			for i := range v {
				if k == "op" && v[i] == "or" {
					logicOr = true
				}
			}
		}
		for k, v := range r.URL.Query() {
			for i, cur := range v {
				if k == "op" && v[i] == "or" {
				} else if k == "type" ||
						  k == "subtype" || 
						  k == "manufacturer" || 
						  k == "value" {
					condition := Condition{key: k + "=", value: v[i]}
					conditions = append(conditions, condition)
					viewOps = viewOps + " <span style=\"color:#ababab\">" + k + "=</span>'" +
						v[i] + "'"
					if logicOr {
						viewOps = viewOps + " <span style=\"color:#2222ff\">or</span>"
					} else {
						viewOps = viewOps + " <span style=\"color:#2222ff\">and</span>"
					}
				} else if k == "search-type" {
					if cur == "" {
						continue
					}
					condition := Condition{
						key: "type like ", value: "%" + v[i] + "%"}
					conditions = append(conditions, condition)
					viewOps = viewOps + " <span style=\"color:#00ff00\">contains(type='" + v[i] + "')</span>"
					if logicOr {
						viewOps = viewOps + " <span style=\"color:#2222ff\">or</span>"
					} else {
						viewOps = viewOps + " <span style=\"color:#2222ff\">and</span>"
					}
				} else if k == "search-subtype" {
					if cur == "" {
						continue
					}
					condition := Condition{
						key: "subtype like ", value: "%" + v[i] + "%"}
					conditions = append(conditions, condition)
					viewOps = viewOps + " <span style=\"color:#00ff00\">contains(subtype='" + v[i] + "')</span>"
					if logicOr {
						viewOps = viewOps + " <span style=\"color:#2222ff\">or</span>"
					} else {
						viewOps = viewOps + " <span style=\"color:#2222ff\">and</span>"
					}
				} else if k == "search-manufacturer" {
					if cur == "" {
						continue
					}
					condition := Condition{
						key: "manufacturer like ", value: "%" + v[i] + "%"}
					conditions = append(conditions, condition)
					viewOps = viewOps + " <span style=\"color:#00ff00\">contains(manufacturer='" + v[i] + "')</span>"
					if logicOr {
						viewOps = viewOps + " <span style=\"color:#2222ff\">or</span>"
					} else {
						viewOps = viewOps + " <span style=\"color:#2222ff\">and</span>"
					}
				} else if k == "search-part-number" {
					if cur == "" {
						continue
					}
					condition := Condition{
						key: "model_number like ", value: "%" + v[i] + "%"}
					conditions = append(conditions, condition)
					viewOps = viewOps + " <span style=\"color:#00ff00\">contains(model_number='" + v[i] + "')</span>"
					if logicOr {
						viewOps = viewOps + " <span style=\"color:#2222ff\">or</span>"
					} else {
						viewOps = viewOps + " <span style=\"color:#2222ff\">and</span>"
					}
				} else if k == "search-description" {
					if cur == "" {
						continue
					}
					condition := Condition{
						key: "descriptive_name like ", value: "%" + v[i] + "%"}
					conditions = append(conditions, condition)
					viewOps = viewOps + " <span style=\"color:#00ff00\">contains(description='" + v[i] + "')</span>"
					if logicOr {
						viewOps = viewOps + " <span style=\"color:#2222ff\">or</span>"
					} else {
						viewOps = viewOps + " <span style=\"color:#2222ff\">and</span>"
					}
				} else if k == "search-phys-description" {
					if cur == "" {
						continue
					}
					condition := Condition{
						key: "phys_description like ", value: "%" + v[i] + "%"}
					conditions = append(conditions, condition)
					viewOps = viewOps + " <span style=\"color:#00ff00\">contains(phys_description='" + v[i] + "')</span>"
					if logicOr {
						viewOps = viewOps + " <span style=\"color:#2222ff\">or</span>"
					} else {
						viewOps = viewOps + " <span style=\"color:#2222ff\">and</span>"
					}
				}  else {
					viewOps = " <span style=\"color:#ff0000\">Error(" + k + "=" + v[i] + ")</span>"
				}
			}				
		}
		if logicOr {
			viewOps = strings.TrimRight(viewOps, " <span style=\"color:#2222ff\">or</span>")
		} else {
			viewOps = strings.TrimRight(viewOps, " <span style=\"color:#2222ff\">and</span>")
		}
		items = getItemsFiltered("type", logicOr, conditions...)
	} else {
		items = getItems("type")
	}
	types, _ := getDistinctCol("type")
	manufacturers, _ := getDistinctCol("manufacturer")
	viewData := ViewPageFields{
			Prefix:			invPrefix,
			UserName: 		userName.(string),
			ViewTitle:		"List Items",
			Data: 			items,
			Types:			types,
			Manufacturers:	manufacturers,
	}
	if viewOps != "" {
		viewData.ViewOps = template.HTML(viewOps)
	} else {
		viewData.ViewOps = "All Items"
	}
	
	t.Execute(w, viewData)
}

func ItemPage(w http.ResponseWriter, r *http.Request) {
	if !checkSession(w, r) {
		return
	}
	s := session.Get(r)
	var itemID string
	if r.Method == "GET" {
		itemID = r.URL.Query().Get("id")
	} else {
		itemID = r.FormValue("id")
	}
	if itemID == "" {
		t, _ := template.ParseFiles(invTemplatePath + "/message.gtpl")
		msg := MessagePage{
				Prefix: invPrefix,
				Header: "Missing Item ID",
				Message: template.HTML(
						"<p>ItemID was not provided in the URL</p>",
				),
		}
		t.Execute(w, msg)
	} else {
		t, _ := template.ParseFiles(invTemplatePath + "/item.gtpl")
		item, invEntries, err := getItem(itemID)
		if item == nil || err != nil {
			errT, _ := template.ParseFiles(invTemplatePath + "/message.gtpl")
			msg := MessagePage{
					Prefix: invPrefix,
					Header: "Failed to Process Item",
					Message: template.HTML(
							"<p>ItemID: " + itemID + "</p>" +
							"<p><a href=\"" + invPrefix + "/edit?id=" + itemID + 
							"\">Add this Item</a> - " + 
							"<a href=\"" + invPrefix + "/browse\">Browse</a> - <a href=\"" + invPrefix + "/list\">View All Items</a></p>",
					),
			}
			errT.Execute(w, msg)
			return
		}
		userName := s.CAttr("UserName")
		itemData := ItemPageFields{
				Prefix:		invPrefix,
				UserName:	userName.(string),
				Info:		item,
				InvEntries:	invEntries,
		}
		t.Execute(w, itemData)
	}
}

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	if !checkSession(w, r) {
		return
	}
	var itemID string
	if r.Method == "GET" {
		itemID = r.URL.Query().Get("id")
	} else {
		itemID = r.FormValue("id")
	}
	if itemID == "" {
		t, _ := template.ParseFiles(invTemplatePath + "/message.gtpl")
		msg := MessagePage{
				Prefix: invPrefix,
				Header: "Missing Item ID",
				Message: template.HTML(
						"<p>ItemID was not provided in the URL</p>",
				),
		}
		t.Execute(w, msg)
	} else {
		item, _, _ := getItem(itemID)
		err := deleteItem(itemID)
		if err != nil {
			errT, _ := template.ParseFiles(invTemplatePath + "/message.gtpl")
			msg := MessagePage{
					Prefix: invPrefix,
					Header: "Failed to Process Item",
					Message: template.HTML(
							"<p>Error: " + err.Error() + "</p>",
					),
			}
			errT.Execute(w, msg)
		} else {
			http.Redirect(w, r, invPrefix + "/list?type=" + item.Type, 301)
		}
	}
}

func AddEditItemPage(w http.ResponseWriter, r *http.Request) {
	if !checkSession(w, r) {
		return
	}
	s := session.Get(r)
	var itemID string
	if r.Method == "GET" {
		itemID = r.URL.Query().Get("id")
	} else {
		itemID = r.FormValue("id")
	}
	t, _ := template.ParseFiles(invTemplatePath + "/edit.gtpl")
	item, invEntries, err := getItem(itemID)
	if itemID == "" || err != nil {
		errT, _ := template.ParseFiles(invTemplatePath + "/message.gtpl")
		msg := MessagePage{
				Prefix: invPrefix,
				Header: "Failed to Process Item",
				Message: template.HTML(
						"<p>Item ID was not provided or was invalid!</p>",
				),
		}
		errT.Execute(w, msg)
		return
	}
	if item != nil {
		strFooter := " <button onclick=\"history.back()\">Cancel</button>"
		fields := AddEditPageFields {
			Prefix:		invPrefix,
			UserName:	s.CAttr("UserName").(string),
			Header:		"Edit Item #" + itemID,
			Info:		item,
			InvEntries:	invEntries,
			Footer:		template.HTML(strFooter),
		}
		t.Execute(w, fields)
	} else {
		strFooter := " <button onclick=\"window.location.href='" + invPrefix + "/browse'; return false\">Cancel</button>"
		intItemID, _ := strconv.Atoi(itemID)
		fields := AddEditPageFields {
			Prefix:		invPrefix,
			UserName:	s.CAttr("UserName").(string),
			Header:		"Add Item #" + itemID,
			Info:		&Item{ItemID: intItemID},
			InvEntries:	make([]InventoryEntry, 0),
			Footer:		template.HTML(strFooter),
		}
		t.Execute(w, fields)
	}
}

func BrowsePage(w http.ResponseWriter, r *http.Request) {
	if !checkSession(w, r) {
		return
	}
	s := session.Get(r)
	userName := s.CAttr("UserName").(string)
	types, err := getDistinctCol("type")	
	if err != nil {
		errT, _ := template.ParseFiles(invTemplatePath + "/message.gtpl")
		msg := MessagePage{
				Prefix: invPrefix,
				Header: "Failed to Get Item Types",
				Message: template.HTML(
						"<p>This should not happen, please contact the administrator</p>" +
						"<p>Error: " + err.Error() + "</p>",
				),
		}
		errT.Execute(w, msg)
		return
	}
	manufacturers, err := getDistinctCol("manufacturer")
	if err != nil {
		errT, _ := template.ParseFiles(invTemplatePath + "/message.gtpl")
		msg := MessagePage{
				Prefix: invPrefix,
				Header: "Failed to Get Item Types",
				Message: template.HTML(
						"<p>This should not happen, please contact the administrator</p>" +
						"<p>Error: " + err.Error() + "</p>",
				),
		}
		errT.Execute(w, msg)
		return
	}
	browsePageFields := BrowsePageFields{
		Prefix:			invPrefix,
		UserName:		userName,
		Types:			make([]Type, len(types)),
		Manufacturers:	manufacturers,
	}		
	for i := range types {
		condition := Condition{
			key: 	"type=",
			value:	types[i]}
		subtypes, err := getDistinctCol("subtype", condition)
		typeManufacturers, _ := getDistinctCol("manufacturer", condition)
		if err != nil {
			errT, _ := template.ParseFiles(invTemplatePath + "/message.gtpl")
			msg := MessagePage{
					Prefix: invPrefix,
					Header: "Failed to Get Subtypes",
					Message: template.HTML(
							"<p>This should not happen, please contact the administrator</p>" +
							"<p>Error: " + err.Error() + "</p>",
					),
			}
			errT.Execute(w, msg)
		}
		browsePageFields.Types[i].Name = types[i]
		browsePageFields.Types[i].Subtypes = subtypes
		browsePageFields.Types[i].Manufacturers = typeManufacturers
	}
	
	t, _ := template.ParseFiles(invTemplatePath + "/browse.gtpl")
	t.Execute(w, browsePageFields)
}

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	if !checkSession(w, r) {
		return
	}
	s := session.Get(r)
	if r.Method == "GET" {
		t, _ := template.ParseFiles(invTemplatePath + "/search.gtpl")
		userName := s.CAttr("UserName").(string)
		searchPageFields := SearchPage{invPrefix, userName}
		t.Execute(w, searchPageFields)
	} else if r.Method == "POST" {
		qType := r.FormValue("type")
		qSubtype := r.FormValue("subtype")
		qManufacturer := r.FormValue("manufacturer")
		qDescription := r.FormValue("description")
		qPartNumber := r.FormValue("part_number")
		qPhysDescription := r.FormValue("phys_description")
		url := "./list?search-type=" + qType +
			   "&search-subtype=" + qSubtype +
			   "&search-manufacturer=" + qManufacturer +
			   "&search-description=" + qDescription +
			   "&search-part-number=" + qPartNumber +
			   "&search-phys-description=" + qPhysDescription
		http.Redirect(w, r, url, 301)
	}
}

func CommitItemHandler(w http.ResponseWriter, r *http.Request) {
	if !checkSession(w, r) {
		return
	}
	if r.Method == "POST" {
		itemID := r.FormValue("id")
		intItemID, err := strconv.Atoi(itemID)
		if err != nil {
			errT, _ := template.ParseFiles(invTemplatePath + "/message.gtpl")
			msg := MessagePage{
					Prefix: invPrefix,
					Header: "Failed to Process Item",
					Message: template.HTML(
							"<p>Invalid Field Values</p>",
					),
			}
			errT.Execute(w, msg)
		}		
		var item Item
		item.ItemID = intItemID
		item.Model_number = r.FormValue("model")
		item.Manufacturer = r.FormValue("manufacturer")
		item.Type = r.FormValue("type")
		item.Subtype = r.FormValue("subtype")
		item.Descriptive_name = r.FormValue("description")
		item.Phys_description = r.FormValue("phys_description")
		item.ProductURL = r.FormValue("productURL")
		item.DatasheetURL = r.FormValue("datasheetURL")
		item.Seller1URL = r.FormValue("seller1URL")
		item.Seller2URL = r.FormValue("seller2URL")
		item.Seller3URL = r.FormValue("seller3URL")
		item.UnitPrice = r.FormValue("unitprice")
		item.Notes = r.FormValue("notes")
		item.Value = r.FormValue("value")
		err = addUpdateItem(item)
		if err != nil {
			errT, _ := template.ParseFiles(invTemplatePath + "/message.gtpl")
			msg := MessagePage{
					Prefix: invPrefix,
					Header: "Failed to Process Item",
					Message: template.HTML(
							"<p>Error: " + err.Error() + "</p>",
					),
			}
			errT.Execute(w, msg)
		} else {
			http.Redirect(w, r, "./item?id=" + itemID, 301)			
		}
	}
}

func QtyPostHandler(w http.ResponseWriter, r *http.Request) {
	if !checkSession(w, r) {
		return
	}
	if r.Method == "POST" {
		itemID := r.FormValue("id")
		location := r.FormValue("location")
		quantity := r.FormValue("quantity")
		err := updateInventory(itemID, location, quantity)
		if err != nil {
			errT, _ := template.ParseFiles(invTemplatePath + "/message.gtpl")
			msg := MessagePage{
					Prefix: invPrefix,
					Header: "Failed to Process Item",
					Message: template.HTML(
							"<p>Error: " + err.Error() + "</p>",
					),
			}
			errT.Execute(w, msg)
		} else {
			http.Redirect(w, r, "./item?id=" + itemID, 301)			
		}
	}
}

func DeleteEntryHandler(w http.ResponseWriter, r *http.Request) {
	if !checkSession(w, r) {
		return
	}
	var itemID string
	var location string
	if r.Method == "GET" {
		itemID = r.URL.Query().Get("id")
		location = r.URL.Query().Get("location")
	} else {
		itemID = r.FormValue("id")
		location = r.FormValue("location")
	}
	fmt.Println("itemID =", itemID, "location =", location)
	if itemID == "" || location == "" {
		t, _ := template.ParseFiles(invTemplatePath + "/message.gtpl")
		msg := MessagePage{
				Prefix: invPrefix,
				Header: "Missing Fields",
				Message: template.HTML(
						"<p>Some fields missing</p>",
				),
		}
		t.Execute(w, msg)
	} else {
		err := deleteInventoryEntry(itemID, location)
		if err != nil {
			errT, _ := template.ParseFiles(invTemplatePath + "/message.gtpl")
			msg := MessagePage{
					Prefix: invPrefix,
					Header: "Failed to Process Item",
					Message: template.HTML(
							"<p>Error: " + err.Error() + "</p>",
					),
			}
			errT.Execute(w, msg)
		} else {
			http.Redirect(w, r, "./item?id=" + itemID, 301)
		}
	}
}

func AddEntryHandler(w http.ResponseWriter, r *http.Request) {
	if !checkSession(w, r) {
		return
	}
	var itemID string
	var location string
	if r.Method == "GET" {
		itemID = r.URL.Query().Get("id")
		location = r.URL.Query().Get("location")
	} else {
		itemID = r.FormValue("id")
		location = r.FormValue("location")
	}
	fmt.Println("itemID =", itemID, "location =", location)
	if itemID == "" || location == "" {
		t, _ := template.ParseFiles(invTemplatePath + "/message.gtpl")
		msg := MessagePage{
				Prefix: invPrefix,
				Header: "Missing Fields",
				Message: template.HTML(
						"<p>Some fields missing</p>",
				),
		}
		t.Execute(w, msg)
	} else {
		err := addInventoryEntry(itemID, location, "0")
		if err != nil {
			errT, _ := template.ParseFiles(invTemplatePath + "/message.gtpl")
			msg := MessagePage{
					Prefix: invPrefix,
					Header: "Failed to Process Item",
					Message: template.HTML(
							"<p>Error: " + err.Error() + "</p>",
					),
			}
			errT.Execute(w, msg)
		} else {
			http.Redirect(w, r, "./item?id=" + itemID, 301)
		}
	}
}

func DownloadItemsHandler(w http.ResponseWriter, r *http.Request) {
	if !checkSession(w, r) {
		return
	}
	now := time.Now()
	timestamp := fmt.Sprintf("%d-%02d-%02d-%02dh%02dm%02ds",
		now.Year(), int(now.Month()), now.Day(),
		now.Hour(), now.Minute(), now.Second())
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment;filename=items-" +
		timestamp + ".csv")
	err := ExportItems(w)
	if err != nil {
		// oh wow
	}
}

func DownloadInventoryHandler(w http.ResponseWriter, r *http.Request) {
	if !checkSession(w, r) {
		return
	}
	now := time.Now()
	timestamp := fmt.Sprintf("%d-%02d-%02d-%02dh%02dm%02ds",
		now.Year(), int(now.Month()), now.Day(),
		now.Hour(), now.Minute(), now.Second())
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment;filename=inventory-" +
		timestamp + ".csv")
	err := ExportInventory(w)
	if err != nil {
		// oh wow
	}
}

func checkSession(w http.ResponseWriter, r *http.Request) bool {
	s := session.Get(r)
	if s == nil {
		http.Redirect(w, r, invPrefix + "/login", 301)
		return false
	}
	if s.CAttr("Inventory") == nil {
		http.Redirect(w, r, invPrefix + "/login", 301)
		return false
	}
	inventoryToken := s.CAttr("Inventory").(string)
	if inventoryToken != "Yes" {
		http.Redirect(w, r, invPrefix + "/login", 301)
		return false
	}
	return true
}

func Install(prefix string, usersFile string, templateDir string, staticDir string,
		dbConf string) (error) {
	if prefix != "" {
		invPrefix = "/" + prefix
	} else {
		invPrefix = ""
	}
	invTemplatePath = templateDir
	fmt.Println("Installing inventory handlers to '" + prefix + "'")
	// check for template files
	templates := []string{"browse.gtpl", "edit.gtpl", "item.gtpl",
						  "list.gtpl", "login.gtpl", "message.gtpl",
						  "search.gtpl"}
	for _, fileName := range templates {
		if _, err := os.Stat(templateDir + "/" + fileName); os.IsNotExist(err) {
			return errors.New(templateDir + "/" + fileName + " does not exist")
		}
	}
	
	http.HandleFunc(invPrefix+"/login", LoginPage)
	http.HandleFunc(invPrefix+"/list", ListingPage)
	http.HandleFunc(invPrefix+"/logout", LogoutPage)
	http.HandleFunc(invPrefix+"/item", ItemPage)
	http.HandleFunc(invPrefix+"/edit", AddEditItemPage)
	http.HandleFunc(invPrefix+"/delete", DeleteHandler)
	http.HandleFunc(invPrefix+"/commit", CommitItemHandler)
	http.HandleFunc(invPrefix+"/modify-qty", QtyPostHandler)
	http.HandleFunc(invPrefix+"/delete-entry", DeleteEntryHandler)
	http.HandleFunc(invPrefix+"/add-entry", AddEntryHandler)
	http.HandleFunc(invPrefix+"/download-items", DownloadItemsHandler)		
	http.HandleFunc(invPrefix+"/download-inventory", DownloadInventoryHandler)
	http.HandleFunc(invPrefix+"/browse", BrowsePage)
	http.HandleFunc(invPrefix+"/search", SearchHandler)
	fs := http.FileServer(http.Dir(staticDir))
	http.Handle(invPrefix+"/static/", http.StripPrefix(invPrefix+"/static/", fs))
	fmt.Println("Loading users data from '" + usersFile + "'...")
	invUsers = netutil.NewUsers()
	err := invUsers.LoadFromFile(usersFile)
	if err != nil {
		return err 
	}
	fmt.Println("Connecting to database defined by '" +
				os.Args[3] + "'...")
	invDB, err = netutil.OpenPostgresDBFromConfig(dbConf)
	if err != nil {
		return err
	}
	return nil
}
