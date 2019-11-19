package mail

//nolint:lll
const (
	//ExportedBody is the html for the exported info email, http://pojo.sodhanalibrary.com/string.html
	ExportedBody = "<!DOCTYPE html PUBLIC \"-//W3C//DTD XHTML 1.0 Transitional //EN\" \"http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd\">" +
		"" +
		"<html xmlns=\"http://www.w3.org/1999/xhtml\" xmlns:o=\"urn:schemas-microsoft-com:office:office\" xmlns:v=\"urn:schemas-microsoft-com:vml\">" +
		"<head>" +
		"<!--[if gte mso 9]><xml><o:OfficeDocumentSettings><o:AllowPNG/><o:PixelsPerInch>96</o:PixelsPerInch></o:OfficeDocumentSettings></xml><![endif]-->" +
		"<meta content=\"text/html; charset=utf-8\" http-equiv=\"Content-Type\"/>" +
		"<meta content=\"width=device-width\" name=\"viewport\"/>" +
		"<!--[if !mso]><!-->" +
		"<meta content=\"IE=edge\" http-equiv=\"X-UA-Compatible\"/>" +
		"<!--<![endif]-->" +
		"<title></title>" +
		"<!--[if !mso]><!-->" +
		"<link href=\"https://fonts.googleapis.com/css?family=Open+Sans\" rel=\"stylesheet\" type=\"text/css\"/>" +
		"<!--<![endif]-->" +
		"<style type=\"text/css\">" +
		"		body {" +
		"			margin: 0;" +
		"			padding: 0;" +
		"		}" +
		"" +
		"		table," +
		"		td," +
		"		tr {" +
		"			vertical-align: top;" +
		"			border-collapse: collapse;" +
		"		}" +
		"" +
		"		* {" +
		"			line-height: inherit;" +
		"		}" +
		"" +
		"		a[x-apple-data-detectors=true] {" +
		"			color: inherit !important;" +
		"			text-decoration: none !important;" +
		"		}" +
		"	</style>" +
		"<style id=\"media-query\" type=\"text/css\">" +
		"		@media (max-width: 520px) {" +
		"" +
		"			.block-grid," +
		"			.col {" +
		"				min-width: 320px !important;" +
		"				max-width: 100% !important;" +
		"				display: block !important;" +
		"			}" +
		"" +
		"			.block-grid {" +
		"				width: 100% !important;" +
		"			}" +
		"" +
		"			.col {" +
		"				width: 100% !important;" +
		"			}" +
		"" +
		"			.col>div {" +
		"				margin: 0 auto;" +
		"			}" +
		"" +
		"			img.fullwidth," +
		"			img.fullwidthOnMobile {" +
		"				max-width: 100% !important;" +
		"			}" +
		"" +
		"			.no-stack .col {" +
		"				min-width: 0 !important;" +
		"				display: table-cell !important;" +
		"			}" +
		"" +
		"			.no-stack.two-up .col {" +
		"				width: 50% !important;" +
		"			}" +
		"" +
		"			.no-stack .col.num4 {" +
		"				width: 33% !important;" +
		"			}" +
		"" +
		"			.no-stack .col.num8 {" +
		"				width: 66% !important;" +
		"			}" +
		"" +
		"			.no-stack .col.num4 {" +
		"				width: 33% !important;" +
		"			}" +
		"" +
		"			.no-stack .col.num3 {" +
		"				width: 25% !important;" +
		"			}" +
		"" +
		"			.no-stack .col.num6 {" +
		"				width: 50% !important;" +
		"			}" +
		"" +
		"			.no-stack .col.num9 {" +
		"				width: 75% !important;" +
		"			}" +
		"" +
		"			.video-block {" +
		"				max-width: none !important;" +
		"			}" +
		"" +
		"			.mobile_hide {" +
		"				min-height: 0px;" +
		"				max-height: 0px;" +
		"				max-width: 0px;" +
		"				display: none;" +
		"				overflow: hidden;" +
		"				font-size: 0px;" +
		"			}" +
		"" +
		"			.desktop_hide {" +
		"				display: block !important;" +
		"				max-height: none !important;" +
		"			}" +
		"		}" +
		"	</style>" +
		"</head>" +
		"<body class=\"clean-body\" style=\"margin: 0; padding: 0; -webkit-text-size-adjust: 100%; background-color: #F6F6FF;\">" +
		"<!--[if IE]><div class=\"ie-browser\"><![endif]-->" +
		"<table bgcolor=\"#F6F6FF\" cellpadding=\"0\" cellspacing=\"0\" class=\"nl-container\" role=\"presentation\" style=\"table-layout: fixed; vertical-align: top; min-width: 320px; Margin: 0 auto; border-spacing: 0; border-collapse: collapse; mso-table-lspace: 0pt; mso-table-rspace: 0pt; background-color: #F6F6FF; width: 100%;\" valign=\"top\" width=\"100%\">" +
		"<tbody>" +
		"<tr style=\"vertical-align: top;\" valign=\"top\">" +
		"<td style=\"word-break: break-word; vertical-align: top;\" valign=\"top\">" +
		"<!--[if (mso)|(IE)]><table width=\"100%\" cellpadding=\"0\" cellspacing=\"0\" border=\"0\"><tr><td align=\"center\" style=\"background-color:#F6F6FF\"><![endif]-->" +
		"<div style=\"background-color:#F6F6FF;\">" +
		"<div class=\"block-grid\" style=\"Margin: 0 auto; min-width: 320px; max-width: 500px; overflow-wrap: break-word; word-wrap: break-word; word-break: break-word; background-color: #FFFFFF;\">" +
		"<div style=\"border-collapse: collapse;display: table;width: 100%;background-color:#FFFFFF;\">" +
		"<!--[if (mso)|(IE)]><table width=\"100%\" cellpadding=\"0\" cellspacing=\"0\" border=\"0\" style=\"background-color:#F6F6FF;\"><tr><td align=\"center\"><table cellpadding=\"0\" cellspacing=\"0\" border=\"0\" style=\"width:500px\"><tr class=\"layout-full-width\" style=\"background-color:#FFFFFF\"><![endif]-->" +
		"<!--[if (mso)|(IE)]><td align=\"center\" width=\"500\" style=\"background-color:#FFFFFF;width:500px; border-top: 0px solid transparent; border-left: 0px solid transparent; border-bottom: 0px solid transparent; border-right: 0px solid transparent;\" valign=\"top\"><table width=\"100%\" cellpadding=\"0\" cellspacing=\"0\" border=\"0\"><tr><td style=\"padding-right: 0px; padding-left: 0px; padding-top:5px; padding-bottom:5px;\"><![endif]-->" +
		"<div class=\"col num12\" style=\"min-width: 320px; max-width: 500px; display: table-cell; vertical-align: top; width: 500px;\">" +
		"<div style=\"width:100% !important;\">" +
		"<!--[if (!mso)&(!IE)]><!-->" +
		"<div style=\"border-top:0px solid transparent; border-left:0px solid transparent; border-bottom:0px solid transparent; border-right:0px solid transparent; padding-top:5px; padding-bottom:5px; padding-right: 0px; padding-left: 0px;\">" +
		"<!--<![endif]-->" +
		"<div align=\"center\" class=\"img-container center autowidth fullwidth\" style=\"padding-right: 0px;padding-left: 0px;\">" +
		"<!--[if mso]><table width=\"100%\" cellpadding=\"0\" cellspacing=\"0\" border=\"0\"><tr style=\"line-height:0px\"><td style=\"padding-right: 0px;padding-left: 0px;\" align=\"center\"><![endif]--><a href=\"www.eco.com\" tabindex=\"-1\" target=\"_blank\"> <img align=\"center\" alt=\"Linked Up\" border=\"0\" class=\"center autowidth fullwidth\" src=\"https://linkedup-email.s3-us-west-2.amazonaws.com/images/Email+Header.png\" style=\"text-decoration: none; -ms-interpolation-mode: bicubic; height: auto; border: none; width: 100%; max-width: 500px; display: block;\" title=\"Linked Up\" width=\"500\"/></a>" +
		"<!--[if mso]></td></tr></table><![endif]-->" +
		"</div>" +
		"<div align=\"center\" class=\"img-container center autowidth\" style=\"padding-right: 0px;padding-left: 0px;\">" +
		"<!--[if mso]><table width=\"100%\" cellpadding=\"0\" cellspacing=\"0\" border=\"0\"><tr style=\"line-height:0px\"><td style=\"padding-right: 0px;padding-left: 0px;\" align=\"center\"><![endif]--><img align=\"center\" alt=\"LInkedUp Logo\" border=\"0\" class=\"center autowidth\" src=\"https://linkedup-email.s3-us-west-2.amazonaws.com/images/Linked+Up+Graphic+Hexagon.png\" style=\"text-decoration: none; -ms-interpolation-mode: bicubic; border: 0; height: auto; width: 100%; max-width: 196px; display: block;\" title=\"LInkedUp Logo\" width=\"196\"/>" +
		"<!--[if mso]></td></tr></table><![endif]-->" +
		"</div>" +
		"<!--[if mso]><table width=\"100%\" cellpadding=\"0\" cellspacing=\"0\" border=\"0\"><tr><td style=\"padding-right: 10px; padding-left: 10px; padding-top: 10px; padding-bottom: 10px; font-family: Arial, sans-serif\"><![endif]-->" +
		"<div style=\"color:#555555;font-family:Arial, 'Helvetica Neue', Helvetica, sans-serif;line-height:1.2;padding-top:10px;padding-right:10px;padding-bottom:10px;padding-left:10px;\">" +
		"<div style=\"line-height: 1.2; font-family: Arial, 'Helvetica Neue', Helvetica, sans-serif; font-size: 12px; color: #555555; mso-line-height-alt: 14px;\">" +
		"<p style=\"line-height: 1.2; font-size: 18px; text-align: center; mso-line-height-alt: 22px; margin: 0;\"><span style=\"font-size: 18px;\">You can find the requested contact export attached to this email. </span></p>" +
		"<p style=\"line-height: 1.2; font-size: 12px; text-align: center; mso-line-height-alt: 14px; margin: 0;\"> </p>" +
		"<p style=\"line-height: 1.2; font-size: 18px; text-align: center; mso-line-height-alt: 22px; margin: 0;\"><span style=\"font-size: 18px;\">If you haven't already, share your thoughts on the game in our <a href=\"https://docs.google.com/forms/d/e/1FAIpQLSeazQOpq5qGO-SDXRJOcr0BWqnUHoSP80r0RD9Z-NCinjT7Mw/viewform?usp=sf_link\" rel=\"noopener\" style=\"text-decoration: underline; color: #0068A5;\" target=\"_blank\">user survey</a>!</span></p>" +
		"</div>" +
		"</div>" +
		"<!--[if mso]></td></tr></table><![endif]-->" +
		"<table border=\"0\" cellpadding=\"0\" cellspacing=\"0\" class=\"divider\" role=\"presentation\" style=\"table-layout: fixed; vertical-align: top; border-spacing: 0; border-collapse: collapse; mso-table-lspace: 0pt; mso-table-rspace: 0pt; min-width: 100%; -ms-text-size-adjust: 100%; -webkit-text-size-adjust: 100%;\" valign=\"top\" width=\"100%\">" +
		"<tbody>" +
		"<tr style=\"vertical-align: top;\" valign=\"top\">" +
		"<td class=\"divider_inner\" style=\"word-break: break-word; vertical-align: top; min-width: 100%; -ms-text-size-adjust: 100%; -webkit-text-size-adjust: 100%; padding-top: 10px; padding-right: 10px; padding-bottom: 10px; padding-left: 10px;\" valign=\"top\">" +
		"<table align=\"center\" border=\"0\" cellpadding=\"0\" cellspacing=\"0\" class=\"divider_content\" height=\"0\" role=\"presentation\" style=\"table-layout: fixed; vertical-align: top; border-spacing: 0; border-collapse: collapse; mso-table-lspace: 0pt; mso-table-rspace: 0pt; border-top: 1px solid #BBBBBB; height: 0px; width: 100%;\" valign=\"top\" width=\"100%\">" +
		"<tbody>" +
		"<tr style=\"vertical-align: top;\" valign=\"top\">" +
		"<td height=\"0\" style=\"word-break: break-word; vertical-align: top; -ms-text-size-adjust: 100%; -webkit-text-size-adjust: 100%;\" valign=\"top\"><span></span></td>" +
		"</tr>" +
		"</tbody>" +
		"</table>" +
		"</td>" +
		"</tr>" +
		"</tbody>" +
		"</table>" +
		"<!--[if mso]><table width=\"100%\" cellpadding=\"0\" cellspacing=\"0\" border=\"0\"><tr><td style=\"padding-right: 10px; padding-left: 10px; padding-top: 10px; padding-bottom: 10px; font-family: Arial, sans-serif\"><![endif]-->" +
		"<div style=\"color:#555555;font-family:Arial, 'Helvetica Neue', Helvetica, sans-serif;line-height:1.2;padding-top:10px;padding-right:10px;padding-bottom:10px;padding-left:10px;\">" +
		"<div style=\"line-height: 1.2; font-family: Arial, 'Helvetica Neue', Helvetica, sans-serif; font-size: 12px; color: #555555; mso-line-height-alt: 14px;\">" +
		"<p style=\"line-height: 1.2; text-align: center; font-size: 18px; mso-line-height-alt: 22px; margin: 0;\"><span style=\"font-size: 18px;\">To learn more about how Eco is building the future of spendable currencies, head over to our <a href=\"https://eco.com/\" rel=\"noopener\" style=\"text-decoration: underline; color: #0068A5;\" target=\"_blank\">website</a> to sign up for early access and see how you can <a href=\"https://eco.com/jobs\" rel=\"noopener\" style=\"text-decoration: underline; color: #0068A5;\" target=\"_blank\">join us</a>.<br/><br/></span></p>" +
		"<p style=\"font-size: 18px; line-height: 1.2; text-align: center; mso-line-height-alt: 22px; margin: 0;\"><span style=\"font-size: 18px;\">Don’t forget to check out the attached infographic as well as our game write-ups on <a href=\"https://medium.com/@eco\" rel=\"noopener\" style=\"text-decoration: underline; color: #0068A5;\" target=\"_blank\" title=\"Eco on Medium\">Medium</a>.</span><span style=\"font-size: 18px;\"> </span></p>" +
		"</div>" +
		"</div>" +
		"<!--[if mso]></td></tr></table><![endif]-->" +
		"<table border=\"0\" cellpadding=\"0\" cellspacing=\"0\" class=\"divider\" role=\"presentation\" style=\"table-layout: fixed; vertical-align: top; border-spacing: 0; border-collapse: collapse; mso-table-lspace: 0pt; mso-table-rspace: 0pt; min-width: 100%; -ms-text-size-adjust: 100%; -webkit-text-size-adjust: 100%;\" valign=\"top\" width=\"100%\">" +
		"<tbody>" +
		"<tr style=\"vertical-align: top;\" valign=\"top\">" +
		"<td class=\"divider_inner\" style=\"word-break: break-word; vertical-align: top; min-width: 100%; -ms-text-size-adjust: 100%; -webkit-text-size-adjust: 100%; padding-top: 10px; padding-right: 10px; padding-bottom: 10px; padding-left: 10px;\" valign=\"top\">" +
		"<table align=\"center\" border=\"0\" cellpadding=\"0\" cellspacing=\"0\" class=\"divider_content\" height=\"0\" role=\"presentation\" style=\"table-layout: fixed; vertical-align: top; border-spacing: 0; border-collapse: collapse; mso-table-lspace: 0pt; mso-table-rspace: 0pt; border-top: 1px solid #BBBBBB; height: 0px; width: 100%;\" valign=\"top\" width=\"100%\">" +
		"<tbody>" +
		"<tr style=\"vertical-align: top;\" valign=\"top\">" +
		"<td height=\"0\" style=\"word-break: break-word; vertical-align: top; -ms-text-size-adjust: 100%; -webkit-text-size-adjust: 100%;\" valign=\"top\"><span></span></td>" +
		"</tr>" +
		"</tbody>" +
		"</table>" +
		"</td>" +
		"</tr>" +
		"</tbody>" +
		"</table>" +
		"<table cellpadding=\"0\" cellspacing=\"0\" class=\"social_icons\" role=\"presentation\" style=\"table-layout: fixed; vertical-align: top; border-spacing: 0; border-collapse: collapse; mso-table-lspace: 0pt; mso-table-rspace: 0pt;\" valign=\"top\" width=\"100%\">" +
		"<tbody>" +
		"<tr style=\"vertical-align: top;\" valign=\"top\">" +
		"<td style=\"word-break: break-word; vertical-align: top; padding-top: 10px; padding-right: 10px; padding-bottom: 10px; padding-left: 10px;\" valign=\"top\">" +
		"<table activate=\"activate\" align=\"center\" alignment=\"alignment\" cellpadding=\"0\" cellspacing=\"0\" class=\"social_table\" role=\"presentation\" style=\"table-layout: fixed; vertical-align: top; border-spacing: 0; border-collapse: undefined; mso-table-tspace: 0; mso-table-rspace: 0; mso-table-bspace: 0; mso-table-lspace: 0;\" to=\"to\" valign=\"top\">" +
		"<tbody>" +
		"<tr align=\"center\" style=\"vertical-align: top; display: inline-block; text-align: center;\" valign=\"top\">" +
		"<td style=\"word-break: break-word; vertical-align: top; padding-bottom: 5px; padding-right: 3px; padding-left: 3px;\" valign=\"top\"><a href=\"https://twitter.com/eco\" target=\"_blank\"><img alt=\"Twitter\" height=\"32\" src=\"https://linkedup-email.s3-us-west-2.amazonaws.com/images/twitter%402x.png\" style=\"text-decoration: none; -ms-interpolation-mode: bicubic; height: auto; border: none; display: block;\" title=\"Twitter\" width=\"32\"/></a></td>" +
		"<td style=\"word-break: break-word; vertical-align: top; padding-bottom: 5px; padding-right: 3px; padding-left: 3px;\" valign=\"top\"><a href=\"https://medium.com/@eco\" target=\"_blank\"><img alt=\"Medium\" height=\"32\" src=\"https://linkedup-email.s3-us-west-2.amazonaws.com/images/medium%402x.png\" style=\"text-decoration: none; -ms-interpolation-mode: bicubic; height: auto; border: none; display: block;\" title=\"Medium\" width=\"32\"/></a></td>" +
		"</tr>" +
		"</tbody>" +
		"</table>" +
		"</td>" +
		"</tr>" +
		"</tbody>" +
		"</table>" +
		"<table border=\"0\" cellpadding=\"0\" cellspacing=\"0\" class=\"divider\" role=\"presentation\" style=\"table-layout: fixed; vertical-align: top; border-spacing: 0; border-collapse: collapse; mso-table-lspace: 0pt; mso-table-rspace: 0pt; min-width: 100%; -ms-text-size-adjust: 100%; -webkit-text-size-adjust: 100%;\" valign=\"top\" width=\"100%\">" +
		"<tbody>" +
		"<tr style=\"vertical-align: top;\" valign=\"top\">" +
		"<td class=\"divider_inner\" style=\"word-break: break-word; vertical-align: top; min-width: 100%; -ms-text-size-adjust: 100%; -webkit-text-size-adjust: 100%; padding-top: 10px; padding-right: 10px; padding-bottom: 10px; padding-left: 10px;\" valign=\"top\">" +
		"<table align=\"center\" border=\"0\" cellpadding=\"0\" cellspacing=\"0\" class=\"divider_content\" height=\"0\" role=\"presentation\" style=\"table-layout: fixed; vertical-align: top; border-spacing: 0; border-collapse: collapse; mso-table-lspace: 0pt; mso-table-rspace: 0pt; border-top: 1px solid #BBBBBB; height: 0px; width: 100%;\" valign=\"top\" width=\"100%\">" +
		"<tbody>" +
		"<tr style=\"vertical-align: top;\" valign=\"top\">" +
		"<td height=\"0\" style=\"word-break: break-word; vertical-align: top; -ms-text-size-adjust: 100%; -webkit-text-size-adjust: 100%;\" valign=\"top\"><span></span></td>" +
		"</tr>" +
		"</tbody>" +
		"</table>" +
		"</td>" +
		"</tr>" +
		"</tbody>" +
		"</table>" +
		"<div align=\"center\" class=\"img-container center autowidth fullwidth\" style=\"padding-right: 0px;padding-left: 0px;\">" +
		"<!--[if mso]><table width=\"100%\" cellpadding=\"0\" cellspacing=\"0\" border=\"0\"><tr style=\"line-height:0px\"><td style=\"padding-right: 0px;padding-left: 0px;\" align=\"center\"><![endif]--><a href=\"https://eco.com/\" tabindex=\"-1\" target=\"_blank\"> <img align=\"center\" alt=\"Powered by Eco\" border=\"0\" class=\"center autowidth fullwidth\" src=\"https://linkedup-email.s3-us-west-2.amazonaws.com/images/Email+Footer.png\" style=\"text-decoration: none; -ms-interpolation-mode: bicubic; height: auto; border: none; width: 100%; max-width: 500px; display: block;\" title=\"Powered by Eco\" width=\"500\"/></a>" +
		"<!--[if mso]></td></tr></table><![endif]-->" +
		"</div>" +
		"<!--[if mso]><table width=\"100%\" cellpadding=\"0\" cellspacing=\"0\" border=\"0\"><tr><td style=\"padding-right: 10px; padding-left: 10px; padding-top: 10px; padding-bottom: 10px; font-family: Arial, sans-serif\"><![endif]-->" +
		"<div style=\"color:#555555;font-family:Arial, 'Helvetica Neue', Helvetica, sans-serif;line-height:1.2;padding-top:10px;padding-right:10px;padding-bottom:10px;padding-left:10px;\">" +
		"<div style=\"line-height: 1.2; font-family: Arial, 'Helvetica Neue', Helvetica, sans-serif; font-size: 12px; color: #555555; mso-line-height-alt: 14px;\">" +
		"<p style=\"line-height: 1.2; text-align: center; font-size: 12px; mso-line-height-alt: 14px; margin: 0;\"> </p>" +
		"</div>" +
		"</div>" +
		"<!--[if mso]></td></tr></table><![endif]-->" +
		"<!--[if (!mso)&(!IE)]><!-->" +
		"</div>" +
		"<!--<![endif]-->" +
		"</div>" +
		"</div>" +
		"<!--[if (mso)|(IE)]></td></tr></table><![endif]-->" +
		"<!--[if (mso)|(IE)]></td></tr></table></td></tr></table><![endif]-->" +
		"</div>" +
		"</div>" +
		"</div>" +
		"<!--[if (mso)|(IE)]></td></tr></table><![endif]-->" +
		"</td>" +
		"</tr>" +
		"</tbody>" +
		"</table>" +
		"<!--[if (IE)]></div><![endif]-->" +
		"</body>" +
		"</html>"
)
