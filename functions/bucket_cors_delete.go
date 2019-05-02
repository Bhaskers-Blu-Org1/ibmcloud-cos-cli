package functions

import (
	"strings"

	"github.com/IBM/ibmcloud-cos-cli/config/fields"
	. "github.com/IBM/ibmcloud-cos-cli/i18n"

	"github.com/IBM/ibmcloud-cos-cli/config/flags"

	"github.com/IBM/ibmcloud-cos-cli/config"
	"github.com/IBM/ibmcloud-cos-cli/utils"

	"github.com/IBM/ibm-cos-sdk-go/service/s3"

	"github.com/urfave/cli"
)

// BucketCorsDelete - Deletes CORS confuration from a bucket
// Parameter:
//     	CLI Context Application
// Returns:
//  	Error if triggered
func BucketCorsDelete(c *cli.Context) error {
	// Load COS Context
	cosContext := c.App.Metadata[config.CosContextKey].(*utils.CosContext)

	// Load COS Context UI and Config
	ui := cosContext.UI
	conf := cosContext.Config

	// Set DeleteBucketCorsInput
	input := new(s3.DeleteBucketCorsInput)

	// Required parameter and no optional parameter for DeleteBucketCors
	mandatory := map[string]string{
		fields.Bucket: flags.Bucket,
	}

	// Validate User Inputs and Retrieve Region
	region, err := ValidateUserInputsAndSetRegion(c, input, mandatory, map[string]string{}, conf)
	if err != nil {
		ui.Failed(err.Error())
		cli.ShowCommandHelp(c, c.Command.Name)
		return cli.NewExitError(err.Error(), 1)
	}

	// Setting client to do the call
	client := cosContext.GetClient(region)

	// Alert User that we are performing the call
	ui.Say(T("Deleting CORS from the bucket..."))

	// DeleteBucketCors API
	_, err = client.DeleteBucketCors(input)
	// Error handling
	if err != nil {
		if strings.Contains(err.Error(), "EmptyStaticCreds") {
			ui.Failed(err.Error() + "\n" + T("Try logging in using 'ibmcloud login'."))
		} else {
			ui.Failed(err.Error())
		}
		return cli.NewExitError("", 1)
	}

	// Success
	ui.Ok()
	ui.Say(T("Successfully deleted CORS configuration on bucket: {{.Bucket}}",
		map[string]interface{}{"Bucket": utils.EntityNameColor(*input.Bucket)}))

	// Return
	return nil
}
