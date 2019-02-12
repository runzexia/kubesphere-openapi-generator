package install

import (
	"github.com/runzexia/kubesphere-crd-sample/pkg/apis/iam/v1alpha2"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
)

func Install(scheme *runtime.Scheme) {
	utilruntime.Must(v1alpha2.AddToScheme(scheme))
	utilruntime.Must(scheme.SetVersionPriority(v1alpha2.SchemeGroupVersion))
}
