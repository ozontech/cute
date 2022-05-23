package cute

var commonBuilder *HTTPTestMaker

func init() {
	commonBuilder = NewHTTPTestMaker()
}

func NewTestBuilder() AllureBuilder {
	return commonBuilder.NewTestBuilder()
}
