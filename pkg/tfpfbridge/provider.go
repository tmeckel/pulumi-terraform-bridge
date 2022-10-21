// Copyright 2016-2022, Pulumi Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tfbridge

import (
	"context"
	"fmt"
	"sync"

	"github.com/blang/semver"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	tfsdkprovider "github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	tfsdkresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource/plugin"
	"github.com/pulumi/pulumi/sdk/v3/go/common/tokens"
	"github.com/pulumi/pulumi/sdk/v3/go/common/workspace"
	pulumirpc "github.com/pulumi/pulumi/sdk/v3/proto/go"

	"github.com/pulumi/pulumi-terraform-bridge/pkg/tfpfbridge/info"
)

// Provider implements the Pulumi resource provider operations for any
// Terraform plugin built with Terraform Plugin Framework.
//
// https://www.terraform.io/plugin/framework
type Provider struct {
	tfProvider      tfsdkprovider.Provider
	tfServer        tfprotov6.ProviderServer
	resourcesByType resourcesByType
	info            info.ProviderInfo
	resourcesCache  resources
	resourcesOnce   sync.Once
}

var _ plugin.Provider = &Provider{}

func NewProvider(info info.ProviderInfo) plugin.Provider {
	p := info.P()
	server6 := providerserver.NewProtocol6(p)
	return &Provider{
		tfProvider: p,
		tfServer:   server6(),
		info:       info,
	}
}

func NewProviderServer(info info.ProviderInfo) pulumirpc.ResourceProviderServer {
	return plugin.NewProviderServer(NewProvider(info))
}

// Closer closes any underlying OS resources associated with this provider (like processes, RPC channels, etc).
func (p *Provider) Close() error {
	panic("TODO")
}

// Pkg fetches this provider's package.
func (p *Provider) Pkg() tokens.Package {
	panic("TODO")
}

// GetSchema returns the schema for the provider.
func (p *Provider) GetSchema(version int) ([]byte, error) {
	panic("TODO")
}

// CheckConfig validates the configuration for this resource provider.
func (p *Provider) CheckConfig(urn resource.URN,
	olds, news resource.PropertyMap, allowUnknowns bool) (resource.PropertyMap, []plugin.CheckFailure, error) {
	// TODO proper implementation here.
	return news, []plugin.CheckFailure{}, nil
}

// DiffConfig checks what impacts a hypothetical change to this provider's configuration will have on the provider.
func (p *Provider) DiffConfig(urn resource.URN, olds, news resource.PropertyMap,
	allowUnknowns bool, ignoreChanges []string) (plugin.DiffResult, error) {

	// TODO proper implementation here.
	return plugin.DiffResult{}, nil
}

// Configure configures the resource provider with "globals" that control its behavior.
func (p *Provider) Configure(inputs resource.PropertyMap) error {
	// TODO actually configure
	return nil
}

func PropertyMapToValue(schema tfsdk.Schema, props resource.PropertyMap) (tftypes.Value, diag.Diagnostics) {
	panic("TODO")
}

func ValueToPropertyMap(schema tfsdk.Schema, value tftypes.Value) (resource.PropertyMap, diag.Diagnostics) {
	panic("TODO")
}

// Create allocates a new instance of the provided resource and returns its unique resource.ID.
func (p *Provider) Create(urn resource.URN, news resource.PropertyMap,
	timeout float64, preview bool) (resource.ID, resource.PropertyMap, resource.Status, error) {

	ctx := context.TODO()

	// TODO handle preview=true that should not call Create

	var diags diag.Diagnostics

	res, err := p.resourcesByType.ByURN(urn)
	if err != nil {
		return "", nil, 0, err
	}

	schema, diag1 := res.GetSchema(ctx)
	diags.Append(diag1...)

	plannedValue, diag2 := PropertyMapToValue(schema, news)
	diags.Append(diag2...)

	req := tfsdkresource.CreateRequest{
		Plan: tfsdk.Plan{
			Raw:    plannedValue,
			Schema: schema,
		},
	}

	// TODO set req.ProviderMeta
	//
	// See https://www.terraform.io/internals/provider-meta

	// TODO set req.Config: tfsdk.Config.
	//
	// See https://www.terraform.io/plugin/framework/accessing-values
	//
	// Provider may want to read resource configuration separately from the Plan. Need to clarify how these can be
	// different (perhaps .Config is as-written and excludes any computations performed by executing the program).
	// Currently it is not obvious where to find this data in Pulumi protocol.

	resp := tfsdkresource.CreateResponse{
		State: tfsdk.State{
			Raw:    req.Plan.Raw,
			Schema: req.Plan.Schema,
		},
	}

	res.Create(ctx, req, &resp)

	diags.Append(resp.Diagnostics...)

	createdState, diag3 := ValueToPropertyMap(resp.State.Schema, resp.State.Raw)
	diags.Append(diag3...)

	if diags.HasError() {
		// TODO error out
	}

	// TODO handle resp.Private field to save that state inside Pulumi state.

	var createdID resource.ID // TODO allocate ID

	return createdID, createdState, resource.StatusOK, nil
}

// Read the current live state associated with a resource. Enough state must be include in the inputs to uniquely
// identify the resource; this is typically just the resource ID, but may also include some properties. If the resource
// is missing (for instance, because it has been deleted), the resulting property map will be nil.
func (p *Provider) Read(urn resource.URN, id resource.ID,
	inputs, state resource.PropertyMap) (plugin.ReadResult, resource.Status, error) {
	panic("TODO")
}

// Update updates an existing resource with new values.
func (p *Provider) Update(urn resource.URN, id resource.ID, olds resource.PropertyMap, news resource.PropertyMap,
	timeout float64, ignoreChanges []string, preview bool) (resource.PropertyMap, resource.Status, error) {
	panic("TODO")
}

func (p *Provider) Delete(urn resource.URN, id resource.ID,
	props resource.PropertyMap, timeout float64) (resource.Status, error) {
	panic("TODO")
}

// Construct creates a new component resource.
func (p *Provider) Construct(info plugin.ConstructInfo, typ tokens.Type, name tokens.QName, parent resource.URN,
	inputs resource.PropertyMap, options plugin.ConstructOptions) (plugin.ConstructResult, error) {
	return plugin.ConstructResult{},
		fmt.Errorf("Construct is not implemented for Terraform Plugin Framework bridged providers")
}

// Invoke dynamically executes a built-in function in the provider.
func (p *Provider) Invoke(tok tokens.ModuleMember,
	args resource.PropertyMap) (resource.PropertyMap, []plugin.CheckFailure, error) {
	panic("TODO")
}

// StreamInvoke dynamically executes a built-in function in the provider, which returns a stream of responses.
func (p *Provider) StreamInvoke(tok tokens.ModuleMember, args resource.PropertyMap,
	onNext func(resource.PropertyMap) error) ([]plugin.CheckFailure, error) {
	panic("TODO")
}

// Call dynamically executes a method in the provider associated with a component resource.
func (p *Provider) Call(tok tokens.ModuleMember, args resource.PropertyMap, info plugin.CallInfo,
	options plugin.CallOptions) (plugin.CallResult, error) {
	return plugin.CallResult{},
		fmt.Errorf("Call is not implemented for Terraform Plugin Framework bridged providers")
}

// GetPluginInfo returns this plugin's information.
func (p *Provider) GetPluginInfo() (workspace.PluginInfo, error) {
	ver, err := semver.Parse(p.info.Version)
	if err != nil {
		return workspace.PluginInfo{}, err
	}
	info := workspace.PluginInfo{
		Name:    p.info.Name,
		Version: &ver,
		Kind:    workspace.ResourcePlugin,
	}
	return info, nil
}

// SignalCancellation asks all resource providers to gracefully shut down and abort any ongoing operations. Operation
// aborted in this way will return an error (e.g., `Update` and `Create` will either a creation error or an
// initialization error. SignalCancellation is advisory and non-blocking; it is up to the host to decide how long to
// wait after SignalCancellation is called before (e.g.) hard-closing any gRPC connection.
func (p *Provider) SignalCancellation() error {
	return nil
}

func (p *Provider) terraformResourceName(resourceToken tokens.Type) (string, error) {
	for tfname, v := range p.info.Resources {
		if v.Tok == resourceToken {
			return tfname, nil
		}
	}
	return "", fmt.Errorf("Unkonwn resource: %v", resourceToken)
}

// func (p *Provider) findResource(ctx context.Context, token tokens.Type) {
// 	for _, makeResource := range p.tfProvider.Resources(ctx) {
// 		res := makeResource()
// 		schema := res.GetSchema(ctx)
// 		schema.
// 	}
// }

type resourcesByType map[tokens.Type]tfsdkresource.Resource

func (rbt resourcesByType) ByURN(urn resource.URN) (tfsdkresource.Resource, error) {
	r, ok := rbt[urn.Type()]
	if !ok {
		return nil, fmt.Errorf("unrecognized resource type: %s", urn.Type())
	}
	return r, nil
}

func newResourcesByType(ctx context.Context, prov *tfsdkprovider.Provider) resourcesByType {
	panic("TODO")
}
