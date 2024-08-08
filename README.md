### goFuzzyMind
goFuzzyMind is a Go library for implementing fuzzy logic systems. It provides a framework for defining fuzzy sets, creating fuzzy rules, and performing inference using an inference engine. The library also includes methods for defuzzification using various techniques. This library can be used to build fuzzy logic systems for applications such as risk assessment, decision-making, or control systems.

#### FuzzySet
Represents a fuzzy set with a membership function and provides methods for operations on fuzzy sets.

```
```go
type FuzzySet struct {
    Name               string
    MembershipFunction func(float64) float64
}
// NewFuzzySet creates a new FuzzySet.
func NewFuzzySet(name string, membershipFunc func(float64) float64) *FuzzySet
```

##### Methods
- MembershipDegree(x float64) float64: Returns the membership degree of x.
- **Union(other FuzzySet) FuzzySet: Returns a new fuzzy set representing the union of this set and other.
- **Intersection(other FuzzySet) FuzzySet: Returns a new fuzzy set representing the intersection of this set and other.
- *Complement() FuzzySet: Returns a new fuzzy set representing the complement of this set.
- *Normalize() FuzzySet: Returns a new fuzzy set with a normalized membership function.
- Centroid(min, max, step float64) float64: Calculates the centroid for defuzzification.
#### FuzzyRule
Represents a fuzzy rule consisting of a condition and a consequence.

```go
type FuzzyRule struct {
    Condition   func(map[string]float64) bool
    Consequence interface{}
    Weight      float64
}

// NewFuzzyRule creates a new FuzzyRule.
func NewFuzzyRule(condition func(map[string]float64) bool, consequence interface{}, weight float64) *FuzzyRule
```
##### Methods
Evaluate(inputs map[string]float64) map[string]interface{}: Evaluates the rule against the given inputs and returns the result and weight if the condition is satisfied.

#### InferenceEngine
Uses a set of fuzzy rules to perform inference and defuzzification.

```go
type InferenceEngine struct {
    Rules []*FuzzyRule
}

// NewInferenceEngine creates a new InferenceEngine.
func NewInferenceEngine(rules []*FuzzyRule) *InferenceEngine
```
##### Methods
- Infer(inputs map[string]float64) string: Performs inference based on the input values and returns the aggregated result.
- AggregateResults(results []map[string]interface{}) string: Aggregates results from the fuzzy rules.
- DefuzzifyCentroid(min, max, step float64) float64: Performs defuzzification using the centroid method.
- DefuzzifyMOM(min, max, step float64) float64: Performs defuzzification using the Mean of Maxima (MOM) method.
- DefuzzifyBisector(min, max, step float64) float64: Performs defuzzification using the Bisector method.
- *GetFuzzySetConsequences() []FuzzySet: Returns a list of fuzzy sets as consequences of the rules.

##### Example Usage
Here's a sample usage of the goFuzzyMind library:

```go
package main

import (
    "fmt"
    "github.com/FadeDreams/gofuzzymind"
)

func main() {
    // Define fuzzy sets for urgency and complexity
    urgencySet := gofuzzymind.NewFuzzySet("Urgency", func(urgency float64) float64 {
        switch {
        case urgency < 3:
            return 0
        case urgency < 7:
            return (urgency - 3) / 4
        default:
            return 1
        }
    })

    complexitySet := gofuzzymind.NewFuzzySet("Complexity", func(complexity float64) float64 {
        switch {
        case complexity < 2:
            return 0
        case complexity < 5:
            return (complexity - 2) / 3
        default:
            return 1
        }
    })

    // Define fuzzy rules
    rules := []*gofuzzymind.FuzzyRule{
        gofuzzymind.NewFuzzyRule(
            func(inputs map[string]float64) bool {
                return urgencySet.MembershipDegree(inputs["urgency"]) > 0.7 &&
                    complexitySet.MembershipDegree(inputs["complexity"]) > 0.7
            },
            gofuzzymind.NewFuzzySet("Urgent", func(x float64) float64 {
                if x >= 7 {
                    return 1
                }
                return x / 7
            }),
            1,
        ),
        gofuzzymind.NewFuzzyRule(
            func(inputs map[string]float64) bool {
                return urgencySet.MembershipDegree(inputs["urgency"]) > 0.5
            },
            func(_ map[string]float64) string { return "High Priority" },
            1,
        ),
        gofuzzymind.NewFuzzyRule(
            func(inputs map[string]float64) bool {
                return complexitySet.MembershipDegree(inputs["complexity"]) > 0.5
            },
            func(_ map[string]float64) string { return "Medium Priority" },
            1,
        ),
        gofuzzymind.NewFuzzyRule(
            func(inputs map[string]float64) bool {
                return urgencySet.MembershipDegree(inputs["urgency"]) <= 0.5 &&
                    complexitySet.MembershipDegree(inputs["complexity"]) <= 0.5
            },
            func(_ map[string]float64) string { return "Low Priority" },
            1,
        ),
    }

    engine := gofuzzymind.NewInferenceEngine(rules)

    // Example ticket
    ticket := map[string]float64{"urgency": 8, "complexity": 6}
    priority := engine.Infer(ticket)
    fmt.Printf("Ticket Priority: %s\n", priority)

    // Defuzzification examples
    centroid := urgencySet.Centroid(0, 10, 0.01)
    fmt.Printf("Centroid defuzzification: %.2f\n", centroid)

    defuzzifiedCentroid := engine.DefuzzifyCentroid(0, 10, 0.01)
    fmt.Printf("Defuzzified Centroid: %.2f\n", defuzzifiedCentroid)

    defuzzifiedMOM := engine.DefuzzifyMOM(0, 10, 0.01)
    fmt.Printf("Defuzzified MOM: %.2f\n", defuzzifiedMOM)

    defuzzifiedBisector := engine.DefuzzifyBisector(0, 10, 0.01)
    fmt.Printf("Defuzzified Bisector: %.2f\n", defuzzifiedBisector)
}
```
