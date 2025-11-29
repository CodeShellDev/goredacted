package service

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"syscall"
	"time"

	"github.com/codeshelldev/goredacted/internals/parser"
	"github.com/codeshelldev/goredacted/utils/clipboardutils"
	"github.com/codeshelldev/goredacted/utils/fsutils"
)

var lastClipboard string

func Spawn(flags ...string) error {
	execPath, err := os.Executable()

	if err != nil {
		return err
	}

	cmd := exec.Command(execPath, flags...)

	cmd.Env = append(os.Environ(), "DAEMON_MODE=1")

	cmd.Stdin = nil
	cmd.Stdout = nil
	cmd.Stderr = nil

	err = cmd.Start()
	if err != nil {
		return err
	}

	pid := cmd.Process.Pid

	err = fsutils.Write("/tmp/goredacted.pid", strconv.Itoa(pid))

	if err != nil {
		return err
	}

	return nil
}

func Despawn() error {
	data, err := fsutils.Read("/tmp/goredacted.pid")

	if err != nil {
		return err
	}

	pid, err := strconv.Atoi(data)
    if err != nil {
        return errors.New("invalid PID in file: " + err.Error())
    }

    process, err := os.FindProcess(pid)
    if err != nil {
        return errors.New("cannot find process: " + err.Error())
    }

    err = process.Signal(syscall.SIGTERM)
    if err != nil {
        return errors.New("failed to signal process: " + err.Error())
    }

	return nil
}

func Hook() {
	for {
		clip := clipboardutils.Read()
		
		if lastClipboard != clip {
			redacted := Parse(clip)

			clipboardutils.Write(redacted)

			lastClipboard = redacted
		}

		time.Sleep(time.Millisecond * 100)
	}
}

func Parse(content string) string {
	return parser.Parse(content)
}

func Redact() error {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return err
	}

	if (stat.Mode() & os.ModeCharDevice) != 0 {
		return errors.New("no input provided on stdin")
	}

	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		return err
	}

	redacted := Parse(string(data))

	_, err = os.Stdout.Write([]byte(redacted))
	return err
}

func Listen(addr string) {
	http.HandleFunc("/parse", func(w http.ResponseWriter, req *http.Request) {
    	log.Print(req.Method, " ", req.URL.Path, req.URL.RawQuery, " from ", req.RemoteAddr, "\n")

        if req.Method != http.MethodGet {
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
            return
        }

		redacted := Parse(clipboardutils.Read())

		response := map[string]any{
			"output": redacted,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
    })

    log.Print("Server listening on http://", addr, "\n")
    log.Fatal(http.ListenAndServe(addr, nil))
}