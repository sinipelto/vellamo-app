$env:GOOS = 'linux'
$env:GOARCH = 'amd64'

$target = 'peltoloc'
$bin = 'linux_amd64'

$bint = 'vellamo'
$bintdev = 'vellamo_dev'

$confdir = 'config'
$confdev = 'config.dev.json'

$binp = ".\out\$bin"
$confp = ".\$confdir\$conf"
$confpdev = "$confdir/$confdev"

$tgtp = "/opt/vellamo"
$bintgp = "$tgtp/$bint"
$bintgpdev = "$tgtp/$bintdev"

# Wrapper to debug print / actually execute commands
# function Exec {
# 	$ArgsJ = $Args -Join ' '
# 	Write-Output $ArgsJ
# }

function Exec {
	$ArgsJ = $Args -Join ' '
	Invoke-Expression $ArgsJ
}

# Determine if dev mode
$devMode = $false
if ($Args.Count -gt 0) {
	$devMode = $Args[0].ToLower() -eq "dev"
}

if ($devMode) {
	go build -ldflags="-X 'main.DEBUGS=true' -X 'main.CONFIGS=${confpdev}'" -o $binp .
} else {
	go build -o $binp .
}

if ($LastExitCode -ne 0) {
	Write-Output "ERROR: Build failed."
	exit $LastExitCode
}

Exec scp -r .\$confdir ${target}:"$tgtp/$confdir"

if ($LastExitCode -ne 0) {
	Write-Output "ERROR: Failed to copy configs over SCP."
	exit $LastExitCode
}

if ($devMode) {
	Exec scp $binp ${target}:"$bintgpdev"
} else {
	Exec scp $binp ${target}:"$bintgp"
}

if ($LastExitCode -ne 0) {
	Write-Output "ERROR: Failed to copy binary over SCP."
	exit $LastExitCode
}

if ($devMode) {
	Exec ssh ${target} -C "chmod 0750 $bintgpdev"
} else {
	Exec ssh ${target} -C "chmod 0750 $bintgp"
}
	
if ($LastExitCode -ne 0) {
	Write-Output "ERROR: Failed to chmod binary."
	exit $LastExitCode
}

# Execute on remote
# UPDATE: (RE)LAUNCHED BY INOTIFYWAIT SCRIPT
#Exec ssh ${target} -tC "$bintgp"
