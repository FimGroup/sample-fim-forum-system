package forumcore

import (
	"errors"
	"fmt"

	"github.com/FimGroup/fim/fimapi/pluginapi"
	"github.com/FimGroup/fim/fimapi/providers"
	"github.com/FimGroup/fim/fimapi/rule"
)

type CustomFunctions struct {
	_logger providers.Logger
}

func (c *CustomFunctions) FnPrintObject(params []interface{}) (pluginapi.Fn, error) {
	key := params[0].(string)
	if !rule.ValidateFullPath(key) {
		return nil, errors.New("invalid path:" + key)
	}
	paths := rule.SplitFullPath(key)
	return func(m pluginapi.Model) error {
		//FIXME have to handle object/array properly
		o := m.GetFieldUnsafe0(paths)
		if c._logger.IsInfoEnabled() {
			c._logger.InfoF("print object=[%s] of key=[%s]", fmt.Sprint(o), key)
		}
		return nil
	}, nil
}

func (c *CustomFunctions) FnPanic(params []interface{}) (pluginapi.Fn, error) {
	key := params[0].(string)
	return func(m pluginapi.Model) error {
		panic(errors.New(key))
	}, nil
}
