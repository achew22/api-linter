// Copyright 2022 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 		https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/googleapis/api-linter/lint"
	"google.golang.org/protobuf/types/descriptorpb"
)

func TestFormatGitHubActionOutput(t *testing.T) {
	tests := []struct {
		name string
		data []lint.Response
		want string
	}{
		{
			name: "Empty input",
			data: []lint.Response{},
			want: "",
		},
		{
			name: "Example with location specifics",
			data: []lint.Response{
				{
					FilePath: "example.proto",
					Problems: []lint.Problem{
						{
							RuleID: "core::naming_formats::field_names",
							Location: &descriptorpb.SourceCodeInfo_Location{
								Span: []int32{1, 2, 3, 4},
							},
						},
						{
							RuleID:  "core::naming_formats::field_names",
							Message: "multi\nline\ncomment",
							Location: &descriptorpb.SourceCodeInfo_Location{
								Span: []int32{5, 6, 7, 8},
							},
						},
					},
				},
			},
			want: `::error file=example.proto line=1 col=2 endLine=3 endColumn=4 title=core։։naming_formats։։field_names::
::error file=example.proto line=5 col=6 endLine=7 endColumn=8 title=core։։naming_formats։։field_names::multi\nline\ncomment
`,
		},
		{
			name: "Example with a couple of responses",
			data: []lint.Response{
				{
					FilePath: "example.proto",
					Problems: []lint.Problem{
						{RuleID: "core::naming_formats::field_names"},
						{RuleID: "core::naming_formats::field_names"},
					},
				},
				{
					FilePath: "example2.proto",
					Problems: []lint.Problem{
						{RuleID: "core::0131::request_message::name"},
						{RuleID: "core::0132::response_message::name"},
					},
				},
				{
					FilePath: "example3.proto",
					Problems: []lint.Problem{
						{RuleID: "core::naming_formats::field_names"},
					},
				},
				{
					FilePath: "example4.proto",
					Problems: []lint.Problem{
						{RuleID: "core::naming_formats::field_names"},
						{RuleID: "core::0132::response_message::name"},
					},
				},
			},
			want: `::error file=example.proto title=core։։naming_formats։։field_names::
::error file=example.proto title=core։։naming_formats։։field_names::
::error file=example2.proto title=core։։0131։։request_message։։name::
::error file=example2.proto title=core։։0132։։response_message։։name::
::error file=example3.proto title=core։։naming_formats։։field_names::
::error file=example4.proto title=core։։naming_formats։։field_names::
::error file=example4.proto title=core։։0132։։response_message։։name::
`,
		},
	}
	for _, test := range tests {
		got := formatGitHubActionOutput(test.data)
		if diff := cmp.Diff(string(test.want), string(got)); diff != "" {
			t.Errorf("formatGitHubActionOutput() mismatch (-want +got):\n%s", diff)
		}
	}
}
