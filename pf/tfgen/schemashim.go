// Copyright 2016-2023, Pulumi Corporation.
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

package tfgen

import (
	"context"

	"github.com/pulumi/pulumi-terraform-bridge/pf/internal/schemashim"
	pf "github.com/pulumi/pulumi-terraform-bridge/pf/tfbridge"
	"github.com/pulumi/pulumi-terraform-bridge/v3/pkg/tfbridge"
)

func shimSchemaOnlyProviderInfo(ctx context.Context, provider pf.ProviderInfo) tfbridge.ProviderInfo {
	shimProvider := schemashim.ShimSchemaOnlyProvider(ctx, provider.NewProvider())
	var copy tfbridge.ProviderInfo = provider.ProviderInfo
	copy.P = shimProvider
	return copy
}
