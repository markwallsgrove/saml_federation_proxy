package models

import "encoding/xml"

type ChidleyRoot314159 struct {
	EntitiesDescriptor *EntitiesDescriptor `xml:"urn:oasis:names:tc:SAML:2.0:metadata EntitiesDescriptor,omitempty" json:"EntitiesDescriptor,omitempty"`
}

type AccountableUsers struct {
	XMLName xml.Name `xml:"http://ukfederation.org.uk/2006/11/label AccountableUsers,omitempty" json:"AccountableUsers,omitempty"`
}

type UKFederationMember struct {
	XMLName xml.Name `xml:"http://ukfederation.org.uk/2006/11/label UKFederationMember,omitempty" json:"UKFederationMember,omitempty"`
}

type CanonicalizationMethod struct {
	Algorithm string   `xml:" Algorithm,attr"  json:",omitempty"`
	XMLName   xml.Name `xml:"http://www.w3.org/2000/09/xmldsig# CanonicalizationMethod,omitempty" json:"CanonicalizationMethod,omitempty"`
}

type DigestMethod struct {
	Algorithm string `xml:" Algorithm,attr"  json:",omitempty"`
}

type DigestValue struct {
	Text    string   `xml:",chardata" json:",omitempty"`
	XMLName xml.Name `xml:"http://www.w3.org/2000/09/xmldsig# DigestValue,omitempty" json:"DigestValue,omitempty"`
}

type Exponent struct {
	Text    string   `xml:",chardata" json:",omitempty"`
	XMLName xml.Name `xml:"http://www.w3.org/2000/09/xmldsig# Exponent,omitempty" json:"Exponent,omitempty"`
}

type KeyInfo struct {
	KeyValue *KeyValue `xml:"http://www.w3.org/2000/09/xmldsig# KeyValue,omitempty" json:"KeyValue,omitempty"`
	X509Data *X509Data `xml:"http://www.w3.org/2000/09/xmldsig# X509Data,omitempty" json:"X509Data,omitempty"`
	XMLName  xml.Name  `xml:"http://www.w3.org/2000/09/xmldsig# KeyInfo,omitempty" json:"KeyInfo,omitempty"`
}

type KeyValue struct {
	RSAKeyValue *RSAKeyValue `xml:"http://www.w3.org/2000/09/xmldsig# RSAKeyValue,omitempty" json:"RSAKeyValue,omitempty"`
	XMLName     xml.Name     `xml:"http://www.w3.org/2000/09/xmldsig# KeyValue,omitempty" json:"KeyValue,omitempty"`
}

type Modulus struct {
	Text    string   `xml:",chardata" json:",omitempty"`
	XMLName xml.Name `xml:"http://www.w3.org/2000/09/xmldsig# Modulus,omitempty" json:"Modulus,omitempty"`
}

type RSAKeyValue struct {
	Exponent *Exponent `xml:"http://www.w3.org/2000/09/xmldsig# Exponent,omitempty" json:"Exponent,omitempty"`
	Modulus  *Modulus  `xml:"http://www.w3.org/2000/09/xmldsig# Modulus,omitempty" json:"Modulus,omitempty"`
	XMLName  xml.Name  `xml:"http://www.w3.org/2000/09/xmldsig# RSAKeyValue,omitempty" json:"RSAKeyValue,omitempty"`
}

type Reference struct {
	URI          string        `xml:" URI,attr"  json:",omitempty"`
	DigestMethod *DigestMethod `xml:"http://www.w3.org/2000/09/xmldsig# DigestMethod,omitempty" json:"DigestMethod,omitempty"`
	DigestValue  *DigestValue  `xml:"http://www.w3.org/2000/09/xmldsig# DigestValue,omitempty" json:"DigestValue,omitempty"`
	Transforms   *Transforms   `xml:"http://www.w3.org/2000/09/xmldsig# Transforms,omitempty" json:"Transforms,omitempty"`
	XMLName      xml.Name      `xml:"http://www.w3.org/2000/09/xmldsig# Reference,omitempty" json:"Reference,omitempty"`
}

type Signature struct {
	KeyInfo        *KeyInfo        `xml:"http://www.w3.org/2000/09/xmldsig# KeyInfo,omitempty" json:"KeyInfo,omitempty"`
	SignatureValue *SignatureValue `xml:"http://www.w3.org/2000/09/xmldsig# SignatureValue,omitempty" json:"SignatureValue,omitempty"`
	SignedInfo     *SignedInfo     `xml:"http://www.w3.org/2000/09/xmldsig# SignedInfo,omitempty" json:"SignedInfo,omitempty"`
	XMLName        xml.Name        `xml:"http://www.w3.org/2000/09/xmldsig# Signature,omitempty" json:"Signature,omitempty"`
}

type SignatureMethod struct {
	Algorithm string   `xml:" Algorithm,attr"  json:",omitempty"`
	XMLName   xml.Name `xml:"http://www.w3.org/2000/09/xmldsig# SignatureMethod,omitempty" json:"SignatureMethod,omitempty"`
}

type SignatureValue struct {
	Text    string   `xml:",chardata" json:",omitempty"`
	XMLName xml.Name `xml:"http://www.w3.org/2000/09/xmldsig# SignatureValue,omitempty" json:"SignatureValue,omitempty"`
}

type SignedInfo struct {
	CanonicalizationMethod *CanonicalizationMethod `xml:"http://www.w3.org/2000/09/xmldsig# CanonicalizationMethod,omitempty" json:"CanonicalizationMethod,omitempty"`
	Reference              *Reference              `xml:"http://www.w3.org/2000/09/xmldsig# Reference,omitempty" json:"Reference,omitempty"`
	SignatureMethod        *SignatureMethod        `xml:"http://www.w3.org/2000/09/xmldsig# SignatureMethod,omitempty" json:"SignatureMethod,omitempty"`
	XMLName                xml.Name                `xml:"http://www.w3.org/2000/09/xmldsig# SignedInfo,omitempty" json:"SignedInfo,omitempty"`
}

type Transform struct {
	Algorithm string   `xml:" Algorithm,attr"  json:",omitempty"`
	XMLName   xml.Name `xml:"http://www.w3.org/2000/09/xmldsig# Transform,omitempty" json:"Transform,omitempty"`
}

type Transforms struct {
	Transform []*Transform `xml:"http://www.w3.org/2000/09/xmldsig# Transform,omitempty" json:"Transform,omitempty"`
	XMLName   xml.Name     `xml:"http://www.w3.org/2000/09/xmldsig# Transforms,omitempty" json:"Transforms,omitempty"`
}

type X509Certificate struct {
	Text    string   `xml:",chardata" json:",omitempty"`
	XMLName xml.Name `xml:"http://www.w3.org/2000/09/xmldsig# X509Certificate,omitempty" json:"X509Certificate,omitempty"`
}

type X509Data struct {
	X509Certificate *X509Certificate `xml:"http://www.w3.org/2000/09/xmldsig# X509Certificate,omitempty" json:"X509Certificate,omitempty"`
	XMLName         xml.Name         `xml:"http://www.w3.org/2000/09/xmldsig# X509Data,omitempty" json:"X509Data,omitempty"`
}

type KeySize struct {
	Xmlns   string   `xml:" xmlns,attr"  json:",omitempty"`
	Text    string   `xml:",chardata" json:",omitempty"`
	XMLName xml.Name `xml:"http://www.w3.org/2001/04/xmlenc# KeySize,omitempty" json:"KeySize,omitempty"`
}

type Scope struct {
	Regexp  string   `xml:" regexp,attr"  json:",omitempty"`
	Text    string   `xml:",chardata" json:",omitempty"`
	XMLName xml.Name `xml:"urn:mace:shibboleth:metadata:1.0 Scope,omitempty" json:"Scope,omitempty"`
}

type Attribute struct {
	FriendlyName   string            `xml:" FriendlyName,attr"  json:",omitempty"`
	Name           string            `xml:" Name,attr"  json:",omitempty"`
	NameFormat     string            `xml:" NameFormat,attr"  json:",omitempty"`
	AttributeValue []*AttributeValue `xml:"urn:oasis:names:tc:SAML:2.0:assertion AttributeValue,omitempty" json:"AttributeValue,omitempty"`
	XMLName        xml.Name          `xml:"urn:oasis:names:tc:SAML:2.0:assertion Attribute,omitempty" json:"Attribute,omitempty"`
}

type AttributeValue struct {
	Text    string   `xml:",chardata" json:",omitempty"`
	XMLName xml.Name `xml:"urn:oasis:names:tc:SAML:2.0:assertion AttributeValue,omitempty" json:"AttributeValue,omitempty"`
}

type ArtifactResolutionService struct {
	Binding   string   `xml:" Binding,attr"  json:",omitempty"`
	Index     string   `xml:" index,attr"  json:",omitempty"`
	IsDefault string   `xml:" isDefault,attr"  json:",omitempty"`
	Location  string   `xml:" Location,attr"  json:",omitempty"`
	XMLName   xml.Name `xml:"urn:oasis:names:tc:SAML:2.0:metadata ArtifactResolutionService,omitempty" json:"ArtifactResolutionService,omitempty"`
}

type AssertionConsumerService struct {
	Binding   string   `xml:" Binding,attr"  json:",omitempty"`
	Index     string   `xml:" index,attr"  json:",omitempty"`
	IsDefault string   `xml:" isDefault,attr"  json:",omitempty"`
	Location  string   `xml:" Location,attr"  json:",omitempty"`
	XMLName   xml.Name `xml:"urn:oasis:names:tc:SAML:2.0:metadata AssertionConsumerService,omitempty" json:"AssertionConsumerService,omitempty"`
}

type AssertionIDRequestService struct {
	Binding  string   `xml:" Binding,attr"  json:",omitempty"`
	Location string   `xml:" Location,attr"  json:",omitempty"`
	XMLName  xml.Name `xml:"urn:oasis:names:tc:SAML:2.0:metadata AssertionIDRequestService,omitempty" json:"AssertionIDRequestService,omitempty"`
}

type AttributeAuthorityDescriptor struct {
	ProtocolSupportEnumeration string              `xml:" protocolSupportEnumeration,attr"  json:",omitempty"`
	Attribute                  []*Attribute        `xml:"urn:oasis:names:tc:SAML:2.0:assertion Attribute,omitempty" json:"Attribute,omitempty"`
	AttributeService           []*AttributeService `xml:"urn:oasis:names:tc:SAML:2.0:metadata AttributeService,omitempty" json:"AttributeService,omitempty"`
	Extensions                 *Extensions         `xml:"urn:oasis:names:tc:SAML:2.0:metadata Extensions,omitempty" json:"Extensions,omitempty"`
	KeyDescriptor              []*KeyDescriptor    `xml:"urn:oasis:names:tc:SAML:2.0:metadata KeyDescriptor,omitempty" json:"KeyDescriptor,omitempty"`
	NameIDFormat               []*NameIDFormat     `xml:"urn:oasis:names:tc:SAML:2.0:metadata NameIDFormat,omitempty" json:"NameIDFormat,omitempty"`
	XMLName                    xml.Name            `xml:"urn:oasis:names:tc:SAML:2.0:metadata AttributeAuthorityDescriptor,omitempty" json:"AttributeAuthorityDescriptor,omitempty"`
}

type AttributeConsumingService struct {
	Index              string                `xml:" index,attr"  json:",omitempty"`
	IsDefault          string                `xml:" isDefault,attr"  json:",omitempty"`
	RequestedAttribute []*RequestedAttribute `xml:"urn:oasis:names:tc:SAML:2.0:metadata RequestedAttribute,omitempty" json:"RequestedAttribute,omitempty"`
	ServiceDescription []*ServiceDescription `xml:"urn:oasis:names:tc:SAML:2.0:metadata ServiceDescription,omitempty" json:"ServiceDescription,omitempty"`
	ServiceName        []*ServiceName        `xml:"urn:oasis:names:tc:SAML:2.0:metadata ServiceName,omitempty" json:"ServiceName,omitempty"`
	XMLName            xml.Name              `xml:"urn:oasis:names:tc:SAML:2.0:metadata AttributeConsumingService,omitempty" json:"AttributeConsumingService,omitempty"`
}

type AttributeService struct {
	Binding  string   `xml:" Binding,attr"  json:",omitempty"`
	Location string   `xml:" Location,attr"  json:",omitempty"`
	XMLName  xml.Name `xml:"urn:oasis:names:tc:SAML:2.0:metadata AttributeService,omitempty" json:"AttributeService,omitempty"`
}

type Company struct {
	Text    string   `xml:",chardata" json:",omitempty"`
	XMLName xml.Name `xml:"urn:oasis:names:tc:SAML:2.0:metadata Company,omitempty" json:"Company,omitempty"`
}

type ContactPerson struct {
	ContactType     string             `xml:"contactType,attr"  json:",omitempty"`
	XmlnsIcmd       string             `xml:"xmlns icmd,attr"  json:",omitempty"`
	XmlnsRefeds     string             `xml:"xmlns refeds,attr"  json:",omitempty"`
	Company         *Company           `xml:"urn:oasis:names:tc:SAML:2.0:metadata Company,omitempty" json:"Company,omitempty"`
	EmailAddress    []*EmailAddress    `xml:"urn:oasis:names:tc:SAML:2.0:metadata EmailAddress,omitempty" json:"EmailAddress,omitempty"`
	GivenName       *GivenName         `xml:"urn:oasis:names:tc:SAML:2.0:metadata GivenName,omitempty" json:"GivenName,omitempty"`
	SurName         *SurName           `xml:"urn:oasis:names:tc:SAML:2.0:metadata SurName,omitempty" json:"SurName,omitempty"`
	TelephoneNumber []*TelephoneNumber `xml:"urn:oasis:names:tc:SAML:2.0:metadata TelephoneNumber,omitempty" json:"TelephoneNumber,omitempty"`
	XMLName         xml.Name           `xml:"urn:oasis:names:tc:SAML:2.0:metadata ContactPerson,omitempty" json:"ContactPerson,omitempty"`
}

type EmailAddress struct {
	Text    string   `xml:",chardata" json:",omitempty"`
	XMLName xml.Name `xml:"urn:oasis:names:tc:SAML:2.0:metadata EmailAddress,omitempty" json:"EmailAddress,omitempty"`
}

type EncryptionMethod struct {
	Algorithm string   `xml:" Algorithm,attr"  json:",omitempty"`
	KeySize   *KeySize `xml:"http://www.w3.org/2001/04/xmlenc# KeySize,omitempty" json:"KeySize,omitempty"`
	XMLName   xml.Name `xml:"urn:oasis:names:tc:SAML:2.0:metadata EncryptionMethod,omitempty" json:"EncryptionMethod,omitempty"`
}

type EntitiesDescriptor struct {
	XmlnsDs          string              `xml:"xmlns ds,attr"  json:",omitempty"`
	ID               string              `xml:" ID,attr"  json:",omitempty"`
	XmlnsIdpdisc     string              `xml:"xmlns idpdisc,attr"  json:",omitempty"`
	XmlnsInit        string              `xml:"xmlns init,attr"  json:",omitempty"`
	XmlnsMdattr      string              `xml:"xmlns mdattr,attr"  json:",omitempty"`
	XmlnsMdrpi       string              `xml:"xmlns mdrpi,attr"  json:",omitempty"`
	XmlnsMdui        string              `xml:"xmlns mdui,attr"  json:",omitempty"`
	Name             string              `xml:" Name,attr"  json:",omitempty"`
	XmlnsRemd        string              `xml:"xmlns remd,attr"  json:",omitempty"`
	XmlnsSaml        string              `xml:"xmlns saml,attr"  json:",omitempty"`
	XmlnsShibmd      string              `xml:"xmlns shibmd,attr"  json:",omitempty"`
	XmlnsUkfedlabel  string              `xml:"xmlns ukfedlabel,attr"  json:",omitempty"`
	ValidUntil       string              `xml:" validUntil,attr"  json:",omitempty"`
	Xmlns            string              `xml:" xmlns,attr"  json:",omitempty"`
	XmlnsXsi         string              `xml:"xmlns xsi,attr"  json:",omitempty"`
	EntityDescriptor []*EntityDescriptor `xml:"urn:oasis:names:tc:SAML:2.0:metadata EntityDescriptor,omitempty" json:"EntityDescriptor,omitempty"`
	Extensions       *Extensions         `xml:"urn:oasis:names:tc:SAML:2.0:metadata Extensions,omitempty" json:"Extensions,omitempty"`
	Signature        *Signature          `xml:"http://www.w3.org/2000/09/xmldsig# Signature,omitempty" json:"Signature,omitempty"`
	XMLName          xml.Name            `xml:"urn:oasis:names:tc:SAML:2.0:metadata EntitiesDescriptor,omitempty" json:"EntitiesDescriptor,omitempty"`
}

type EntityDescriptor struct {
	FederationID                 string                        `json:",omitempty"`
	Checksum                     []byte                        `json:",omitempty"`
	EntityID                     string                        `xml:" entityID,attr"  json:",omitempty"`
	ID                           string                        `xml:" ID,attr"  json:",omitempty"`
	AttributeAuthorityDescriptor *AttributeAuthorityDescriptor `xml:"urn:oasis:names:tc:SAML:2.0:metadata AttributeAuthorityDescriptor,omitempty" json:"AttributeAuthorityDescriptor,omitempty"`
	ContactPerson                []*ContactPerson              `xml:"urn:oasis:names:tc:SAML:2.0:metadata ContactPerson,omitempty" json:"ContactPerson,omitempty"`
	Extensions                   *Extensions                   `xml:"urn:oasis:names:tc:SAML:2.0:metadata Extensions,omitempty" json:"Extensions,omitempty"`
	IDPSSODescriptor             *IDPSSODescriptor             `xml:"urn:oasis:names:tc:SAML:2.0:metadata IDPSSODescriptor,omitempty" json:"IDPSSODescriptor,omitempty"`
	Organization                 *Organization                 `xml:"urn:oasis:names:tc:SAML:2.0:metadata Organization,omitempty" json:"Organization,omitempty"`
	SPSSODescriptor              *SPSSODescriptor              `xml:"urn:oasis:names:tc:SAML:2.0:metadata SPSSODescriptor,omitempty" json:"SPSSODescriptor,omitempty"`
	XMLName                      xml.Name                      `xml:"urn:oasis:names:tc:SAML:2.0:metadata EntityDescriptor,omitempty" json:"EntityDescriptor,omitempty"`
}

type Extensions struct {
	AccountableUsers   *AccountableUsers    `xml:"http://ukfederation.org.uk/2006/11/label AccountableUsers,omitempty" json:"AccountableUsers,omitempty"`
	DigestMethod       []*DigestMethod      `xml:"urn:oasis:names:tc:SAML:metadata:algsupport DigestMethod,omitempty" json:"DigestMethod,omitempty"`
	DiscoHints         *DiscoHints          `xml:"urn:oasis:names:tc:SAML:metadata:ui DiscoHints,omitempty" json:"DiscoHints,omitempty"`
	DiscoveryResponse  []*DiscoveryResponse `xml:"urn:oasis:names:tc:SAML:profiles:SSO:idp-discovery-protocol DiscoveryResponse,omitempty" json:"DiscoveryResponse,omitempty"`
	EntityAttributes   *EntityAttributes    `xml:"urn:oasis:names:tc:SAML:metadata:attribute EntityAttributes,omitempty" json:"EntityAttributes,omitempty"`
	PublicationInfo    *PublicationInfo     `xml:"urn:oasis:names:tc:SAML:metadata:rpi PublicationInfo,omitempty" json:"PublicationInfo,omitempty"`
	RegistrationInfo   *RegistrationInfo    `xml:"urn:oasis:names:tc:SAML:metadata:rpi RegistrationInfo,omitempty" json:"RegistrationInfo,omitempty"`
	RequestInitiator   []*RequestInitiator  `xml:"urn:oasis:names:tc:SAML:profiles:SSO:request-init RequestInitiator,omitempty" json:"RequestInitiator,omitempty"`
	Scope              []*Scope             `xml:"urn:mace:shibboleth:metadata:1.0 Scope,omitempty" json:"Scope,omitempty"`
	SigningMethod      []*SigningMethod     `xml:"urn:oasis:names:tc:SAML:metadata:algsupport SigningMethod,omitempty" json:"SigningMethod,omitempty"`
	UIInfo             *UIInfo              `xml:"urn:oasis:names:tc:SAML:metadata:ui UIInfo,omitempty" json:"UIInfo,omitempty"`
	UKFederationMember *UKFederationMember  `xml:"http://ukfederation.org.uk/2006/11/label UKFederationMember,omitempty" json:"UKFederationMember,omitempty"`
	XMLName            xml.Name             `xml:"urn:oasis:names:tc:SAML:2.0:metadata Extensions,omitempty" json:"Extensions,omitempty"`
}

type GivenName struct {
	Text    string   `xml:",chardata" json:",omitempty"`
	XMLName xml.Name `xml:"urn:oasis:names:tc:SAML:2.0:metadata GivenName,omitempty" json:"GivenName,omitempty"`
}

type IDPSSODescriptor struct {
	ErrorURL                   string                       `xml:" errorURL,attr"  json:",omitempty"`
	ID                         string                       `xml:" ID,attr"  json:",omitempty"`
	ProtocolSupportEnumeration string                       `xml:" protocolSupportEnumeration,attr"  json:",omitempty"`
	WantAuthnRequestsSigned    string                       `xml:" WantAuthnRequestsSigned,attr"  json:",omitempty"`
	ArtifactResolutionService  []*ArtifactResolutionService `xml:"urn:oasis:names:tc:SAML:2.0:metadata ArtifactResolutionService,omitempty" json:"ArtifactResolutionService,omitempty"`
	AssertionIDRequestService  []*AssertionIDRequestService `xml:"urn:oasis:names:tc:SAML:2.0:metadata AssertionIDRequestService,omitempty" json:"AssertionIDRequestService,omitempty"`
	Attribute                  []*Attribute                 `xml:"urn:oasis:names:tc:SAML:2.0:assertion Attribute,omitempty" json:"Attribute,omitempty"`
	Extensions                 *Extensions                  `xml:"urn:oasis:names:tc:SAML:2.0:metadata Extensions,omitempty" json:"Extensions,omitempty"`
	KeyDescriptor              []*KeyDescriptor             `xml:"urn:oasis:names:tc:SAML:2.0:metadata KeyDescriptor,omitempty" json:"KeyDescriptor,omitempty"`
	ManageNameIDService        []*ManageNameIDService       `xml:"urn:oasis:names:tc:SAML:2.0:metadata ManageNameIDService,omitempty" json:"ManageNameIDService,omitempty"`
	NameIDFormat               []*NameIDFormat              `xml:"urn:oasis:names:tc:SAML:2.0:metadata NameIDFormat,omitempty" json:"NameIDFormat,omitempty"`
	NameIDMappingService       *NameIDMappingService        `xml:"urn:oasis:names:tc:SAML:2.0:metadata NameIDMappingService,omitempty" json:"NameIDMappingService,omitempty"`
	SingleLogoutService        []*SingleLogoutService       `xml:"urn:oasis:names:tc:SAML:2.0:metadata SingleLogoutService,omitempty" json:"SingleLogoutService,omitempty"`
	SingleSignOnService        []*SingleSignOnService       `xml:"urn:oasis:names:tc:SAML:2.0:metadata SingleSignOnService,omitempty" json:"SingleSignOnService,omitempty"`
	XMLName                    xml.Name                     `xml:"urn:oasis:names:tc:SAML:2.0:metadata IDPSSODescriptor,omitempty" json:"IDPSSODescriptor,omitempty"`
}

type KeyDescriptor struct {
	Use              string              `xml:" use,attr"  json:",omitempty"`
	EncryptionMethod []*EncryptionMethod `xml:"urn:oasis:names:tc:SAML:2.0:metadata EncryptionMethod,omitempty" json:"EncryptionMethod,omitempty"`
	KeyInfo          *KeyInfo            `xml:"http://www.w3.org/2000/09/xmldsig# KeyInfo,omitempty" json:"KeyInfo,omitempty"`
	XMLName          xml.Name            `xml:"urn:oasis:names:tc:SAML:2.0:metadata KeyDescriptor,omitempty" json:"KeyDescriptor,omitempty"`
}

type ManageNameIDService struct {
	Binding          string   `xml:" Binding,attr"  json:",omitempty"`
	Location         string   `xml:" Location,attr"  json:",omitempty"`
	ResponseLocation string   `xml:" ResponseLocation,attr"  json:",omitempty"`
	XMLName          xml.Name `xml:"urn:oasis:names:tc:SAML:2.0:metadata ManageNameIDService,omitempty" json:"ManageNameIDService,omitempty"`
}

type NameIDFormat struct {
	Text    string   `xml:",chardata" json:",omitempty"`
	XMLName xml.Name `xml:"urn:oasis:names:tc:SAML:2.0:metadata NameIDFormat,omitempty" json:"NameIDFormat,omitempty"`
}

type NameIDMappingService struct {
	Binding  string   `xml:" Binding,attr"  json:",omitempty"`
	Location string   `xml:" Location,attr"  json:",omitempty"`
	XMLName  xml.Name `xml:"urn:oasis:names:tc:SAML:2.0:metadata NameIDMappingService,omitempty" json:"NameIDMappingService,omitempty"`
}

type Organization struct {
	OrganizationDisplayName []*OrganizationDisplayName `xml:"urn:oasis:names:tc:SAML:2.0:metadata OrganizationDisplayName,omitempty" json:"OrganizationDisplayName,omitempty"`
	OrganizationName        []*OrganizationName        `xml:"urn:oasis:names:tc:SAML:2.0:metadata OrganizationName,omitempty" json:"OrganizationName,omitempty"`
	OrganizationURL         []*OrganizationURL         `xml:"urn:oasis:names:tc:SAML:2.0:metadata OrganizationURL,omitempty" json:"OrganizationURL,omitempty"`
	XMLName                 xml.Name                   `xml:"urn:oasis:names:tc:SAML:2.0:metadata Organization,omitempty" json:"Organization,omitempty"`
}

type OrganizationDisplayName struct {
	NamespaceLang string   `xml:"http://www.w3.org/XML/1998/namespace lang,attr"  json:",omitempty"`
	Text          string   `xml:",chardata" json:",omitempty"`
	XMLName       xml.Name `xml:"urn:oasis:names:tc:SAML:2.0:metadata OrganizationDisplayName,omitempty" json:"OrganizationDisplayName,omitempty"`
}

type OrganizationName struct {
	NamespaceLang string   `xml:"http://www.w3.org/XML/1998/namespace lang,attr"  json:",omitempty"`
	Text          string   `xml:",chardata" json:",omitempty"`
	XMLName       xml.Name `xml:"urn:oasis:names:tc:SAML:2.0:metadata OrganizationName,omitempty" json:"OrganizationName,omitempty"`
}

type OrganizationURL struct {
	NamespaceLang string   `xml:"http://www.w3.org/XML/1998/namespace lang,attr"  json:",omitempty"`
	Text          string   `xml:",chardata" json:",omitempty"`
	XMLName       xml.Name `xml:"urn:oasis:names:tc:SAML:2.0:metadata OrganizationURL,omitempty" json:"OrganizationURL,omitempty"`
}

type RequestedAttribute struct {
	FriendlyName   string            `xml:" FriendlyName,attr"  json:",omitempty"`
	IsRequired     string            `xml:" isRequired,attr"  json:",omitempty"`
	Name           string            `xml:" Name,attr"  json:",omitempty"`
	NameFormat     string            `xml:" NameFormat,attr"  json:",omitempty"`
	AttributeValue []*AttributeValue `xml:"urn:oasis:names:tc:SAML:2.0:assertion AttributeValue,omitempty" json:"AttributeValue,omitempty"`
	XMLName        xml.Name          `xml:"urn:oasis:names:tc:SAML:2.0:metadata RequestedAttribute,omitempty" json:"RequestedAttribute,omitempty"`
}

type SPSSODescriptor struct {
	AuthnRequestsSigned        string                       `xml:" AuthnRequestsSigned,attr"  json:",omitempty"`
	ErrorURL                   string                       `xml:" errorURL,attr"  json:",omitempty"`
	ProtocolSupportEnumeration string                       `xml:" protocolSupportEnumeration,attr"  json:",omitempty"`
	WantAssertionsSigned       string                       `xml:" WantAssertionsSigned,attr"  json:",omitempty"`
	ArtifactResolutionService  []*ArtifactResolutionService `xml:"urn:oasis:names:tc:SAML:2.0:metadata ArtifactResolutionService,omitempty" json:"ArtifactResolutionService,omitempty"`
	AssertionConsumerService   []*AssertionConsumerService  `xml:"urn:oasis:names:tc:SAML:2.0:metadata AssertionConsumerService,omitempty" json:"AssertionConsumerService,omitempty"`
	AttributeConsumingService  *AttributeConsumingService   `xml:"urn:oasis:names:tc:SAML:2.0:metadata AttributeConsumingService,omitempty" json:"AttributeConsumingService,omitempty"`
	Extensions                 *Extensions                  `xml:"urn:oasis:names:tc:SAML:2.0:metadata Extensions,omitempty" json:"Extensions,omitempty"`
	KeyDescriptor              []*KeyDescriptor             `xml:"urn:oasis:names:tc:SAML:2.0:metadata KeyDescriptor,omitempty" json:"KeyDescriptor,omitempty"`
	ManageNameIDService        []*ManageNameIDService       `xml:"urn:oasis:names:tc:SAML:2.0:metadata ManageNameIDService,omitempty" json:"ManageNameIDService,omitempty"`
	NameIDFormat               []*NameIDFormat              `xml:"urn:oasis:names:tc:SAML:2.0:metadata NameIDFormat,omitempty" json:"NameIDFormat,omitempty"`
	SingleLogoutService        []*SingleLogoutService       `xml:"urn:oasis:names:tc:SAML:2.0:metadata SingleLogoutService,omitempty" json:"SingleLogoutService,omitempty"`
	XMLName                    xml.Name                     `xml:"urn:oasis:names:tc:SAML:2.0:metadata SPSSODescriptor,omitempty" json:"SPSSODescriptor,omitempty"`
}

type ServiceDescription struct {
	NamespaceLang string   `xml:"http://www.w3.org/XML/1998/namespace lang,attr"  json:",omitempty"`
	Text          string   `xml:",chardata" json:",omitempty"`
	XMLName       xml.Name `xml:"urn:oasis:names:tc:SAML:2.0:metadata ServiceDescription,omitempty" json:"ServiceDescription,omitempty"`
}

type ServiceName struct {
	NamespaceLang string   `xml:"http://www.w3.org/XML/1998/namespace lang,attr"  json:",omitempty"`
	Text          string   `xml:",chardata" json:",omitempty"`
	XMLName       xml.Name `xml:"urn:oasis:names:tc:SAML:2.0:metadata ServiceName,omitempty" json:"ServiceName,omitempty"`
}

type SingleLogoutService struct {
	Binding          string   `xml:" Binding,attr"  json:",omitempty"`
	Location         string   `xml:" Location,attr"  json:",omitempty"`
	ResponseLocation string   `xml:" ResponseLocation,attr"  json:",omitempty"`
	XMLName          xml.Name `xml:"urn:oasis:names:tc:SAML:2.0:metadata SingleLogoutService,omitempty" json:"SingleLogoutService,omitempty"`
}

type SingleSignOnService struct {
	Binding  string   `xml:" Binding,attr"  json:",omitempty"`
	Location string   `xml:" Location,attr"  json:",omitempty"`
	XMLName  xml.Name `xml:"urn:oasis:names:tc:SAML:2.0:metadata SingleSignOnService,omitempty" json:"SingleSignOnService,omitempty"`
}

type SurName struct {
	Text    string   `xml:",chardata" json:",omitempty"`
	XMLName xml.Name `xml:"urn:oasis:names:tc:SAML:2.0:metadata SurName,omitempty" json:"SurName,omitempty"`
}

type TelephoneNumber struct {
	Text    string   `xml:",chardata" json:",omitempty"`
	XMLName xml.Name `xml:"urn:oasis:names:tc:SAML:2.0:metadata TelephoneNumber,omitempty" json:"TelephoneNumber,omitempty"`
}

type SigningMethod struct {
	Algorithm string   `xml:" Algorithm,attr"  json:",omitempty"`
	Xmlns     string   `xml:" xmlns,attr"  json:",omitempty"`
	XMLName   xml.Name `xml:"urn:oasis:names:tc:SAML:metadata:algsupport SigningMethod,omitempty" json:"SigningMethod,omitempty"`
}

type EntityAttributes struct {
	Attribute []*Attribute `xml:"urn:oasis:names:tc:SAML:2.0:assertion Attribute,omitempty" json:"Attribute,omitempty"`
	XMLName   xml.Name     `xml:"urn:oasis:names:tc:SAML:metadata:attribute EntityAttributes,omitempty" json:"EntityAttributes,omitempty"`
}

type PublicationInfo struct {
	CreationInstant string   `xml:" creationInstant,attr"  json:",omitempty"`
	Publisher       string   `xml:" publisher,attr"  json:",omitempty"`
	XMLName         xml.Name `xml:"urn:oasis:names:tc:SAML:metadata:rpi PublicationInfo,omitempty" json:"PublicationInfo,omitempty"`
}

type RegistrationInfo struct {
	RegistrationAuthority string                `xml:" registrationAuthority,attr"  json:",omitempty"`
	RegistrationInstant   string                `xml:" registrationInstant,attr"  json:",omitempty"`
	RegistrationPolicy    []*RegistrationPolicy `xml:"urn:oasis:names:tc:SAML:metadata:rpi RegistrationPolicy,omitempty" json:"RegistrationPolicy,omitempty"`
	XMLName               xml.Name              `xml:"urn:oasis:names:tc:SAML:metadata:rpi RegistrationInfo,omitempty" json:"RegistrationInfo,omitempty"`
}

type RegistrationPolicy struct {
	NamespaceLang string   `xml:"http://www.w3.org/XML/1998/namespace lang,attr"  json:",omitempty"`
	Text          string   `xml:",chardata" json:",omitempty"`
	XMLName       xml.Name `xml:"urn:oasis:names:tc:SAML:metadata:rpi RegistrationPolicy,omitempty" json:"RegistrationPolicy,omitempty"`
}

type Description struct {
	NamespaceLang string   `xml:"http://www.w3.org/XML/1998/namespace lang,attr"  json:",omitempty"`
	Text          string   `xml:",chardata" json:",omitempty"`
	XMLName       xml.Name `xml:"urn:oasis:names:tc:SAML:metadata:ui Description,omitempty" json:"Description,omitempty"`
}

type DiscoHints struct {
	DomainHint      []*DomainHint      `xml:"urn:oasis:names:tc:SAML:metadata:ui DomainHint,omitempty" json:"DomainHint,omitempty"`
	GeolocationHint []*GeolocationHint `xml:"urn:oasis:names:tc:SAML:metadata:ui GeolocationHint,omitempty" json:"GeolocationHint,omitempty"`
	IPHint          []*IPHint          `xml:"urn:oasis:names:tc:SAML:metadata:ui IPHint,omitempty" json:"IPHint,omitempty"`
	XMLName         xml.Name           `xml:"urn:oasis:names:tc:SAML:metadata:ui DiscoHints,omitempty" json:"DiscoHints,omitempty"`
}

type DisplayName struct {
	NamespaceLang string   `xml:"http://www.w3.org/XML/1998/namespace lang,attr"  json:",omitempty"`
	Text          string   `xml:",chardata" json:",omitempty"`
	XMLName       xml.Name `xml:"urn:oasis:names:tc:SAML:metadata:ui DisplayName,omitempty" json:"DisplayName,omitempty"`
}

type DomainHint struct {
	Text    string   `xml:",chardata" json:",omitempty"`
	XMLName xml.Name `xml:"urn:oasis:names:tc:SAML:metadata:ui DomainHint,omitempty" json:"DomainHint,omitempty"`
}

type GeolocationHint struct {
	Text    string   `xml:",chardata" json:",omitempty"`
	XMLName xml.Name `xml:"urn:oasis:names:tc:SAML:metadata:ui GeolocationHint,omitempty" json:"GeolocationHint,omitempty"`
}

type IPHint struct {
	Text    string   `xml:",chardata" json:",omitempty"`
	XMLName xml.Name `xml:"urn:oasis:names:tc:SAML:metadata:ui IPHint,omitempty" json:"IPHint,omitempty"`
}

type InformationURL struct {
	NamespaceLang string   `xml:"http://www.w3.org/XML/1998/namespace lang,attr"  json:",omitempty"`
	Text          string   `xml:",chardata" json:",omitempty"`
	XMLName       xml.Name `xml:"urn:oasis:names:tc:SAML:metadata:ui InformationURL,omitempty" json:"InformationURL,omitempty"`
}

type Keywords struct {
	NamespaceLang string   `xml:"http://www.w3.org/XML/1998/namespace lang,attr"  json:",omitempty"`
	Text          string   `xml:",chardata" json:",omitempty"`
	XMLName       xml.Name `xml:"urn:oasis:names:tc:SAML:metadata:ui Keywords,omitempty" json:"Keywords,omitempty"`
}

type Logo struct {
	Height        string   `xml:" height,attr"  json:",omitempty"`
	NamespaceLang string   `xml:"http://www.w3.org/XML/1998/namespace lang,attr"  json:",omitempty"`
	Width         string   `xml:" width,attr"  json:",omitempty"`
	Text          string   `xml:",chardata" json:",omitempty"`
	XMLName       xml.Name `xml:"urn:oasis:names:tc:SAML:metadata:ui Logo,omitempty" json:"Logo,omitempty"`
}

type PrivacyStatementURL struct {
	NamespaceLang string   `xml:"http://www.w3.org/XML/1998/namespace lang,attr"  json:",omitempty"`
	Text          string   `xml:",chardata" json:",omitempty"`
	XMLName       xml.Name `xml:"urn:oasis:names:tc:SAML:metadata:ui PrivacyStatementURL,omitempty" json:"PrivacyStatementURL,omitempty"`
}

type UIInfo struct {
	Description         []*Description         `xml:"urn:oasis:names:tc:SAML:metadata:ui Description,omitempty" json:"Description,omitempty"`
	DisplayName         []*DisplayName         `xml:"urn:oasis:names:tc:SAML:metadata:ui DisplayName,omitempty" json:"DisplayName,omitempty"`
	InformationURL      []*InformationURL      `xml:"urn:oasis:names:tc:SAML:metadata:ui InformationURL,omitempty" json:"InformationURL,omitempty"`
	Keywords            []*Keywords            `xml:"urn:oasis:names:tc:SAML:metadata:ui Keywords,omitempty" json:"Keywords,omitempty"`
	Logo                []*Logo                `xml:"urn:oasis:names:tc:SAML:metadata:ui Logo,omitempty" json:"Logo,omitempty"`
	PrivacyStatementURL []*PrivacyStatementURL `xml:"urn:oasis:names:tc:SAML:metadata:ui PrivacyStatementURL,omitempty" json:"PrivacyStatementURL,omitempty"`
	XMLName             xml.Name               `xml:"urn:oasis:names:tc:SAML:metadata:ui UIInfo,omitempty" json:"UIInfo,omitempty"`
}

type DiscoveryResponse struct {
	Binding   string   `xml:" Binding,attr"  json:",omitempty"`
	Index     string   `xml:" index,attr"  json:",omitempty"`
	IsDefault string   `xml:" isDefault,attr"  json:",omitempty"`
	Location  string   `xml:" Location,attr"  json:",omitempty"`
	XMLName   xml.Name `xml:"urn:oasis:names:tc:SAML:profiles:SSO:idp-discovery-protocol DiscoveryResponse,omitempty" json:"DiscoveryResponse,omitempty"`
}

type RequestInitiator struct {
	Binding  string   `xml:" Binding,attr"  json:",omitempty"`
	Location string   `xml:" Location,attr"  json:",omitempty"`
	XMLName  xml.Name `xml:"urn:oasis:names:tc:SAML:profiles:SSO:request-init RequestInitiator,omitempty" json:"RequestInitiator,omitempty"`
}
