# Backend Bootstrap

## Purpose

Version is an implementation detail of `cmd/server`. The build embeds the version string via ldflags directly into the server binary — no separate module is needed or should be created. Air hosting is wired via a `just backend` recipe pointing at `cmd/server`.
