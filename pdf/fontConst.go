package pdf

var (
	// CMapIdentityH : The horizontal identity mapping for 2-byte CIDs;
	// may be used with CIDFonts using any Registry, Ordering, and Supplement values.
	// It maps 2-byte character codes ranging from 0 to 65,535 to the same 2-byte CID value,
	// interpreted high-order byte first
	CMapIdentityH = newPredefinedCMap("Identity-H")
	// CMapIdentityV : Vertical version of Identity−H. The mapping is the same as for Identity−H.
	CMapIdentityV = newPredefinedCMap("Identity-V")
	// CIDSystemInfoAdobeIdentity0 : Special CIDSystemInfo for not CID-based fonts.
	CIDSystemInfoAdobeIdentity0 = newCIDSystemInfo("Adobe", "Identity", 0)
	// CIDSystemInfoAdobeCNS1 : Traditional Chinese
	CIDSystemInfoAdobeCNS1 = newCIDSystemInfo("Adobe", "CNS1", 1)
	// CIDSystemInfoAdobeCNS2 : Traditional Chinese
	CIDSystemInfoAdobeCNS2 = newCIDSystemInfo("Adobe", "CNS1", 2)
	// CIDSystemInfoAdobeCNS3 : Traditional Chinese
	CIDSystemInfoAdobeCNS3 = newCIDSystemInfo("Adobe", "CNS1", 3)
	// CIDSystemInfoAdobeCNS4 : Traditional Chinese
	CIDSystemInfoAdobeCNS4 = newCIDSystemInfo("Adobe", "CNS1", 4)
	// CIDSystemInfoAdobeCNS5 : Traditional Chinese
	CIDSystemInfoAdobeCNS5 = newCIDSystemInfo("Adobe", "CNS1", 5)
	// CIDSystemInfoAdobeCNS6 : Traditional Chinese
	CIDSystemInfoAdobeCNS6 = newCIDSystemInfo("Adobe", "CNS1", 6)
	// CIDSystemInfoAdobeCNS7 : Traditional Chinese
	CIDSystemInfoAdobeCNS7 = newCIDSystemInfo("Adobe", "CNS1", 7)
	// CIDSystemInfoAdobeGB1 : Simplified Chinese
	CIDSystemInfoAdobeGB1 = newCIDSystemInfo("Adobe", "GB1", 1)
	// CIDSystemInfoAdobeGB2 : Simplified Chinese
	CIDSystemInfoAdobeGB2 = newCIDSystemInfo("Adobe", "GB1", 2)
	// CIDSystemInfoAdobeGB3 : Simplified Chinese
	CIDSystemInfoAdobeGB3 = newCIDSystemInfo("Adobe", "GB1", 3)
	// CIDSystemInfoAdobeGB4 : Simplified Chinese
	CIDSystemInfoAdobeGB4 = newCIDSystemInfo("Adobe", "GB1", 4)
	// CIDSystemInfoAdobeGB5 : Simplified Chinese
	CIDSystemInfoAdobeGB5 = newCIDSystemInfo("Adobe", "GB1", 5)
	// CIDSystemInfoAdobeJapan1 : Japanese
	CIDSystemInfoAdobeJapan1 = newCIDSystemInfo("Adobe", "Japan1", 1)
	// CIDSystemInfoAdobeJapan2 : Japanese
	CIDSystemInfoAdobeJapan2 = newCIDSystemInfo("Adobe", "Japan1", 2)
	// CIDSystemInfoAdobeJapan3 : Japanese
	CIDSystemInfoAdobeJapan3 = newCIDSystemInfo("Adobe", "Japan1", 3)
	// CIDSystemInfoAdobeJapan4 : Japanese
	CIDSystemInfoAdobeJapan4 = newCIDSystemInfo("Adobe", "Japan1", 4)
	// CIDSystemInfoAdobeJapan5 : Japanese
	CIDSystemInfoAdobeJapan5 = newCIDSystemInfo("Adobe", "Japan1", 5)
	// CIDSystemInfoAdobeJapan6 : Japanese
	CIDSystemInfoAdobeJapan6 = newCIDSystemInfo("Adobe", "Japan1", 6)
	// CIDSystemInfoAdobeKorea1 : Korean
	CIDSystemInfoAdobeKorea1 = newCIDSystemInfo("Adobe", "Korea1", 1)
	// CIDSystemInfoAdobeKorea2 : Korean
	CIDSystemInfoAdobeKorea2 = newCIDSystemInfo("Adobe", "Korea1", 2)
)
