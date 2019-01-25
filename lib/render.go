package lib

import (
	"encoding/json"
	"fmt"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apiserver/pkg/registry/rest"
	"net"

	"github.com/go-openapi/spec"
	"github.com/golang/glog"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	apiopenapi "k8s.io/apiserver/pkg/endpoints/openapi"
	genericapiserver "k8s.io/apiserver/pkg/server"
	genericoptions "k8s.io/apiserver/pkg/server/options"

	"k8s.io/kube-openapi/pkg/builder"
	"k8s.io/kube-openapi/pkg/common"
)

type Config struct {
	Scheme *runtime.Scheme
	Codecs serializer.CodecFactory

	Info               spec.InfoProps
	OpenAPIDefinitions []common.GetOpenAPIDefinitions
	Resources          []schema.GroupVersionResource
	Mapper             *meta.DefaultRESTMapper
}

func (c *Config) GetOpenAPIDefinitions(ref common.ReferenceCallback) map[string]common.OpenAPIDefinition {
	out := map[string]common.OpenAPIDefinition{}
	for _, def := range c.OpenAPIDefinitions {
		for k, v := range def(ref) {
			out[k] = v
		}
	}
	return out
}

func RenderOpenAPISpec(cfg Config) (string, error) {
	// we need to add the options to empty v1
	// TODO fix the server code to avoid this
	metav1.AddToGroupVersion(cfg.Scheme, schema.GroupVersion{Version: "v1"})

	// TODO: keep the generic API server from wanting this
	unversioned := schema.GroupVersion{Group: "", Version: "v1"}
	cfg.Scheme.AddUnversionedTypes(unversioned,
		&metav1.Status{},
		&metav1.APIVersions{},
		&metav1.APIGroupList{},
		&metav1.APIGroup{},
		&metav1.APIResourceList{},
	)

	recommendedOptions := genericoptions.NewRecommendedOptions("/registry/foo.com", cfg.Codecs.LegacyCodec())
	recommendedOptions.SecureServing.BindPort = 8443
	recommendedOptions.Etcd = nil
	recommendedOptions.Authentication = nil
	recommendedOptions.Authorization = nil
	recommendedOptions.CoreAPI = nil
	recommendedOptions.Admission = nil

	// TODO have a "real" external address
	if err := recommendedOptions.SecureServing.MaybeDefaultWithSelfSignedCerts("localhost", nil, []net.IP{net.ParseIP("127.0.0.1")}); err != nil {
		glog.Fatal(fmt.Errorf("error creating self-signed certificates: %v", err))
	}

	serverConfig := genericapiserver.NewRecommendedConfig(cfg.Codecs)

	if err := recommendedOptions.ApplyTo(serverConfig, cfg.Scheme); err != nil {
		glog.Fatal(err)
		return "", err
	}
	serverConfig.OpenAPIConfig = genericapiserver.DefaultOpenAPIConfig(cfg.GetOpenAPIDefinitions, apiopenapi.NewDefinitionNamer(cfg.Scheme))
	serverConfig.OpenAPIConfig.Info.InfoProps = cfg.Info

	genericServer, err := serverConfig.Complete().New("stash-server", genericapiserver.NewEmptyDelegate()) // completion is done in Complete, no need for a second time
	if err != nil {
		glog.Fatal(err)
		return "", err
	}

	{
		// api router map
		table := map[schema.GroupVersion]map[string]ResourceInfo{}
		for _, gvr := range cfg.Resources {
			var resmap map[string]ResourceInfo
			// init ResourceInfo map
			if m, found := table[gvr.GroupVersion()]; found {
				resmap = m
			} else {
				resmap = map[string]ResourceInfo{}
				table[gvr.GroupVersion()] = resmap
			}

			gvk, err := cfg.Mapper.KindFor(gvr)
			if err != nil {
				glog.Fatal(err)
				return "", err
			}
			obj, err := cfg.Scheme.New(gvk)
			if err != nil {
				return "", err
			}
			list, err := cfg.Scheme.New(gvk.GroupVersion().WithKind(gvk.Kind + "List"))
			if err != nil {
				glog.Fatal(err)
				return "", err
			}

			resmap[gvr.Resource] = ResourceInfo{
				gvk:  gvk,
				obj:  obj,
				list: list,
			}
		}

		for gv, resmap := range table {
			apiGroupInfo := genericapiserver.NewDefaultAPIGroupInfo(gv.Group, cfg.Scheme, metav1.ParameterCodec, cfg.Codecs)
			apiGroupInfo.MetaGroupVersion = &gv
			storage := map[string]rest.Storage{}
			for r, stuff := range resmap {
				storage[r] = NewREST(stuff)
			}
			apiGroupInfo.VersionedResourcesStorageMap[gv.Version] = storage

			if err := genericServer.InstallAPIGroup(&apiGroupInfo); err != nil {
				glog.Fatal(err)
				return "", err
			}
		}
	}

	spec, err := builder.BuildOpenAPISpec(genericServer.Handler.GoRestfulContainer.RegisteredWebServices(), serverConfig.OpenAPIConfig)
	if err != nil {
		glog.Fatal(err)
		return "", err
	}
	data, err := json.MarshalIndent(spec, "", "  ")
	if err != nil {
		glog.Fatal(err)
		return "", err
	}
	return string(data), nil
}