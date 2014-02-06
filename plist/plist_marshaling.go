// This file extends the plist package with the ability to marshal a struct into plist byte data

package plist

import (
    "reflect"
    "strconv"
)

func Marshal(v interface{}) ([]byte, error) {
    marshaled, err := marshalValue(reflect.ValueOf(v), []byte("\n<plist version=\"1.0\">\n"), "", 0)
    if err != nil {
        return nil, err
    }

    return append(marshaled, []byte("</plist>\n")...), nil
}

func marshalValue(v reflect.Value, startData []byte, fieldName string, indentLevel int) (data []byte, err error) {
    data = append(startData)
    if fieldName != "" {
        data = appendKey(data, fieldName, indentLevel)
    }
    switch v.Kind() {
    case reflect.Ptr:
        {
            return marshalValue(v.Elem(), startData, fieldName, indentLevel)
        }
    case reflect.Struct:
        {
            data = appendLineWithIndent(data, []byte("<dict>"), indentLevel)

            t := v.Type()
            for i := 0; i < t.NumField(); i++ {
                f := t.Field(i)
                fieldName := ""
                tag := f.Tag.Get("plist")
                if tag != "" {
                    fieldName = tag
                } else {
                    fieldName = f.Name
                }

                data, err = marshalValue(v.Field(i), data, fieldName, indentLevel+1)

                if err != nil {
                    return []byte{}, err
                }

            }

            data = appendLineWithIndent(data, []byte("</dict>"), indentLevel)
        }
    case reflect.Bool:
        {
            boolVal := v.Bool()
            if boolVal {
                data = appendLineWithIndent(data, []byte("<true/>"), indentLevel)
            } else {
                data = appendLineWithIndent(data, []byte("<false/>"), indentLevel)
            }
        }
    case reflect.Int:
        {
            intVal := v.Int()

            tempData := append([]byte("<integer>"), []byte(strconv.Itoa(int(intVal)))...)
            tempData = append(tempData, []byte("</integer>")...)

            data = appendLineWithIndent(data, tempData, indentLevel)
        }
    case reflect.String:
        {
            strVal := v.String()

            tempData := append([]byte("<string>"), []byte(strVal)...)
            tempData = append(tempData, []byte("</string>")...)

            data = appendLineWithIndent(data, tempData, indentLevel)
        }
    case reflect.Slice:
        {
            data = appendLineWithIndent(data, []byte("<array>"), indentLevel)

            for i := 0; i < v.Len(); i++ {
                data, err = marshalValue(v.Index(i), data, "", indentLevel+1)

                if err != nil {
                    return []byte{}, err
                }
            }

            data = appendLineWithIndent(data, []byte("</array>"), indentLevel)
        }
    default:
        // ignore anything not listed
        return
    }
    return
}

func appendLineWithIndent(startData []byte, appendData []byte, indentLevel int) []byte {
    data := append(startData)
    for i := 0; i < indentLevel; i++ {
        data = append(data, '\t')
    }
    data = append(data, appendData...)
    data = append(data, '\n')

    return data
}

func appendKey(startData []byte, key string, indentLevel int) []byte {
    data := append([]byte{}, []byte("<key>")...)

    data = append(data, []byte(key)...)

    data = append(data, []byte("</key>")...)

    return appendLineWithIndent(startData, data, indentLevel)
}
