package schemas

import "database/sql/driver"

type MyString string

const MyStringNull MyString = "\x00"

func (s MyString) Value() (driver.Value, error) {
	if s == MyStringNull {
		return nil, nil
	}
	return []byte(s), nil
}

func (s *MyString) Scan(src interface{}) error {
	switch v := src.(type) {
	case string:
		*s = MyString(v)
	case []byte:
		*s = MyString(v)
	case nil:
		*s = MyStringNull
	}
	return nil
}
