// Code generated by mockery v1.0.0. DO NOT EDIT.

package cache

import mock "github.com/stretchr/testify/mock"
import types "github.com/aquasecurity/fanal/types"

// MockImageCache is an autogenerated mock type for the ImageCache type
type MockImageCache struct {
	mock.Mock
}

type ImageCacheMissingLayersArgs struct {
	ImageID          string
	ImageIDAnything  bool
	LayerIDs         []string
	LayerIDsAnything bool
}

type ImageCacheMissingLayersReturns struct {
	MissingImage    bool
	MissingLayerIDs []string
	Err             error
}

type ImageCacheMissingLayersExpectation struct {
	Args    ImageCacheMissingLayersArgs
	Returns ImageCacheMissingLayersReturns
}

func (_m *MockImageCache) ApplyMissingLayersExpectation(e ImageCacheMissingLayersExpectation) {
	var args []interface{}
	if e.Args.ImageIDAnything {
		args = append(args, mock.Anything)
	} else {
		args = append(args, e.Args.ImageID)
	}
	if e.Args.LayerIDsAnything {
		args = append(args, mock.Anything)
	} else {
		args = append(args, e.Args.LayerIDs)
	}
	_m.On("MissingLayers", args...).Return(e.Returns.MissingImage, e.Returns.MissingLayerIDs, e.Returns.Err)
}

func (_m *MockImageCache) ApplyMissingLayersExpectations(expectations []ImageCacheMissingLayersExpectation) {
	for _, e := range expectations {
		_m.ApplyMissingLayersExpectation(e)
	}
}

// MissingLayers provides a mock function with given fields: imageID, layerIDs
func (_m *MockImageCache) MissingLayers(imageID string, layerIDs []string) (bool, []string, error) {
	ret := _m.Called(imageID, layerIDs)

	var r0 bool
	if rf, ok := ret.Get(0).(func(string, []string) bool); ok {
		r0 = rf(imageID, layerIDs)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 []string
	if rf, ok := ret.Get(1).(func(string, []string) []string); ok {
		r1 = rf(imageID, layerIDs)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).([]string)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(string, []string) error); ok {
		r2 = rf(imageID, layerIDs)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

type ImageCachePutImageArgs struct {
	ImageID           string
	ImageIDAnything   bool
	ImageInfo         types.ImageInfo
	ImageInfoAnything bool
}

type ImageCachePutImageReturns struct {
	Err error
}

type ImageCachePutImageExpectation struct {
	Args    ImageCachePutImageArgs
	Returns ImageCachePutImageReturns
}

func (_m *MockImageCache) ApplyPutImageExpectation(e ImageCachePutImageExpectation) {
	var args []interface{}
	if e.Args.ImageIDAnything {
		args = append(args, mock.Anything)
	} else {
		args = append(args, e.Args.ImageID)
	}
	if e.Args.ImageInfoAnything {
		args = append(args, mock.Anything)
	} else {
		args = append(args, e.Args.ImageInfo)
	}
	_m.On("PutImage", args...).Return(e.Returns.Err)
}

func (_m *MockImageCache) ApplyPutImageExpectations(expectations []ImageCachePutImageExpectation) {
	for _, e := range expectations {
		_m.ApplyPutImageExpectation(e)
	}
}

// PutImage provides a mock function with given fields: imageID, imageInfo
func (_m *MockImageCache) PutImage(imageID string, imageInfo types.ImageInfo) error {
	ret := _m.Called(imageID, imageInfo)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, types.ImageInfo) error); ok {
		r0 = rf(imageID, imageInfo)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type ImageCachePutLayerArgs struct {
	LayerID                     string
	LayerIDAnything             bool
	DecompressedLayerID         string
	DecompressedLayerIDAnything bool
	LayerInfo                   types.LayerInfo
	LayerInfoAnything           bool
}

type ImageCachePutLayerReturns struct {
	Err error
}

type ImageCachePutLayerExpectation struct {
	Args    ImageCachePutLayerArgs
	Returns ImageCachePutLayerReturns
}

func (_m *MockImageCache) ApplyPutLayerExpectation(e ImageCachePutLayerExpectation) {
	var args []interface{}
	if e.Args.LayerIDAnything {
		args = append(args, mock.Anything)
	} else {
		args = append(args, e.Args.LayerID)
	}
	if e.Args.DecompressedLayerIDAnything {
		args = append(args, mock.Anything)
	} else {
		args = append(args, e.Args.DecompressedLayerID)
	}
	if e.Args.LayerInfoAnything {
		args = append(args, mock.Anything)
	} else {
		args = append(args, e.Args.LayerInfo)
	}
	_m.On("PutLayer", args...).Return(e.Returns.Err)
}

func (_m *MockImageCache) ApplyPutLayerExpectations(expectations []ImageCachePutLayerExpectation) {
	for _, e := range expectations {
		_m.ApplyPutLayerExpectation(e)
	}
}

// PutLayer provides a mock function with given fields: layerID, decompressedLayerID, layerInfo
func (_m *MockImageCache) PutLayer(layerID string, decompressedLayerID string, layerInfo types.LayerInfo) error {
	ret := _m.Called(layerID, decompressedLayerID, layerInfo)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, types.LayerInfo) error); ok {
		r0 = rf(layerID, decompressedLayerID, layerInfo)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
