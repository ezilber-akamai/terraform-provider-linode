package lketypes

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

func NewDataSource() datasource.DataSource {
	return &DataSource{
		BaseDataSource: helper.NewBaseDataSource(
			helper.BaseDataSourceConfig{
				Name:   "linode_lke_types",
				Schema: &frameworkDataSourceSchema,
			},
		),
	}
}

type DataSource struct {
	helper.BaseDataSource
}

func (r *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	tflog.Debug(ctx, "Read data.linode_lke_types")

	var data LKETypeFilterModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, d := filterConfig.GenerateID(data.Filters)
	if d != nil {
		resp.Diagnostics.Append(d)
		return
	}
	data.ID = id

	result, d := filterConfig.GetAndFilter(
		ctx, r.Meta.Client, data.Filters, listLKETypes, data.Order, data.OrderBy)
	if d != nil {
		resp.Diagnostics.Append(d)
		return
	}

	data.parseLKETypes(ctx, helper.AnySliceToTyped[linodego.LKEType](result))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func listLKETypes(ctx context.Context, client *linodego.Client, filter string) ([]any, error) {
	tflog.Debug(ctx, "Listing LKE types", map[string]any{
		"filter_header": filter,
	})

	types, err := client.ListLKETypes(ctx, &linodego.ListOptions{
		Filter: filter,
	})
	if err != nil {
		return nil, err
	}

	return helper.TypedSliceToAny(types), nil
}
