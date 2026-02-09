/*
	Trend Micro Vision One API SDK
	(c) 2026 by Mikhail Kondrashin (mkondrashin@gmail.com)

	Sandbox API capabilities

	generate.go - all generate commands
*/

package vone

//go:generate enum -package=vone -type=RiskLevel -names=high,medium,low,noRisk
//go:generate enum -package=vone -type=Action -names=analyzeFile,analyzeUrl
//go:generate enum -package=vone -type=Status -names=succeeded,running,failed
//go:generate enum -package=vone -type=SO -names=Domain,IP,SenderMailAddress,FileSha1,FileSha256
//go:generate enum -package=vone -type=ErrorCode -names=OK,AccessDenied,BadRequest,ConditionNotMet,InternalServerError,InvalidCredentials,NotFound,ParameterNotAccepted,RequestEntityTooLarge,TooManyRequests,Unsupported
//go:generate enum -package=vone -type=AlertStatus -names=Open,InProgress,Closed
//go:generate enum -package=vone -type=InvestigationResult -names "No Findings,Noteworthy,True Positive,False Positive,Benign True Positive,Other Findings"

