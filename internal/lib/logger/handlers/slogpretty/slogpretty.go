package slogpretty

import (
	"context"
	"encoding/json"
	"io"
	stdLog "log"
	"log/slog"
	"sync"

	"github.com/fatih/color"
)

type PrettyHandlerOptions struct {
	SlogOpts   *slog.HandlerOptions
	ForceColor bool
}

type PrettyHandler struct {
	opts   PrettyHandlerOptions
	mu     sync.Mutex
	l      *stdLog.Logger
	attrs  []slog.Attr
	groups []string
}

func (opts PrettyHandlerOptions) NewPrettyHandler(out io.Writer) *PrettyHandler {
	if opts.SlogOpts == nil {
		opts.SlogOpts = &slog.HandlerOptions{}
	}

	color.NoColor = false
	if opts.ForceColor {
		color.NoColor = false
	}

	return &PrettyHandler{
		opts:  opts,
		l:     stdLog.New(out, "", 0),
		attrs: make([]slog.Attr, 0),
	}
}

func (h *PrettyHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.opts.SlogOpts.Level.Level()
}

func (h *PrettyHandler) Handle(_ context.Context, r slog.Record) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	level := r.Level.String() + ":"

	switch r.Level {
	case slog.LevelDebug:
		level = color.MagentaString(level)
	case slog.LevelInfo:
		level = color.BlueString(level)
	case slog.LevelWarn:
		level = color.YellowString(level)
	case slog.LevelError:
		level = color.RedString(level)
	}

	fields := make(map[string]any)
	r.Attrs(func(a slog.Attr) bool {
		h.addAttr(fields, a)
		return true
	})

	for _, a := range h.attrs {
		h.addAttr(fields, a)
	}

	var b []byte
	if len(fields) > 0 {
		var err error
		b, err = json.MarshalIndent(fields, "", "  ")
		if err != nil {
			return err
		}
	}

	timeStr := r.Time.Format("[15:04:05.000]")
	msg := color.CyanString(r.Message)
	jsonStr := color.WhiteString(string(b))

	h.l.Println(
		timeStr,
		level,
		msg,
		jsonStr,
	)

	return nil
}

func (h *PrettyHandler) addAttr(m map[string]any, a slog.Attr) {
	key := a.Key
	if len(h.groups) > 0 {
		key = h.groups[len(h.groups)-1] + "." + key
	}
	m[key] = a.Value.Any()
}

func (h *PrettyHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) == 0 {
		return h
	}
	return &PrettyHandler{
		opts:   h.opts,
		l:      h.l,
		attrs:  append(h.attrs, attrs...),
		groups: h.groups,
	}
}

func (h *PrettyHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}
	groups := make([]string, len(h.groups)+1)
	copy(groups, h.groups)
	groups[len(groups)-1] = name
	return &PrettyHandler{
		opts:   h.opts,
		l:      h.l,
		attrs:  h.attrs,
		groups: groups,
	}
}
