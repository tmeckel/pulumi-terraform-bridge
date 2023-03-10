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

package tfbridgetests

import (
	"testing"

	"github.com/pulumi/pulumi-terraform-bridge/pf/tests/internal/testprovider"
	testutils "github.com/pulumi/pulumi-terraform-bridge/testing/x"
)

func TestCreateWithComputedOptionals(t *testing.T) {
	server := newProviderServer(t, testprovider.SyntheticTestBridgeProvider())
	testCase := `
        {
          "method": "/pulumirpc.ResourceProvider/Create",
          "request": {
            "urn": "urn:pulumi:test-stack::basicprogram::testbridge:index/testres:Testcompres::r1",
            "properties": {
              "ecdsacurve": "P384"
            },
            "preview": false
          },
          "response": {
            "id": "r1",
            "properties": {
              "ecdsacurve": "P384",
              "id": "r1"
            }
          }
        }
        `
	testutils.Replay(t, server, testCase)
}

func TestCreateWritesSchemaVersion(t *testing.T) {
	server := newProviderServer(t, testprovider.RandomProvider())

	testutils.Replay(t, server, `
	{
	  "method": "/pulumirpc.ResourceProvider/Create",
	  "request": {
	    "urn": "urn:pulumi:dev::repro-pulumi-random::random:index/randomString:RandomString::s",
	    "properties": {
	      "length": 1
	    }
	  },
	  "response": {
	    "id": "*",
	    "properties": {
	      "__meta": "{\"schema_version\":\"2\"}",
	      "id": "*",
	      "result": "*",
	      "length": 1,
	      "lower": true,
	      "minLower": 0,
	      "minNumeric": 0,
	      "minSpecial": 0,
	      "minUpper": 0,
	      "number": true,
	      "numeric": true,
	      "special": true,
	      "upper": true
	    }
	  }
	}
        `)
}

func TestPreviewCreate(t *testing.T) {
	server := newProviderServer(t, testprovider.RandomProvider())

	testCase := `
	{
	  "method": "/pulumirpc.ResourceProvider/Create",
	  "request": {
	    "urn": "urn:pulumi:dev::repro::random:index/randomInteger:RandomInteger::k",
	    "properties": {
	      "max": 10,
	      "min": 0
	    },
	    "preview": true
	  },
	  "response": {
	    "properties": {
	      "id": "04da6b54-80e4-46f7-96ec-b56ff0331ba9",
	      "max": 10,
	      "min": 0,
	      "result": "04da6b54-80e4-46f7-96ec-b56ff0331ba9"
	    }
	  },
	  "metadata": {
	    "kind": "resource",
	    "mode": "client",
	    "name": "random"
	  }
	}
`
	testutils.Replay(t, server, testCase)
}

func TestCreateWithFirstClassSecrets(t *testing.T) {
	server := newProviderServer(t, testprovider.RandomProvider())
	testCase := `
	{
	  "method": "/pulumirpc.ResourceProvider/Create",
	  "request": {
	    "urn": "urn:pulumi:dev::pulumi-terraform-bridge-812::random:index/randomPet:RandomPet::pet",
	    "properties": {
	      "separator": {
		"4dabf18193072939515e22adb298388d": "1b47061264138c4ac30d75fd1eb44270",
		"value": "BbAXG:}h"
	      }
	    },
	    "preview": true
	  },
	  "response": {
            "properties": {
              "id": "*",
              "length": 2,
              "separator": "BbAXG:}h"
            }
          }
	}`
	testutils.Replay(t, server, testCase)
}
