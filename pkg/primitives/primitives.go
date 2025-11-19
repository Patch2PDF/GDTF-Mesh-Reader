package Primitives

import (
	"embed"

	Types "github.com/Patch2PDF/GDTF-Mesh-Reader/pkg/MeshTypes"
	FileHandlers "github.com/Patch2PDF/GDTF-Mesh-Reader/pkg/file_handlers"
)

// embedded primitives
//
//go:embed assets/1.0/*.3ds
//go:embed assets/1.1/*.3ds
var modelFS embed.FS

// TODO: add primitives for Cube, Cylinder and Sphere, possibly Pigtail (just a scaled cube?)

// paths to embedded primitives
var primitivePaths = map[string]string{
	"Cube":            "",
	"Cylinder":        "",
	"Sphere":          "",
	"Base":            "assets/1.0/primitivetype_base.3ds",
	"Yoke":            "assets/1.0/primitivetype_yoke.3ds",
	"Head":            "assets/1.0/primitivetype_head.3ds",
	"Scanner":         "assets/1.0/primitivetype_scanner.3ds",
	"Conventional":    "assets/1.0/primitivetype_conventional.3ds",
	"Pigtail":         "",
	"Base1_1":         "assets/1.1/primitivetype_base_1.1.3ds",
	"Scanner1_1":      "assets/1.1/primitivetype_scanner_1.1.3ds",
	"Conventional1_1": "assets/1.1/primitivetype_conventional_1.1.3ds",
}

// Map of GDTF Primitive meshes.
// Requries LoadPrimitives() to be executed before accessing
var Primitives = map[string]*Types.Mesh{}

// Load GDTF Primitive meshes into Primitives Map.
// Required to be run on app startup, before accessing Primitives Map
func LoadPrimitives() error {
	for primitiveType, path := range primitivePaths {
		if path == "" {
			continue
		}
		data, err := modelFS.ReadFile(path)
		if err != nil {
			return err
		}
		Primitives[primitiveType], err = FileHandlers.Load3DS(&data, nil)
		if err != nil {
			return err
		}
	}
	return nil
}
