package gofuzzymind

import "math"

// FuzzySet represents a fuzzy set with a membership function.
type FuzzySet struct {
	Name               string
	MembershipFunction func(float64) float64
}

// NewFuzzySet creates a new FuzzySet.
func NewFuzzySet(name string, membershipFunc func(float64) float64) *FuzzySet {
	return &FuzzySet{
		Name:               name,
		MembershipFunction: membershipFunc,
	}
}

// MembershipDegree returns the membership degree of x.
func (fs *FuzzySet) MembershipDegree(x float64) float64 {
	return fs.MembershipFunction(x)
}

// Union returns a new FuzzySet representing the union of this set and another set.
func (fs *FuzzySet) Union(other *FuzzySet) *FuzzySet {
	return NewFuzzySet(
		"Union("+fs.Name+", "+other.Name+")",
		func(x float64) float64 {
			return math.Max(fs.MembershipFunction(x), other.MembershipFunction(x))
		},
	)
}

// Intersection returns a new FuzzySet representing the intersection of this set and another set.
func (fs *FuzzySet) Intersection(other *FuzzySet) *FuzzySet {
	return NewFuzzySet(
		"Intersection("+fs.Name+", "+other.Name+")",
		func(x float64) float64 {
			return math.Min(fs.MembershipFunction(x), other.MembershipFunction(x))
		},
	)
}

// Complement returns a new FuzzySet representing the complement of this set.
func (fs *FuzzySet) Complement() *FuzzySet {
	return NewFuzzySet(
		"Complement("+fs.Name+")",
		func(x float64) float64 {
			return 1 - fs.MembershipFunction(x)
		},
	)
}

// Normalize returns a new FuzzySet with a normalized membership function.
func (fs *FuzzySet) Normalize() *FuzzySet {
	return NewFuzzySet(
		"Normalized("+fs.Name+")",
		func(x float64) float64 {
			return fs.MembershipFunction(x) / math.Max(1, fs.MembershipFunction(x))
		},
	)
}

// Centroid calculates the centroid of the fuzzy set for defuzzification.
func (fs *FuzzySet) Centroid(min, max, step float64) float64 {
	numerator, denominator := 0.0, 0.0
	for x := min; x <= max; x += step {
		mu := fs.MembershipFunction(x)
		numerator += x * mu
		denominator += mu
	}
	if denominator == 0 {
		return 0
	}
	return numerator / denominator
}

// FuzzyRule represents a fuzzy rule consisting of a condition and a consequence.
type FuzzyRule struct {
	Condition   func(map[string]float64) bool
	Consequence interface{}
	Weight      float64
}

// NewFuzzyRule creates a new FuzzyRule.
func NewFuzzyRule(condition func(map[string]float64) bool, consequence interface{}, weight float64) *FuzzyRule {
	return &FuzzyRule{
		Condition:   condition,
		Consequence: consequence,
		Weight:      weight,
	}
}

// Evaluate evaluates the rule against the given inputs and returns the result and weight if the condition is satisfied.
func (fr *FuzzyRule) Evaluate(inputs map[string]float64) map[string]interface{} {
	if fr.Condition(inputs) {
		var result interface{}
		if consequence, ok := fr.Consequence.(*FuzzySet); ok {
			result = consequence
		} else if consequenceFunc, ok := fr.Consequence.(func(map[string]float64) string); ok {
			result = consequenceFunc(inputs)
		}
		return map[string]interface{}{"result": result, "weight": fr.Weight}
	}
	return nil
}

// InferenceEngine uses a set of fuzzy rules to perform inference and defuzzification.
type InferenceEngine struct {
	Rules []*FuzzyRule
}

// NewInferenceEngine creates a new InferenceEngine.
func NewInferenceEngine(rules []*FuzzyRule) *InferenceEngine {
	return &InferenceEngine{Rules: rules}
}

// Infer performs inference based on the input values and returns the aggregated result.
func (ie *InferenceEngine) Infer(inputs map[string]float64) string {
	var results []map[string]interface{}
	for _, rule := range ie.Rules {
		if evaluation := rule.Evaluate(inputs); evaluation != nil {
			results = append(results, evaluation)
		}
	}
	return ie.AggregateResults(results)
}

func (ie *InferenceEngine) AggregateResults(results []map[string]interface{}) string {
	if len(results) == 0 {
		return "Low Priority"
	}

	totalWeight, weightedSum := 0.0, 0.0
	for _, result := range results {
		weight := result["weight"].(float64)
		weightedSum += ie.PriorityMapping(result["result"]) * weight
		totalWeight += weight
	}

	if totalWeight > 0 {
		return ie.ReversePriorityMapping(weightedSum / totalWeight)
	}
	return "Low Priority"
}

func (ie *InferenceEngine) PriorityMapping(priority interface{}) float64 {
	if p, ok := priority.(string); ok {
		switch p {
		case "Urgent":
			return 3
		case "High Priority":
			return 2
		case "Medium Priority":
			return 1
		}
	}
	return 0
}

func (ie *InferenceEngine) ReversePriorityMapping(score float64) string {
	switch {
	case score >= 2.5:
		return "Urgent"
	case score >= 1.5:
		return "High Priority"
	case score >= 0.5:
		return "Medium Priority"
	default:
		return "Low Priority"
	}
}

func (ie *InferenceEngine) GetFuzzySetConsequences() []*FuzzySet {
	var fuzzySets []*FuzzySet
	for _, rule := range ie.Rules {
		if fs, ok := rule.Consequence.(*FuzzySet); ok {
			fuzzySets = append(fuzzySets, fs)
		}
	}
	return fuzzySets
}

func (ie *InferenceEngine) DefuzzifyCentroid(min, max, step float64) float64 {
	numerator, denominator := 0.0, 0.0
	fuzzySets := ie.GetFuzzySetConsequences()
	for x := min; x <= max; x += step {
		mu := 0.0
		for _, fuzzySet := range fuzzySets {
			mu = math.Max(mu, fuzzySet.MembershipDegree(x))
		}
		numerator += x * mu
		denominator += mu
	}
	if denominator == 0 {
		return 0
	}
	return numerator / denominator
}

func (ie *InferenceEngine) DefuzzifyMOM(min, max, step float64) float64 {
	maxMu, sumX, count := 0.0, 0.0, 0.0
	fuzzySets := ie.GetFuzzySetConsequences()
	for x := min; x <= max; x += step {
		mu := 0.0
		for _, fuzzySet := range fuzzySets {
			mu = math.Max(mu, fuzzySet.MembershipDegree(x))
		}
		if mu > maxMu {
			maxMu = mu
			sumX = x
			count = 1
		} else if mu == maxMu {
			sumX += x
			count++
		}
	}
	if count == 0 {
		return 0
	}
	return sumX / count
}

func (ie *InferenceEngine) DefuzzifyBisector(min, max, step float64) float64 {
	totalArea, leftArea := 0.0, 0.0
	bisector := min
	fuzzySets := ie.GetFuzzySetConsequences()
	for x := min; x <= max; x += step {
		mu := 0.0
		for _, fuzzySet := range fuzzySets {
			mu = math.Max(mu, fuzzySet.MembershipDegree(x))
		}
		totalArea += mu * step
	}
	for x := min; x <= max; x += step {
		mu := 0.0
		for _, fuzzySet := range fuzzySets {
			mu = math.Max(mu, fuzzySet.MembershipDegree(x))
		}
		leftArea += mu * step
		if leftArea >= totalArea/2 {
			bisector = x
			break
		}
	}
	return bisector
}

