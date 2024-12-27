package resources

import (
	"bytes"
	"html/template"
	"os"
	"path/filepath"
)

func ParseTemplate(templateName string, data interface{}) (string, error) {
	resourcesDir := filepath.Join("internal", "resources")
	templatePath := filepath.Join(resourcesDir, templateName)

	content, err := os.ReadFile(templatePath)
	if err != nil {
		return "", err
	}

	tmpl, err := template.New(templateName).Parse(string(content))
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func GenerateOrderEmail(isShop bool, data OrderEmailData) (string, error) {
	templateName := "user_order_template.html"
	if isShop {
		templateName = "shop_order_template.html"
	}

	return ParseTemplate(templateName, data)
}
