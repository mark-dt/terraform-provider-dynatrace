/**
* @license
* Copyright 2020 Dynatrace LLC
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
 */

package service

import (
	"strings"

	common "github.com/dynatrace-oss/terraform-provider-dynatrace/datasources"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/export"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/config"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/logging"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		Read: logging.EnableDS(DataSourceRead),
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func DataSourceRead(d *schema.ResourceData, m any) (err error) {
	var name string
	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	}
	creds, err := config.Credentials(m, config.CredValDefault)
	if err != nil {
		return err
	}

	service := export.Service(creds, export.ResourceTypes.CalculatedServiceMetric)
	var stubs api.Stubs
	if stubs, err = service.List(); err != nil {
		return err
	}
	if len(stubs) > 0 {
		for _, stub := range stubs {
			if name == stub.Name {
				d.SetId(stub.ID)
				return nil
			}
		}
	}
	d.SetId(common.NotFoundID(strings.ToLower(strings.ReplaceAll(name, " ", ""))))
	return nil
}
