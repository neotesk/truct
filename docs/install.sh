#!/bin/bash

__dependencies=("lscpu awk grep curl")
for cmd in $__dependencies; do
    if ! [ -x "$( command -v $cmd )" ]; then
        echo "Error: $cmd is not installed." >&2
        exit 1
    fi
done;

__arch=$( lscpu | grep Architecture | awk {'print $2'} );
if [[ "$__arch" == "x86_64" ]]; then
    __arch="amd64";
elif [[ "$__arch" == "aarch64" ]]; then
    __arch="arm64";
elif [[ "$__arch" == "arm64" ]]; then
    __arch="arm64";
elif [[ "$__arch" == "i686" ]]; then
    __arch="i386";
elif [[ "$__arch" == arm* ]]; then
    __arch="arm";
fi;
__arch=$( echo $__arch | awk '{print tolower($0)}' )
__os=$( uname -s | awk '{print tolower($0)}' )

__filename="truct.$__os-$__arch"

echo "Fetching release for $__os-$__arch..."
__link=$( curl -s https://api.github.com/repos/neotesk/truct/releases/latest | grep "browser_download_url.*$__filename" | cut -d ":" -f 2,3 | tr -d \" | tr -d " " )

if [[ "$__link" == "" ]]; then
    echo "OS and architecture cannot be found in the releases. Cannot download." >&2
    exit 1;
fi;

_comm=""
if [ -x "$( command -v sudo )" ]; then
    echo "'sudo' is found in the system"
    _comm="sudo"
elif [ -x "$( command -v doas )" ]; then
    echo "'doas' is found in the system"
    _comm="doas"
else
    echo "Falling back to su binary"
    _comm="su -c"
fi

if [[ "$_comm" == "su -c" ]]; then
    eval "$_comm 'curl -Lo /usr/bin/truct \"$__link\" && chmod +x /usr/bin/truct'"
else
    eval "$_comm curl -Lo /usr/bin/truct \"$__link\" && $_comm chmod +x /usr/bin/truct"
fi

echo Done.