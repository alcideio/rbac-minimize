package analysis

import (
	"encoding/json"
	"fmt"
	v1 "k8s.io/api/rbac/v1"
	"reflect"
	"time"

	"github.com/alcideio/rbac-tool/pkg/rbac"
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker/decls"
	"google.golang.org/protobuf/proto"
	"k8s.io/klog"
)

type Analyzer interface {
	Analyze() (*AnalysisReport, error)
}

func CreateAnalyzer(config *AnalysisConfig, policies []rbac.SubjectPolicyList) Analyzer {
	analyzer := analyzer{
		config:           *config,
		policies:         policies,
		rules:            []*analysisRule{},
		globalExclusions: []*exclusion{},
	}

	if err := analyzer.initialize(); err != nil {
		klog.Errorf("Failed to initialize Analyzer - %v", err)
		return nil
	}

	return &analyzer
}

type analysisRule struct {
	rule *Rule

	//Internal State
	compiledAnalysisExpr cel.Program

	//Internal State
	compiledRecommendationExpr cel.Program

	exclusions []*exclusion
}

func createAnalysisExpr(expr string) (cel.Program, error) {

	d := cel.Declarations(
		decls.NewVar("subjects", decls.Dyn),
	)

	env, err := cel.NewEnv(d)
	if err != nil {
		return nil, err
	}

	ast, iss := env.Compile(expr)
	// Check iss for compilation errors.
	if iss.Err() != nil {
		return nil, iss.Err()
	}

	// Type-check the expression for correctness.
	checked, iss := env.Check(ast)

	// Report semantic errors, if present.
	if iss.Err() != nil {
		return nil, iss.Err()
	}

	// Check the result type is a string.
	if !proto.Equal(checked.ResultType(), decls.NewListType(decls.Dyn)) {
		return nil, fmt.Errorf("Got %v, wanted %v result type", checked.ResultType(), decls.NewListType(decls.Dyn))
	}

	prg, err := env.Program(checked)
	if err != nil {
		return nil, err
	}

	return prg, nil
}

func createRecommendationExpr(expr string) (cel.Program, error) {

	d := cel.Declarations(
		decls.NewVar("subject", decls.Dyn),
	)

	env, err := cel.NewEnv(d)
	if err != nil {
		return nil, err
	}

	ast, iss := env.Compile(expr)
	// Check iss for compilation errors.
	if iss.Err() != nil {
		return nil, iss.Err()
	}

	// Type-check the expression for correctness.
	checked, iss := env.Check(ast)

	// Report semantic errors, if present.
	if iss.Err() != nil {
		return nil, iss.Err()
	}

	// Check the result type is a string.
	if !proto.Equal(checked.ResultType(), decls.String) {
		return nil, fmt.Errorf("Got %v, wanted %v result type", checked.ResultType(), decls.String)
	}

	prg, err := env.Program(checked)
	if err != nil {
		return nil, err
	}

	return prg, nil
}

func newAnalysisRule(rule *Rule) (*analysisRule, error) {
	r := &analysisRule{
		rule:       rule,
		exclusions: []*exclusion{},
	}

	compiledAnalysisExpr, err := createAnalysisExpr(rule.AnalysisExpr)
	if err != nil {
		return nil, err
	}
	r.compiledAnalysisExpr = compiledAnalysisExpr

	compiledRecommendationExpr, err := createRecommendationExpr(rule.Recommendation)
	if err != nil {
		return nil, err
	}
	r.compiledRecommendationExpr = compiledRecommendationExpr

	for i := range rule.Exclusions {
		anExclusion, err := newExclusion(&rule.Exclusions[i])
		if err != nil {
			return nil, err
		}
		klog.V(5).Infof("Initialized Rule Exclusion '%v'", rule.Exclusions[i].Comment)
		r.exclusions = append(r.exclusions, anExclusion)
	}

	return r, nil
}

type exclusion struct {
	exclusion *Exclusion

	//Internal State
	compiledExceptionExpr cel.Program
}

func createExclusionExpr(expr string) (cel.Program, error) {

	d := cel.Declarations(
		decls.NewVar("subject", decls.Dyn),
	)

	env, err := cel.NewEnv(d)
	if err != nil {
		return nil, err
	}

	ast, iss := env.Compile(expr)
	// Check iss for compilation errors.
	if iss.Err() != nil {
		return nil, iss.Err()
	}

	// Type-check the expression for correctness.
	checked, iss := env.Check(ast)

	// Report semantic errors, if present.
	if iss.Err() != nil {
		return nil, iss.Err()
	}

	// Check the result type is a string.
	if !proto.Equal(checked.ResultType(), decls.Bool) {
		return nil, fmt.Errorf("Got %v, wanted %v result type", checked.ResultType(), decls.Bool)
	}

	prg, err := env.Program(checked)
	if err != nil {
		return nil, err
	}

	return prg, nil
}

func newExclusion(exclusionInfo *Exclusion) (*exclusion, error) {
	r := &exclusion{
		exclusion: exclusionInfo,
	}

	compiledExclusionExpr, err := createExclusionExpr(exclusionInfo.Expression)
	if err != nil {
		return nil, err
	}
	r.compiledExceptionExpr = compiledExclusionExpr

	return r, nil
}

type analyzer struct {
	config   AnalysisConfig
	policies []rbac.SubjectPolicyList

	Findings []AnalysisReportFinding

	rules []*analysisRule

	globalExclusions []*exclusion

	policiesObj interface{}
}

func (a *analyzer) initialize() error {

	b, err := json.Marshal(map[string]interface{}{"subjects": a.policies})
	if err != nil {
		return err
	}

	for i := range a.config.GlobalExclusions {
		anExclusion, err := newExclusion(&a.config.GlobalExclusions[i])
		if err != nil {
			return err
		}
		klog.V(5).Infof("Initialized Global Exclusion '%v'", a.config.GlobalExclusions[i].Comment)
		a.globalExclusions = append(a.globalExclusions, anExclusion)
	}

	m := make(map[string]interface{})
	if err := json.Unmarshal(b, &m); err != nil {
		return err
	}
	a.policiesObj = m["subjects"]

	for i, _ := range a.config.Rules {
		aRule, err := newAnalysisRule(&a.config.Rules[i])
		if err != nil {
			return err
		}
		klog.V(5).Infof("Initialized Rule '%v'", a.config.Rules[i].Name)
		a.rules = append(a.rules, aRule)
	}

	return nil
}

func (a *analyzer) shouldExclude(subject map[string]interface{}, exclusions []*exclusion) (bool, error) {
	for _, exclusion := range exclusions {
		recommendationOutput, _, err := exclusion.compiledExceptionExpr.Eval(map[string]interface{}{
			"subject": subject,
		})

		if err != nil {
			return false, err
		}

		exclude, ok := recommendationOutput.Value().(bool)
		if !ok {
			return false, fmt.Errorf("Failed to cast exclusion result '%v'", exclusion.exclusion.Comment)
		}

		if exclude {
			return true, nil
		}
	}

	return false, nil
}

func (a *analyzer) Analyze() (*AnalysisReport, error) {
	analysisStats := AnalysisStats{
		RuleCount: len(a.config.Rules),
	}
	report := AnalysisReport{
		AnalysisConfigInfo: AnalysisConfigInfo{
			Name:        a.config.Name,
			Description: a.config.Description,
			Uuid:        a.config.Uuid,
		},
		CreatedOn: time.Now().Format(time.RFC3339),
		Findings:  []AnalysisReportFinding{},
	}

	errs := []error{}
	for _, rule := range a.rules {

		out, _, err := rule.compiledAnalysisExpr.Eval(map[string]interface{}{
			"subjects": a.policiesObj,
		})

		if err != nil {
			klog.Errorf("Failed to evaluate rule '%v' - %v", rule.rule.Name, err)
			errs = append(errs, err)
			continue
		}

		outObj, err := out.ConvertToNative(reflect.TypeOf([]interface{}{}))
		if err != nil {
			klog.Fatalf("Failed to evaluate rule '%v' - %v", rule.rule.Name, err)
		}

		subjects, ok := outObj.([]interface{})

		if !ok {
			klog.Fatalf("Failed to cast - %v", reflect.TypeOf(outObj).Name())
		}

		if len(subjects) == 0 {
			klog.V(4).Infof("Rule - '%v' - no match", rule.rule.Name)
			continue
		}

		klog.V(5).Infof("Rule - '%v' - matched \n%v\n", rule.rule.Name, subjects)

		for _, subject := range subjects {
			sub := subject.(map[string]interface{})

			exclude, err := a.shouldExclude(sub, rule.exclusions)
			if err != nil {
				klog.Errorf("Failed to check exclusion for rule '%v' and subject %v - %v", rule.rule.Name, sub, err)
				errs = append(errs, err)
				//Continue on error - assume malformed exception expression
			}

			if exclude {
				analysisStats.ExclusionCount++
				klog.V(5).Infof("Skipping subject '%v' from rule exclusion - %v", sub, rule.rule.Name)
				continue
			}

			exclude, err = a.shouldExclude(sub, a.globalExclusions)
			if err != nil {
				klog.Errorf("Failed to check global exclusion for rule '%v' and subject %v - %v", rule.rule.Name, sub, err)
				errs = append(errs, err)
				//Continue on error - assume malformed exception expression
			}

			if exclude {
				analysisStats.ExclusionCount++
				klog.V(5).Infof("Skipping subject '%v' from rule exclusion - %v", sub, rule.rule.Name)
				continue
			}

			recommendationOutput, _, err := rule.compiledRecommendationExpr.Eval(map[string]interface{}{
				"subject": sub,
			})

			if err != nil {
				klog.Errorf("Failed to render recommendation for rule '%v' and subject %v - %v", rule.rule.Name, sub, err)
				errs = append(errs, err)
				continue
			}

			recommendation, ok := recommendationOutput.Value().(string)
			if !ok {
				klog.Fatalf("Failed to evaluate rule '%v' - %v", rule.rule.Name, err)
			}

			info := AnalysisFinding{
				Severity:       rule.rule.Severity,
				Message:        rule.rule.Description,
				Recommendation: recommendation,
				RuleName:       rule.rule.Name,
				RuleUuid:       rule.rule.Uuid,
				References:     rule.rule.References,
			}

			s := v1.Subject{}
			if kind, exist := sub["kind"]; exist {
				s.Kind = kind.(string)
			}
			if apiGroup, exist := sub["apiGroup"]; exist {
				s.APIGroup = apiGroup.(string)
			}
			if name, exist := sub["name"]; exist {
				s.Name = name.(string)
			}
			if namespace, exist := sub["namespace"]; exist {
				s.Namespace = namespace.(string)
			}

			finding := AnalysisReportFinding{
				Subject: &s,
				Finding: info,
			}
			report.Findings = append(report.Findings, finding)
		}

	}

	report.Stats = analysisStats

	return &report, nil
}
