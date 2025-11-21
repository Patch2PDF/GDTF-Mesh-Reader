package GDTFMeshReader

import (
	"fmt"
	"io"
	"path/filepath"

	"github.com/Patch2PDF/GDTF-Mesh-Reader/pkg/MeshTypes"
	FileHandlers "github.com/Patch2PDF/GDTF-Mesh-Reader/pkg/file_handlers"
	Primitives "github.com/Patch2PDF/GDTF-Mesh-Reader/pkg/primitives"
)

type ModelReaderConf struct {
	File          io.Reader
	Filename      *string
	PrimitiveType string
}

// Load Primitive Models
func LoadPrimitives() error {
	return Primitives.LoadPrimitives()
}

// Get Model by mesh file or PrimitiveType.
//
// Note: requires LoadPrimitives() to be run beforehand if you want to get a primitive
func GetModel(conf ModelReaderConf, desiredSize MeshTypes.Vector) (*MeshTypes.Mesh, error) {
	var mesh *MeshTypes.Mesh

	if conf.PrimitiveType == "Undefined" && conf.File != nil && conf.Filename != nil && *conf.Filename != "" {
		filetype := filepath.Ext(*conf.Filename)
		switch filetype {
		case ".gltf", ".glb":
			meshes, err := FileHandlers.LoadGLTF(conf.File, desiredSize)
			if err != nil {
				return nil, err
			}
			mesh = meshes[0]
		case ".3ds":
			data, err := io.ReadAll(conf.File)
			if err != nil {
				return nil, err
			}
			mesh, err = FileHandlers.Load3DS(&data, &desiredSize)
			if err != nil {
				return nil, err
			}
		default:
			return nil, fmt.Errorf("unknown model type %s", filetype)
		}
	} else if conf.PrimitiveType != "Undefined" && Primitives.Primitives[conf.PrimitiveType] != nil {
		tempMesh := Primitives.Primitives[conf.PrimitiveType].Copy()
		mesh = &tempMesh
		mesh.ScaleToDimensions(&desiredSize)
	}

	return mesh, nil
}
