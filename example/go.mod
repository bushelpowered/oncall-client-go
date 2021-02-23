module main

go 1.15

require (
	github.com/bushelpowered/oncall-client-go/oncall v0.0.0
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.8.0
)

replace github.com/bushelpowered/oncall-client-go/oncall v0.0.0 => ../oncall
