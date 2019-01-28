package main

import (
	"github.com/runzexia/kubesphere-openapi-generator/lib"
	"go/build"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/api/meta"
	"os"
	"path/filepath"

	"github.com/go-openapi/spec"
	"github.com/golang/glog"
	devopsinstall "github.com/runzexia/kubesphere-crd-sample/pkg/apis/devops/install"
	devopsv1alpha2 "github.com/runzexia/kubesphere-crd-sample/pkg/apis/devops/v1alpha2"
	iaminstall "github.com/runzexia/kubesphere-crd-sample/pkg/apis/iam/install"
	iamv1alpha2 "github.com/runzexia/kubesphere-crd-sample/pkg/apis/iam/v1alpha2"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/kube-openapi/pkg/common"
)

func main() {

	var (
		Scheme = runtime.NewScheme()
		Codecs = serializer.NewCodecFactory(Scheme)
	)

	devopsinstall.Install(Scheme)
	iaminstall.Install(Scheme)
	mapper := meta.NewDefaultRESTMapper(nil)
	mapper.AddSpecific(devopsv1alpha2.SchemeGroupVersion.WithKind(devopsv1alpha2.ResourceKindDevOpsProject),
		devopsv1alpha2.SchemeGroupVersion.WithResource(devopsv1alpha2.ResourcePluralDevOpsProject),
		devopsv1alpha2.SchemeGroupVersion.WithResource(devopsv1alpha2.ResourceSingularDevOpsProject), meta.RESTScopeRoot)
	mapper.AddSpecific(iamv1alpha2.SchemeGroupVersion.WithKind(iamv1alpha2.ResourceKindWorkspace),
		iamv1alpha2.SchemeGroupVersion.WithResource(iamv1alpha2.ResourcePluralWorkspace),
		iamv1alpha2.SchemeGroupVersion.WithResource(iamv1alpha2.ResourceSingularWorkspace), meta.RESTScopeRoot)

	spec, err := lib.RenderOpenAPISpec(lib.Config{
		Scheme: Scheme,
		Codecs: Codecs,
		Info: spec.InfoProps{
			Title:   "KubeSphere Advanced",
			Version: "v2.0.0",
			Contact: &spec.ContactInfo{
				Name:  "KubeSphere",
				URL:   "https://kubesphere.io/",
				Email: "kubesphere@yunify.com",
			},
			License: &spec.License{
				Name: "Apache 2.0",
				URL:  "https://www.apache.org/licenses/LICENSE-2.0.html",
			},
		},
		OpenAPIDefinitions: []common.GetOpenAPIDefinitions{
			devopsv1alpha2.GetOpenAPIDefinitions,
			iamv1alpha2.GetOpenAPIDefinitions,
		},
		Resources: []schema.GroupVersionResource{
			devopsv1alpha2.SchemeGroupVersion.WithResource(devopsv1alpha2.ResourcePluralDevOpsProject),
			iamv1alpha2.SchemeGroupVersion.WithResource(iamv1alpha2.ResourcePluralWorkspace),
		},
		Mapper: mapper,
	})
	if err != nil {
		glog.Fatal(err)
	}

	filename := build.Default.GOPATH + "/src/github.com/runzexia/kubesphere-crd-sample/api/openapi-spec/swagger.json"
	err = os.MkdirAll(filepath.Dir(filename), 0755)
	if err != nil {
		glog.Fatal(err)
	}
	err = ioutil.WriteFile(filename, []byte(spec), 0644)
	if err != nil {
		glog.Fatal(err)
	}
}
