package cute

var commonBuilder *HTTPTestMaker

func init() {
	commonBuilder = NewHTTPTestMaker()
}

// NewTestBuilder is function for create base test builder,
// For create custom test builder use NewHTTPTestMaker()
func NewTestBuilder() AllureBuilder {
	return commonBuilder.NewTestBuilder()
}
