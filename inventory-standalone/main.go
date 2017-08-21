package main

import (
	"os"
	"fmt"
	"log"
	"strconv"
	"net/http"
	"encoding/csv"
	"bbi/netutil"
	"bbi/inventory"	
)

func main() {
	if len(os.Args) < 2 {
		usage()
		return
	}
	switch os.Args[1] {
	case "users":
		if len(os.Args) < 4 {
			usersUsage()
			return
		}
		handleUserOps()
	case "db":
		if len(os.Args) < 4 {
			dbUsage()
			return
		}
		handleDbOps()
	case "serve":
		if len(os.Args) != 7 {
			serveUsage()
			return
		}

		err := inventory.Install("", os.Args[2], "templates", os.Args[6], os.Args[3])
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Listening for connections on port 44443...")
		err = http.ListenAndServeTLS(":44443", os.Args[4], os.Args[5], nil)
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	}
}

func usersUsage() {
	fmt.Println("  users [users-file] list")
	fmt.Println("                     add [username] [password]")
	fmt.Println("                     delete [username]")
	fmt.Println("                     test-login [username] [password]")
}

func serveUsage() {
	fmt.Println("  serve [users-file] [db-config] [cert] [key] [static-dir]")
}

func dbUsage() {
	fmt.Println("  db [db-config] create-default-config")
	fmt.Println("                 create-tables")
	fmt.Println("                 delete-tables")
	fmt.Println("                 export-items [output-file]")
	fmt.Println("                 export-inventory [output-file]")
	fmt.Println("                 import-items [input-file]")
	fmt.Println("                 import-inventory [input-file]")
	fmt.Println("                 list-items")
	fmt.Println("                 list-inventory")
}

func usage() {
	fmt.Println("usage: inventory [command]\n")
	fmt.Println("commands:")
	usersUsage()
	fmt.Println()
	dbUsage()
	fmt.Println()	
	serveUsage()
	fmt.Println()
}

func handleUserOps() {
	path := os.Args[2]
	command := os.Args[3]
	params := os.Args[4:]
	users := netutil.NewUsers()
	err := users.LoadFromFile(path)
	if err != nil {
		log.Fatal(err)
	}

	switch command {
	case "list":
		userList := users.GetList()
		fmt.Println(len(userList), "users: ")
		for i := range userList {
			fmt.Print(i, " " + userList[i] + "\n")
		}
	case "add":
		if len(params) != 2 {
			fmt.Println("usage: inventory users [users-file] add [name] [password]")
			return
		}
		fmt.Println("Adding user '" + params[0] + "'")
		err := users.Add(params[0], params[1])
		if err != nil {
			log.Fatal(err)
		}
	case "test-login":
		if len(params) != 2 {
			fmt.Println("usage: inventory users [users-file] test-login [name] [password]")
			return
		}
		err := users.Login(params[0], params[1])
		if err != nil {
			log.Fatal(err)
		} else {
			fmt.Println("Login Successful")
		}
	case "delete":
		if len(params) != 1 {
			fmt.Println("usage: inventory users [users-file] delete [name]")
			return
		}
		users.Delete(params[0])
	default:
		fmt.Println(command, "is an invalid subcommand\n")
		usersUsage()
		return
	}
	
	err = users.SaveToFile(path)
	if err != nil {
		log.Fatal(err)
	}
}

func handleDbOps() {
	path := os.Args[2]
	command := os.Args[3]
	params := os.Args[4:]
	
	if command == "create-default-config" {
		if len(params) != 0 {
			fmt.Println("usage: inventory db [db-config] create-default-config")
			return
		}
		defaultStr := []byte("bbiinv bbipassword localhost 5432 bbiinvdb disable\n")
		f, err := os.Create(path)
		if err != nil {
			log.Fatal(err)
		}
		_, err = f.Write(defaultStr)
		if err != nil {
			log.Fatal(err)
		}
		f.Close()
		return;
	}
	var err error
	db, err := netutil.OpenPostgresDBFromConfig(path)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	
	switch command {
	case "create-tables":		
		stmt := `
			CREATE TABLE items (  
				id SERIAL PRIMARY KEY,
				itemID INTEGER,
				descriptive_name TEXT,
				model_number TEXT,
				manufacturer TEXT,
				type TEXT,
				subtype TEXT,
				phys_description TEXT,
				datasheetURL TEXT,
				productURL TEXT,
				seller1URL TEXT,
				seller2URL TEXT,
				seller3URL TEXT,
				unitPrice NUMERIC(8,2),
				notes TEXT,
				value NUMERIC(12,6)
			);`
		_, err = db.Exec(stmt)
		if err != nil {
			fmt.Println(err)
		}
		stmt = `
			CREATE TABLE inventory (
				id SERIAL PRIMARY KEY,
				itemID INTEGER,
				location TEXT,
				quantity INTEGER
			);`
		_, err = db.Exec(stmt)
		if err != nil {
			fmt.Println(err)
		}
	case "delete-tables":
		stmt := `DROP TABLE items`
		_, err = db.Exec(stmt)
		if err != nil {
			fmt.Println(err)
		}
		stmt = `DROP TABLE inventory`
		_, err = db.Exec(stmt)
		if err != nil {
			fmt.Println(err)
		}
	case "import-inventory":
		if len(params) != 1 {
			fmt.Println("usage: inventory db [db-config] import-inventory [input-file]")
			return
		}
		file, err := os.Open(params[0])
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		lines, err := csv.NewReader(file).ReadAll()
		if err != nil {
			log.Fatal(err)
		}
		for i := range lines {
			tokens := lines[i]
			if err != nil || len(tokens) != 3 {
				fmt.Println("Invalid Inventory File Format tokens=", len(tokens))
				return
			}
			itemID, err := strconv.Atoi(tokens[0])
			if err != nil {
				fmt.Println("Failed to parse item ID column")
				return
			}
			number, err := strconv.Atoi(tokens[2])
			if err != nil {
				fmt.Println("Failed to parse number column")
				return
			}
			// check if we just need to change the total number
			query := "select count(*) from inventory where " +
				"(itemID=" + tokens[0] + " and location='" +
				tokens[1] + "')"
			var rows int
			err = db.QueryRow(query).Scan(&rows)
			if rows > 0 {
				fmt.Println(tokens[0] + " in '" + tokens[1] + "' " +
					"already defined, changing number only")
				stmt, err := db.Prepare("update inventory set " +
					"quantity=$1 where (itemID=" + tokens[0] + " and " +
					"location='" + tokens[1] + "')")
				if err != nil {
					fmt.Println("Inventory Import Failed")
					return
				}
				_, err = stmt.Exec(number)
				if err != nil {
					fmt.Println("Inventory Import Failed")
					return
				}
			} else {
				stmt, err := db.Prepare("insert into inventory(" +
					"itemID, location, quantity) values(" +
					"$1, $2, $3)")
				if err != nil {
					fmt.Println("Inventory Import Failed")
					return
				}
				_, err = stmt.Exec(itemID, tokens[1], number)
				if err != nil {
					fmt.Println("Inventory Import Failed")
					return
				}
			}
		}
	case "import-items":
		if len(params) != 1 {
			fmt.Println("usage: inventory db [db-config] import-items [input-file]")
			return
		}
		file, err := os.Open(params[0])
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		lines, err := csv.NewReader(file).ReadAll()
		if err != nil {
			log.Fatal(err)
		}
		for i := range lines {
			tokens := lines[i]
			if err != nil || len(tokens) != 15 {
				fmt.Println("Invalid Items File Format tokens=", len(tokens))
				return
			}
			fPrice, err := strconv.ParseFloat(tokens[12], 32)
			if err != nil {
				fmt.Println("Failed to parse unit price column")
				return
			}
			// check if item is already defined
			query := "select count(*) from items where itemID=" + tokens[0]
			var rows int
			err = db.QueryRow(query).Scan(&rows)
			if err != nil {
				fmt.Println("Items Import failed:", err.Error())
				return
			}
			if rows > 0 {
				fmt.Println(tokens[0], "is already defined")
				continue
			}
			stmt, err := db.Prepare("insert into items(" +
				"itemID, descriptive_name, model_number, manufacturer, " +
				"type, subtype, " +
				"phys_description, datasheetURL, productURL, seller1URL, " +
				"seller2URL, seller3URL, unitPrice, notes, value) values( " +
				"$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, " +
				"$13, $14, $15)")
			if err != nil {
				fmt.Println("Items Import failed:", err.Error())
				return
			}
			_, err = stmt.Exec(
				tokens[0], tokens[1], tokens[2], tokens[3], tokens[4],
				tokens[5], tokens[6], tokens[7], tokens[8], tokens[9],
				tokens[10], tokens[11], fPrice, tokens[13], tokens[14])
			if err != nil {
				fmt.Println("Items Import failed:", err.Error())
				return
			}
		}
	case "list-inventory":
		rows, err := db.Query("select * from inventory")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("itemID\t Qty\t Location")
		fmt.Println("------\t ---\t --------")
		for rows.Next() {
			var entry inventory.InventoryEntry
			err := rows.Scan(&entry.Serial,
				&entry.ItemID, &entry.Location, &entry.Quantity)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(entry.ItemID, "\t", entry.Quantity, "\t",
						entry.Location)
		}
	case "list-items":
		rows, err := db.Query("select * from items")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
		fmt.Println("itemID\t descriptive_name")
		fmt.Println("------\t ----------------")
		for rows.Next() {
			var item inventory.Item
			err := rows.Scan(
				&item.Serial, &item.ItemID, &item.Descriptive_name,
				&item.Model_number, &item.Manufacturer, 
				&item.Type, &item.Subtype,
				&item.Phys_description,
				&item.DatasheetURL, &item.ProductURL,
				&item.Seller1URL, &item.Seller2URL,
				&item.Seller3URL, &item.UnitPrice, &item.Notes,
				&item.Value,
				)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(item.ItemID, "\t",
						item.Descriptive_name)
		}
	case "export-inventory":
		if len(params) != 1 {
			fmt.Println("usage: inventory db [db-config] export-inventory [output-file]")
			return
		}
		f, err := os.Create(params[0])
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		err = inventory.ExportInventory(f)
		if err != nil {
			log.Fatal(err)
		}
	case "export-items":
		if len(params) != 1 {
			fmt.Println("usage: inventory db [db-config] export-items [output-file]")
			return
		}
		f, err := os.Create(params[0])
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		err = inventory.ExportItems(f)				
		if err != nil {
			log.Fatal(err)
		}
	default:
		fmt.Println(command, "is an invalid subcommand\n")
		dbUsage()
	}
}
