// Copyright 2023 Woodpecker Authors
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

package secret

import (
	"context"
	"fmt"
	"html/template"
	"os"

	"github.com/urfave/cli/v3"

	"go.woodpecker-ci.org/woodpecker/v2/cli/common"
	"go.woodpecker-ci.org/woodpecker/v2/cli/internal"
)

var secretShowCmd = &cli.Command{
	Name:      "show",
	Usage:     "show secret information",
	ArgsUsage: "[repo-id|repo-full-name]",
	Action:    secretShow,
	Flags: []cli.Flag{
		common.OrgFlag,
		&cli.StringFlag{
			Name:  "name",
			Usage: "secret name",
		},
		common.FormatFlag(tmplSecretList, true),
	},
}

func secretShow(ctx context.Context, c *cli.Command) error {
	var (
		secretName = c.String("name")
		format     = c.String("format") + "\n"
	)

	if secretName == "" {
		return fmt.Errorf("secret name is missing")
	}

	client, err := internal.NewClient(ctx, c)
	if err != nil {
		return err
	}

	orgID, err := parseTargetArgs(client, c)
	if err != nil {
		return err
	}

	secret, err := client.OrgSecret(orgID, secretName)
	if err != nil {
		return err
	}

	tmpl, err := template.New("_").Funcs(secretFuncMap).Parse(format)
	if err != nil {
		return err
	}
	return tmpl.Execute(os.Stdout, secret)
}
