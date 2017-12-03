package prettyprint

import (
	"os"
	"github.com/olekukonko/tablewriter"
	"fmt"
	"github.com/user/numb/utils"
)

func stringify(values []interface{}) []string {
	strs := make([]string, len(values))
	for idx, value := range values {
		strs[idx] = fmt.Sprint(value)
	}
	return strs
}

func TablePrint(data map[string]interface{})  {
	keys, values := utils.MapKeyValue(data)
	valuesStr := stringify(values)
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(keys)
	table.Append(valuesStr)
	table.Render()
}