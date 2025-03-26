/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/spf13/cobra"
)

// grantCmd represents the grant command
var grantCmd = &cobra.Command{
	Use:   "grant",
	Short: "Grant user as project role",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		project, _ := cmd.Flags().GetString("project")
		role, _ := cmd.Flags().GetString("role")
		member, _ := cmd.Flags().GetString("member")

		// Get persistent flags from root command
		baseURL, _ := rootCmd.PersistentFlags().GetString("base")
		username, _ := rootCmd.PersistentFlags().GetString("user")
		password, _ := rootCmd.PersistentFlags().GetString("pwd")

		// Build URL
		url := fmt.Sprintf("%s/rest/user/%s/project/%s/role/%s",
			baseURL,
			url.PathEscape(member),
			url.PathEscape(project),
			url.PathEscape(role))

		// Create request
		req, err := http.NewRequest(http.MethodPut, url, nil)
		if err != nil {
			fmt.Printf("Error creating request: %v\n", err)
			return
		}

		// Add basic auth
		req.SetBasicAuth(username, password)

		// Send request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("Error sending request: %v\n", err)
			return
		}
		defer resp.Body.Close()

		// Check response
		if resp.StatusCode == http.StatusOK {
			fmt.Printf("OK")
		} else {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Printf("Failed. Status code: %d, error reading response: %v\n", resp.StatusCode, err)
				return
			}
			fmt.Printf("Failed. Status code: %d, message: %s\n", resp.StatusCode, string(body))
		}
	},
}

func init() {
	projectCmd.AddCommand(grantCmd)

	// Add required flags for role and member
	grantCmd.Flags().String("project", "", "Specify the project name (required)")
	grantCmd.Flags().String("role", "", "Role to grant (required)")
	grantCmd.Flags().String("member", "", "Member to grant role to (required)")

	grantCmd.MarkFlagRequired("project")
	grantCmd.MarkFlagRequired("role")
	grantCmd.MarkFlagRequired("member")
}
