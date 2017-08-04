package main

import (
	"os"
	"fmt"
	"log"
	"strings"
	"strconv"
	"encoding/csv"
	"bbi/netutil"
)

type Item struct {
	serial				int
	ItemID 				int
	Descriptive_name 	string
	Model_number 		string
	Manufacturer		string
	Type				string
	Subtype				string	
	Phys_description	string
	DatasheetURL		string
	ProductURL			string
	Seller1URL			string
	Seller2URL			string
	Seller3URL			string
	UnitPrice			float64
	Notes				string
	TotalQty			int
}

type InventoryEntry struct {
	serial		int
	ItemID		int
	Location	string
	Quantity	int
}

type FetchFilter struct {
	key		string
	value	string
}

func getItem(itemID string) (*Item, []InventoryEntry, error) {
	var count int
	rows, err := db.Query("select * from items where itemID=$1", itemID)	
	if err != nil {
		return nil, nil, err
	}
	rows.Scan(&count)
	defer rows.Close()
	for rows.Next() {
		var item Item
		item.TotalQty = 0
		err := rows.Scan(
			&item.serial, &item.ItemID, &item.Descriptive_name,
			&item.Model_number, &item.Manufacturer, 
			&item.Type, &item.Subtype,
			&item.Phys_description,
			&item.DatasheetURL, &item.ProductURL,
			&item.Seller1URL, &item.Seller2URL,
			&item.Seller3URL, &item.UnitPrice, &item.Notes,
			)
		if err != nil {			
			return nil, nil, err
		}
		entries := getInventoryEntries(item.ItemID)
		for i := range entries {
			if entries[i].ItemID == item.ItemID {
				item.TotalQty = item.TotalQty + entries[i].Quantity
			}
		}
		return &item, entries, nil
	}
	return nil, nil, nil
}

func getItems() []Item {
	return getItemsFiltered()
}

func getItemsFiltered(filters...FetchFilter) []Item {
	var count int
	values := make([]interface{}, len(filters))
	stmt := "select * from items"
	if filters != nil {
		stmt = stmt + " where ("
		for i := range filters {
			count := strconv.Itoa(i+1)
			stmt = stmt + filters[i].key + "=$" + count + " and "
			values[i] = filters[i].value
		}
		stmt = strings.TrimRight(stmt, " and ")
		stmt = stmt + ")"
	}
	rows, err := db.Query(stmt, values...)	
	if err != nil {
		fmt.Println(err)
		return nil
	}
	rows.Scan(&count)
	list := make([]Item, count)
	defer rows.Close()
	for rows.Next() {
		var item Item
		item.TotalQty = 0
		err := rows.Scan(
			&item.serial, &item.ItemID, &item.Descriptive_name,
			&item.Model_number, &item.Manufacturer, 
			&item.Type, &item.Subtype,
			&item.Phys_description,
			&item.DatasheetURL, &item.ProductURL,
			&item.Seller1URL, &item.Seller2URL,
			&item.Seller3URL, &item.UnitPrice, &item.Notes,
			)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		entries := getInventoryEntries(item.ItemID)
		for i := range entries {
			if entries[i].ItemID == item.ItemID {
				item.TotalQty = item.TotalQty + entries[i].Quantity
			}
		}
		list = append(list, item)
	}
	return list
}

func getInventoryEntries(id int) []InventoryEntry {
	itemID := strconv.Itoa(id)
	list := make([]InventoryEntry, 0)
	invrows, err := db.Query("select * from inventory where " +
		"itemID=" + itemID)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer invrows.Close()
	for invrows.Next() {
		var entry InventoryEntry
		err = invrows.Scan(&entry.serial, &entry.ItemID,
			&entry.Location, &entry.Quantity)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		list = append(list, entry)
	}
	return list
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
				unitPrice float(2),
				notes TEXT
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
			if err != nil || len(tokens) != 14 {
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
				fmt.Println("Items Import failed")
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
				"seller2URL, seller3URL, unitPrice, notes) values( " +
				"$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, " +
				"$13, $14)")
			if err != nil {
				fmt.Println("Items Import failed")
				return
			}
			_, err = stmt.Exec(
				tokens[0], tokens[1], tokens[2], tokens[3], tokens[4],
				tokens[5], tokens[6], tokens[7], tokens[8], tokens[9],
				tokens[10], tokens[11], fPrice, tokens[13])
			if err != nil {
				fmt.Println("Items Import failed")
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
			var entry InventoryEntry
			err := rows.Scan(&entry.serial,
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
			var item Item
			err := rows.Scan(
				&item.serial, &item.ItemID, &item.Descriptive_name,
				&item.Model_number, &item.Manufacturer, 
				&item.Type, &item.Subtype,
				&item.Phys_description,
				&item.DatasheetURL, &item.ProductURL,
				&item.Seller1URL, &item.Seller2URL,
				&item.Seller3URL, &item.UnitPrice, &item.Notes,
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
		w := csv.NewWriter(f)
		defer w.Flush()
		rows, err := db.Query("select * from inventory")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
		for rows.Next() {
			var entry InventoryEntry
			err := rows.Scan(&entry.serial,
				&entry.ItemID, &entry.Location, &entry.Quantity)
			if err != nil {
				log.Fatal(err)
			}
			line := make([]string, 3)
			line[0] = strconv.Itoa(entry.ItemID)
			line[1] = entry.Location
			line[2] = strconv.Itoa(entry.Quantity)
			w.Write(line)
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
		w := csv.NewWriter(f)
		defer w.Flush()
		rows, err := db.Query("select * from items")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
		for rows.Next() {
			var item Item
			err := rows.Scan(
				&item.serial, &item.ItemID, &item.Descriptive_name,
				&item.Model_number, &item.Manufacturer, 
				&item.Type, &item.Subtype,
				&item.Phys_description,
				&item.DatasheetURL, &item.ProductURL,
				&item.Seller1URL, &item.Seller2URL,
				&item.Seller3URL, &item.UnitPrice, &item.Notes,
				)
			if err != nil {
				log.Fatal(err)
			}
			line := make([]string, 14)
			line[0] = strconv.Itoa(item.ItemID)
			line[1] = item.Descriptive_name
			line[2] = item.Model_number
			line[3] = item.Manufacturer
			line[4] = item.Type
			line[5] = item.Subtype
			line[6] = item.Phys_description
			line[7] = item.DatasheetURL
			line[8] = item.ProductURL
			line[9] = item.Seller1URL
			line[10] = item.Seller2URL
			line[11] = item.Seller3URL
			line[12] = strconv.FormatFloat(item.UnitPrice, 'E', -1, 64)
			line[13] = item.Notes
			w.Write(line)
		}
	default:
		fmt.Println(command, "is an invalid subcommand\n")
		dbUsage()
	}
}