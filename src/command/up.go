package command

type Up struct {
}

func (u Up) UsageInfoShort() string {
	return "publishes the page using the page definition file provided"
}

func (u Up) UsageInfoExpanded() string {
	return ""
}

func (u Up) UsageCategory() int {
	return 1
}

func (u Up) Execute() bool {
	// Parse yaml file

	// Resolve template url, is it valid?
	// Download template from url, build static assets as needed,
	// then read build files into memory. Take into consideration the size
	// of the built assets - will it be ok to store in memory until deploy?
	// or maybe copy these one by one into a deploy directory (zip if needed)?
	// maintaining a flag that signals deploy step once assets are ready.

	// Get default host for host_value on yaml file. Does infrastructure
	// exist to deploy assets? If not create infrastructure with message
	// 'Creating infrastructure on [host_value]...'
	// Infrastructure could potentially be defined and created with
	// Infrastructure as Code tool e.g. terraform (this logic)
	// may need to be done 'page conf host...' command

	// Get default registrar for registrar_value on yaml file,
	// does domain exist on registrar? if not register with message
	// 'Registering domain.com with [registrar_value]...'
	// configure dns records as needed so that the custom domain
	// points to the host infrastructure

	// Take assets from deploy directory, and execute depoyment via host
	// cli

	return true
}

func (u Up) Output() string {
	return ""
}

func (u Up) ValidArgs() bool {
	return true
}
