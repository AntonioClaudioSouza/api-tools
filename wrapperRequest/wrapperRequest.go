package wrapperRequest

import (
	"encoding/json"
	"fmt"
	"strings"
)

type WrapperDataRequest struct {
	Data     map[string]interface{}
	BookMark string
	PageSize int32
	Fields   []WrapperField
	DataOut  map[string]interface{}
}

type WrapperField struct {
	Required bool
	FieldAPI string
	TagCC    string
	Asset    string
}

func (w *WrapperDataRequest) addField(fieldName string, value interface{}) {
	w.DataOut[fieldName] = value
}

func (w *WrapperDataRequest) addAsset(fieldNameChaincode, assetNameChaincode, fieldNameRequest string) {

	value, has := w.Data[fieldNameRequest]
	if !has {
		return
	}

	if strings.Contains(assetNameChaincode, "[]->") {
		nameAsset := strings.Replace(assetNameChaincode, "[]->", "", 1)

		// ** Check is content is not nil
		arrayString := value.([]interface{})
		arrayAsset := []map[string]interface{}{}
		for _, v := range arrayString {

			if v.(string) == "" {
				return
			}

			arrayAsset = append(arrayAsset, map[string]interface{}{
				"@assetType": nameAsset,
				"@key":       v,
			})
		}

		w.DataOut[fieldNameChaincode] = arrayAsset
	} else {

		if value.(string) == "" {
			return
		}

		w.DataOut[fieldNameChaincode] = map[string]interface{}{
			"@assetType": assetNameChaincode,
			"@key":       value,
		}
	}
}

func (w *WrapperDataRequest) GetJson() ([]byte, error) {

	w.DataOut = map[string]interface{}{}

	for _, f := range w.Fields {

		// ** Check required
		if f.Required {
			_, has := w.Data[f.FieldAPI]
			if !has {
				return nil, fmt.Errorf("%s failed is required field not found", f.FieldAPI)
			}
		}

		if f.Asset != "" {
			w.addAsset(f.TagCC, f.Asset, f.FieldAPI)
		} else {

			valueField, has := w.Data[f.FieldAPI]
			if has {
				w.addField(f.TagCC, valueField)
			}
		}
	}

	if len(w.BookMark) > 0 {
		w.DataOut["bookmark"] = w.BookMark
	}

	if w.PageSize > 0 {
		w.DataOut["pageSize"] = w.PageSize
	}

	pJson, err := json.Marshal(w.DataOut)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	return pJson, nil
}
