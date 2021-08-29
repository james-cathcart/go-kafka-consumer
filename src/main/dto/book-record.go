package dto

import "github.com/linkedin/goavro/v2"

type Author struct {
	FirstName string
	LastName  string
}

type Book struct {
	ID     string
	Title  string
	Author Author
	Errors []string
}

func (book *Book) ToStringMap() map[string]interface{} {

	datumIn := map[string]interface{}{
		"ID":    string(book.ID),
		"Title": string(book.Title),
	}

	if len(book.Errors) > 0 {
		datumIn["Errors"] = goavro.Union("array", book.Errors)
	} else {
		datumIn["Errors"] = goavro.Union("null", nil)
	}

	authorDatum := map[string]interface{}{
		"FirstName": string(book.Author.FirstName),
		"LastName":  string(book.Author.LastName),
	}

	datumIn["Author"] = goavro.Union("my.namespace.com.author", authorDatum)

	return datumIn
}

func StringMapToUser(data map[string]interface{}) *Book {

	bookPlaceHolder := &Book{}
	for k, v := range data {
		switch k {
		case "ID":
			if value, ok := v.(string); ok {
				bookPlaceHolder.ID = value
			}
		case "Title":
			if value, ok := v.(string); ok {
				bookPlaceHolder.Title = value
			}
		case "Errors":
			if value, ok := v.(map[string]interface{}); ok {
				for _, item := range value["array"].([]interface{}) {
					bookPlaceHolder.Errors = append(bookPlaceHolder.Errors, item.(string))
				}
			}
		case "Author":
			if vmap, ok := v.(map[string]interface{}); ok {
				if cookieSMap, ok := vmap["my.namespace.com.author"].(map[string]interface{}); ok {
					authorPlaceHolder := &Author{}
					for k, v := range cookieSMap {
						switch k {
						case "FirstName":
							if value, ok := v.(string); ok {
								authorPlaceHolder.FirstName = value
							}
						case "LastName":
							if value, ok := v.(string); ok {
								authorPlaceHolder.LastName = value
							}
						}
					}
					bookPlaceHolder.Author = *authorPlaceHolder
				}
			}
		}
	}
	return bookPlaceHolder
}
