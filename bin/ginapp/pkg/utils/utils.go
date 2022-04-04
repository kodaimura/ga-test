package utils

import (
	"strconv"
)

func AtoiSlice(sl []string) ([]int, error) {
	isl := make([]int, len(sl))
	for i, v := range sl {
		x , err := strconv.Atoi(v)
		isl[i] = x
		if err != nil {
			return []int{}, err
		}
	}
	return isl, nil
}


func ItoaSlice(sl []int) []string {
	asl := make([]string, len(sl))
	for i, v := range sl {
		asl[i] = strconv.Itoa(v)
	}

	return asl
} 


func Combinations(sl []string, n int) [][]string{
    combs := [][]string{}

    if len(sl) <= n{
        return [][]string{sl} 
    } 

    if n == 1 {
        for _, x := range sl {
            combs = append(combs, []string{x})
        }
        return combs
    }

    for _, c := range Combinations(sl[1:], n - 1) {
        combs = append(combs, append(c, sl[0]))
    }    
    return append(Combinations(sl[1:], n), combs...)
}