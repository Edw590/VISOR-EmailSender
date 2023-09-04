module EmailSender

// Keep it on 1.20, so that it can be compiled for Windows 7 too if it's compiled with Go 1.20 (it's the last version
// supporting it).
go 1.20

require VISOR_Server/Utils v0.0.0-00010101000000-000000000000

require (
	github.com/dchest/jsmin v0.0.0-20220218165748-59f39799265f // indirect
	github.com/ztrue/tracerr v0.4.0 // indirect
)

replace VISOR_Server/Utils => ./Utils
