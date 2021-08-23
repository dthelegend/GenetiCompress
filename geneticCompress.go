package main

import (
	"math"
	"math/bits"
	"math/rand"
	"os"
)

type FitnessScore int;

type Population []IndividualCompressor

func (population Population) TournamentSelection(selectionSize int) IndividualCompressor {
	var maxIndividual IndividualCompressor;

	for i := 0; i < selectionSize; i++ {
		var entrant IndividualCompressor = population[rand.Intn(len(population))];
		if entrant.EvaluateFitness() >= maxIndividual.EvaluateFitness() {
			maxIndividual = entrant;
		}
	}

	return maxIndividual;
}

func (population Population) NextGeneration() {
	var populationSize int = len(population);
	var parentPoolSize int = int(0.5 + math.Sqrt(float64(populationSize) + 0.25));

	var parentPool []IndividualCompressor = make([]IndividualCompressor, parentPoolSize)
	for i := 0; i < parentPoolSize; i++ {
		parentPool[i] = population.TournamentSelection(10);
	}

	var k int = 0;
	for i := 0; i < parentPoolSize; i++ {
		for j := i + 1; j < parentPoolSize; j++ {
			population[k], population[k + 1] = UniformCrossover(parentPool[i], parentPool[j])
			if(rand.Intn(5) == 0) {
				population[k].mutate();
			}
			k+=2;
		}
	}

	for ; k < populationSize; k++ {
		population[k] = IndividualCompressor{Chromosome{}};
		for i := 0; i < len(population[k].Chromosome); i++ {
			population[k].Chromosome[i] = Gene(rand.Int())
		} 
	}
}

func (population Population) MaxIndividual() IndividualCompressor {
	var maxIndividual IndividualCompressor;

	for i := 0; i < len(population); i++ {
		var entrant IndividualCompressor = population[i];
		if entrant.EvaluateFitness() >= maxIndividual.EvaluateFitness() {
			maxIndividual = entrant;
		}
	}

	return maxIndividual
}

type Gene int64

type Chromosome [2]Gene;

type IndividualCompressor struct {
	Chromosome Chromosome;
}

type CompressorReader struct {
	index int;
	individualCompressor IndividualCompressor;
}

func (I IndividualCompressor) NewReader() CompressorReader {
	return CompressorReader{index: 0, individualCompressor: I}
}

func (I CompressorReader) Read(b []byte) (int, error) {
	// sin(Ax + B)
	var i int;
	for i = 0; i < len(b); i++ {
		b[i] = 0x0;
		for j := 0; j < 8; j++ {
			var p byte = 0x1;
			I.index += i;
			if(math.Round(math.Sin(float64((I.index * 8 + j) * int(I.individualCompressor.Chromosome[0])) * math.Pi * 0.5) + float64(I.individualCompressor.Chromosome[1])) == 0) {
				p = 0x0;
			}
			b[i] = b[i] | (p << j);
		}
	}
	return i, nil;
}

func (I IndividualCompressor) EvaluateFitness() FitnessScore {
	f, err := os.Open("./test.txt");
    if err != nil { 
		panic(err);
	}
	defer f.Close();

	p := I.NewReader();

	var original []byte = make([]byte, 1);
	var prediction []byte = make([]byte, 1);
	var errorCount int = 0;
	for i := 0; ; i++ {
		_, err := f.Read(original[:]);
		if err != nil { 
			break;
		}
		p.Read(prediction[:]);
		if err != nil { 
			break;
		}
		errorCount += bits.OnesCount8(uint8(prediction[0] & original[0]));
	}

	return FitnessScore(errorCount);
}

func (I IndividualCompressor) mutate() {
	I.Chromosome[rand.Intn(len(I.Chromosome))] = Gene(rand.Int63());
}

func UniformCrossover(parent1 IndividualCompressor, parent2 IndividualCompressor) (IndividualCompressor, IndividualCompressor) {
	var child1 IndividualCompressor = IndividualCompressor{};
	var child2 IndividualCompressor = IndividualCompressor{};

	for i := 0; i < len(parent1.Chromosome); i++ {
		var mask uint64 = rand.Uint64();

		child1.Chromosome[i] = Gene((mask & uint64(parent1.Chromosome[i])) | (^mask & uint64(parent2.Chromosome[i])));
		child2.Chromosome[i] = Gene((mask & uint64(parent2.Chromosome[i])) | (^mask & uint64(parent1.Chromosome[i])));
	}
	

	return child1, child2;
}