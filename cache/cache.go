package cache

import "errors"

type Cache interface {
	Set(k, v interface{}) error
	Remove(k interface{}) error
	Get(k interface{}) (interface{}, error)
	Len() int
}

var ErrorMissingKey = errors.New("missing key")
