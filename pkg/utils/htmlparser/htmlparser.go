package htmlparser

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

func ParseHTML(HTMLContent string, model any) error {
	root, err := htmlquery.Parse(strings.NewReader(HTMLContent))
	if err != nil {
		return fmt.Errorf("failed to parse html: %w", err)
	}

	return parseRoot(root, model)
}

func parseRoot(root *html.Node, model any) error {
	refValue := reflect.ValueOf(model)

	if refValue.Kind() != reflect.Ptr {
		return fmt.Errorf("model should be a pointer")
	}
	refValue = refValue.Elem()

	if refValue.Kind() != reflect.Struct {
		return fmt.Errorf("model should be a struct")
	}

	// traverse the fields of the struct
	for i := 0; i < refValue.NumField(); i++ {
		field := refValue.Field(i)
		fieldType := refValue.Type().Field(i)
		tag := fieldType.Tag.Get("xpath")
		if tag == "" {
			continue
		}

		switch field.Kind() {
		case reflect.String:
			node := htmlquery.FindOne(root, tag)
			if node == nil {
				continue
			}
			field.SetString(htmlquery.InnerText(node))
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			node := htmlquery.FindOne(root, tag)
			if node == nil {
				continue
			}
			value, err := strconv.ParseInt(strings.TrimSpace(htmlquery.InnerText(node)), 10, 64)
			if err != nil {
				return fmt.Errorf("failed to parse int: %w", err)
			}
			field.SetInt(value)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			node := htmlquery.FindOne(root, tag)
			if node == nil {
				continue
			}
			value, err := strconv.ParseUint(htmlquery.InnerText(node), 10, 64)
			if err != nil {
				return fmt.Errorf("failed to parse uint: %w", err)
			}
			field.SetUint(value)
		case reflect.Float32, reflect.Float64:
			node := htmlquery.FindOne(root, tag)
			if node == nil {
				continue
			}
			value, err := strconv.ParseFloat(htmlquery.InnerText(node), 64)
			if err != nil {
				return fmt.Errorf("failed to parse float: %w", err)
			}
			field.SetFloat(value)
		case reflect.Bool:
			node := htmlquery.FindOne(root, tag)
			if node == nil {
				continue
			}
			value, err := strconv.ParseBool(htmlquery.InnerText(node))
			if err != nil {
				return fmt.Errorf("failed to parse bool: %w", err)
			}
			field.SetBool(value)
		case reflect.Slice:
			nodes := htmlquery.Find(root, tag)
			if len(nodes) == 0 {
				continue
			}

			elemType := field.Type().Elem()
			for _, node := range nodes {
				// if the element type is ptr, we need to create a new instance
				if elemType.Kind() == reflect.Ptr {
					elemValue := reflect.New(elemType.Elem())
					err := parseRoot(node, elemValue.Interface())
					if err != nil {
						return err
					}
					field.Set(reflect.Append(field, elemValue))
				} else {
					elemValue := reflect.New(elemType).Elem()
					err := parseRoot(node, elemValue.Addr().Interface())
					if err != nil {
						return err
					}

					field.Set(reflect.Append(field, elemValue))
				}
			}
		case reflect.Struct:
			node := htmlquery.FindOne(root, tag)
			if node == nil {
				continue
			}

			err := parseRoot(node, field.Addr().Interface())
			if err != nil {
				return err
			}

		default:
			return fmt.Errorf("unsupported field type: %s", field.Kind())
		}
	}

	return nil
}
