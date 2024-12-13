package rules

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/hyperjumptech/grule-rule-engine/ast"
	"github.com/hyperjumptech/grule-rule-engine/builder"
	"github.com/hyperjumptech/grule-rule-engine/engine"
	"github.com/hyperjumptech/grule-rule-engine/pkg"
)

var rule = []byte(
	`
rule amount_a "Extact A" salience 1 {
  WHEN 
    In.GetTransType() == "A"
  THEN
    In.Extract("a", "$.data.a.amount");
  	Retract("amount_a");
}

rule amount_a2 "Extact A" salience 1 {
  WHEN 
    In.GetTransType() == "A"
  THEN
    In.Extract("a2", "$.data.a.amount");
  	Retract("amount_a2");
}

rule amount_b "Extact B" salience 2 {
  WHEN 
    In.GetTransType() == "B"
  THEN
    In.Extract("amount_b", "$.data.b.amount");
  	Retract("amount_b");
}

rule amount_c "Extact amount c" salience 4 {
  WHEN 
    In.GetTransType() == "C"
  THEN
    In.Extract("amount_c", "$.data.c.amount");
  	Retract("amount_c");
}
`)

var MainEngine IRule

func init() {
	MainEngine = &_rule{
		script: rule,
		ID:     0,
	}
}

type IRule interface {
	Execute(ctx context.Context, in any, log logr.Logger) error
	GetID() int64
}

type _rule struct {
	script []byte
	ID     int64
}

func (r *_rule) GetID() int64 {
	return r.ID
}

func (r *_rule) Execute(ctx context.Context, in any, log logr.Logger) error {
	lib := ast.NewKnowledgeLibrary()
	ruleBuilder := builder.NewRuleBuilder(lib)
	dataContext := ast.NewDataContext()
	err := dataContext.Add("In", in)
	if err != nil {
		log.Error(err, "buildDataContextError")
		return err
	}

	// Build normally
	err = ruleBuilder.BuildRuleFromResource("CalcRisk", "0.1.1", pkg.NewBytesResource(r.script))
	if err != nil {
		if reporter, ok := err.(*pkg.GruleErrorReporter); ok {
			for i, er := range reporter.Errors {
				log.Error(er, "Error", "Index", i)
				return err
			}
		} else {
			log.Error(err, "There should be GruleErrorReporter")
		}

		return err
	}

	kb, err := lib.NewKnowledgeBaseInstance("CalcRisk", "0.1.1")
	if err != nil {
		log.Error(err, "NewKnowledgeBaseInstanceError")
		return err
	}

	eng1 := &engine.GruleEngine{
		MaxCycle:  999_999_999,
		Listeners: []engine.GruleEngineListener{},
	}
	err = eng1.Execute(dataContext, kb)
	if err != nil {
		log.Error(err, "ExecuteError")
		return err
	}

	return nil
}
