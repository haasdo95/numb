package prettyprint

import (
	"os"
	"github.com/olekukonko/tablewriter"
	"fmt"
	"github.com/haasdo95/numb//utils"
)

func stringify(values []interface{}) []string {
	strs := make([]string, len(values))
	for idx, value := range values {
		strs[idx] = fmt.Sprint(value)
	}
	return strs
}

func TablePrint(data map[string]interface{}, headerColor int) *tablewriter.Table {
	keys, values := utils.MapKeyValue(data)
	valuesStr := stringify(values)
	table := tablewriter.NewWriter(os.Stdout)
	table.SetColWidth(110)
	table.SetHeader(keys)
	hcolors := makeHeaderColor(data, headerColor)
	table.SetHeaderColor(hcolors...)
	table.Append(valuesStr)
	return table
}

func makeHeaderColor(data map[string]interface{}, bgColor int) []tablewriter.Colors {
	cs := make([]tablewriter.Colors, len(data))
	for i := 0; i < len(cs); i++ {
		cs[i] = tablewriter.Color(bgColor, tablewriter.FgBlackColor)
	}
	return cs
}

func DisplayTestFailure() {
	table := TablePrint(map[string]interface{}{
		"Failed to display testing result": "Perhaps they are not even recorded",
	}, tablewriter.BgRedColor)
	table.SetColWidth(110)
	table.Render()
}

func DisplayHyperParamFailure() {
	table := TablePrint(map[string]interface{}{
		"Failed to display hyper-parameters": "Perhaps they are not even recorded",
	}, tablewriter.BgRedColor)
	table.SetColWidth(110)
	table.Render()
}