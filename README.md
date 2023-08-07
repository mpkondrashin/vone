# VOne - Trend Micro Vision One Web API CLI tool

**Command line utility to use Vision One**

--------

## Parameters
Each parameter can be provided in three ways: 
1. Configuration file config.yaml. ```vone``` seeks for this file in its current folder or in the folder of ```vone``` executable itself.
2. Environment variables
3. Command line parameters (not all commands support all parameters):

| YAML option<br/>Command line<br/>Env variable | Description | 
| --------------------------------------------- | ----------- | 
| address<br/>--address<br/>VONE_ADDRESS | Vision One URL (See https://automation.trendmicro.com/xdr/Guides/Regional-Domains) |
| token<br/>--token<br/>VONE_TOKEN | Vision One Token (Generate using Vision One console) |
| filename<br/>--filename<br/>VONE_FILENAME | Path to file |
| mask<br/>--mask<br/>VONE_MASK | Files mask |
| url<br/>--url<br/>VONE_URL | URL |
| urlfile<br/>--urlfile<br/>VONE_URLFILE | Text file with URLs (one per line). Use "-" name to use stdin |
| timeout<br>--timeout<br>VONE_TIMEOUT | Timeout for sample analysis |
| log<br>--log<br>VONE_LOG | Log file path |
| query<br>--query<br>VONE_QUERY | Query expression |
| top<br>--top<br>VONE_TOP | Limit mount of downloaded data (50, 100, or 200) |

Any combination of parameters can be used with ```vone```. For example, creating following configuration file (config.yaml):
```yaml
address: api.xdr.trendmicro.com
token: eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJjaWQiOiJiZTg0OTU0NS1lNjc0LTQwZjAtOTlkYy1mYjU2NWYzMjQ3NjAiLCJjcGlkIjoic3ZwIiwicHBpZCI6ImN1cyIsIml0IjoxNjU5MzU3ODYzLCJ1aWQiOiJta29uZHJhc2hpbkBnbWFpbC5jb20iLCJwbCI6IqIsImV0IjoxNjkwODkzODYzfQ.AhWwdZEWp4BwEXl4Mukd3baVIAm848c6Y3TdhvIyhxjsAPMxqdmOV0RXYxeItdoFWt5ljxIS5LdsPtjERYq8QaB9CYD-tVd886KknUpxQ8llceo_wDKcKGRDIkrQU6UkHJsI4yeYvEZCKrkMPHTLG5-1xjClOK1IfzGHA-t_nNLYx3pFJS_VohKEDaPmKRM9Lnc6OQPju6k8wt-QxQ0ksq_qNu0ba0XL_cTe02lkLTt3TGYZgPwhkVPrH7_4Pe_vsIuF3r-r9VVYIPGmfqYuddnkLJopZ8heNOoal1WdtlFp_p-ckzcSAjWS9mxZDVp6W4HIr3heONzyebGVXMbTttWAe-V_b75VjcN6HLAjI4OxGiiU9Pm_ZOntlBIBNldncOsxl29WpZShIli_qh4PJilXPmpHRW4pxL9soSIMTRI7H5ALqVEK_6QxEEKR2dexvoB4uYG0wss5e1c9RMQveJqQ8soYfB0y0WyJ5vS2KzeU5EOlIR3Ql4XDIphxZkGMtfUKK3AKPY2J7QSHnyBKiJYo12Q03ZdDJAtveDwr0ADyWkwrmDqaHB86_PEbyWJtfIIBgG848g1R0YcRAow76_944U_mGcomU1N5PK2_SZOr6n9-HQz_99vmn23S2TPHB-R2oEN2snB3aXaI9VTdQWNqrtwQBQOFIcTJgIEwS_8
log: vone.log
```
Following command can be used to upload file, files by mask, url and urls from file (example for Linux command line):
```
./vone submit --filename file_name.exe --mask collection/*.exe --url http://www.google.com --urlfile url_list.txt
```

In this example, address and token are taken from config.yaml file.

**Note:** If the same parameter is provided in two ways, command line parameters have higher priority than environment variable and the latter, higher priority than configuration file.

## Commands

### Submit Single File
Submit file for analisys

Required parameters: address, token, filename
```commandline
./vone submit --filename <file_path> <options>
```

### Submit Files By Mask
Submit multiply for analisys

Required parameters: address, token, filename
```commandline
./vone submit --mask <mask> <options>
```

### Submit Single URL
Submit URL for analisys

Required parameters: address, token, filename
```commandline
./vone submit --url <url> <options>
```

### Submit URLs From File
Submit multiply urls for analisys

Required parameters: address, token, filename
```commandline
./vone submit --urlfile <filename> <options>
```

### Get Quota
Get amount of files to submit

Required parameters: address, token
```commandline
./vone quota <options>
```

# Go Library

If this repo is treated as go package to use Vision One Web API (github.com/mpkondrashin/vone), followig functions are supported:
- Get daily reserve
- Submit URLs to sandbox
- Submit file to sandbox
- List submissions
- Get submission status
- Get list of analysis results - not supported!
- Get analysis results
- Download analysis results
- Download Investigation Package
- Download suspicious object list

For package usage examples, please check cmd/vone folder and cmd/sandbox/main.go file.