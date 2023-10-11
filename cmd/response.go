package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

func NewResponseStruct(success bool, message, err string) CLIResponse {
	return CLIResponse{
		Success: success,
		Message: message,
		Error:   err,
	}

}

func NewResponse(cmd *cobra.Command, success bool, message, err string) {
	resp := CLIResponse{
		Success: success,
		Message: message,
		Error:   err,
	}

	switch {
	case returnJSON:
		uJSON, err := json.Marshal(resp)
		if err != nil {
			NewResponse(cmd, false, "", err.Error())
			return
		}
		fmt.Println(string(uJSON))
	default:
		if !resp.Success {
			fmt.Println(resp.Error)
			cmd.Help()
		} else {
			fmt.Println(resp.Message)
		}
	}
}
