package testprovider

import (
	"unicode"

	testproviderdata "github.com/pulumi/pulumi-terraform-bridge/v3/internal/testprovider"
	"github.com/pulumi/pulumi-terraform-bridge/v3/pkg/tfbridge"
	shimv2 "github.com/pulumi/pulumi-terraform-bridge/v3/pkg/tfshim/sdk-v2"
	"github.com/pulumi/pulumi-terraform-bridge/x/muxer"
	"github.com/pulumi/pulumi/pkg/v3/codegen/schema"
	"github.com/pulumi/pulumi/pkg/v3/resource/provider"
	"github.com/pulumi/pulumi/sdk/v3/go/common/tokens"
	pulumirpc "github.com/pulumi/pulumi/sdk/v3/proto/go"
)

func ProviderMiniMuxedReplace() tfbridge.ProviderInfo {
	minimuxedPkg := "minimuxed"
	minimuxedMod := "index"

	minimuxedMember := func(mod string, mem string) tokens.ModuleMember {
		return tokens.ModuleMember(minimuxedPkg + ":" + mod + ":" + mem)
	}

	minimuxedType := func(mod string, typ string) tokens.Type {
		return tokens.Type(minimuxedMember(mod, typ))
	}

	minimuxedResource := func(mod string, res string) tokens.Type {
		fn := string(unicode.ToLower(rune(res[0]))) + res[1:]
		return minimuxedType(mod+"/"+fn, res)
	}

	return tfbridge.ProviderInfo{
		P:           shimv2.NewProvider(testproviderdata.ProviderMiniMuxed()),
		Name:        "minimuxed",
		Description: "A Pulumi package to safely use minimuxed resources in Pulumi programs.",
		Keywords:    []string{"pulumi", "minimuxed"},
		License:     "Apache-2.0",
		Homepage:    "https://pulumi.io",
		Repository:  "https://github.com/pulumi/pulumi-minimuxed",
		Resources: map[string]*tfbridge.ResourceInfo{
			"minimuxed_integer": {Tok: minimuxedResource(minimuxedMod, "MinimuxedInteger")},
		},
		MuxWith: []muxer.Provider{
			newMuxReplaceProvider(),
		},
	}
}

func newMuxReplaceProvider() muxer.Provider {
	return &muxReplaceProvider{
		packageSchema: schema.PackageSpec{
			Name: "minimuxed",
			Resources: map[string]schema.ResourceSpec{
				"minimuxed:index/minimuxedInteger:MinimuxedInteger": {
					ObjectTypeSpec: schema.ObjectTypeSpec{
						Properties: map[string]schema.PropertySpec{
							"max": {
								TypeSpec: schema.TypeSpec{
									Type: "integer",
								},
							},
							"min": {
								TypeSpec: schema.TypeSpec{
									Type: "integer",
								},
							},
							"result": {
								TypeSpec: schema.TypeSpec{
									Type: "integer",
								},
							},
						},
						Required: []string{
							"max",
							"min",
							"result",
						},
					},
					InputProperties: map[string]schema.PropertySpec{
						"max": {
							TypeSpec: schema.TypeSpec{
								Type: "integer",
							},
							WillReplaceOnChanges: true,
						},
						"min": {
							TypeSpec: schema.TypeSpec{
								Type: "integer",
							},
							WillReplaceOnChanges: true,
						},
					},
					RequiredInputs: []string{
						"max",
						"min",
					},
					StateInputs: &schema.ObjectTypeSpec{
						Properties: map[string]schema.PropertySpec{
							"max": {
								TypeSpec: schema.TypeSpec{
									Type: "integer",
								},
								WillReplaceOnChanges: true,
							},
							"min": {
								TypeSpec: schema.TypeSpec{
									Type: "integer",
								},
								WillReplaceOnChanges: true,
							},
							"result": {
								TypeSpec: schema.TypeSpec{
									Type: "integer",
								},
							},
						},
						Type: "object",
					},
				},
			},
		},
	}
}

type muxReplaceProvider struct {
	pulumirpc.UnimplementedResourceProviderServer

	packageSchema schema.PackageSpec
}

func (m *muxReplaceProvider) GetSpec() (schema.PackageSpec, error) {
	return m.packageSchema, nil
}

func (m *muxReplaceProvider) GetInstance(*provider.HostClient) (pulumirpc.ResourceProviderServer, error) {
	return m, nil
}
