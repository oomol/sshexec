#! /usr/bin/env bash
set -e
USER=ihexon
IP_ADDR=192.168.1.250
PORT=22
arg0="$(basename "$0")"

expand-q() {
	for i; do echo -n " ${i@Q} "; done
}

output_args() {
	echo -n \"
	expand-q "$@"
	echo -n \"
	echo
}

write_cmd() {
	echo "#! /usr/bin/env bash"
	echo -n ssh -q -o StrictHostKeyChecking=no $USER@$IP_ADDR -p $PORT
	echo -n ' '
	output_args "$arg0" "$@"
}

write_cmd "$@" >/tmp/.cmd

chmod +x /tmp/.cmd && /tmp/.cmd
