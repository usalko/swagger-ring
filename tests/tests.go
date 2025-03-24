package tests

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/getkin/kin-openapi/openapi3"
)

func main() {
	// Загрузка первого swagger.json
	swagger1, err := loadSwagger("swagger1.json")
	if err != nil {
		log.Fatalf("Ошибка загрузки swagger1.json: %v", err)
	}

	// Загрузка второго swagger.json
	swagger2, err := loadSwagger("swagger2.json")
	if err != nil {
		log.Fatalf("Ошибка загрузки swagger2.json: %v", err)
	}

	// Слияние двух спецификаций
	mergedSwagger := mergeSwagger(swagger1, swagger2)

	// Сохранение результата в новый файл
	err = saveSwagger("merged_swagger.json", mergedSwagger)
	if err != nil {
		log.Fatalf("Ошибка сохранения merged_swagger.json: %v", err)
	}

	fmt.Println("Слияние завершено. Результат сохранен в merged_swagger.json")
}

func loadSwagger(filename string) (*openapi3.T, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	loader := openapi3.NewLoader()
	swagger, err := loader.LoadFromData(data)
	if err != nil {
		return nil, err
	}

	return swagger, nil
}

func mergeSwagger(swagger1, swagger2 *openapi3.T) *openapi3.T {
	mergedPaths := openapi3.Paths{}
	// Копируем базовые поля из первой спецификации
	merged := &openapi3.T{
		OpenAPI:      swagger1.OpenAPI,
		Info:         swagger1.Info,
		Servers:      swagger1.Servers,
		Paths:        &mergedPaths,
		Components:   swagger1.Components,
		Security:     swagger1.Security,
		Tags:         swagger1.Tags,
		ExternalDocs: swagger1.ExternalDocs,
	}

	// Добавляем пути из первой спецификации
	for path, pathItem := range swagger1.Paths.Map() {
		merged.Paths.Set(path, pathItem)
	}

	// Добавляем пути из второй спецификации
	for path, pathItem := range swagger2.Paths.Map() {
		if existingPathItem := merged.Paths.Find(path); existingPathItem != nil {
			// Если путь уже существует, объединяем операции
			for method, operation := range pathItem.Operations() {
				existingPathItem.SetOperation(method, operation)
			}
		} else {
			// Если путь не существует, добавляем его
			merged.Paths.Set(path, pathItem)
		}
	}

	// Объединяем компоненты
	if swagger2.Components != nil {
		if merged.Components == nil {
			merged.Components = &openapi3.Components{}
		}

		// Объединяем схемы
		for name, schema := range swagger2.Components.Schemas {
			merged.Components.Schemas[name] = schema
		}

		// Объединяем параметры
		for name, parameter := range swagger2.Components.Parameters {
			merged.Components.Parameters[name] = parameter
		}

		// Объединяем ответы
		for name, response := range swagger2.Components.Responses {
			merged.Components.Responses[name] = response
		}

		// Объединяем примеры
		for name, example := range swagger2.Components.Examples {
			merged.Components.Examples[name] = example
		}

		// Объединяем запросы
		for name, requestBody := range swagger2.Components.RequestBodies {
			merged.Components.RequestBodies[name] = requestBody
		}

		// Объединяем заголовки
		for name, header := range swagger2.Components.Headers {
			merged.Components.Headers[name] = header
		}

		// Объединяем security схемы
		for name, securityScheme := range swagger2.Components.SecuritySchemes {
			merged.Components.SecuritySchemes[name] = securityScheme
		}

		// Объединяем ссылки
		for name, link := range swagger2.Components.Links {
			merged.Components.Links[name] = link
		}

		// Объединяем обратные вызовы
		for name, callback := range swagger2.Components.Callbacks {
			merged.Components.Callbacks[name] = callback
		}
	}

	return merged
}

func saveSwagger(filename string, swagger *openapi3.T) error {
	data, err := json.MarshalIndent(swagger, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, data, 0644)
}
