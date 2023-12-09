package commands

import (
	"encoding/json"
	"mime"
	"path/filepath"
)

type Infra struct {
}

type InfraArgs struct {
	InfraFuncName  string
	InfraFuncInput string
}

var infra CommandInfo = CommandInfo{
	DisplayName:              "infra",
	ExecutionOutput:          "",
	ExecutionOk:              true,
	MinimumExpectedArgs:      1,
	MaximumExpectedArguments: 99,
	OrderedArgLabel:          []string{"infraFuncName", "infraFuncInput"},
	ArgValues: map[string]string{
		"infraFuncName":  "",
		"infraFuncInput": "",
	},
}

var infraArgs InfraArgs = InfraArgs{
	InfraFuncName:  "",
	InfraFuncInput: "",
}

func (inf Infra) UsageInfoShort() string {
	return ""
}

func (inf Infra) UsageInfoExpanded() string {
	return ""
}

func (inf Infra) UsageCategory() int {
	return -1
}

func (inf Infra) Execute() {
	if !infra.ExecutionOk {
		return
	}

	switch infraArgs.InfraFuncName {
	case "mimetype":
		mimeType := contentType(infraArgs.InfraFuncInput)
		output := map[string]string{"mimetype": mimeType}
		jsonOutput, _ := json.Marshal(output)
		infra.ExecutionOutput = string(jsonOutput)
	default:
		infra.ExecutionOutput = "Invalid infra function name"
	}
}

func (inf Infra) Output() string {
	// output arguments submitted to the infra command
	return infra.ExecutionOutput
}

func (inf Infra) BindArgs() {
	if !conf.ExecutionOk {
		return
	}

	infraArgs.InfraFuncName = infra.ArgValues["infraFuncName"]
	infraArgs.InfraFuncInput = infra.ArgValues["infraFuncInput"]
}

func contentType(filePath string) string {
	ext := filepath.Ext(filePath)
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		return "application/octet-stream" // Default MIME type
	}
	return mimeType
}
