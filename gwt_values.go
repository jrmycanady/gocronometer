package gocronometer

// The following constants contain header values required for GWT requests. These values are found by inspecting a
// request from the web app. When the web app is updated these values can change. The values provided here are the
// default that will be used by the library if new values are not provided.
const (
	GWTContentType = "text/x-gwt-rpc; charset=UTF-8"
	GWTModuleBase  = "https://cronometer.com/cronometer/"
	GWTPermutation = "7B121DC5483BF272B1BC1916DA9FA963"

	// GWTHeader is what appears to be a hash value that is provided at the beginning of every GWT request. As it
	// changes with app updates it appears to be related to validating the version the requester is expecting.
	//GWTHeader = "3B6C5196158464C5643BA376AF05E7F1"
	GWTHeader = "2D6A926E3729946302DC68073CB0D550"
)

// The following are the GWT procedure calls as found from inspection of the app.
const (

	// GWTGenerateAuthToken will generate a GWT auth token. The only known use case is for accessing non GWT API calls
	// such as data export.
	// The first parameter in the string should be the sesnonce and the second is the users ID.
	GWTGenerateAuthToken = "7|0|8|https://cronometer.com/cronometer/|" + GWTHeader + "|com.cronometer.shared.rpc.CronometerService|generateAuthorizationToken" +
		"|java.lang.String/2004016611|I|com.cronometer.shared.user.AuthScope/2065601159|%s|1|2|3|4|4|5|6|6|7|8|%s|3600|7|2|"

	// GWTAuthenticate will authenticate with the GWT api. The sesnonce should be set in the cookies.
	GWTAuthenticate = "7|0|5|https://cronometer.com/cronometer/|" + GWTHeader + "|com.cronometer.shared.rpc.CronometerService|authenticate|java.lang.Integer/3438268394|1|2|3|4|1|5|5|-300|"

	// GWTLogout will log the session out.
	// The only parameter should be the sesnonce.
	GWTLogout = "7|0|6|https://cronometer.com/cronometer/|" + GWTHeader + "|com.cronometer.shared.rpc.CronometerService|logout|java.lang.String/2004016611|%s|1|2|3|4|1|5|6|"
)
