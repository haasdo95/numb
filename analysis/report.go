package analysis

import (
	"github.com/nasyxx/numb/database"
	"strconv"
	"gopkg.in/mgo.v2/bson"
	"os"
	"os/exec"
	"text/template"
	"fmt"
	"github.com/nasyxx/numb/utils"
	"gopkg.in/mgo.v2"
)

// reportData is used to make rendering cleaner
type reportData struct {
	TestResult map[string]interface{} `json:"test-result"`
	ParamsUsed map[string]interface{} `json:"hyperparams"`
	SourceCode string				  `json:"code"`
}

// Report enables user to run "numb report <ID>"
func Report(collection *mgo.Collection, timestamp string) {
	timestampNum, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		fmt.Println("Failed to parse input")
		return
	}
	query := collection.Find(bson.M{"timestamp": timestampNum})
	if cnt, err := query.Count(); cnt != 1 || err != nil {
		fmt.Println("Failed to find record")
		return
	}
	var result database.Schema
	query.One(&result)
	if result.Test == "" {
		fmt.Println("Cannot report an untest model")
		return
	}

	// create report dir
	stateDictName := strconv.FormatInt(result.Timestamp, 10)
	reportDirName := "report-" + stateDictName
	_, err = os.Stat(reportDirName)
	if !os.IsNotExist(err) { // remove if existing
		utils.Check(os.RemoveAll(reportDirName))
	}
	err = os.Mkdir(reportDirName, 0755)
	utils.Check(err)
	err = os.Chdir(reportDirName)
	utils.Check(err)

	utils.Check(exec.Command("cp", "-rf", "github.com/nasyxx/numb/.nmb/" + stateDictName, "statedict.pkl").Run())

	paramObj, err := utils.Str2Obj(result.Params)
	utils.Check(err)
	testObj, err := utils.Str2Obj(result.Test)
	utils.Check(err)
	reportData := reportData {
		TestResult: testObj,
		ParamsUsed: paramObj,
		SourceCode: result.Code,
	}

	tmplt, err := template.New("report").Parse(`
		<!DOCTYPE html>
		<html lang="en">
			<head>
				<title>report</title>
				<meta charset="UTF-8">
				<meta name="viewport" content="width=device-width, initial-scale=1">
				<script src="https://cdn.rawgit.com/google/code-prettify/master/loader/run_prettify.js"></script>
			</head>
			<body>
				<h2>Source Code</h2>
				<pre class="prettyprint">
					<code class="language-python">
						<br/>{{.SourceCode}}
					</code>
				</pre>
				<h2>Hyper-parameters</h2>
				<table>
					<tr>
						{{range $k,$v := .ParamsUsed}}
							<th>{{$k}}</th>
						{{end}}
					</tr>
					<tr>
						{{range $k,$v := .ParamsUsed}}
							<td>{{$v}}</td>
						{{end}}
					</tr>
				</table>
				<h2>Test Results</h2>
				<table>
					<tr>
						{{range $k,$v := .TestResult}}
							<th>{{$k}}</th>
						{{end}}
					</tr>
					<tr>
						{{range $k,$v := .TestResult}}
							<td>{{$v}}</td>
						{{end}}
					</tr>
				</table>
			</body>
		</html>
		`)
	reportHTML, err := os.OpenFile("report.html", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	utils.Check(err)
	tmplt.Execute(reportHTML, reportData)
	
}