package trace

import (
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

const (
	callBackBeforeName = "tracing:before"
	callBackAfterName  = "tracing:after"
)

var gormTracer = otel.Tracer("gorm.io/gorm")

var _ gorm.Plugin = &GormTracingPlugin{}

type GormTracingPlugin struct{}

func (op *GormTracingPlugin) Name() string {
	return "tracking_plugin"
}

func (op *GormTracingPlugin) Initialize(db *gorm.DB) (err error) {
	if err = db.Callback().Create().Before("gorm:before_create").Register(callBackBeforeName, before); err != nil {
		return errors.Wrapf(err, "gorm before_create register tracing error.")
	}
	if err = db.Callback().Query().Before("gorm:query").Register(callBackBeforeName, before); err != nil {
		return errors.Wrapf(err, "gorm query register tracing error.")
	}
	if err = db.Callback().Delete().Before("gorm:before_delete").Register(callBackBeforeName, before); err != nil {
		return errors.Wrapf(err, "gorm before_delete register tracing error.")
	}
	if err = db.Callback().Update().Before("gorm:setup_reflect_value").Register(callBackBeforeName, before); err != nil {
		return errors.Wrapf(err, "gorm setup_reflect_value register tracing error.")
	}
	if err = db.Callback().Row().Before("gorm:row").Register(callBackBeforeName, before); err != nil {
		return errors.Wrapf(err, "gorm row register tracing error.")
	}
	if err = db.Callback().Raw().Before("gorm:raw").Register(callBackBeforeName, before); err != nil {
		return errors.Wrapf(err, "gorm raw_ register tracing error.")
	}

	if err = db.Callback().Create().After("gorm:after_create").Register(callBackAfterName, after); err != nil {
		return errors.Wrapf(err, "gorm after_create register tracing error.")
	}
	if err = db.Callback().Query().After("gorm:after_query").Register(callBackAfterName, after); err != nil {
		return errors.Wrapf(err, "gorm after_query register tracing error.")
	}
	if err = db.Callback().Delete().After("gorm:after_delete").Register(callBackAfterName, after); err != nil {
		return errors.Wrapf(err, "gorm after_delete register tracing error.")
	}
	if err = db.Callback().Update().After("gorm:after_update").Register(callBackAfterName, after); err != nil {
		return errors.Wrapf(err, "gorm after_update register tracing error.")
	}
	if err = db.Callback().Row().After("gorm:row").Register(callBackAfterName, after); err != nil {
		return errors.Wrapf(err, "gorm row register tracing error.")
	}
	if err = db.Callback().Raw().After("gorm:raw").Register(callBackAfterName, after); err != nil {
		return errors.Wrapf(err, "gorm raw_ register tracing error.")
	}
	return nil
}

func before(db *gorm.DB) {
	ctx := db.Statement.Context
	if !trace.SpanFromContext(ctx).IsRecording() {
		return
	}

	var span trace.Span
	db.Statement.Context, span = gormTracer.Start(ctx, "gorm")
	span.SetAttributes(
		attribute.String("db.system", db.Statement.Name()),
	)
}

func after(db *gorm.DB) {
	span := trace.SpanFromContext(db.Statement.Context)
	defer span.End()

	if err := db.Error; err != nil && err != gorm.ErrRecordNotFound {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
}
