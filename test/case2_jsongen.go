package test

import (
	"fmt"
)

func (this *RecursiveStruct) FastUnmarshalJSON(data string) error {
	pos := 0
	dataLen := len(data)
	if dataLen == 0 {
		return fmt.Errorf("Empty JSON")
	}

	var res RecursiveStruct
	var parser_RecursiveStruct func(*RecursiveStruct) error
	parser_RecursiveStruct = func(ptr *RecursiveStruct) error {
		var res RecursiveStruct
	L1:
		for {
			if pos >= dataLen {
				break L1
			}

			switch data[pos] {
			case ' ', '\t', '\r', '\n':
				pos++
				continue
			default:
				break L1
			}
		}

		if pos >= dataLen {
			return fmt.Errorf("Waited {, got EOF")
		}

		if data[pos] != '{' {
			return fmt.Errorf("Waited {, got %c", data[pos])
		}
		pos++

	L2:
		for {
		L3:
			for {
				if pos >= dataLen {
					break L3
				}

				switch data[pos] {
				case ' ', '\t', '\r', '\n':
					pos++
					continue
				default:
					break L3
				}
			}

			if pos >= dataLen {
				return fmt.Errorf("Invalid JSON")
			}

			switch data[pos] {
			case ',':
				pos++
				continue
			case '}':
				pos++
				break L2
			}

			var fieldName string
		L4:
			for {
				if pos >= dataLen {
					break L4
				}

				switch data[pos] {
				case ' ', '\t', '\r', '\n':
					pos++
					continue
				default:
					break L4
				}
			}

			if pos >= dataLen {
				return fmt.Errorf("Waited \", got EOF")
			}

			if data[pos] != '"' {
				return fmt.Errorf("Waited \", got %c", data[pos])
			}
			pos++
			r := pos
			for data[r] != '"' {
				r++
				if pos >= dataLen {
					return fmt.Errorf("Invalid JSON")
				}
			}
			fieldName = data[pos:r]
			pos = r + 1
		L5:
			for {
				if pos >= dataLen {
					break L5
				}

				switch data[pos] {
				case ' ', '\t', '\r', '\n':
					pos++
					continue
				default:
					break L5
				}
			}

			if pos >= dataLen {
				return fmt.Errorf("Waited :, got EOF")
			}

			if data[pos] != ':' {
				return fmt.Errorf("Waited :, got %c", data[pos])
			}

			pos++

			switch fieldName {
			case "IntField":
			L6:
				for {
					if pos >= dataLen {
						break L6
					}

					switch data[pos] {
					case ' ', '\t', '\r', '\n':
						pos++
						continue
					default:
						break L6
					}
				}

				if pos >= dataLen {
					return fmt.Errorf("Waited digit, got EOF")
				}
				r := pos
				for data[r] >= '0' && data[r] <= '9' {
					r++
					if pos >= dataLen {
						return fmt.Errorf("Waited digit, got EOF")
					}
				}
				if pos == r {
					return fmt.Errorf("Waited digit, got %c", data[pos])
				}
				res.IntField = 0
				for ; pos < r; pos++ {
					res.IntField *= 10
					res.IntField += int(data[pos] - '0')
				}
			case "StrField":
			L7:
				for {
					if pos >= dataLen {
						break L7
					}

					switch data[pos] {
					case ' ', '\t', '\r', '\n':
						pos++
						continue
					default:
						break L7
					}
				}

				if pos >= dataLen {
					return fmt.Errorf("Waited \", got EOF")
				}

				if data[pos] != '"' {
					return fmt.Errorf("Waited \", got %c", data[pos])
				}
				pos++
				r := pos
				for data[r] != '"' {
					r++
					if pos >= dataLen {
						return fmt.Errorf("Invalid JSON")
					}
				}
				res.StrField = data[pos:r]
				pos = r + 1
			case "RecursiveField":
				var tmp RecursiveStruct
				parser_RecursiveStruct(&tmp)
				res.RecursiveField = &tmp

			}

		}
		*ptr = res
		return nil
	}
	parser_RecursiveStruct(&res)
	*this = res
	return nil
}
