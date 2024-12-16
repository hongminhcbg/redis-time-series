package rules

import (
	"github.com/go-logr/logr"
	"github.com/hyperjumptech/grule-rule-engine/ast"
)

type logListener struct {
	l logr.Logger
}

func (l *logListener) EvaluateRuleEntry(cycle uint64, entry *ast.RuleEntry, candidate bool) {
	l.l.Info("EvaluateRuleEntry", "cycle", cycle, "entry", entry, "candidate", candidate)
}

// ExecuteRuleEntry will be called by the engine if it execute a rule entry in a cycle
func (l *logListener) ExecuteRuleEntry(cycle uint64, entry *ast.RuleEntry) {
	l.l.Info("ExecuteRuleEntry", "cycle", cycle, "entry", entry)
}

// BeginCycle will be called by the engine every time it start a new evaluation cycle
func (l *logListener) BeginCycle(cycle uint64) {
	l.l.Info("BeginCycle", "cycle", cycle)
}
