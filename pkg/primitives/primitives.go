package Primitives

import (
	"embed"

	Types "github.com/Patch2PDF/GDTF-Mesh-Reader/pkg/MeshTypes"
	FileHandlers "github.com/Patch2PDF/GDTF-Mesh-Reader/pkg/file_handlers"
)

//go:embed assets/1.0/*.3ds
//go:embed assets/1.1/*.3ds
var modelFS embed.FS

var primitivePaths = map[string]string{
	"Undefined":       "",
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

var Primitives = map[string]*Types.Mesh{}

func LoadPrimitives() error {
	for primitiveType, path := range primitivePaths {
		switch primitiveType {
		case "Cube":
			mesh := NewCube()
			Primitives[primitiveType] = &mesh
		case "Cylinder":
			mesh := NewCylinder(10, true)
			Primitives[primitiveType] = &mesh
		case "Sphere":
			mesh := NewSphere(2)
			Primitives[primitiveType] = &mesh
		case "Pigtail", "Undefined":
			Primitives[primitiveType] = &Types.Mesh{} // empty mesh
		default:
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
	}
	return nil
}
