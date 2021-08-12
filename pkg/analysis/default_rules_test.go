package analysis

import (
	"testing"
	"time"

	"github.com/alcideio/rbac-tool/pkg/rbac"
	"github.com/kr/pretty"
	v1 "k8s.io/api/rbac/v1"
	"k8s.io/klog"
)

func Test__DefultRules(t *testing.T) {
	for i, rule := range defaultRules {
		if r, err := newAnalysisRule(&defaultRules[i]); err != nil || r == nil {
			t.Fatalf("Rule '%v' - %v - failed to initialize\n%v\n", i, rule.Name, err)
			t.Fail()
		}
	}
}

func Test__Analyzer(t *testing.T) {
	defer klog.Flush()

	config := DefaultAnalysisConfig()

	analyzer := CreateAnalyzer(
		config,
		[]rbac.SubjectPolicyList{
			{Subject: v1.Subject{
				Kind:      "ServiceAccount",
				APIGroup:  "",
				Name:      "test-sa",
				Namespace: "test",
			}, AllowedTo: []rbac.NamespacedPolicyRule{
				{Namespace: "test", Verb: "get", APIGroup: "*", Resource: "*", ResourceNames: nil, NonResourceURLs: nil},
			}},
		},
	)

	if analyzer == nil {
		t.Fail()
	}

	report, err := analyzer.Analyze()
	if err != nil {
		t.Fatalf("Analysis failed - %v", err)
		t.Fail()
	}

	if len(report.Findings) == 0 {
		t.Fatalf("Expecting findings")
		t.Fail()
	}

	t.Logf("%v", pretty.Sprint(report))
}

func Test__GlobalExclusion(t *testing.T) {
	defer klog.Flush()

	config := DefaultAnalysisConfig()

	analyzer := CreateAnalyzer(
		config,
		[]rbac.SubjectPolicyList{
			{Subject: v1.Subject{
				Kind:      "ServiceAccount",
				APIGroup:  "",
				Name:      "test-sa",
				Namespace: "kube-system",
			}, AllowedTo: []rbac.NamespacedPolicyRule{
				{Namespace: "test", Verb: "get", APIGroup: "*", Resource: "*", ResourceNames: nil, NonResourceURLs: nil},
			}},
		},
	)

	if analyzer == nil {
		t.Fail()
	}

	report, err := analyzer.Analyze()
	if err != nil {
		t.Fatalf("Analysis failed - %v", err)
		t.Fail()
	}

	if len(report.Findings) != 0 {
		t.Fatalf("Expecting no findings")
		t.Fail()
	}

	t.Logf("%v", pretty.Sprint(report))
}

func Test__RuleExclusion(t *testing.T) {
	defer klog.Flush()

	config := DefaultAnalysisConfig()

	config.Rules[0].Exclusions = []Exclusion{
		{
			Disabled:     false,
			Comment:      "Exclude test from analysis",
			AddedBy:      "tester",
			LastModified: time.Now().Format(time.RFC3339),
			SnoozeUntil:  0,
			Expression:   `subject.namespace == "test"`,
		},
	}

	analyzer := CreateAnalyzer(
		config,
		[]rbac.SubjectPolicyList{
			{Subject: v1.Subject{
				Kind:      "ServiceAccount",
				APIGroup:  "",
				Name:      "test-sa",
				Namespace: "test",
			}, AllowedTo: []rbac.NamespacedPolicyRule{
				{Namespace: "test", Verb: "get", APIGroup: "*", Resource: "*", ResourceNames: nil, NonResourceURLs: nil},
			}},
		},
	)

	if analyzer == nil {
		t.Fail()
	}

	report, err := analyzer.Analyze()
	if err != nil {
		t.Fatalf("Analysis failed - %v", err)
		t.Fail()
	}

	if len(report.Findings) != 0 {
		t.Fatalf("Expecting no findings")
		t.Fail()
	}

	t.Logf("%v", pretty.Sprint(report))
}