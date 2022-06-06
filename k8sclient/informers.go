package k8sclient

import (
	"bufio"
	"fmt"
	"log"
	"os"
	// "k8s.io/utils/strings/slices"
)

type NSLabels struct {
	Namespace string
	Labels    map[string]string
}

var watchlist []NSLabels

var nsLabel map[string]string

func CreateInformer() {

	// nsList := readNSList_file()

	systemDefNS := []string{"default", "kube-", "kubernetes-", "ingress-"}
	userDefNS := []string{"client-namespace", "my-namespace-podb", "client-namespace-podc"}

	// Get Namespaces
	nsList := GetNamespaces()

	// Filtered Namespaces without System Defined NS
	for i := 0; i < len(nsList); i++ {
		for _, sD := range systemDefNS {
			// Break from Loop if
			if len(sD) <= len(nsList[i]) && sD == nsList[i][:len(sD)] {

				copy(nsList[i:], nsList[i+1:]) // Shift a[i+1:] left one index.
				nsList[len(nsList)-1] = ""     // Erase last element (write zero value).
				nsList = nsList[:len(nsList)-1]
				i--
				break
			}
		}
	}
	// fmt.Println("List All Namespaces without System Defined : ", nsList)

	// Filtered Namespaces without User Defined
	for i := 0; i < len(nsList); i++ {
		for _, uD := range userDefNS {
			// Break from Loop if
			if len(uD) <= len(nsList[i]) && uD == nsList[i][:len(uD)] {

				copy(nsList[i:], nsList[i+1:]) // Shift a[i+1:] left one index.
				nsList[len(nsList)-1] = ""     // Erase last element (write zero value).
				nsList = nsList[:len(nsList)-1]
				i--
				break
			}
		}
	}
	fmt.Println("List All Namespaces without User Defined : ", nsList)

	// Delete Non Listed Namespaces
	if len(nsList) > 0 {
		for _, ns := range nsList {

			isDeleted := DeleteNamespace(ns)

			if isDeleted {
				log.Println("Sucessfully Deleted Namespace ")
			}
		}
	} else {
		log.Println("No additional namespaces to be deleted")
	}

	// Check Labels on Namespaces
	if len(userDefNS) > 0 {

		for _, val := range userDefNS {
			nsLabelList := ListUserDefNSLabels(val)
			watchlist = append(watchlist, NSLabels{val, nsLabelList})
		}
	}

	fmt.Println(watchlist)

}

func readNSList_file() []string {

	of, err := os.Open("./nsList.txt")
	if err != nil {
		log.Fatal("NS List Open File Error : ", err)
	}

	defer of.Close()

	var nsList []string

	scanner := bufio.NewScanner(of)
	for scanner.Scan() {
		nsList = append(nsList, scanner.Text())
	}
	return nsList

}
