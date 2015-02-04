package resource

import (
	"fmt"
	"reflect"

	"github.com/jinzhu/gorm"
	"github.com/qor/qor"
)

type Resource struct {
	Name         string
	primaryField *gorm.Field
	Value        interface{}
	Searcher     func(interface{}, *qor.Context) error
	Finder       func(interface{}, *MetaValues, *qor.Context) error
	Saver        func(interface{}, *qor.Context) error
	Deleter      func(interface{}, *qor.Context) error
	validators   []func(interface{}, *MetaValues, *qor.Context) error
	processors   []func(interface{}, *MetaValues, *qor.Context) error
}

type Resourcer interface {
	GetResource() *Resource
	GetMetas(...[]string) []Metaor
	CallSearcher(interface{}, *qor.Context) error
	CallFinder(interface{}, *MetaValues, *qor.Context) error
	CallSaver(interface{}, *qor.Context) error
	CallDeleter(interface{}, *qor.Context) error
	NewSlice() interface{}
	NewStruct() interface{}
}

func New(value interface{}, names ...string) *Resource {
	name := reflect.Indirect(reflect.ValueOf(value)).Type().Name()
	for _, n := range names {
		name = n
	}

	return &Resource{Value: value, Name: name}
}

func (res *Resource) GetResource() *Resource {
	return res
}

func (res *Resource) PrimaryField() *gorm.Field {
	if res.primaryField == nil {
		scope := gorm.Scope{Value: res.Value}
		res.primaryField = scope.PrimaryKeyField()
	}
	return res.primaryField
}

func (res *Resource) PrimaryFieldName() (name string) {
	field := res.PrimaryField()
	if field != nil {
		name = field.Name
	}
	return
}

func (res *Resource) CallSearcher(result interface{}, context *qor.Context) error {
	if res.Searcher != nil {
		return res.Searcher(result, context)
	} else {
		return context.GetDB().Order(fmt.Sprintf("%v DESC", res.PrimaryField().DBName)).Find(result).Error
	}
}

func (res *Resource) CallSaver(result interface{}, context *qor.Context) error {
	if res.Saver != nil {
		return res.Saver(result, context)
	} else {
		return context.GetDB().Save(result).Error
	}
}

func (res *Resource) CallDeleter(result interface{}, context *qor.Context) error {
	if res.Deleter != nil {
		return res.Deleter(result, context)
	} else {
		db := context.GetDB().Delete(result, context.ResourceID)
		if db.Error != nil {
			return db.Error
		} else if db.RowsAffected == 0 {
			return gorm.RecordNotFound
		}
		return nil
	}
}

func (res *Resource) CallFinder(result interface{}, metaValues *MetaValues, context *qor.Context) error {
	if res.Finder != nil {
		return res.Finder(result, metaValues, context)
	} else {
		if metaValues == nil {
			return context.GetDB().First(result, context.ResourceID).Error
		}
		return nil
	}
}

func (res *Resource) AddValidator(fc func(interface{}, *MetaValues, *qor.Context) error) {
	res.validators = append(res.validators, fc)
}

func (res *Resource) AddProcessor(fc func(interface{}, *MetaValues, *qor.Context) error) {
	res.processors = append(res.processors, fc)
}

func (res *Resource) NewSlice() interface{} {
	sliceType := reflect.SliceOf(reflect.ValueOf(res.Value).Type())
	slice := reflect.MakeSlice(sliceType, 0, 0)
	slicePtr := reflect.New(sliceType)
	slicePtr.Elem().Set(slice)
	return slicePtr.Interface()
}

func (res *Resource) NewStruct() interface{} {
	return reflect.New(reflect.Indirect(reflect.ValueOf(res.Value)).Type()).Interface()
}

func (res *Resource) GetMetas(...[]string) []Metaor {
	panic("not defined")
}
