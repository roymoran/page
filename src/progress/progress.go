package progress

import "time"

type ProgressHostState string
type ProgressCertificateState string
type ProgressDomainState string
type ProgressWebsiteFilesState string

const (
	HostCheck         string = "Checking Host..."
	CertificateCheck  string = "Checking Certificate..."
	DomainCheck       string = "Checking Domain..."
	WebsiteFilesCheck string = "Checking Website Files..."
)

var StandardTimeout time.Duration = 3 * time.Minute
var ProgressSequence = []string{"[\\]", "[|]", "[/]", "[-]"}
var HostProvisioningSequence = []string{"Provisioning...", "Provisioned..."}
var HostWebsiteFilesUploadingSequence = []string{"Uploading...", "Uploaded..."}
var CertificateGeneratingSequence = []string{"Generating...", "Generated..."}
var CertificateRenewingSequence = []string{"Renewing...", "Renewed..."}
var DomainUpdatingSequence = []string{"Updating...", "Updated..."}
var ValidatingSequence = []string{"Validating...", "Valid..."}
var ProgressFailedMessage = "Failed..."
var ProgressFailed = "[x]"
var ProgressComplete = "[✓]"
