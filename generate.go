/*
	Trend Micro Vision One API SDK
	(c) 2023 by Mikhail Kondrashin (mkondrashin@gmail.com)

	Sandbox API capabilities

	generate.go - all generate commands
*/

package vone

//go:generate enum -package=vone -type=RiskLevel -names=high,medium,low,noRisk
//go:generate enum -package=vone -type=Action -names=analyzeFile,analyzeUrl
//go:generate enum -package=vone -type=Status -names=succeeded,running,failed
//go:generate enum -package=vone -type=ErrorCode -names=AccessDenied,BadRequest,ConditionNotMet,InternalServerError,InvalidCredentials,NotFound,ParameterNotAccepted,RequestEntityTooLarge,TooManyRequests,Unsupported
