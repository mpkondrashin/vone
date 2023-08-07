/*
	Trend Micro Vision One API SDK
	(c) 2023 by Mikhail Kondrashin (mkondrashin@gmail.com)

	Sandbox API capabilities

	generate.go - all generate commands
*/

package vone

//go:generate enum -package=vone -type=RiskLevel -values=high,medium,low,noRisk
//go:generate enum -package=vone -type=Action -values=analyzeFile,analyzeUrl
//go:generate enum -package=vone -type=Status -values=succeeded,running,failed
//go:generate enum -package=vone -type=ErrorCode -values=AccessDenied,BadRequest,ConditionNotMet,InternalServerError,InvalidCredentials,NotFound,ParameterNotAccepted,RequestEntityTooLarge,TooManyRequests,Unsupported
