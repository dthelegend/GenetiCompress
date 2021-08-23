package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	var popSize = 100;
	var initialPopulation Population = make(Population, popSize);

	for i := 0; i < popSize; i++ {
		initialPopulation[i] = IndividualCompressor{Chromosome{1, 0}}
	}

	for i := 0; i < 10; i++ {
		initialPopulation.NextGeneration()
		go fmt.Printf("%d: %v\n", i, initialPopulation.MaxIndividual().EvaluateFitness())
	}

	var maxIndividual IndividualCompressor = initialPopulation.MaxIndividual();

	for i := 0; i < len(maxIndividual.Chromosome); i++ {
		fmt.Printf("Chromosome %d: %s\n", i, strconv.FormatUint(uint64(maxIndividual.Chromosome[i]), 2));
	}

	var outstring []byte = make([]byte, 8);
	maxIndividual.NewReader().Read(outstring);
	fmt.Printf("outstring:")
	for i := 0; i < len(outstring); i++ {
		fmt.Printf(" %s", strconv.FormatInt(int64(outstring[i]), 2))
	}
	fmt.Printf("\nstring(outstring): %v\n", string(outstring))
}