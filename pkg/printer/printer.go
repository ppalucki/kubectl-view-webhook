/*
Copyright © 2020 Trendyol Tech

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package printer

import (
	"github.com/olekukonko/tablewriter"
	"github.com/pterm/pterm"
	"io"
	"os"
	"strconv"
	"strings"
)

// Printer formats and prints check results and warnings.
type Printer struct {
	out io.Writer
}

// NewPrinter constructs a new Printer with the specified output io.Writer
// and output format.
func NewPrinter(out io.Writer) *Printer {
	return &Printer{
		out: out,
	}
}

func convertStringArrayToBulletListItem(s []string) []pterm.BulletListItem {
	var bulletItems []pterm.BulletListItem
	for _, t := range s {
		bulletItems = append(bulletItems, pterm.BulletListItem{
			Level: 0,
			Text:  t,
		})
	}
	return bulletItems
}

func (p *Printer) Print(model *PrintModel) {
	var data [][]string

	for _, item := range model.Items {
		operationsData, _ := pterm.DefaultBulletList.WithItems(convertStringArrayToBulletListItem(item.Operations)).Srender()
		resourcesData, _ := pterm.DefaultBulletList.WithItems(convertStringArrayToBulletListItem(item.Resources)).Srender()
		namespacesData, _ := pterm.DefaultBulletList.WithItems(convertStringArrayToBulletListItem(item.ActiveNamespaces)).Srender()

		var valid string
		if item.ValidUntil < 4000 {
			valid = pterm.Red(strconv.FormatInt(item.ValidUntil, 10) + "d")
		} else if item.ValidUntil < 60000 {
			valid = pterm.Yellow(strconv.FormatInt(item.ValidUntil, 10) + "d")
		} else {
			valid = pterm.Green(strconv.FormatInt(item.ValidUntil, 10) + "d")
		}

		webhookTreeList := pterm.NewTreeFromLeveledList(pterm.LeveledList{
			pterm.LeveledListItem{Level: 0, Text: item.Webhook.ServiceName},
			pterm.LeveledListItem{Level: 1, Text: "NS  : " + item.Webhook.ServiceNamespace},
			pterm.LeveledListItem{Level: 1, Text: "Path: " + *item.Webhook.ServicePath},
			pterm.LeveledListItem{Level: 1, Text: "Port: " + strconv.Itoa(int(*item.Webhook.ServicePort))},
		})

		s, _ := pterm.DefaultTree.WithRoot(webhookTreeList).Srender()

		data = append(data, []string{item.Kind, item.Name, item.Webhook.Name, strings.TrimSuffix(s, "\n"), resourcesData, operationsData, valid, namespacesData})
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Kind", "Name", "Webhook", "Service", "Resources", "Operations", "Remaining Day", "Active NS"})
	table.SetRowLine(true)
	table.SetAutoMergeCells(true)
	table.SetHeaderLine(true)
	table.SetBorder(true)

	//table.SetReflowDuringAutoWrap(true)
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	//table.SetCenterSeparator("")
	//table.SetColumnSeparator("")
	//table.SetRowSeparator("")

	//table.SetTablePadding(" ")
	//table.SetNoWhiteSpace(true)
	table.AppendBulk(data)

	table.Render()
}
