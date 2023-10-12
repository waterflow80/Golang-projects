package main

import (
	"fmt"
	"math"
	"sync"
	"time"
)

func bubbleSort(arr []int16, channel chan any, wg *sync.WaitGroup){
	fmt.Printf("============== SORTING %v ==============\n", arr)
	var arr_copy []int16
	arr_copy = append(arr_copy, arr...)
	//fmt.Printf("Array_copy b4: %v\n", arr_copy)
	for i := range arr_copy {
			for j:= 0; j<len(arr_copy)-i-1; j++ {
				if arr_copy[j] > arr_copy[j+1] {
					
					temp := arr_copy[j]
					arr_copy[j] = arr_copy[j+1]
					arr_copy[j+1] = temp
				}
			}
	}
	//fmt.Printf("Array_copy after: %v\n", arr_copy)
	channel <- arr_copy
	wg.Done()
}

/**
Splits the given array into num sub-arrays*/
func split_array(num int, original_arr []int16) [][]int16{
	fmt.Printf("================ SPLITTING ================\n")
	fmt.Printf("original array: %v\n", original_arr)
	
	var local_original_arr = original_arr
	len_arr := len(original_arr)
	result_arr := make([][]int16, 0, 3)
	// if len(array) % num == 0 (eg: num=3, array=[1,2,3,4,5,6])
	// the subarrays will be in equal size
	if len_arr % num == 0 {
		num_elem := len_arr/num // the number of elements in each sub array
		for i:=0; i<num; i++ {
			result_arr = append(result_arr, local_original_arr[i*num_elem:num_elem+(i*num_elem)])
		}
	} else { // The array cannot be divided into equal sizes
		num_common_elts := len_arr / num
		num_remaining_elts := len_arr % num

		// distributing the equal elements
		last_indx :=0 // keeps track of the last indx distributed + 1 (the index to process next)
		fmt.Printf("---- Common elements ----\n")
		for i:=0; i<num_common_elts; i++ {
			result_arr = append(result_arr, local_original_arr[i*num_common_elts:num_common_elts+(i*num_common_elts)])
			last_indx = last_indx + num_common_elts   
			fmt.Printf("result array: %v\n", result_arr)
		}
		
		// Now distributing the rest of the elements (note: some sub-arrays will not receive any element here)
		fmt.Printf("---- Remaining elements ----\n")
		fmt.Printf("result array: %v\n", result_arr)
		fmt.Printf("Last index: %v\n", last_indx)
		for i:=0;i<num_remaining_elts;i++ {
			//result_arr[i] = append(result_arr[i], local_original_arr[last_indx])
			var temp []int16
			_ = copy(temp, result_arr[i])
			temp = append(temp, local_original_arr[last_indx])
			var result_arr_temp [][]int16
			_ = copy(result_arr_temp, result_arr)
			result_arr_temp[i] = temp
			fmt.Printf("temp=%v\n", temp)
			last_indx++
			fmt.Printf("result array: %v\n", result_arr)
			fmt.Printf("Last index: %v\n", last_indx)
		}

	}
	return result_arr
}

/**
Merging the given num_sub_arrays arrays into one sorted array with size array_size (the original array)
*/
func merge(sub_arrays [][]int16, array_size int) []int16 {
	tracker := make([]int, int(len((sub_arrays)))) // This will contain the indexes of the next test in each array. Eg: [2, 1, 3]
	final_array := make([]int16, int(array_size))
	n_sub_arrays := len(sub_arrays)

	// Initialize indexes with zeros
	for i:=0;i<n_sub_arrays;i++ {
		tracker[i] = 0
	}

	// Initialize the final array with +INF
	for i:=0; i<array_size;i++ {
		final_array[i] = math.MaxInt16
	}

	//Merge
	for i:=0; i<array_size; i++ {
		min_in := 0 // The index of the array from which the current min value has been taken 
		for j:=0; j<n_sub_arrays;j++ {
			if (tracker[j] < len(sub_arrays[j])) { // The case of the index out of bound
				if sub_arrays[j][tracker[j]] < final_array[i] {
					final_array[i] = sub_arrays[j][tracker[j]]
					min_in = j
				}
			}
		}
		tracker[min_in]++
	}

	return final_array

}

func main() {
	start := time.Now()
	// Array
	var arr = []int16 {10,50,12,5,30,6,15,88,36,44, 1, 7}

	// Distributed
	var num_walkers int = 3
	wlakers_channel := make(chan any, num_walkers)
	wg := &sync.WaitGroup{}
	wg.Add(num_walkers)

	// The list of sub arrays to be distributed	
	var sub_arrays_list [][]int16 = split_array(num_walkers, arr)

	fmt.Printf("Sub arrays list= %v\n", sub_arrays_list)
	for i:=0; i<num_walkers;i++ {
		go bubbleSort(sub_arrays_list[i], wlakers_channel, wg)
	}

	wg.Wait()
	close(wlakers_channel)

	var sorted_sub_arrays [][]int16
	for item := range wlakers_channel {
		sorted_sub_arrays = append(sorted_sub_arrays, item.([]int16))
	}	
	fmt.Printf("Sorted sub arrays %v\n", sorted_sub_arrays)
	//var sorted_array []int16 // The final sorted array
	sorted_array := merge(sorted_sub_arrays, len(arr)) // The final sorted array
	
	fmt.Printf("The final sorted array: %v\n",sorted_array)
	fmt.Printf("Total time of exeution: %v\n", time.Since(start))
}