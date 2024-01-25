// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &JsonFileResource{}
var _ resource.ResourceWithImportState = &JsonFileResource{}

func NewJsonFileResource() resource.Resource {
	return &JsonFileResource{}
}

type JsonFileResource struct {
}

type JsonFileResourceModel struct {
	Value  types.String `tfsdk:"value"`
	Nested types.List   `tfsdk:"nested"`
}

type JsonFileStruct struct {
	Value string
}

func writeJson(data JsonFileStruct) {
	file, _ := json.MarshalIndent(data, "", "")
	os.WriteFile("data.json", file, 0644)
}

func readJson() JsonFileStruct {
	var ret JsonFileStruct
	file, _ := os.ReadFile("data.json")
	json.Unmarshal(file, &ret)
	return ret
}

func (r *JsonFileResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_data"
}

func (r *JsonFileResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"value": schema.StringAttribute{
				Required: true,
			},
			"nested": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"fixed": schema.StringAttribute{
							Computed: true,
						},
						"time": schema.StringAttribute{
							Computed: true,
							Optional: true,
							Default:  stringdefault.StaticString("banana"),
						},
					},
				},
			},
		},
	}
}

func (r *JsonFileResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
}

func (r *JsonFileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data JsonFileResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	writeJson(JsonFileStruct{
		Value: data.Value.ValueString(),
	})
	fillNested(&data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func fillNested(data *JsonFileResourceModel) {
	nestedType := map[string]attr.Type{"time": types.StringType, "fixed": types.StringType}
	data.Nested = types.ListValueMust(types.ObjectType{AttrTypes: nestedType},
		[]attr.Value{
			types.ObjectValueMust(nestedType,
				map[string]attr.Value{
					"time":  types.StringValue(time.Now().String()),
					"fixed": types.StringValue("fixed"),
				}),
		})
}

func (r *JsonFileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data JsonFileResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	fileData := readJson()
	data.Value = types.StringValue(fileData.Value)
	fillNested(&data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *JsonFileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data JsonFileResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	fileData := readJson()
	fileData.Value = data.Value.ValueString()
	writeJson(fileData)
	fillNested(&data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *JsonFileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data JsonFileResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	os.Remove("data.json")
}

func (r *JsonFileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
