package enums

type OAuthProvider string
type UserEntityLoginType string
type NotificationType string

const (
	OauthProviderGoogle OAuthProvider = "google"
	OauthProviderGithub OAuthProvider = "github"
)

const (
	UserEntityLoginTypeEmailAndPassword UserEntityLoginType = "email_pass"
	UserEntityLoginTypeOauthAndOtp      UserEntityLoginType = "oauth_otp"
	UserEntityLoginTypeOauthOnly        UserEntityLoginType = "oauth_only"
)

const (
	NotificationTypeEmail NotificationType = "EMAIL"
)
