/*
	Trend Micro Vision One API SDK
	(c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)

	Sandbox API capabilities

	it_add_exception.go - add exceptions for Threat Intelligance
*/

package vone

import (
	"context"
	"fmt"
	"iter"
	"strconv"
	"strings"
)

type EventID int

const (
	EVENT_PROCESS EventID = iota + 1
	EVENT_FILE
	EVENT_CONNECTIO
	EVENT_DNS
	EVENT_REGISTRY
	EVENT_ACCOUNT
	EVENT_INTERNET
	XDR_EVENT_MODIFIED_PROCESS
	EVENT_WINDOWS_HOOK
	EVENT_WINDOWS_EVENT
	EVENT_AMSI
	EVENT_WMI
	TELEMETRY_MEMORY
	TELEMETRY_BM
)

//go:generate stringer -type EventID -trimprefix EventID

// Implement Unmarshaler interface for EventID
func (e *EventID) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), `"`)
	v, err := strconv.Atoi(s)
	if err != nil {
		return err
	}
	*e = EventID(v)
	return nil
}

// Implement Marshaler interface for EventID
func (e EventID) MarshalJSON() ([]byte, error) {
	return []byte(e.String()), nil
}

type ObjectSubTrueType int

const (
	MSOfficeOrUnknownOLE         ObjectSubTrueType = 1000
	Winword20                    ObjectSubTrueType = 1001
	Winword10                    ObjectSubTrueType = 1002
	HangulWordProcessor          ObjectSubTrueType = 1003
	Ichitaro                     ObjectSubTrueType = 1004
	JungUmGlobal                 ObjectSubTrueType = 1005
	OutlookItem                  ObjectSubTrueType = 1006
	MSPublisher                  ObjectSubTrueType = 1007
	MSOffice                     ObjectSubTrueType = 4000
	HancomHancell                ObjectSubTrueType = 4001
	UnknownTypeMSOffice2007      ObjectSubTrueType = 4045000
	WordMSOffice2007             ObjectSubTrueType = 4045001
	ExcelMSOffice2007            ObjectSubTrueType = 4045002
	PTTMSOffice2007              ObjectSubTrueType = 4045003
	XPSOffice2007                ObjectSubTrueType = 4045004
	HwpStandardOWPMLDocument     ObjectSubTrueType = 1004000
	UnknownXMLDocument           ObjectSubTrueType = 1003000
	Word2003XMLDocument          ObjectSubTrueType = 1003001
	XMLSpreadsheet2003           ObjectSubTrueType = 1003002
	PowerPointXMLPresentation    ObjectSubTrueType = 1003003
	FLI                          ObjectSubTrueType = 22000
	FLC                          ObjectSubTrueType = 22001
	FLIC                         ObjectSubTrueType = 22002
	WindowsWrite                 ObjectSubTrueType = 16000
	WordForDOS                   ObjectSubTrueType = 16001
	DOS_EXE                      ObjectSubTrueType = 7000
	WIN16_EXE                    ObjectSubTrueType = 7001
	WIN32_EXE                    ObjectSubTrueType = 7002
	OS2_EXE                      ObjectSubTrueType = 7003
	WIN16_DLL                    ObjectSubTrueType = 7004
	Win32_DLL                    ObjectSubTrueType = 7005
	WindowsVxD                   ObjectSubTrueType = 7006
	OS2_2xVxD                    ObjectSubTrueType = 7007
	NT_MIPS_EXE                  ObjectSubTrueType = 7008
	PKLITE_EXE                   ObjectSubTrueType = 7009
	LZEXE                        ObjectSubTrueType = 7010
	DIET_EXE                     ObjectSubTrueType = 7011
	PKZIP_EXE                    ObjectSubTrueType = 7012
	ARJ_EXE                      ObjectSubTrueType = 7013
	LZH_EXE                      ObjectSubTrueType = 7014
	LZH_EXE_UsedByZipMail        ObjectSubTrueType = 7015
	ASPACK                       ObjectSubTrueType = 7016
	UPX_EXE                      ObjectSubTrueType = 7017
	MSIL                         ObjectSubTrueType = 7018
	ASPACK2x                     ObjectSubTrueType = 7019
	WWPACK                       ObjectSubTrueType = 7020
	PETITE                       ObjectSubTrueType = 7021
	PEPACK                       ObjectSubTrueType = 7022
	MEW11                        ObjectSubTrueType = 7023
	MEW05                        ObjectSubTrueType = 7024
	MEW10                        ObjectSubTrueType = 7025
	AMD64_EXE                    ObjectSubTrueType = 7026
	AMD64_DLL                    ObjectSubTrueType = 7027
	ARM_EXE                      ObjectSubTrueType = 7028
	THUMB_EXE                    ObjectSubTrueType = 7029
	Miscellaneous_EXE            ObjectSubTrueType = 7030
	UPX64                        ObjectSubTrueType = 7031
	DLLNotProgram                ObjectSubTrueType = 7032
	DOS_COM                      ObjectSubTrueType = 5000
	PKLITE_COM                   ObjectSubTrueType = 5001
	DIET_COM                     ObjectSubTrueType = 5002
	LZH_COM                      ObjectSubTrueType = 5003
	ELFFile                      ObjectSubTrueType = 19000
	ELFReloactableFile           ObjectSubTrueType = 19001
	ELFEXEFile                   ObjectSubTrueType = 19002
	ELFLibraryFile               ObjectSubTrueType = 19003
	ELFCoredumpFile              ObjectSubTrueType = 19004
	EPOCBinaryFile               ObjectSubTrueType = 126000
	EPOCEXEFile                  ObjectSubTrueType = 126001
	EPOCLibraryFile              ObjectSubTrueType = 126002
	EPOCCompressedFile           ObjectSubTrueType = 126003
	SISFile                      ObjectSubTrueType = 128000
	SISFileAlt                   ObjectSubTrueType = 128001
	PDF                          ObjectSubTrueType = 6015000
	ISO9660File                  ObjectSubTrueType = 1005000
	UDFFile                      ObjectSubTrueType = 1005001
	FrameMakerDocumentFile       ObjectSubTrueType = 6004000
	FrameMakerMIFFile            ObjectSubTrueType = 6004001
	FrameMakerMMLFile            ObjectSubTrueType = 6004002
	FrameMakerBookFile           ObjectSubTrueType = 6004003
	FrameMakerDictionaryFile     ObjectSubTrueType = 6004004
	FrameMakerFontFile           ObjectSubTrueType = 6004005
	FrameMakerIPL                ObjectSubTrueType = 6004006
	AdobeFontMetrics             ObjectSubTrueType = 6001000
	AdobeFontBits                ObjectSubTrueType = 6001001
	AVI                          ObjectSubTrueType = 4020000
	WAV                          ObjectSubTrueType = 4020001
	BND                          ObjectSubTrueType = 4020002
	RMI                          ObjectSubTrueType = 4020003
	RDI                          ObjectSubTrueType = 4020004
	CDA                          ObjectSubTrueType = 4020005
	ANI                          ObjectSubTrueType = 4020006
	CMX                          ObjectSubTrueType = 4020007
	Compressed16Bits             ObjectSubTrueType = 2004000
	PackedData                   ObjectSubTrueType = 2004001
	CompackedData                ObjectSubTrueType = 2004002
	SCOCompressed                ObjectSubTrueType = 2004003
	ScriptVBScriptJavascript     ObjectSubTrueType = 28000
	HTMLFile                     ObjectSubTrueType = 28001
	PALMResourceCodeFile         ObjectSubTrueType = 28002
	ASPFile                      ObjectSubTrueType = 28003
	GeneralTextFile              ObjectSubTrueType = 28004
	ActionScript                 ObjectSubTrueType = 28005
	XMLDataPackage               ObjectSubTrueType = 28006
	MACBIN                       ObjectSubTrueType = 115000
	AutoCADDWG                   ObjectSubTrueType = 58000
	AutoCADR2000                 ObjectSubTrueType = 58001
	MBX_OUTLOOK4                 ObjectSubTrueType = 116000
	MBX_Unix                     ObjectSubTrueType = 116001
	MBX_FOXMAIL                  ObjectSubTrueType = 116002
	MSAccessMDB                  ObjectSubTrueType = 26000
	MSAccess2000XP               ObjectSubTrueType = 26001
	MSAccessMDB20                ObjectSubTrueType = 26002
	MSAccess2007                 ObjectSubTrueType = 26003
	OutlookMSOFile               ObjectSubTrueType = 6016000
	ExchangeMSOData              ObjectSubTrueType = 6016001
	UnknownOpenDocumentFile      ObjectSubTrueType = 4048000
	OpenDocumentTextFile         ObjectSubTrueType = 4048001
	OpenDocumentGraphicFile      ObjectSubTrueType = 4048002
	OpenDocumentPresentationFile ObjectSubTrueType = 4048003
	OpenDocumentSpreadsheetFile  ObjectSubTrueType = 4048004
	OpenDocumentFormulaFile      ObjectSubTrueType = 4048005
	OpenDocumentDatabaseFile     ObjectSubTrueType = 4048006
	UncompressedSWF              ObjectSubTrueType = 4038000
	CompressedSWF                ObjectSubTrueType = 4038001
	LZMACompressedSWF            ObjectSubTrueType = 4038002
	UnknownMatchObjectFile       ObjectSubTrueType = 4052000
	X86MatchObjectFile           ObjectSubTrueType = 4052001
	X64MatchObjectFile           ObjectSubTrueType = 4052002
	ONENOTE_ONE                  ObjectSubTrueType = 4055000
	ONENOTE_ONETOC2              ObjectSubTrueType = 4055001
	ONENOTE_ONE_FSSHTTP          ObjectSubTrueType = 4055002
	ONENOTE_ONETOC2_FSSHTTP      ObjectSubTrueType = 4055003
	FILESS_REG                   ObjectSubTrueType = 8000001
	FILESS_WMI                   ObjectSubTrueType = 8000002
	FILESS_SCHEDULE_TASK         ObjectSubTrueType = 8000003
	FILESS_AMSI                  ObjectSubTrueType = 8000004
	FILESS_MIP3                  ObjectSubTrueType = 8000005
	FILESS_CMD                   ObjectSubTrueType = 8000006
	FILESS_MIP3_X86_64           ObjectSubTrueType = 8000007
)

//go:generate stringer -type ObjectSubTrueType -trimprefix ObjectSubTrueType

// Implement Unmarshaler interface for EventID
func (e *ObjectSubTrueType) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), `"`)
	v, err := strconv.Atoi(s)
	if err != nil {
		return err
	}
	*e = ObjectSubTrueType(v)
	return nil
}

// Implement Marshaler interface for EventID
func (e ObjectSubTrueType) MarshalJSON() ([]byte, error) {
	s := e.String()
	if len(s) == 0 {
		return []byte(strconv.Itoa(int(e))), nil
	}
	return []byte(s), nil
}

type DataSource int

const (
	Detections DataSource = iota
	EndpointActivityData
	CloudActivityData
	EmailActivityData
	MobileActivityData
	NetworkActivityData
	ContainerActivityData
	IdentityActivityData
	ThirdPartyLogData
)

//go:generate stringer -type DataSource -trimprefix DataSource

// Implement Unmarshaler interface for DataSource
func (d *DataSource) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), `"`)
	for i := range ThirdPartyLogData + 1 {
		if s == DataSource(i).String() {
			*d = DataSource(i)
			return nil
		}
	}
	return nil
}

type (
	ObservedAttackTechniquesEventsResponse struct {
		TotalCount int                                  `json:"totalCount"`
		Count      int                                  `json:"count"`
		Items      []ObservedAttackTechniquesEventsItem `json:"items"`
		NextLink   string                               `json:"nextLink"`
	}

	OATFilter struct {
		ID                 string   `json:"id"`
		Name               string   `json:"name"`
		Description        string   `json:"description"`
		MitreTacticIds     []string `json:"mitreTacticIds"`
		MitreTechniqueIds  []string `json:"mitreTechniqueIds"`
		HighlightedObjects []struct {
			Type  string `json:"type"`
			Field string `json:"field"`
			Value any    `json:"value"`
		} `json:"highlightedObjects"`
		RiskLevel string `json:"riskLevel"`
		Type      string `json:"type"`
	}

	OATFilters   []OATFilter
	StringsSlice []string

	ObservedAttackTechniquesEventsItem struct {
		Source   DataSource `json:"source"`
		UUID     string     `json:"uuid"`
		Filters  OATFilters `json:"filters"`
		Endpoint struct {
			EndpointName string       `json:"endpointName"`
			AgentGUID    string       `json:"agentGuid"`
			Ips          StringsSlice `json:"ips"`
		} `json:"endpoint"`
		EntityType       string        `json:"entityType"`
		EntityName       any           `json:"entityName"` // Can be list, so not string
		DetectedDateTime VisionOneTime `json:"detectedDateTime"`
		IngestedDateTime VisionOneTime `json:"ingestedDateTime"`
		Detail           Detail        `json:"detail"`
	}

	Detail struct {
		Dpt                             int
		Dst                             string
		EndpointGuid                    string
		EndpointHostName                string
		EndpointIp                      StringsSlice
		EventId                         string
		EventSubId                      int
		ObjectIntegrityLevel            int
		ObjectTrueType                  int
		ObjectSubTrueType               int
		WinEventID                      int
		WinEventId                      int
		EventTime                       string
		EventTimeDT                     VisionOneTime
		HostName                        string
		LogonUser                       StringsSlice
		ObjectCmd                       string
		ObjectFileHashSha1              string
		ObjectFilePath                  string
		ObjectHostName                  string
		ObjectIp                        string
		ObjectIps                       StringsSlice
		ObjectPort                      int
		ObjectRegistryData              string
		ObjectRegistryKeyHandle         string
		ObjectRegistryValue             string
		ObjectSigner                    StringsSlice
		ObjectSignerValid               []bool
		ObjectUser                      string
		Os                              string
		ParentCmd                       string
		ParentFileHashSha1              string
		ParentFilePath                  string
		ProcessCmd                      string
		ProcessFileHashSha1             string
		ProcessFilePath                 string
		ProductCode                     string
		SearchDL                        string
		Spt                             int
		Src                             string
		SrcFileHashSha1                 string
		SrcFilePath                     string
		Tags                            StringsSlice
		Uuid                            string
		Act                             string // "Unknown" "N/A" "Clean" "Delete" "Move" "Rename" "Pass/Log" "Strip" "Drop" "Quarantine" "Insert/Replace" "Archive" "Stamp" "Block" "Redirect mail for approval" "Encrypted" "Reset" "Pass" "Terminate" "Assessment"
		ActResult                       string // "Unknown" "N/A" "File blocked(Deny Access)" "File quarantined" "File blocked(Defined By Product)" "File cleaned" "File deleted" "File quarantined" "File renamed" "File passed" "Unable to clean file. Passed" "Unable to clean file. Deleted" "Unable to clean file. Renamed" "Unable to clean file. Quarantined" "File stripped" "Unable to clean file. Stripped" "File replaced" "File dropped" "File archived" "Block successfully" "Quarantine successfully" "Stamp successfully" "File uploaded" "Unable to clean file. Quarantined" "Unable to clean file. Passed" "Access denied" "No action" "Reboot system successfully" "Spyware/Grayware unsafe to clean" "Stop scan manually successfully" "Redirect Mail for Approval successfully" "Encrypted" "Detect" "Unable to block file(Deny Access)" "Unable to quarantine file" "Unable to block file(Defined By Product)" "Unable to clean file" "Unable to delete file" "Unable to quarantine file" "Unable to rename file" "Unable to pass file" "Unable to clean or pass file" "Unable to clean or delete file" "Unable to clean or rename file" "Unable to clean or quarantine file" "Unable to strip file" "Unable to clean or strip file" "Unable to replace file" "Unable to drop file" "Unable to archive file" "Unable to block file" "Unable to quarantine file" "Unable to stamp file" "Unable to upload file" "Unable to clean or quarantine file" "Unable to clean or pass file" "Unable to deny access" "Unable to perform no action" "Action Required - Restart the endpoint to finish cleaning the security threat" "Unsafe to Clean File" "Unable to stop scan manually" "Unable to Redirect Mail for Approval" "Action required – Perform a full system scan" "Action required - Use the \\\"Rescue Disk\\\" tool in the OfficeScan ToolBox to remove this threat. If the problem persists" "contact Support" "Action required - Use the \\\"Rootkit Buster\\\" tool in the OfficeScan ToolBox to remove this threat. If the problem persists" "contact Support" "Action required - Use the \\\"Clean Boot\\\" tool in the OfficeScan ToolBox to remove this threat. If the problem persists" "contact Support"
		App                             string // "HTTP" "HTTPS" "KERBEROS" "TCP" "SMB" "SMB2" "DCE-RPC" "LDAP" "DNS Response"
		AppGroup                        string // "HTTP" "AUTH" "CIFS" "TCP" "LDAP" "DNS Response"
		AptRelated                      bool
		BehaviorCat                     string
		Blocking                        string //"Suspicious behavior" "Fraud" "Suspicious spyware/grayware" "Suspicious virus/malware" "IntelliTrap" "Network virus" "Spyware/Grayware" "Virus/Malware" "Unknown" "File Name" "Webmail site" "Web Server" "URL Pattern" "Java/VB Script" "True File Type" "User-defined" "Server-defined" "Web Policy" "Phish" "Phish/Spyware/Grayware" "Phish/Virus/Malware accomplice" "Phish/Forged signature" "Phish/Disease vector" "Phish/Malicious applet" "Phishing reputation" "Policy IP translation" "Policy Java Scanning" "Policy Malicious Mobile Code" "Pharming" "URL Blocking" "URL Filtering" "Client IP Blocking" "Destination Port Blocking" "Web reputation" "Unsupported file type" "Exceeds total file count limit" "Exceeds file size limit" "Exceeds decompression layer limit" "Exceeds decompression time limit" "Exceeds compression ratio limit" "Password protected file" "Restricted spyware/grayware type" "String pattern"
		Cat                             string
		CccaDetection                   string
		CccaDetectionSource             string // "GLOBAL_INTELLIGENCE" "VIRTUAL_ANALYZER" "USER_DEFINED" "RELEVANCE_RULE"
		CccaRiskLevel                   int    // SLF_CCCA_RISKLEVEL_UNKNOWN (0) SLF_CCCA_RISKLEVEL_LOW (1) SLF_CCCA_RISKLEVEL_MEDIUM (2) SLF_CCCA_RISKLEVEL_HIGH (3)
		ClientFlag                      string // Enum: "Unknown" "src" "dst" 0:Unknown 1:src 2:dst
		Cnt                             string
		Component                       StringsSlice
		CompressedFileSize              string
		DetectionType                   string
		DeviceDirection                 string // Enum: "inbound" "outbound" "unknown" 		0: inbound 1: outbound 2: unknown (If cannot be parsed correctly, 2 is assigned)
		DeviceGUID                      string
		DeviceProcessName               string
		DeviceMacAddress                string
		Dhost                           string
		DomainName                      string
		DstGroup                        string
		End                             string
		EndpointGUID                    string
		EndpointMacAddress              StringsSlice
		EngType                         string
		EngVer                          string
		EventName                       string
		EventSubName                    string
		FileHash                        string
		FileName                        StringsSlice
		FileOperation                   string
		FilePath                        string
		FilePathName                    string
		FileSize                        string
		FirstAct                        string
		FirstActResult                  string
		FullPath                        string
		HttpReferer                     string
		InterestedHost                  string
		InterestedIp                    StringsSlice
		InterestedMacAddress            string
		MalName                         string
		MalType                         string
		MDevice                         StringsSlice
		MDeviceGUID                     string
		MitreMapping                    StringsSlice
		MitreVersion                    string
		Mpname                          string
		Mpver                           string
		ObjectFileHashMd5               string
		ObjectFileHashSha256            string
		ObjectFileName                  string
		ObjectName                      string
		ObjectPid                       int
		ParentFileHashSha256            string
		PeerHost                        string
		PeerIp                          StringsSlice
		Pname                           string
		ProcessFileHashMd5              string
		ProcessFileHashSha256           string
		ProcessName                     string
		ProcessPid                      int
		ProcessSigner                   StringsSlice
		Pver                            string
		Request                         string
		RequestClientApplication        string
		Rt                              string
		Rt_utc                          string
		SrcGroup                        string
		TacticId                        StringsSlice
		ThreatName                      string
		ClusterId                       string
		ClusterName                     string
		K8sNamespace                    string
		K8sPodName                      string
		K8sPodId                        string
		ContainerName                   string
		ContainerId                     string
		ContainerImage                  string
		ContainerImageDigest            string
		RuleIdStr                       string
		RuleSetName                     string
		RuleSetId                       string
		AwsRegion                       string
		CustomerId                      string
		EventCategory                   string //Management Data Insights
		EventID                         string
		EventSource                     string
		EventCase                       string
		EventType                       string // AwsApiCall AwsServiceEvent AwsConsoleAction AwsConsoleSignIn AwsCloudTrailInsight
		EventVersion                    string
		ManagementEvent                 bool
		PackageTraceId                  string
		ReadOnly                        bool
		RecipientAccountId              string
		RequestID                       string
		ResponseElements                string
		SourceIPAddress                 string
		PrincipalName                   string
		ServerTls                       string
		ServerProtocol                  string
		UserAgent                       string
		TenantGuid                      string
		Application                     string
		RuleName                        string
		ClientIp                        string
		RequestBase                     string
		Score                           int
		UserDomain                      StringsSlice
		Suid                            string
		Duration                        int
		FileHashSha256                  string
		FFileSize                       string
		FileType                        string
		MimeType                        string
		Sender                          string
		Profile                         string
		UserDepartment                  string
		RequestMethod                   string
		RequestMimeType                 string
		RuleType                        string
		RuleUuid                        string
		ObjectId                        string
		PolicyUuid                      string
		CompanyName                     string
		Start                           int
		EndpointModel                   string
		OsName                          string
		OsVer                           string
		FirstSeen                       string
		UserType                        string
		ObjectAppPackageName            string
		ObjectAppInstalledTime          string
		ObjectAppLabel                  string
		ObjectAppSize                   string
		ObjectAppIsSystemApp            bool
		ObjectAppVerCode                string
		ObjectAppSha256                 string
		ObjectAppPublicKeySha1          string
		AppLabel                        string
		AppPkgName                      string
		AppPublicKeySha1                string
		AppSize                         string
		AppIsSystem                     string
		AppVerCode                      string
		FilterRiskLevel                 string
		MailMsgSubject                  string
		MailMsgId                       string
		MsgUuid                         string
		Mailbox                         string
		MailSenderIp                    string
		MailFromAddresses               StringsSlice
		MailWholeHeader                 StringsSlice
		MailToAddresses                 StringsSlice
		MailSourceDomain                string
		ScanType                        string
		Org_id                          string
		MailUrlsVisibleLink             StringsSlice
		MailUrlsRealLink                StringsSlice
		Version                         string
		EventSourceType                 int
		ReceivedTime                    string
		GroupId                         string
		BitwiseFilterRiskLevel          int
		CustomFilterTags                string
		MgmtInstanceId                  string
		IdpName                         string
		IdpId                           string
		LocationCountry                 string
		LocationCity                    string
		LocationState                   string
		LocationLongitude               string
		LocationLatitude                string
		ClientId                        string
		ClientDisplayName               string
		ClientOS                        string
		ClientBrowser                   string
		ClientApp                       string
		IpAddress                       string
		UserId                          string
		UserDisplayName                 string
		StatusDetail                    string
		Status                          string
		StatusReason                    string
		TargetResourceId                string
		TargetResourceDisplayName       string
		EventAdditionalDetails          string
		InitiatedByAppId                string
		InitiatedByAppDisplayName       string
		InitiatedByServicePrincipalId   string
		InitiatedByServicePrincipalName string
		InitiatedByUserId               string
		InitiatedByUserDisplayName      string
		InitiatedByUserHomeTenantId     string
		InitiatedByUserHomeTenantName   string
		InitiatedByUserIpAddress        string
		InitiatedByUserPrincipalName    string
		LoggedByService                 string
		OperationType                   string
		Result                          string
		ResultReason                    string
		TargetResources                 StringsSlice
		PartitionKey                    string
		CustomFilterRiskLevel           string
		TmFilterRiskLevel               string
		ActionName                      string
		CorrelationId                   string
		PolicyTreePath                  string
		UUID                            string `json:"uuid"`
		LogReceivedTime                 string // int
		Vendor                          string
		Severity                        int
		PolicyName                      string
		SUser1                          string
		DUser1                          string
		VsysName                        string
		SrcZone                         string
		DstZone                         string
		FlowId                          string
		Category                        string
		ReqDataSize                     uint64
		RespDataSize                    uint64
		SessionStart                    int
		SessionEnd                      int
		SessionEndReason                string
		UrlCat                          StringsSlice
		SrcLocation                     string
		DstLocation                     string
		Dvchost                         string
		HttpXForwardedFor               string
		SOSName                         string
		Shost                           string
		Smac                            string
		DOSName                         string
		Dmac                            string
		Requests                        StringsSlice
		RuleId                          int
		Rating                          string
		Direction                       string
		HttpRespContentType             string
		FileHashMd5                     string
		OldFileHash                     string
		Suser                           StringsSlice
		Duser                           StringsSlice
		Dvc                             StringsSlice
		Proto                           string
		VendorParsed                    string
		VendorRaw                       string
		CloudAccountId                  string
		NetworkInterfaceId              string
		IpProto                         int
		Packets                         int
		Bytes                           int
		Action                          string
		LogStatus                       string
		VpcId                           string
		SubnetId                        string
		InstanceId                      string
		TcpFlags                        int
		FlowType                        string
		PktSrcAddr                      string
		PktDstAddr                      string
		RegionCode                      string
		AzId                            string
		SubLocationType                 string
		SubLocationId                   string
		PktSrcCloudServiceName          string
		PktDstCloudServiceName          string
		FlowDirection                   string
		TrafficPath                     int
		RequestParameters               string
		Resources                       StringsSlice
		UserIdentity                    string
		VpcEndpointId                   string
		AdditionalEventData             string
		ApiVersion                      string
		ErrorCode                       string
		ErrorMessage                    string
		SharedEventID                   string
		ServiceEventDetails             string
		TlsDetails                      string
	}
)

// Convert OATFilter the internal date as CSV string
func (f OATFilters) MarshalCSV() (string, error) {
	var sb strings.Builder
	for _, filter := range f {
		sb.WriteString(filter.ID)
		sb.WriteString(",")
		sb.WriteString(filter.Name)
		sb.WriteString(",")
		sb.WriteString(filter.Description)
		sb.WriteString(",")
		sb.WriteString(strings.Join(filter.MitreTacticIds, "|"))
		sb.WriteString(",")
		sb.WriteString(strings.Join(filter.MitreTechniqueIds, "|"))
		sb.WriteString(",")
		for _, o := range filter.HighlightedObjects {
			sb.WriteString(o.Type)
			sb.WriteString(":")
			sb.WriteString(o.Field)
			sb.WriteString(":")
			sb.WriteString(fmt.Sprintf("%v", o.Value))
			sb.WriteString(";")
		}
		sb.WriteString(",")
		sb.WriteString(filter.RiskLevel)
		sb.WriteString(",")
		sb.WriteString(filter.Type)
		sb.WriteString(";")
	}
	return sb.String(), nil
}

// Conver IPs to string
func (s StringsSlice) MarshalCSV() (string, error) {
	return strings.Join(s, "|"), nil
}

type GetOATEventsFunc struct {
	baseFunc
	Response ObservedAttackTechniquesEventsResponse
	top      int
}

var _ vOneFunc = &GetOATEventsFunc{}

func (v *VOne) GetOATEvents() *GetOATEventsFunc {
	f := &GetOATEventsFunc{}
	f.baseFunc.init(v)
	return f
}

func (f *GetOATEventsFunc) DetectedStart(vOneTime VisionOneTime) *GetOATEventsFunc {
	f.setParameter("detectedStartDateTime", vOneTime.String())
	return f
}

func (f *GetOATEventsFunc) DetectedEnd(vOneTime VisionOneTime) *GetOATEventsFunc {
	f.setParameter("detectedEndDateTime", vOneTime.String())
	return f
}

func (f *GetOATEventsFunc) IngestedStart(vOneTime VisionOneTime) *GetOATEventsFunc {
	f.setParameter("ingestedStartDateTime", vOneTime.String())
	return f
}

func (f *GetOATEventsFunc) IngestedEnd(vOneTime VisionOneTime) *GetOATEventsFunc {
	f.setParameter("ingestedEndDateTime", vOneTime.String())
	return f
}

// Filter - filter events
func (f *GetOATEventsFunc) Filter(filter string) *GetOATEventsFunc {
	f.setHeader("TMV1-Filter", filter)
	f.Response.NextLink = ""
	return f
}

// Top - set limit for returned amount of items
func (f *GetOATEventsFunc) Top(t Top) *GetOATEventsFunc {
	f.setParameter("top", t.String())
	f.top = t.Int()
	return f
}

// Iterate - get all events matching query one by one. If callback returns
// non nil error, iteration is aborted and this error is returned
func (f *GetOATEventsFunc) Iterate(ctx context.Context,
	callback func(item *ObservedAttackTechniquesEventsItem) error) error {
	for {
		response, err := f.Do(ctx)
		if err != nil {
			return err
		}
		for n := range response.Items {
			if err := callback(&response.Items[n]); err != nil {
				return err
			}
		}
		if response.NextLink == "" {
			break
		}
		if response.Count != f.top {
			break
		}
	}
	return nil
}

// Range - iterator for all endpoints matching query (go 1.23 and later)
func (f *GetOATEventsFunc) Range(ctx context.Context) iter.Seq2[*ObservedAttackTechniquesEventsItem, error] {
	return func(yield func(*ObservedAttackTechniquesEventsItem, error) bool) {
		for {
			response, err := f.Do(ctx)
			if err != nil {
				yield(nil, err)
				return
			}
			for n := range response.Items {
				if !yield(&response.Items[n], nil) {
					return
				}
			}
			if response.NextLink == "" {
				break
			}
			if response.Count != f.top {
				break
			}
		}
	}
}

// Do - execute the API call
func (f *GetOATEventsFunc) Do(ctx context.Context) (*ObservedAttackTechniquesEventsResponse, error) {
	if err := f.vone.call(ctx, f); err != nil {
		return nil, err
	}
	return &f.Response, nil
}

func (f *GetOATEventsFunc) method() string {
	return methodGet
}

func (s *GetOATEventsFunc) url() string {
	return "/v3.0/oat/detections"
}

func (f *GetOATEventsFunc) responseStruct() any {
	return &f.Response
}
