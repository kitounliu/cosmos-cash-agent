package model

// SchemaField is a field within the credential schema
type SchemaField struct {
	Name        string
	Title       string
	Description string
	Value       string
	ReadOnly    bool
}

func NewSchemaField(name, title, description string) SchemaField {
	return SchemaField{
		Name:        name,
		Title:       title,
		Description: description,
		ReadOnly:    false,
	}
}
func NewSchemaLabel(name, title, description, value string) SchemaField {
	return SchemaField{
		Name:        name,
		Title:       title,
		Description: description,
		Value:       value,
		ReadOnly:    true,
	}
}

// CredentialSchema represent a verifiable credential schema
type CredentialSchema struct {
	// Should match the Verifiable Credentials that belongs to
	Name   string
	Fields []SchemaField
}

// LicenseSchema returns a license schema
func LicenseSchema(licenseType, authority string) CredentialSchema {
	//"license_type": "MICAEMI",
	//"country": "EU",
	//"authority": "Another Financial Services Body (AFFB)",
	//"circulation_limit": {
	//	"denom": "sEUR",
	//	"amount": "1000000000"
	//}
	return CredentialSchema{
		Name: "LicenseCredential",
		Fields: []SchemaField{
			NewSchemaLabel("license_type", "License Type", "license acronym", licenseType),
			NewSchemaLabel("authority", "Authority", "authority issuing the license", authority),
			NewSchemaField("country", "Country", "country where the license applies"),
			NewSchemaField("denom", "Coin denomination", "the coin symbol"),
			NewSchemaField("amount", "Amount", "amount approved by the license"),
		},
	}
}
