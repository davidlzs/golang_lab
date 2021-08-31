package main

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/gocql/gocql"
	"github.com/rivo/tview"
)

const (
	batchSize  = 80         // The number of rows loaded per batch.
	finderPage = "*finder*" // The name of the Finder page.
)

var (
	app         *tview.Application // The tview application.
	pages       *tview.Pages       // The application pages.
	finderFocus tview.Primitive    // The primitive in the Finder that last had focus.
)

// Main entry point.
func main() {
	// Get connect string from the command line.
	if len(os.Args) < 2 {
		fmt.Println(`Please provide a Cassandra initial node IP.`)
		fmt.Println("See https://pkg.go.dev/github.com/gocql/gocql#pkg-overview for details.")
		return
	}

	// Start the application.
	app = tview.NewApplication()
	finder(os.Args[1])
	if err := app.Run(); err != nil {
		fmt.Printf("Error running application: %s\n", err)
	}
}

// Sets up a "Finder" used to navigate the databases, tables, and columns.
func finder(connString string) {
	// Create the basic objects.
	databases := tview.NewList().ShowSecondaryText(false)
	databases.SetBorder(true).SetTitle("Databases")
	columns := tview.NewTable().SetBorders(true)
	columns.SetBorder(true).SetTitle("Columns")
	tables := tview.NewList()
	tables.ShowSecondaryText(false).
		SetDoneFunc(func() {
			tables.Clear()
			columns.Clear()
			app.SetFocus(databases)
		})
	tables.SetBorder(true).SetTitle("Tables")

	// Create the layout.
	flex := tview.NewFlex().
		AddItem(databases, 0, 1, true).
		AddItem(tables, 0, 1, false).
		AddItem(columns, 0, 3, false)

	cluster := gocql.NewCluster(connString)
	cluster.Consistency = gocql.Quorum
	session, err := cluster.CreateSession()
	if err != nil {
		panic(err)
	}
	GetSession := func(keyspace string) *gocql.Session {
		return session
	}

	// Get a list of all databases.
	generalKeyspace, err := cluster.CreateSession()
	if err != nil {
		panic(err)
	}
	defer generalKeyspace.Close() // We really close the DB because we only use it for this one query.
	ctx := context.Background()
	rows := generalKeyspace.Query(`select keyspace_name from system_schema.keyspaces`).
		WithContext(ctx).Iter().Scanner()

	for rows.Next() {
		var keyspaceName string
		if err := rows.Scan(&keyspaceName); err != nil {
			panic(err)
		}
		databases.AddItem(keyspaceName, "", 0, func() {
			// A database was selected. Show all of its tables.
			columns.Clear()
			tables.Clear()
			db := GetSession(keyspaceName)

			t := db.Query(`select table_name from system_schema.tables where keyspace_name=?`,
				keyspaceName).
				WithContext(ctx).Iter().Scanner()
			for t.Next() {
				var tableName string
				if err := t.Scan(&tableName); err != nil {
					panic(err)
				}
				tables.AddItem(tableName, "", 0, nil)
			}
			if err := t.Err(); err != nil {
				panic(err)
			}
			app.SetFocus(tables)

			// When the user navigates to a table, show its columns.
			tables.SetChangedFunc(func(i int, tableName string, t string, s rune) {
				// A table was selected. Show its columns.
				columns.Clear()

				// stat, err := db.Prepare(`
				// select c.column_name,
				// 	c.is_nullable,
				// 	c.data_type,
				// 	c.character_maximum_length,
				// 	c.numeric_precision,
				// 	c.numeric_scale,
				// 	c.ordinal_position,
				// 	k.constraint_name
				// from INFORMATION_SCHEMA.columns c
				// left join information_schema.key_column_usage k
				// 	on c.table_schema = k.table_schema and c.table_name = k.table_name and c.column_name = k.column_name AND k.constraint_name='PRIMARY'
				// where c.table_schema = ?
				// and c.table_name = ?
				// 		`)
				// if err != nil {
				// 	panic(err)
				// }

				// c, err := stat.Query(keyspaceName, tableName)
				// if err != nil {
				// 	panic(err)
				// }
				// defer c.Close()
				// columns.SetCell(0, 0, &tview.TableCell{Text: "Name", Align: tview.AlignCenter, Color: tcell.ColorYellow}).
				// 	SetCell(0, 1, &tview.TableCell{Text: "Type", Align: tview.AlignCenter, Color: tcell.ColorYellow}).
				// 	SetCell(0, 2, &tview.TableCell{Text: "Size", Align: tview.AlignCenter, Color: tcell.ColorYellow}).
				// 	SetCell(0, 3, &tview.TableCell{Text: "Null", Align: tview.AlignCenter, Color: tcell.ColorYellow}).
				// 	SetCell(0, 4, &tview.TableCell{Text: "Constraint", Align: tview.AlignCenter, Color: tcell.ColorYellow})
				// for c.Next() {
				// 	var (
				// 		columnName, isNullable, dataType     string
				// 		constraintType                       sql.NullString
				// 		size, numericPrecision, numericScale sql.NullInt64
				// 		ordinalPosition                      int
				// 	)
				// 	if err := c.Scan(&columnName,
				// 		&isNullable,
				// 		&dataType,
				// 		&size,
				// 		&numericPrecision,
				// 		&numericScale,
				// 		&ordinalPosition,
				// 		&constraintType,
				// 	); err != nil {
				// 		panic(err)
				// 	}
				// 	sizeText := ""
				// 	if size.Valid {
				// 		sizeText = strconv.Itoa(int(size.Int64))
				// 	} else if numericPrecision.Valid {
				// 		sizeText = strconv.Itoa(int(numericPrecision.Int64))
				// 		if numericScale.Valid {
				// 			sizeText += "," + strconv.Itoa(int(numericScale.Int64))
				// 		}
				// 	}
				// 	color := tcell.ColorWhite
				// 	if constraintType.Valid {
				// 		color = map[string]tcell.Color{
				// 			"CHECK":       tcell.ColorGreen,
				// 			"FOREIGN KEY": tcell.ColorDarkMagenta,
				// 			"PRIMARY KEY": tcell.ColorRed,
				// 			"UNIQUE":      tcell.ColorDarkCyan,
				// 		}[constraintType.String]
				// 	}
				// 	columns.SetCell(ordinalPosition, 0, &tview.TableCell{Text: columnName, Color: color}).
				// 		SetCell(ordinalPosition, 1, &tview.TableCell{Text: dataType, Color: color}).
				// 		SetCell(ordinalPosition, 2, &tview.TableCell{Text: sizeText, Align: tview.AlignRight, Color: color}).
				// 		SetCell(ordinalPosition, 3, &tview.TableCell{Text: isNullable, Align: tview.AlignRight, Color: color}).
				// 		SetCell(ordinalPosition, 4, &tview.TableCell{Text: constraintType.String, Align: tview.AlignLeft, Color: color})
				// }
				// if err := c.Err(); err != nil {
				// 	panic(err)
				// }
			})
			tables.SetCurrentItem(0) // Trigger the initial selection.

			// When the user selects a table, show its content.
			tables.SetSelectedFunc(func(i int, tableName string, t string, s rune) {
				session := GetSession(keyspaceName)
				content(session, keyspaceName, tableName)
			})
		})
	}
	if err := rows.Err(); err != nil {
		panic(err)
	}

	// Set up the pages and show the Finder.
	pages = tview.NewPages().
		AddPage(finderPage, flex, true, true)
	app.SetRoot(pages, true)
}

// Shows the contents of the given table.

// TODO: Cassandra pagination: http://www.inanzzz.com/index.php/post/t7fd/cassandra-pagination-example-with-golang
func content(db *gocql.Session, keyspaceName, tableName string) {
	finderFocus = app.GetFocus()

	// If this page already exists, just show it.
	if pages.HasPage(keyspaceName + "." + tableName) {
		pages.SwitchToPage(keyspaceName + "." + tableName)
		return
	}

	// We display the data in a table embedded in a frame.
	table := tview.NewTable().
		SetFixed(1, 0).
		SetSeparator(tview.BoxDrawingsLightHorizontal).
		SetBordersColor(tcell.ColorYellow)
	frame := tview.NewFrame(table).
		SetBorders(0, 0, 0, 0, 0, 0)
	frame.SetBorder(true).
		SetTitle(fmt.Sprintf(`Contents of table "%s"`, tableName))

	// How many rows does this table have?
	var rowCount int
	fullyQualifiedTableName := fmt.Sprintf("%s.%s", keyspaceName, tableName)
	ctx := context.Background()
	err := db.Query(fmt.Sprintf("select count(*) from %s", fullyQualifiedTableName)).
		WithContext(ctx).Consistency(gocql.One).Scan(&rowCount)
	if err != nil {
		panic(err)
	}

	// Load a batch of rows.
	loadRows := func(offset int) {

		iter := db.Query(fmt.Sprintf(`select persistence_id from %s`, fullyQualifiedTableName)).
			WithContext(ctx).Iter()
		columns := iter.Columns()
		for index, columnInfo := range columns {
			name := columnInfo.Name
			// fmt.Fprintf(w, "%s (%s)", columnInfo.Name, columnInfo.TypeInfo)
			table.SetCell(0, index, &tview.TableCell{Text: name, Align: tview.AlignCenter, Color: tcell.ColorYellow})
		}

		// The first row in the table is the list of column names.
		// columnNames, err := rows.Columns()
		// if err != nil {
		// 	panic(err)
		// }
		// for index, name := range columnNames {
		// 	table.SetCell(0, index, &tview.TableCell{Text: name, Align: tview.AlignCenter, Color: tcell.ColorYellow})
		// }

		// Transfer them to the table.
		row := table.GetRowCount()

		for {
			rd, err := iter.RowData()
			if err != nil {
				panic(err)
			}

			if !iter.Scan(rd.Values...) {
				break
			}

			for index, value := range rd.Values {
				switch columns[index].TypeInfo.Type() {
				case gocql.TypeBigInt:
					table.SetCell(row, index, &tview.TableCell{Text: strconv.Itoa(int(*value.(*int64))), Align: tview.AlignRight, Color: tcell.ColorDarkCyan})
				case gocql.TypeFloat:
					table.SetCell(row, index, &tview.TableCell{Text: strconv.FormatFloat(*value.(*float64), 'f', 2, 64), Align: tview.AlignRight, Color: tcell.ColorDarkCyan})
				case gocql.TypeVarchar:
					table.SetCellSimple(row, index, *value.(*string))
				// case 0x000B:
				// 	t := value.Format("2006-01-02")
				// 	table.SetCell(row, index, &tview.TableCell{Text: t, Align: tview.AlignRight, Color: tcell.ColorDarkMagenta})
				// case []uint8:
				// 	str := make([]byte, len(value))
				// 	for index, num := range value {
				// 		str[index] = byte(num)
				// 	}
				// 	table.SetCell(row, index, &tview.TableCell{Text: string(str), Align: tview.AlignRight, Color: tcell.ColorGreen})
				default:
					// We've encountered a type that we don't know yet.
					t := reflect.TypeOf(value)
					str := "?nil?"
					if t != nil {
						str = "?" + t.String() + "?"
					}
					table.SetCellSimple(row, index, str)
				}
			}
		}

		// Show how much we've loaded.
		frame.Clear()
		loadMore := ""
		if table.GetRowCount()-1 < rowCount {
			loadMore = " - press Enter to load more"
		}
		loadMore = fmt.Sprintf("Loaded %d of %d rows%s", table.GetRowCount()-1, rowCount, loadMore)
		frame.AddText(loadMore, false, tview.AlignCenter, tcell.ColorYellow)
	}

	// Load the first batch of rows.
	loadRows(0)

	// Handle key presses.
	table.SetDoneFunc(func(key tcell.Key) {
		switch key {
		case tcell.KeyEscape:
			// Go back to Finder.
			pages.SwitchToPage(finderPage)
			if finderFocus != nil {
				app.SetFocus(finderFocus)
			}
		case tcell.KeyEnter:
			// Load the next batch of rows.
			// loadRows(table.GetRowCount() - 1)
			table.ScrollToEnd()
		}
	})

	// Add a new page and show it.
	pages.AddPage(keyspaceName+"."+tableName, frame, true, true)
}
