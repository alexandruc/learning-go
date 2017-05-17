package main

import (
	"container/list"
	"flag"
	"fmt"
)

func main() {
	options := list.New()

	option1Ptr := flag.String("option1", "null", "the first option")
	option2Ptr := flag.Bool("option2", false, "the second option")
	option3Ptr := flag.Int("option3", 0, "the 3rd option")
	options.PushBack(option1Ptr)
	options.PushBack(option2Ptr)
	options.PushBack(option3Ptr)

	flag.Parse()

	for e := options.Front(); e != nil; e = e.Next() {
		strPtr, bIsStr := e.Value.(*string)
		bPtr, bIsBool := e.Value.(*bool)
		iPtr, bIsInt := e.Value.(*int)
		switch {
		case bIsStr:
			fmt.Println("option str: ", *strPtr)
		case bIsBool:
			fmt.Println("option bool: ", *bPtr)
		case bIsInt:
			fmt.Println("option int: ", *iPtr)
		}
	}

	fmt.Println("option 1: ", *option1Ptr)
	fmt.Println("option 2: ", *option2Ptr)
	fmt.Println("option 3: ", *option3Ptr)
}
