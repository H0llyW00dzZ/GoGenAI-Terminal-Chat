// Copyright (c) 2024 H0llyW00dzZ
// License: MIT License

package fun_stuff

import "os/exec"

// PingIP is a helper function that executes the ping command on the system and returns the result.
//
// Note: This marked as TODO for now, it would be error if you trying exectue this function.
func PingIP(ip string) (string, error) {
	out, err := exec.Command("ping", ip, "-c 4").Output() // '-c 4' is for sending 4 packets.
	if err != nil {
		return "", err
	}
	return string(out), nil
}
