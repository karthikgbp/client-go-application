package k8sclient

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func CreateInformer() {

	nsList := readNSList_file()

	fmt.Println(nsList)

	// nsList := []string{"client-namespace", "client-namespace-podb", "client-namespace-podc"}

	createk8sInformer(nsList)
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

	// bs, err := ioutil.ReadAll(of)
	// if err != nil {
	// 	log.Fatal("NS List Read File Error : ", err)
	// }

	// rf := string(bs)

	// fmt.Println(rf)
}
