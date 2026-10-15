package service

import (
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
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, "ok\n")
	})

	mux.HandleFunc("/parse", func(w http.ResponseWriter, req *http.Request) {
    	log.Print(req.Method, " ", req.URL.Path, req.URL.RawQuery, " from ", req.RemoteAddr, "\n")

		if req.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		body, err := io.ReadAll(http.MaxBytesReader(w, req.Body, 10<<20))
		if err != nil {
			http.Error(w, "Invalid or oversized request body", http.StatusBadRequest)
			return
		}

		input := string(body)
		redacted := Parse(input)

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, redacted)
	})

	server := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	log.Print("Server listening on http://", addr)
	log.Print("Health checkpoint: http://", addr, "/health")

	err := server.ListenAndServe();
	if err != nil && err != http.ErrServerClosed {
		log.Fatal("Server failed: ", err)
	}
}