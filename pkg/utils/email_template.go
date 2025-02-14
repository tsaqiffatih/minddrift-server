package utils

const EmailVerification = `
<!DOCTYPE html>
<html>
<head>
    <title>Email Verification</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            background-color: #f4f4f4;
            text-align: center;
            padding: 20px;
        }
        .container {
            max-width: 500px;
            margin: auto;
            padding: 20px;
            border-radius: 10px;
            box-shadow: 0px 4px 10px rgba(0, 0, 0, 0.1);
            background-color: #ffffff;
            font-size: 16px;
        }
        .button {
            background-color: #0056b3;
            color: white !important;
            padding: 12px 24px;
            text-decoration: none;
            border-radius: 5px;
            display: inline-block;
            margin-top: 20px;
            font-weight: bold;
        }
        .button:hover {
            background-color: #003f7f;
        }
        .footer {
            margin-top: 20px;
            font-size: 12px;
            color: #888888;
        }
    </style>
</head>
<body>
    <div class="container">
        <img src="https://minddrift.com/logo.png" alt="MindDrift Logo" style="width: 150px; margin: 20px auto; display: block;">
        <h2>Hello, {{.Username}}!</h2>
        <p style="font-size: 18px; line-height: 1.6;">Welcome to <b>MindDrift</b>! To complete your registration and secure your account, please verify your email address.</p>
        <table role="presentation" cellspacing="0" cellpadding="0" border="0" align="center">
            <tr>
                <td style="border-radius: 6px; background-color: #007bff; text-align: center;">
                    <a href="{{.VerificationLink}}" 
                       style="display: inline-block; font-size: 16px; font-weight: bold;
                              color: white !important; text-decoration: none; 
                              padding: 12px 24px; border-radius: 6px;">
                       Confirm Your Email
                    </a>
                </td>
            </tr>
        </table>
        <p style="margin-top: 25px;">This link will expire in <b>24 hours</b>.</p>
        <p>If the button above doesn't work, copy and paste this link into your browser:</p>
        <p><a href="{{.VerificationLink}}" style="color: #007bff; word-break: break-all;">{{.VerificationLink}}</a></p>
        <p>If you did not register, please ignore this email.</p>
        <p>Thank you!</p>
        <p>MindDrift Team</p>
        <div class="footer">
            &copy; 2025 MindDrift. All rights reserved. <br>
            Need help? Contact us at <a href="mailto:{{.MindDriftEmail}}">{{.MindDriftEmail}}</a>
        </div>
    </div>
</body>
</html>`

const EmailResetPassword = `
<!DOCTYPE html>
<html>
<head>
    <title>Reset Your Password</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            background-color: #f4f4f4;
            text-align: center;
            padding: 20px;
        }
        .container {
            max-width: 500px;
            margin: auto;
            padding: 20px;
            border-radius: 10px;
            box-shadow: 0px 4px 10px rgba(0, 0, 0, 0.1);
            background-color: #ffffff;
            font-size: 16px;
        }
        .button {
            background-color: #d9534f;
            color: white !important;
            padding: 12px 24px;
            text-decoration: none;
            border-radius: 5px;
            display: inline-block;
            margin-top: 20px;
            font-weight: bold;
        }
        .button:hover {
            background-color: #c9302c;
        }
        .footer {
            margin-top: 20px;
            font-size: 12px;
            color: #888888;
        }
    </style>
</head>
<body>
    <div class="container">
        <img src="https://minddrift.com/logo.png" alt="MindDrift Logo" style="width: 150px; margin: 20px auto; display: block;">
        <h2>Hello, {{.Username}}!</h2>
        <p style="font-size: 18px; line-height: 1.6;">You recently requested to reset your password for <b>MindDrift</b>. Click the button below to reset it:</p>
        <table role="presentation" cellspacing="0" cellpadding="0" border="0" align="center">
            <tr>
                <td style="border-radius: 6px; background-color: #d9534f; text-align: center;">
                    <a href="{{.ResetLink}}" 
                       style="display: inline-block; font-size: 16px; font-weight: bold;
                              color: white !important; text-decoration: none; 
                              padding: 12px 24px; border-radius: 6px;">
                       Reset Your Password
                    </a>
                </td>
            </tr>
        </table>
        <p style="margin-top: 25px;">This link will expire in <b>1 hour</b>.</p>
        <p>If you did not request a password reset, you can safely ignore this email.</p>
        <p>Thank you!</p>
        <p>MindDrift Team</p>
        <div class="footer">
            &copy; 2025 MindDrift. All rights reserved. <br>
            Need help? Contact us at <a href="mailto:{{.MindDriftEmail}}">{{.MindDriftEmail}}</a>
        </div>
    </div>
</body>
</html>`
