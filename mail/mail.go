package mail

import "gopkg.in/gomail.v2"
import "fmt"

func SendWelcomeMail(to string, username string) error {
	m := gomail.NewMessage()

	m.SetHeader("From", "artustech@gmail.com")
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Welcome to Artustech")

	body := fmt.Sprintf(`
			<html lang="en">
			<head>
			  <meta charset="UTF-8">
			  <meta name="viewport" content="width=device-width, initial-scale=1.0">
			  <title>Welcome Email</title>
			  <style>
				body {
				  font-family: Arial, sans-serif;
				  margin: 0;
				  padding: 0;
				  background-color: #f9f9f9;
				}
				.email-container {
				  max-width: 600px;
				  margin: 20px auto;
				  background-color: #ffffff;
				  border-radius: 8px;
				  overflow: hidden;
				  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
				}
				.header {
				  background: linear-gradient(90deg, #4a90e2, #6a11cb);
				  color: white;
				  padding: 20px;
				  text-align: center;
				}
				.header h1 {
				  margin: 0;
				  font-size: 24px;
				}
				.body {
				  padding: 20px;
				  color: #333;
				  line-height: 1.6;
				}
				.body h2 {
				  margin-top: 0;
				  color: #4a90e2;
				}
				.button {
				  display: inline-block;
				  padding: 12px 20px;
				  margin: 20px 0;
				  background: linear-gradient(90deg, #4a90e2, #6a11cb);
				  color: white;
				  text-decoration: none;
				  border-radius: 5px;
				  font-size: 16px;
				}
				.footer {
				  text-align: center;
				  padding: 10px;
				  background-color: #f1f1f1;
				  font-size: 14px;
				  color: #666;
				}
				.social-icons img {
				  width: 24px;
				  margin: 0 5px;
				  vertical-align: middle;
				}
				@media (max-width: 600px) {
				  .header h1, .body h2 {
					font-size: 20px;
				  }
				  .button {
					font-size: 14px;
					padding: 10px 15px;
				  }
				}
			  </style>
			</head>
			<body>
			  <div class="email-container">
				<div class="header">
				  <h1>Welcome to ArtusTehc!</h1>
				</div>
				<div class="body">
				  <h2>Hello %s,</h2>
				  <p>
					We are thrilled to have you with us! At <strong>ArtusTech</strong>, we strive to provide the best experience for our customers. Let’s get started!
				  </p>
				  <a href="#" class="button">Get Started</a>
				  <p>
					If you have any questions, feel free to reply to this email or visit our <a href="#">Help Center</a>.
				  </p>
				</div>
				<div class="footer">
				  <p>Stay connected with us:</p>
				  <div class="social-icons">
					<a href="#"><img src="facebook-icon.png" alt="Facebook"></a>
					<a href="#"><img src="twitter-icon.png" alt="Twitter"></a>
					<a href="#"><img src="instagram-icon.png" alt="Instagram"></a>
				  </div>
				  <p>&copy; 2024 [Your Company Name]. All rights reserved.</p>
				</div>
			  </div>
			</body>
			</html>
			`, username)

	m.SetBody("text/html", body)

	d := gomail.NewDialer("sandbox.smtp.mailtrap.io", 2525, "5bcbb4e8c4661d", "6318cffe27547b")
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil

}
