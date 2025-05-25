package document

func SeedDocuments(repo *Repository) {
	docs := []Document{
		{
			Title: "Return Policy",
			Content: "Our return policy allows returns within 30 days of purchase with a receipt. " +
				"Items must be in original condition with all packaging. " +
				"Refunds are processed to the original payment method within 5-7 business days.",
			Tags: []string{"returns", "policy", "refunds"},
		},
		{
			Title: "Shipping Information",
			Content: "Standard shipping takes 3-5 business days. Express shipping is available for " +
				"an additional fee and delivers within 1-2 business days. International shipping " +
				"may take 7-14 business days and may be subject to customs fees.",
			Tags: []string{"shipping", "delivery", "international"},
		},
		{
			Title: "Account Password Reset",
			Content: "To reset your password, click the 'Forgot Password' link on the login page. " +
				"You will receive an email with a link to create a new password. " +
				"The link expires after 24 hours. If you don't receive the email, please check your spam folder.",
			Tags: []string{"account", "password", "login"},
		},
		{
			Title: "Product Warranty",
			Content: "All electronics come with a 1-year manufacturer warranty covering defects in materials " +
				"and workmanship. The warranty does not cover damage from misuse, accidents, or normal wear and tear. " +
				"Extended warranties are available for purchase.",
			Tags: []string{"warranty", "guarantee", "repairs"},
		},
		{
			Title: "Membership Benefits",
			Content: "Premium members receive free shipping on all orders, early access to sales, " +
				"exclusive discounts, and priority customer support. Membership costs $49.99 per year " +
				"and can be canceled at any time.",
			Tags: []string{"membership", "premium", "benefits"},
		},
	}

	for _, doc := range docs {
		repo.Add(doc)
	}
}
