//+build unit

package functions_test

import (
	"errors"
	"os"
	"testing"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/IBM/ibm-cos-sdk-go/service/s3"
	"github.com/IBM/ibmcloud-cos-cli/config/commands"
	"github.com/IBM/ibmcloud-cos-cli/config/flags"
	"github.com/IBM/ibmcloud-cos-cli/cos"
	"github.com/IBM/ibmcloud-cos-cli/di/providers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/urfave/cli"
)

func TestMPUAbortSunnyPath(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "TargetBucket"
	targetKey := "TargetKey"
	targetUploadID := "TargetUploadID"

	var capturedInput *s3.AbortMultipartUploadInput

	providers.MockS3API.
		On("AbortMultipartUpload", mock.MatchedBy(
			func(input *s3.AbortMultipartUploadInput) bool {
				capturedInput = input
				return true
			})).
		Return(
			new(s3.AbortMultipartUploadOutput), nil).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.AbortMultipartUpload,
		"--" + flags.Bucket, targetBucket,
		"--" + flags.Key, targetKey,
		"--" + flags.UploadID, targetUploadID,

		"--" + flags.Region, "us",
	}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region ( since success is last )
	providers.MockS3API.AssertNumberOfCalls(t, "AbortMultipartUpload", 1)
	//assert exit code is zero
	assert.Equal(t, (*int)(nil), exitCode) // no exit trigger in the cli

	// assert request match cli parameters
	assert.Equal(t, *capturedInput.Bucket, targetBucket)
	assert.Equal(t, *capturedInput.Key, targetKey)
	assert.Equal(t, *capturedInput.UploadId, targetUploadID)

	// capture all output //
	output := providers.FakeUI.Outputs()
	//assert OK
	assert.Contains(t, output, "OK")
	//assert Not Fail
	assert.NotContains(t, output, "FAIL")
}

func TestMPUAbortRainyPath(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "TargetBucket"
	targetKey := "TargetKey"
	targetUploadID := "TargetUploadID"

	var capturedInput *s3.AbortMultipartUploadInput

	providers.MockS3API.
		On("AbortMultipartUpload", mock.MatchedBy(
			func(input *s3.AbortMultipartUploadInput) bool {
				capturedInput = input
				return true
			})).
		Return(
			nil, errors.New("InvalidUploadId")).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.AbortMultipartUpload,
		"--" + flags.Bucket, targetBucket,
		"--" + flags.Key, targetKey,
		"--" + flags.UploadID, targetUploadID,

		"--" + flags.Region, "us",
	}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region ( since success is last )
	providers.MockS3API.AssertNumberOfCalls(t, "AbortMultipartUpload", 1)
	//assert exit code is zero
	assert.Equal(t, 1, *exitCode) // no exit trigger in the cli

	// assert request match cli parameters
	assert.Equal(t, *capturedInput.Bucket, targetBucket)
	assert.Equal(t, *capturedInput.Key, targetKey)
	assert.Equal(t, *capturedInput.UploadId, targetUploadID)

	// capture all output //
	output := providers.FakeUI.Outputs()
	//assert Not OK
	assert.NotContains(t, output, "OK")
	//assert Fail
	assert.Contains(t, output, "FAIL")

}

func TestMPUAbortWithoutUploadID(t *testing.T) {
	defer providers.MocksRESET()

	// --- Arrange ---
	// disable and capture OS EXIT
	var exitCode *int
	cli.OsExiter = func(ec int) {
		exitCode = &ec
	}

	targetBucket := "TargetBucket"
	targetKey := "TargetKey"

	providers.MockS3API.
		On("AbortMultipartUpload", mock.MatchedBy(
			func(input *s3.AbortMultipartUploadInput) bool {
				return true
			})).
		Return(
			nil, errors.New("InvalidUploadId")).
		Once()

	// --- Act ----
	// set os args
	os.Args = []string{"-", commands.AbortMultipartUpload,
		"--" + flags.Bucket, targetBucket,
		"--" + flags.Key, targetKey,
		"--" + flags.Region, "us",
	}
	//call plugin
	plugin.Start(new(cos.Plugin))

	// --- Assert ----
	// assert s3 api called once per region ( since success is last )
	providers.MockS3API.AssertNumberOfCalls(t, "AbortMultipartUpload", 0)
	//assert exit code is zero
	assert.Equal(t, 1, *exitCode) // no exit trigger in the cli

	// capture all output //
	output := providers.FakeUI.Outputs()
	//assert Not OK
	assert.NotContains(t, output, "OK")
	//assert Fail
	assert.Contains(t, output, "FAIL")

}
