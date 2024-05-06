$ErrorActionPreference = 'Stop'
$toolsDir   = "$(Split-Path -parent $MyInvocation.MyCommand.Definition)"

$url        = "https://github.com/thompsonbear/snet-cli/releases/download/v$($env:ChocolateyPackageVersion)/snet-windows-386.zip"
$url64      = "https://github.com/thompsonbear/snet-cli/releases/download/v$($env:ChocolateyPackageVersion)/snet-windows-amd64.zip"

$packageArgs = @{
  packageName   = $env:ChocolateyPackageName
  unzipLocation = $toolsDir
  url           = $url
  url64bit      = $url64

  checksum      = '31790D561290BDA3FAA63AB0D19E80846AEE70D15FB234B50463FD87DA2F3631'
  checksumType  = 'sha256'
  checksum64    = '94DACC33BFE1183414BE4F457BD7021563F444B1DCF520EF316829A9B6A56BC8'
  checksumType64= 'sha256'
}

Install-ChocolateyZipPackage @packageArgs