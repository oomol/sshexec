#! /usr/bin/env bash
set -e
set -u
cmd_file="/tmp/.cmd_$(uuidgen)"
# Default is oomol
USER=oomol
# https://github.com/containers/gvisor-tap-vsock/blob/f0f18025e5b7c7c281a11dfd81034641b40efe18/cmd/gvproxy/main.go#L56
IP_ADDR=192.168.127.254
# https://github.com/oomol/sshexec/blob/f6e0e1583fc874727d68cc5cc3213dff6867dd0e/pkg/define/const.go#L21
PORT=5322

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
	echo -n 'exec'
	echo -n ' '
	echo -n ssh -q -o StrictHostKeyChecking=no $USER@$IP_ADDR -p $PORT
	echo -n ' '
	output_args "$arg0" "$@"
}
write_cmd "$@" >"$cmd_file"

chmod +x "$cmd_file"
sync
exec "$cmd_file"