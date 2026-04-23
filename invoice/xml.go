package invoice

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"text/template"
)

type tmplAC struct {
	Indicator     bool
	Reason        string
	ReasonCode    string
	Amount        string
	VATRate       string
	TaxCategoryID string
}
type Party struct {
	CRN              string
	VATNumber        string
	RegistrationName string
	Street           string
	BuildingNumber   string
	PlotID           string
	District         string
	City             string
	PostalCode       string
}

type tmplLine struct {
	ID            string
	UnitCode      string
	Quantity      float64
	LineTotal     float64
	TaxAmount     float64
	TotalWithVAT  float64
	ItemName      string
	VATPercent    string
	Price         float64
	TaxCategoryID string
	Discounts     []tmplAC
}

type tmplData struct {
	SigningTime              string
	CertificateHash          string
	IssuerName               string
	SerialNumber             string
	X509Certificate          string
	InvoiceDigest            string
	SignedPropsDigest        string
	SignatureValue           string
	ID                       string
	UUID                     string
	IssueDate                string
	IssueTime                string
	ICV                      string
	PreviousInvoiceHash      string
	QRCode                   string
	CRN                      string
	VATNumber                string
	RegistrationName         string
	Street                   string
	BuildingNumber           string
	PlotID                   string
	District                 string
	City                     string
	PostalCode               string
	PaymentMeansCode         string
	LineExtensionAmount      string
	TaxExclusiveAmount       string
	TaxInclusiveAmount       string
	AllowanceTotal           string
	ChargeTotal              string
	TaxAmount                string
	Lines                    []tmplLine
	InvoiceLevelACs          []tmplAC
	BillingReferenceID       string
	InvoiceTypeCode          string
	InstructionNote          string
	SignedPropertiesXML      string
	HasCustomer              bool
	CustomerStreet           string
	CustomerBuildingNumber   string
	CustomerDistrict         string
	CustomerCity             string
	CustomerPostalCode       string
	CustomerVATNumber        string
	CustomerRegistrationName string
	InvoiceTypeName          string
	TaxableAmountS           string
	TaxAmountS               string
	TaxableAmountO           string
}

const invoiceTemplate = `<?xml version="1.0" encoding="UTF-8"?>
<Invoice xmlns="urn:oasis:names:specification:ubl:schema:xsd:Invoice-2" xmlns:cac="urn:oasis:names:specification:ubl:schema:xsd:CommonAggregateComponents-2" xmlns:cbc="urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2" xmlns:ext="urn:oasis:names:specification:ubl:schema:xsd:CommonExtensionComponents-2">
<ext:UBLExtensions>
<ext:UBLExtension>
<ext:ExtensionURI>urn:oasis:names:specification:ubl:dsig:enveloped:xades</ext:ExtensionURI>
<ext:ExtensionContent>
<sig:UBLDocumentSignatures xmlns:sig="urn:oasis:names:specification:ubl:schema:xsd:CommonSignatureComponents-2" xmlns:sac="urn:oasis:names:specification:ubl:schema:xsd:SignatureAggregateComponents-2" xmlns:sbc="urn:oasis:names:specification:ubl:schema:xsd:SignatureBasicComponents-2">
<sac:SignatureInformation>
<cbc:ID>urn:oasis:names:specification:ubl:signature:1</cbc:ID>
<sbc:ReferencedSignatureID>urn:oasis:names:specification:ubl:signature:Invoice</sbc:ReferencedSignatureID>
<ds:Signature xmlns:ds="http://www.w3.org/2000/09/xmldsig#" Id="signature">
<ds:SignedInfo>
<ds:CanonicalizationMethod Algorithm="http://www.w3.org/2006/12/xml-c14n11"/>
<ds:SignatureMethod Algorithm="http://www.w3.org/2001/04/xmldsig-more#ecdsa-sha256"/>
<ds:Reference Id="invoiceSignedData" URI="">
<ds:Transforms>
<ds:Transform Algorithm="http://www.w3.org/TR/1999/REC-xpath-19991116">
<ds:XPath>not(//ancestor-or-self::ext:UBLExtensions)</ds:XPath>
</ds:Transform>
<ds:Transform Algorithm="http://www.w3.org/TR/1999/REC-xpath-19991116">
<ds:XPath>not(//ancestor-or-self::cac:Signature)</ds:XPath>
</ds:Transform>
<ds:Transform Algorithm="http://www.w3.org/TR/1999/REC-xpath-19991116">
<ds:XPath>not(//ancestor-or-self::cac:AdditionalDocumentReference[cbc:ID='QR'])</ds:XPath>
</ds:Transform>
<ds:Transform Algorithm="http://www.w3.org/2006/12/xml-c14n11"/>
</ds:Transforms>
<ds:DigestMethod Algorithm="http://www.w3.org/2001/04/xmlenc#sha256"/>
<ds:DigestValue>{{.InvoiceDigest}}</ds:DigestValue>
</ds:Reference>
<ds:Reference Type="http://uri.etsi.org/01903#SignedProperties" URI="#xadesSignedProperties">
<ds:Transforms>
<ds:Transform Algorithm="http://www.w3.org/2006/12/xml-c14n11"/>
</ds:Transforms>
<ds:DigestMethod Algorithm="http://www.w3.org/2001/04/xmlenc#sha256"/>
<ds:DigestValue>{{.SignedPropsDigest}}</ds:DigestValue>
</ds:Reference>
</ds:SignedInfo>
<ds:SignatureValue>{{.SignatureValue}}</ds:SignatureValue>
<ds:KeyInfo>
<ds:X509Data>
<ds:X509Certificate>{{.X509Certificate}}</ds:X509Certificate>
</ds:X509Data>
</ds:KeyInfo>
<ds:Object>
<xades:QualifyingProperties xmlns:xades="http://uri.etsi.org/01903/v1.3.2#" Target="signature">
{{.SignedPropertiesXML}}
</xades:QualifyingProperties>
</ds:Object>
</ds:Signature>
</sac:SignatureInformation>
</sig:UBLDocumentSignatures>
</ext:ExtensionContent>
</ext:UBLExtension>
</ext:UBLExtensions>
<cbc:ProfileID>reporting:1.0</cbc:ProfileID>
<cbc:ID>{{.ID}}</cbc:ID>
<cbc:UUID>{{.UUID}}</cbc:UUID>
<cbc:IssueDate>{{.IssueDate}}</cbc:IssueDate>
<cbc:IssueTime>{{.IssueTime}}</cbc:IssueTime>
<cbc:InvoiceTypeCode name="{{.InvoiceTypeName}}">{{.InvoiceTypeCode}}</cbc:InvoiceTypeCode>
<cbc:DocumentCurrencyCode>SAR</cbc:DocumentCurrencyCode>
<cbc:TaxCurrencyCode>SAR</cbc:TaxCurrencyCode>
{{if .BillingReferenceID}}
<cac:BillingReference>
  <cac:InvoiceDocumentReference>
    <cbc:ID>{{.BillingReferenceID}}</cbc:ID>
  </cac:InvoiceDocumentReference>
</cac:BillingReference>
{{end}}
<cac:AdditionalDocumentReference>
<cbc:ID>ICV</cbc:ID>
<cbc:UUID>{{.ICV}}</cbc:UUID>
</cac:AdditionalDocumentReference>
<cac:AdditionalDocumentReference>
<cbc:ID>PIH</cbc:ID>
<cac:Attachment>
<cbc:EmbeddedDocumentBinaryObject mimeCode="text/plain">{{.PreviousInvoiceHash}}</cbc:EmbeddedDocumentBinaryObject>
</cac:Attachment>
</cac:AdditionalDocumentReference>
<cac:AdditionalDocumentReference>
<cbc:ID>QR</cbc:ID>
<cac:Attachment>
<cbc:EmbeddedDocumentBinaryObject mimeCode="text/plain">{{.QRCode}}</cbc:EmbeddedDocumentBinaryObject>
</cac:Attachment>
</cac:AdditionalDocumentReference>
<cac:Signature>
<cbc:ID>urn:oasis:names:specification:ubl:signature:Invoice</cbc:ID>
<cbc:SignatureMethod>urn:oasis:names:specification:ubl:dsig:enveloped:xades</cbc:SignatureMethod>
</cac:Signature>
<cac:AccountingSupplierParty>
<cac:Party>
<cac:PartyIdentification>
<cbc:ID schemeID="CRN">{{.CRN}}</cbc:ID>
</cac:PartyIdentification>
<cac:PostalAddress>
<cbc:StreetName>{{.Street}}</cbc:StreetName>
<cbc:BuildingNumber>{{.BuildingNumber}}</cbc:BuildingNumber>
<cbc:PlotIdentification>{{.PlotID}}</cbc:PlotIdentification>
<cbc:CitySubdivisionName>{{.District}}</cbc:CitySubdivisionName>
<cbc:CityName>{{.City}}</cbc:CityName>
<cbc:PostalZone>{{.PostalCode}}</cbc:PostalZone>
<cac:Country>
<cbc:IdentificationCode>SA</cbc:IdentificationCode>
</cac:Country>
</cac:PostalAddress>
<cac:PartyTaxScheme>
<cbc:CompanyID>{{.VATNumber}}</cbc:CompanyID>
<cac:TaxScheme>
<cbc:ID>VAT</cbc:ID>
</cac:TaxScheme>
</cac:PartyTaxScheme>
<cac:PartyLegalEntity>
<cbc:RegistrationName>{{.RegistrationName}}</cbc:RegistrationName>
</cac:PartyLegalEntity>
</cac:Party>
</cac:AccountingSupplierParty>
{{if .HasCustomer}}
<cac:AccountingCustomerParty>
<cac:Party>
<cac:PostalAddress>
<cbc:StreetName>{{.CustomerStreet}}</cbc:StreetName>
<cbc:BuildingNumber>{{.CustomerBuildingNumber}}</cbc:BuildingNumber>
<cbc:CitySubdivisionName>{{.CustomerDistrict}}</cbc:CitySubdivisionName>
<cbc:CityName>{{.CustomerCity}}</cbc:CityName>
<cbc:PostalZone>{{.CustomerPostalCode}}</cbc:PostalZone>
<cac:Country>
<cbc:IdentificationCode>SA</cbc:IdentificationCode>
</cac:Country>
</cac:PostalAddress>
<cac:PartyTaxScheme>
<cbc:CompanyID>{{.CustomerVATNumber}}</cbc:CompanyID>
<cac:TaxScheme>
<cbc:ID>VAT</cbc:ID>
</cac:TaxScheme>
</cac:PartyTaxScheme>
<cac:PartyLegalEntity>
<cbc:RegistrationName>{{.CustomerRegistrationName}}</cbc:RegistrationName>
</cac:PartyLegalEntity>
</cac:Party>
</cac:AccountingCustomerParty>
{{else}}
<cac:AccountingCustomerParty/>
{{end}}
<cac:PaymentMeans>
<cbc:PaymentMeansCode>{{.PaymentMeansCode}}</cbc:PaymentMeansCode>
{{if .InstructionNote}}
<cbc:InstructionNote>{{.InstructionNote}}</cbc:InstructionNote>
{{end}}
</cac:PaymentMeans>
{{range .InvoiceLevelACs}}
<cac:AllowanceCharge>
<cbc:ChargeIndicator>{{.Indicator}}</cbc:ChargeIndicator>
{{if .ReasonCode}}<cbc:AllowanceChargeReasonCode>{{.ReasonCode}}</cbc:AllowanceChargeReasonCode>{{end}}
<cbc:AllowanceChargeReason>{{.Reason}}</cbc:AllowanceChargeReason>
<cbc:Amount currencyID="SAR">{{.Amount}}</cbc:Amount>
<cac:TaxCategory>
<cbc:ID>{{.TaxCategoryID}}</cbc:ID>
<cbc:Percent>{{.VATRate}}</cbc:Percent>
<cac:TaxScheme>
<cbc:ID>VAT</cbc:ID>
</cac:TaxScheme>
</cac:TaxCategory>
</cac:AllowanceCharge>
{{end}}
<cac:TaxTotal>
    <cbc:TaxAmount currencyID="SAR">{{.TaxAmount}}</cbc:TaxAmount>
    <cac:TaxSubtotal>
        <cbc:TaxableAmount currencyID="SAR">{{.TaxableAmountS}}</cbc:TaxableAmount>
        <cbc:TaxAmount currencyID="SAR">{{.TaxAmountS}}</cbc:TaxAmount>
        <cac:TaxCategory>
            <cbc:ID>S</cbc:ID>
            <cbc:Percent>15.00</cbc:Percent>
            <cac:TaxScheme>
                <cbc:ID>VAT</cbc:ID>
            </cac:TaxScheme>
        </cac:TaxCategory>
    </cac:TaxSubtotal>
    <cac:TaxSubtotal>
        <cbc:TaxableAmount currencyID="SAR">{{.TaxableAmountO}}</cbc:TaxableAmount>
        <cbc:TaxAmount currencyID="SAR">0.00</cbc:TaxAmount>
        <cac:TaxCategory>
            <cbc:ID>O</cbc:ID>
            <cbc:Percent>0.00</cbc:Percent>
            <cac:TaxScheme>
                <cbc:ID>VAT</cbc:ID>
            </cac:TaxScheme>
        </cac:TaxCategory>
    </cac:TaxSubtotal>

</cac:TaxTotal>

<!-- Required duplicate total -->
<cac:TaxTotal>
    <cbc:TaxAmount currencyID="SAR">{{.TaxAmount}}</cbc:TaxAmount>
</cac:TaxTotal>
<cac:LegalMonetaryTotal>
<cbc:LineExtensionAmount currencyID="SAR">{{.LineExtensionAmount}}</cbc:LineExtensionAmount>
<cbc:TaxExclusiveAmount currencyID="SAR">{{.TaxExclusiveAmount}}</cbc:TaxExclusiveAmount>
<cbc:TaxInclusiveAmount currencyID="SAR">{{.TaxInclusiveAmount}}</cbc:TaxInclusiveAmount>
<cbc:AllowanceTotalAmount currencyID="SAR">{{.AllowanceTotal}}</cbc:AllowanceTotalAmount>
<cbc:ChargeTotalAmount currencyID="SAR">{{.ChargeTotal}}</cbc:ChargeTotalAmount>
<cbc:PayableAmount currencyID="SAR">{{.TaxInclusiveAmount}}</cbc:PayableAmount>
</cac:LegalMonetaryTotal>
{{- range .Lines}}
<cac:InvoiceLine>
<cbc:ID>{{.ID}}</cbc:ID>
<cbc:InvoicedQuantity unitCode="{{.UnitCode}}">{{printf "%.1f" .Quantity}}</cbc:InvoicedQuantity>
<cbc:LineExtensionAmount currencyID="SAR">{{printf "%.2f" .LineTotal}}</cbc:LineExtensionAmount>
{{range .Discounts}}
<cac:AllowanceCharge>
<cbc:ChargeIndicator>{{.Indicator}}</cbc:ChargeIndicator>
<cbc:AllowanceChargeReason>{{.Reason}}</cbc:AllowanceChargeReason>
<cbc:Amount currencyID="SAR">{{.Amount}}</cbc:Amount>
</cac:AllowanceCharge>
{{end}}
<cac:TaxTotal>
<cbc:TaxAmount currencyID="SAR">{{printf "%.2f" .TaxAmount}}</cbc:TaxAmount>
<cbc:RoundingAmount currencyID="SAR">{{printf "%.2f" .TotalWithVAT}}</cbc:RoundingAmount>
</cac:TaxTotal>
<cac:Item>
<cbc:Name>{{.ItemName}}</cbc:Name>
<cac:ClassifiedTaxCategory>
<cbc:ID>{{.TaxCategoryID}}</cbc:ID>
{{if eq .TaxCategoryID "S"}}<cbc:Percent>{{.VATPercent}}</cbc:Percent>{{end}}
<cac:TaxScheme>
<cbc:ID>VAT</cbc:ID>
</cac:TaxScheme>
</cac:ClassifiedTaxCategory>
</cac:Item>
<cac:Price>
<cbc:PriceAmount currencyID="SAR">{{printf "%.2f" .Price}}</cbc:PriceAmount>
<cbc:BaseQuantity unitCode="{{.UnitCode}}">1</cbc:BaseQuantity>
</cac:Price>
</cac:InvoiceLine>
{{- end}}
</Invoice>`

func encodePIH(raw string) string {
	if raw == "" {
		raw = "0"
	}
	_, err := base64.StdEncoding.DecodeString(raw)
	if err == nil {
		return raw
	}
	return base64.StdEncoding.EncodeToString([]byte(raw))
}

func BuildInvoiceXML(input *InvoiceInput) ([]byte, error) {
	if input == nil {
		return nil, fmt.Errorf("zatca: input is nil")
	}

	baseTime := input.IssueDate
	issueDateStr := baseTime.Format("2006-01-02")
	issueTimeStr := baseTime.Format("15:04:05")

	signingTimeStr := input.SigningTime
	if signingTimeStr == "" {
		signingTimeStr = baseTime.UTC().Format("2006-01-02T15:04:05")
	}

	totals := CalculateTotals(input)

	var lines []tmplLine
	for _, l := range input.Lines {
		t := calcLine(l)

		var itemDiscounts []tmplAC
		for _, d := range l.ItemDiscounts {
			itemDiscounts = append(itemDiscounts, tmplAC{
				Indicator: d.Indicator,
				Reason:    d.Reason,
				Amount:    fmt.Sprintf("%.2f", d.Amount),
			})
		}

		lines = append(lines, tmplLine{
			ID:            l.ID,
			UnitCode:      l.UnitCode,
			Quantity:      l.Quantity,
			LineTotal:     t.LineTotal,
			TaxAmount:     t.TaxAmount,
			TotalWithVAT:  t.TotalWithVAT,
			ItemName:      l.Name,
			VATPercent:    fmt.Sprintf("%.2f", l.VATRate),
			Price:         l.Price,
			TaxCategoryID: taxCategoryID(l.VATRate),
			Discounts:     itemDiscounts,
		})
	}

	var invoiceACs []tmplAC
	for _, ac := range input.InvoiceLevelACs {
		invoiceACs = append(invoiceACs, tmplAC{
			Indicator:     ac.Indicator,
			Reason:        ac.Reason,
			ReasonCode:    ac.ReasonCode,
			Amount:        fmt.Sprintf("%.2f", ac.Amount),
			VATRate:       fmt.Sprintf("%.2f", ac.VATRate),
			TaxCategoryID: ac.TaxCategoryCode,
		})
	}

	signedPropsXML, err := BuildSignedPropertiesXML(input)
	if err != nil {
		return nil, fmt.Errorf("zatca: build signed properties error: %w", err)
	}
	invoiceTypeName := "0100000"

	if input.IsSimplified {
		invoiceTypeName = "0200000"
	}
	hasCustomer := input.Customer != nil &&
		input.Customer.RegistrationName != "" &&
		input.Customer.VATNumber != ""

	var customerStreet string
	var customerBuildingNumber string
	var customerDistrict string
	var customerCity string
	var customerPostalCode string
	var customerVATNumber string
	var customerRegistrationName string

	if input.Customer != nil {
		customerStreet = input.Customer.Street
		customerBuildingNumber = input.Customer.BuildingNumber
		customerDistrict = input.Customer.District
		customerCity = input.Customer.City
		customerPostalCode = input.Customer.PostalCode
		customerVATNumber = input.Customer.VATNumber
		customerRegistrationName = input.Customer.RegistrationName
	}
	taxAmountS := round2(totals.TaxableAmountS * 0.15)
	data := tmplData{
		SigningTime:              signingTimeStr,
		CertificateHash:          input.CertificateHash,
		IssuerName:               input.IssuerName,
		SerialNumber:             input.SerialNumber,
		X509Certificate:          input.X509Certificate,
		InvoiceDigest:            input.InvoiceDigest,
		SignedPropsDigest:        input.SignedPropsDigest,
		SignatureValue:           input.SignatureValue,
		ID:                       input.ID,
		UUID:                     input.UUID,
		IssueDate:                issueDateStr,
		IssueTime:                issueTimeStr,
		ICV:                      fmt.Sprintf("%d", input.ICV),
		PreviousInvoiceHash:      encodePIH(input.PreviousInvoiceHash),
		QRCode:                   input.QRCode,
		CRN:                      input.Supplier.CRN,
		VATNumber:                input.Supplier.VATNumber,
		RegistrationName:         input.Supplier.RegistrationName,
		Street:                   input.Supplier.Street,
		BuildingNumber:           input.Supplier.BuildingNumber,
		PlotID:                   input.Supplier.PlotID,
		District:                 input.Supplier.District,
		City:                     input.Supplier.City,
		PostalCode:               input.Supplier.PostalCode,
		PaymentMeansCode:         input.PaymentMeansCode,
		LineExtensionAmount:      fmt.Sprintf("%.2f", totals.LineExtensionAmount),
		TaxExclusiveAmount:       fmt.Sprintf("%.2f", totals.TaxExclusiveAmount),
		TaxInclusiveAmount:       fmt.Sprintf("%.2f", totals.TaxInclusiveAmount),
		AllowanceTotal:           fmt.Sprintf("%.2f", totals.AllowanceTotal),
		ChargeTotal:              fmt.Sprintf("%.2f", totals.ChargeTotal),
		TaxAmount:                fmt.Sprintf("%.2f", totals.TaxAmount),
		Lines:                    lines,
		InvoiceLevelACs:          invoiceACs,
		BillingReferenceID:       input.BillingReferenceID,
		InvoiceTypeCode:          input.InvoiceTypeCode,
		InstructionNote:          input.InstructionNote,
		SignedPropertiesXML:      string(signedPropsXML),
		HasCustomer:              hasCustomer,
		CustomerStreet:           customerStreet,
		CustomerBuildingNumber:   customerBuildingNumber,
		CustomerDistrict:         customerDistrict,
		CustomerCity:             customerCity,
		CustomerPostalCode:       customerPostalCode,
		CustomerVATNumber:        customerVATNumber,
		CustomerRegistrationName: customerRegistrationName,
		InvoiceTypeName:          invoiceTypeName,
		TaxableAmountS:           fmt.Sprintf("%.2f", totals.TaxableAmountS),
		TaxAmountS:               fmt.Sprintf("%.2f", taxAmountS),
		TaxableAmountO:           fmt.Sprintf("%.2f", totals.TaxableAmountO),
	}

	funcMap := template.FuncMap{
		"printf": fmt.Sprintf,
	}

	tmpl, err := template.New("invoice").Funcs(funcMap).Parse(invoiceTemplate)
	if err != nil {
		return nil, fmt.Errorf("zatca: template parse error: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("zatca: template execute error: %w", err)
	}

	return buf.Bytes(), nil
}

func BuildSignedPropertiesXML(input *InvoiceInput) ([]byte, error) {
	xmlString := fmt.Sprintf(
		`<xades:SignedProperties xmlns:xades="http://uri.etsi.org/01903/v1.3.2#" Id="xadesSignedProperties"><xades:SignedSignatureProperties><xades:SigningTime>%s</xades:SigningTime><xades:SigningCertificate><xades:Cert><xades:CertDigest><ds:DigestMethod xmlns:ds="http://www.w3.org/2000/09/xmldsig#" Algorithm="http://www.w3.org/2001/04/xmlenc#sha256"></ds:DigestMethod><ds:DigestValue xmlns:ds="http://www.w3.org/2000/09/xmldsig#">%s</ds:DigestValue></xades:CertDigest><xades:IssuerSerial><ds:X509IssuerName xmlns:ds="http://www.w3.org/2000/09/xmldsig#">%s</ds:X509IssuerName><ds:X509SerialNumber xmlns:ds="http://www.w3.org/2000/09/xmldsig#">%s</ds:X509SerialNumber></xades:IssuerSerial></xades:Cert></xades:SigningCertificate></xades:SignedSignatureProperties></xades:SignedProperties>`,
		input.SigningTime,
		input.CertificateHash,
		input.IssuerName,
		input.SerialNumber,
	)

	return []byte(xmlString), nil
}
