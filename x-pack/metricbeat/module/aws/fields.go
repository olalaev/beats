// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License;
// you may not use this file except in compliance with the Elastic License.

// Code generated by beats/dev-tools/cmd/asset/asset.go - DO NOT EDIT.

package aws

import (
	"github.com/elastic/beats/libbeat/asset"
)

func init() {
	if err := asset.SetFields("metricbeat", "aws", asset.ModuleFieldsPri, AssetAws); err != nil {
		panic(err)
	}
}

// AssetAws returns asset data.
// This is the base64 encoded gzipped contents of module/aws.
func AssetAws() string {
	return "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"
}