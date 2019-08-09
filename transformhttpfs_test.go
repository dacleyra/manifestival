package manifestival_test

import (
	"os"
	"testing"
	"net/http"
	. "github.com/dacleyra/manifestival"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"github.com/shurcooL/vfsgen"
)

func TestTransformHTTPFileSystem(t *testing.T) {
	
	var fs http.FileSystem = http.Dir("testdata/tree")
	
	err := vfsgen.Generate(fs, vfsgen.Options{})
    if err != nil {
    	t.Errorf("vfsgen.Generate() = %v, wanted no error", err)
    }
	
	f, err := NewManifestHTTPFileSystem(fs, "", true, nil)
	if err != nil {
		t.Errorf("NewManifestHTTPFileSystem(() = %v, wanted no error", err)
	}

	actual := f.Resources
	if len(actual) != 5 {
		t.Errorf("Failed to read all resources: %s", actual)
	}
	f.Transform(func(u *unstructured.Unstructured) error {
		if u.GetKind() == "B" {
			u.SetResourceVersion("69")
		}
		return nil
	})
	transformed := f.Resources
	// Ensure all transformed have a version and B kind
	for _, spec := range transformed {
		if spec.GetResourceVersion() != "69" && spec.GetKind() == "B" {
			t.Errorf("The transform didn't work: %s", transformed)
		}
	}
	// Ensure we didn't change the previous resources
	for _, spec := range actual {
		if spec.GetResourceVersion() != "" {
			t.Errorf("The transform shouldn't affect previous resources: %s", actual)
		}
	}
	
	// remove vfs
	os.Remove("assets_vfsdata.go")
}

func TestTransformHTTPFileSystemCombo(t *testing.T) {
	
	var fs http.FileSystem = http.Dir("testdata/tree")
	
	err := vfsgen.Generate(fs, vfsgen.Options{})
    if err != nil {
    	t.Errorf("vfsgen.Generate() = %v, wanted no error", err)
    }
	
	f, err := NewManifestHTTPFileSystem(fs, "", true, nil)
	if err != nil {
		t.Errorf("NewManifestHTTPFileSystem(() = %v, wanted no error", err)
	}
	
	if len(f.Resources) != 5 {
		t.Errorf("Failed to read all resources: %s", f.Resources)
	}
	fn1 := func(u *unstructured.Unstructured) error {
		if u.GetKind() == "B" {
			u.SetResourceVersion("69")
		}
		return nil
	}
	fn2 := func(u *unstructured.Unstructured) error {
		if u.GetName() == "bar" {
			u.SetResourceVersion("42")
		}
		return nil
	}
	if err := f.Transform(fn1, fn2); err != nil {
		t.Error(err)
	}
	for _, x := range f.Resources {
		if x.GetName() == "bar" && x.GetResourceVersion() != "42" {
			t.Errorf("Failed to transform bar by combo: %s", x)
		}
		if x.GetName() == "B" && x.GetResourceVersion() != "69" {
			t.Errorf("Failed to transform B by combo: %s", x)
		}
	}
	
	// remove vfs
	os.Remove("assets_vfsdata.go")
}
