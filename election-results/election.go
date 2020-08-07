package main

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"

	pretty "github.com/inancgumus/prettyslice"
)

//Party contains all metadata related to votes per party
type Party struct {
	Name            string
	VotesAmount     int
	VotesPercentage float64
}

//ConstituencyElection contains constituency name and metadata for all parties and results of the election
type ConstituencyElection struct {
	Name         string
	Parties      []Party
	VotersAmount int
	ElectedParty Party
}

var dataSet = `Cardiff West, 11014, C, 17803, L, 4923, UKIP, 2069, LD
Islington South & Finsbury, 22547, L, 9389, C, 4829, LD, 3375, UKIP, 3371, G, 309, Ind`

var partyMap = map[string]string{
	"C":    "Conservative Party",
	"L":    "Labour Party",
	"UKIP": "UKIP",
	"LD":   "Liberal Democrats",
	"G":    "Green Party",
	"Ind":  "Independent",
	"SNP":  "SNP",
}

func getConstituencyFromString(data string) ([]ConstituencyElection, error) {
	var electionData []ConstituencyElection

	var dataSet [][]string
	scanner := bufio.NewScanner(strings.NewReader(data))
	for scanner.Scan() {
		dataSet = append(dataSet, strings.Split(scanner.Text(), ","))
	}

	for _, line := range dataSet {
		constituency := ConstituencyElection{
			Name: line[0],
		}

		var parties []Party
		var votersSum int
		for i := 1; i < len(line); i += 2 {
			votes, err := strconv.Atoi(strings.TrimSpace(line[i]))
			if err != nil {
				return electionData, err
			}
			votersSum += votes

			parties = append(parties, Party{
				Name:        partyMap[strings.TrimSpace(line[i+1])],
				VotesAmount: votes,
			})
		}
		constituency.Parties = parties
		constituency.VotersAmount = votersSum
		electionData = append(electionData, constituency)
	}

	return electionData, nil
}

func merge(l, r []Party) []Party {
	res := make([]Party, 0, len(l)+len(r))

	for len(l) > 0 || len(r) > 0 {
		if len(l) == 0 {
			return append(res, r...)
		}
		if len(r) == 0 {
			return append(res, l...)
		}
		if l[0].VotesPercentage >= r[0].VotesPercentage {
			res = append(res, l[0])
			l = l[1:]
		} else {
			res = append(res, r[0])
			r = r[1:]
		}
	}

	return res
}

func mergeSort(data []Party) []Party {
	if len(data) <= 1 {
		return data
	}

	initialDivider := len(data) / 2
	var leftChunk []Party
	var rightChunk []Party

	leftChunk = mergeSort(data[:initialDivider])
	rightChunk = mergeSort(data[initialDivider:])

	return merge(leftChunk, rightChunk)
}

func calculateVotes(data ConstituencyElection) ConstituencyElection {
	fmt.Printf("Calculating percentage of votes in constituency %v...\n", data.Name)
	for pos, party := range data.Parties {
		data.Parties[pos].VotesPercentage = (float64(party.VotesAmount) / (float64(data.VotersAmount)) * float64(100))
	}
	fmt.Printf("Calculating winner in constituency %v...\n", data.Name)
	data.Parties = mergeSort(data.Parties)
	data.ElectedParty = data.Parties[0]

	return data
}

func main() {
	electionMeta, err := getConstituencyFromString(dataSet)
	if err != nil {
		panic(err)
	}

	var electionResults []ConstituencyElection
	for _, v := range electionMeta {
		electionResults = append(electionResults, calculateVotes(v))
	}

	pretty.MaxPerLine = 1
	for _, v := range electionResults {
		fmt.Printf("Region: %v\n", v.Name)
		pretty.Show("Parties", v.Parties)
		fmt.Printf("Winner: %v , voted %v of the citizens\n\n", v.ElectedParty.Name, v.ElectedParty.VotesPercentage)
	}
}
