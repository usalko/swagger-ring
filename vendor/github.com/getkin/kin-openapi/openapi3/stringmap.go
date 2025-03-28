package openapi3

import "encoding/json"

// StringMap is a map[string]string that ignores the origin in the underlying json representation.
type StringMap map[string]string

// UnmarshalJSON sets StringMap to a copy of data.
func (stringMap *StringMap) UnmarshalJSON(data []byte) (err error) {
	*stringMap, _, err = unmarshalStringMap[string](data)
	return
}

// unmarshalStringMapP unmarshals given json into a map[string]*V
func unmarshalStringMapP[V any](data []byte) (map[string]*V, *Origin, error) {
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, nil, err
	}

	origin, err := popOrigin(m, originKey)
	if err != nil {
		return nil, nil, err
	}

	result := make(map[string]*V, len(m))
	for k, v := range m {
		value, err := deepCast[V](v)
		if err != nil {
			return nil, nil, err
		}
		result[k] = value
	}

	return result, origin, nil
}

// unmarshalStringMapP unmarshals given json into a map[string]*V
func unmarshalStringMapPSecuritySchemeRef(data []byte) (map[string]*SecuritySchemeRef, *Origin, error) {
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, nil, err
	}

	origin, err := popOrigin(m, originKey)
	if err != nil {
		return nil, nil, err
	}

	result := make(map[string]*SecuritySchemeRef, len(m))
	for k, v := range m {
		value, err := deepCastSecuritySchemeRef(v)
		if err != nil {
			return nil, nil, err
		}
		result[k] = value
	}

	return result, origin, nil
}

// unmarshalStringMapP unmarshals given json into a map[string]*V
func unmarshalStringMapPSchemaRef(data []byte) (map[string]*SchemaRef, *Origin, error) {
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, nil, err
	}

	origin, err := popOrigin(m, originKey)
	if err != nil {
		return nil, nil, err
	}

	result := make(map[string]*SchemaRef, len(m))
	for k, v := range m {
		value, err := deepCastSchemaRef(v)
		if err != nil {
			return nil, nil, err
		}
		result[k] = value
	}

	return result, origin, nil
}

// unmarshalStringMapP unmarshals given json into a map[string]*V
func unmarshalStringMapPResponseRef(data []byte) (map[string]*ResponseRef, *Origin, error) {
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, nil, err
	}

	origin, err := popOrigin(m, originKey)
	if err != nil {
		return nil, nil, err
	}

	result := make(map[string]*ResponseRef, len(m))
	for k, v := range m {
		value, err := deepCastResponseRef(v)
		if err != nil {
			return nil, nil, err
		}
		result[k] = value
	}

	return result, origin, nil
}

// unmarshalStringMapP unmarshals given json into a map[string]*V
func unmarshalStringMapPRequestBodyRef(data []byte) (map[string]*RequestBodyRef, *Origin, error) {
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, nil, err
	}

	origin, err := popOrigin(m, originKey)
	if err != nil {
		return nil, nil, err
	}

	result := make(map[string]*RequestBodyRef, len(m))
	for k, v := range m {
		value, err := deepCastRequestBodyRef(v)
		if err != nil {
			return nil, nil, err
		}
		result[k] = value
	}

	return result, origin, nil
}

// unmarshalStringMapP unmarshals given json into a map[string]*V
func unmarshalStringMapPParameterRef(data []byte) (map[string]*ParameterRef, *Origin, error) {
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, nil, err
	}

	origin, err := popOrigin(m, originKey)
	if err != nil {
		return nil, nil, err
	}

	result := make(map[string]*ParameterRef, len(m))
	for k, v := range m {
		value, err := deepCastParameterRef(v)
		if err != nil {
			return nil, nil, err
		}
		result[k] = value
	}

	return result, origin, nil
}

// unmarshalStringMapP unmarshals given json into a map[string]*V
func unmarshalStringMapPLinkRef(data []byte) (map[string]*LinkRef, *Origin, error) {
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, nil, err
	}

	origin, err := popOrigin(m, originKey)
	if err != nil {
		return nil, nil, err
	}

	result := make(map[string]*LinkRef, len(m))
	for k, v := range m {
		value, err := deepCastLinkRef(v)
		if err != nil {
			return nil, nil, err
		}
		result[k] = value
	}

	return result, origin, nil
}

// unmarshalStringMapP unmarshals given json into a map[string]*V // usalko:patch
func unmarshalStringMapPHeaderRef(data []byte) (map[string]*HeaderRef, *Origin, error) {
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, nil, err
	}

	origin, err := popOrigin(m, originKey)
	if err != nil {
		return nil, nil, err
	}

	result := make(map[string]*HeaderRef, len(m))
	for k, v := range m {
		value, err := deepCastHeaderRef(v)
		if err != nil {
			return nil, nil, err
		}
		result[k] = value
	}

	return result, origin, nil
}

// unmarshalStringMapP unmarshals given json into a map[string]*V // usalko:patch
func unmarshalStringMapPExampleRef(data []byte) (map[string]*ExampleRef, *Origin, error) {
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, nil, err
	}

	origin, err := popOrigin(m, originKey)
	if err != nil {
		return nil, nil, err
	}

	result := make(map[string]*ExampleRef, len(m))
	for k, v := range m {
		value, err := deepCastExampleRef(v)
		if err != nil {
			return nil, nil, err
		}
		result[k] = value
	}

	return result, origin, nil
}

// unmarshalStringMapP unmarshals given json into a map[string]*V // usalko:patch
func unmarshalStringMapPMediaType(data []byte) (map[string]*MediaType, *Origin, error) {
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, nil, err
	}

	origin, err := popOrigin(m, originKey)
	if err != nil {
		return nil, nil, err
	}

	result := make(map[string]*MediaType, len(m))
	for k, v := range m {
		value, err := deepCastMediaType(v)
		if err != nil {
			return nil, nil, err
		}
		result[k] = value
	}

	return result, origin, nil
}

// unmarshalStringMapP unmarshals given json into a map[string]*V // usalko:patch
func unmarshalStringMapPCallbackRef(data []byte) (map[string]*CallbackRef, *Origin, error) {
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, nil, err
	}

	origin, err := popOrigin(m, originKey)
	if err != nil {
		return nil, nil, err
	}

	result := make(map[string]*CallbackRef, len(m))
	for k, v := range m {
		// value, err := deepCast[CallbackRef](v) // usalko:patch
		value, err := deepCastCallbackRef(v)
		if err != nil {
			return nil, nil, err
		}
		result[k] = value
	}

	return result, origin, nil
}

// unmarshalStringMap unmarshals given json into a map[string]V
func unmarshalStringMap[V any](data []byte) (map[string]V, *Origin, error) {
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, nil, err
	}

	origin, err := popOrigin(m, originKey)
	if err != nil {
		return nil, nil, err
	}

	result := make(map[string]V, len(m))
	for k, v := range m {
		value, err := deepCast[V](v)
		if err != nil {
			return nil, nil, err
		}
		result[k] = *value
	}

	return result, origin, nil
}

// unmarshalStringMap unmarshals given json into a map[string]V
func unmarshalStringMapStringSlice(data []byte) (map[string][]string, *Origin, error) {
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, nil, err
	}

	origin, err := popOrigin(m, originKey)
	if err != nil {
		return nil, nil, err
	}

	result := make(map[string][]string, len(m))
	for k, v := range m {
		value, err := deepCast[[]string](v)
		if err != nil {
			return nil, nil, err
		}
		result[k] = *value
	}

	return result, origin, nil
}

// deepCast casts any value to a value of type V.
func deepCast[V any](value any) (*V, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}

	var result V
	if err = json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// deepCast casts any value to a value of type V.
func deepCastSecuritySchemeRef(value any) (*SecuritySchemeRef, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}

	var result SecuritySchemeRef
	if err = json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// deepCast casts any value to a value of type V.
func deepCastSchemaRef(value any) (*SchemaRef, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}

	var result SchemaRef
	if err = json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// deepCast casts any value to a value of type V.
func deepCastResponseRef(value any) (*ResponseRef, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}

	var result ResponseRef
	if err = json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// deepCast casts any value to a value of type V.
func deepCastRequestBodyRef(value any) (*RequestBodyRef, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}

	var result RequestBodyRef
	if err = json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// deepCast casts any value to a value of type V.
func deepCastParameterRef(value any) (*ParameterRef, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}

	var result ParameterRef
	if err = json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// deepCast casts any value to a value of type V.
func deepCastLinkRef(value any) (*LinkRef, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}

	var result LinkRef
	if err = json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// deepCast casts any value to a value of type V. // usalko:patch
func deepCastHeaderRef(value any) (*HeaderRef, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}

	var result HeaderRef
	if err = json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// deepCast casts any value to a value of type V.
func deepCastExampleRef(value any) (*ExampleRef, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}

	var result ExampleRef
	if err = json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// deepCast casts any value to a value of type V. // usalko:patch
func deepCastMediaType(value any) (*MediaType, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}

	var result MediaType
	if err = json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// deepCast casts any value to a value of type V. // usalko: patch
func deepCastCallbackRef(value any) (*CallbackRef, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}

	var result CallbackRef
	if err = json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// popOrigin removes the origin from the map and returns it.
func popOrigin(m map[string]any, key string) (*Origin, error) {
	if !IncludeOrigin {
		return nil, nil
	}

	origin, err := deepCast[Origin](m[key])
	if err != nil {
		return nil, err
	}
	delete(m, key)
	return origin, nil
}
