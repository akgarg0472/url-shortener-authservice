package utils

func GetSignupSuccessEmailBody(name string) string {
	logoUrl := GetEnvVariable("URL_SHORTENER_LOGO_URL", "https://i.ibb.co/1z4L2sB/url-shortener-logo.png")
	dashboardLink := GetEnvVariable("FRONTEND_BASE_DOMAIN", "http://127.0.0.1:3000/") + GetEnvVariable("FRONTEND_DASHBOARD_PAGE_URL", "dashboard")

	return "<div style='font-family:Arial,sans-serif;margin:0;padding:0;background-color:#f4f4f475!important;max-width:600px;margin:20px auto;padding:20px;background-color:#fff;border-radius:6px;box-shadow:0 0 10px rgba(0,0,0,0.1);color:#333;text-align:center;'><img src='" + logoUrl + "' alt='UrlShortener Logo' style='max-width:33%;height:auto;object-fit:contain;margin-bottom:20px;'/><br /><h2>Welcome to Our Platform!</h2><p style='font-size:16px;text-align:left'>Dear <strong>" + name + "</strong>,</p><p style='text-align:left;font-size:16px;'>Thank you for choosing our URL shortener service. We are thrilled to have you onboard. Your journey with us is about to unfold. Your links, made shorter, smarter, and ready for sharing!</p><p style='text-align:left;font-size:16px;'>Feel free to start shortening your URLs and tracking their performance with some amazing cool features. If you ever need any assistance, our team is here to help you.</p><div style='text-align:left;font-size:14px;'><span>Cheers</span><br/><span>Akhilesh Garg</span><br/><span>Lead Developer</span></div><a href='" + dashboardLink + "' style='display:inline-block;padding:10px 20px;margin-top:20px;text-decoration:none;background-color:#e91e63;color:#fff;border-radius:5px;'>Get Started</a><div style='padding:10px;border-radius:0 0 5px 5px;margin-top:20px;font-size:12px;line-height:18px;'>UrlShortener is a hobby project by Akhilesh Garg<br />&copy; 2024 ÂkHîL, All rights reserved.</div></div>"
}

func GenerateForgotPasswordEmailBody(email string, name string, forgotPasswordUrl string) string {
	return "<div style='font-family:Arial,sans-serif;margin:0;padding:0;background-color:#f4f4f475!important;max-width:600px;margin:20px auto;padding:20px;background-color:#fff;border-radius:5px;box-shadow:0 0 10px rgba(0,0,0,0.1);color:#333;text-align:center;'><p style='font-size:16px;text-align:left'>Hello " + name + ",</p><p style='text-align:left;line-height:24px;font-size:16px;'>We received a request to reset the password for your URLShortener account associated with <span style='color:#15c'>" + email + "</span>.</p><p style='font-size:16px;text-align:left;margin-top:24px;'>To reset your password, click the button below:</p><a href='" + forgotPasswordUrl + "' style='display:inline-block;padding:10px 20px;text-decoration:none;background-color:#5063f0;color:#fff;border-radius:6px;'>Reset Password</a><p style='text-align:left;margin-top:24px;line-height:24px;font-size:16px;'>If you didn't make this request, or if you're having trouble signing in, contact us via our support site. No changes have been made to your account.</p><div style='text-align:left;font-size:16px;margin-top:24px;'>- URLShortener Team</div><div style='padding:10px;border-radius:0 0 5px 5px;margin-top:20px;font-size:12px;line-height:18px;'>UrlShortener is a hobby project by Akhilesh Garg<br />&copy; 2024 ÂkHîL, All rights reserved.</div></div>"
}

func GeneratePasswordChangeSuccessEmailBody(email string) string {
	return "<div style='font-family:Arial,sans-serif;margin:0;padding:0;background-color:#f4f4f475!important;max-width:600px;margin:20px auto;padding:20px;background-color:#fff;border-radius:5px;box-shadow:0 0 10px rgba(0,0,0,0.1);color:#333;text-align:center;'><p style='font-size:16px;text-align:left'>Dear User,</p><p style='text-align:left;line-height:24px;font-size:16px;'>The password of your account associated with <span style='color:#15c'>" + email + "</span> has been successfully changed. If you made this change, no further action is needed.</p><p style='text-align:left;margin-top:24px;line-height:24px;font-size:16px;'>If you didn't request this, please change your password immediately & contact us via our support site. No changes have been made to your account.</p><div style='text-align:left;font-size:16px;margin-top:24px;'>- URLShortener Team</div><div style='padding:10px;border-radius:0 0 5px 5px;margin-top:20px;font-size:12px;line-height:18px;'>UrlShortener is a hobby project by Akhilesh Garg<br />&copy; 2024 ÂkHîL, All rights reserved.</div></div>"
}
