package cloudglog

import (
	"testing"
	"bytes"
	"bufio"
	"github.com/stretchr/testify/assert"
	"strings"
)


type logCase struct {
	name string
	logFunc func(...interface{})
	args []interface{}
	expected string
}

var logCases = []logCase{
	{
		name: "Info",
		logFunc: Info,
		args: []interface{}{"Info"},
		expected: " Info\n",
	},
	{
		name: "Infoln",
		logFunc: Infoln,
		args: []interface{}{"Infoln"},
		expected: " Infoln\n",
	},
	{
		name: "Warning",
		logFunc: Warning,
		args: []interface{}{"Warning"},
		expected: " Warning\n",
	},
	{
		name: "Warningln",
		logFunc: Warningln,
		args: []interface{}{"Warningln"},
		expected: " Warningln\n",
	},
	{
		name: "Error",
		logFunc: Error,
		args: []interface{}{"Error"},
		expected: " Error\n",
	},
	{
		name: "Errorln",
		logFunc: Errorln,
		args: []interface{}{"Errorln"},
		expected: " Errorln\n",
	},
	//{
	// TODO: test os.Exit(1)
	//	name: "Fatal",
	//	logFunc: Fatal,
	//	args: []interface{}{"Fatal"},
	//	expected: " Fatal\n",
	//},
	//{
	//	name: "Fatalln",
	//	logFunc: Fatalln,
	//	args: []interface{}{"Fatalln"},
	//	expected: " Fatalln\n",
	//},
}

func Test_LogFile(t *testing.T) {

	var b bytes.Buffer
	writer := bufio.NewWriter(&b)

	LogFile(writer)

	for _, testCase := range logCases {

		testCase.logFunc(testCase.args...)
		writer.Flush()

		t.Log(string(b.Bytes()))

		assert.Equal(t, testCase.expected, strings.Split(string(b.Bytes()), ":")[5], "%s wrong assertion for %s", t.Name(), testCase.name)

		b.Reset()
	}
}