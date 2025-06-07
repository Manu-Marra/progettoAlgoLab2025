package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// tipo Dizionario rappresenta l’intero Dizionario
type Dizionario struct {
	parole map[string]struct{}
	schemi map[string]struct{}
}

type dizionario *Dizionario

// crea() - crea un nuovo dizionario vuoto (azzera il dizionario)
func crea(d dizionario) {
	d.parole = make(map[string]struct{})
	d.schemi = make(map[string]struct{})
}

// che implementa l’operazione crea(), ovvero che crea un nuovo dizionario, lo inizializza e lo restituisce.
func newDizionario() dizionario {
	return &Dizionario{
		parole: make(map[string]struct{}),
		schemi: make(map[string]struct{}),
	}
}

func isValida(w string) bool {
	regex := regexp.MustCompile(`^[a-zA-Z]+$`)
	return regex.MatchString(w)
}

func esisteParola(d dizionario, w string) bool {
	_, esiste := d.parole[w]
	return esiste
}

func esisteSchema(d dizionario, w string) bool {
	_, esiste := d.schemi[w]
	return esiste
}

func inserisci(d dizionario, w string) {

	if contieneMaiuscola(w) {
		// è uno schema
		// _, esiste := d.schemi[w] // verifica se w è già presente in schemi
		// if !esiste {
		// 	d.schemi[w] = struct{}{}
		// }
		if !esisteSchema(d, w) {
			d.schemi[w] = struct{}{}
		}
	} else {
		// è una parola
		// _, esiste := d.parole[w] // verifica se w è già presente in parole
		// if !esiste {
		// 	d.parole[w] = struct{}{}
		// }
		if !esisteParola(d, w) {
			d.parole[w] = struct{}{}
		}
	}
}

func carica(d dizionario, file string) {
	f, err := os.Open(file)
	if err != nil {
		// file non esiste o non può essere aperto, non fare nulla
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	// scanner default split è ScanLines, cambiamo in ScanWords per token su spazi bianchi
	scanner.Split(bufio.ScanWords)

	for scanner.Scan() {
		w := scanner.Text()
		if isValida(w) {
			inserisci(d, w)
		} else {
			fmt.Printf("formato errato per la parola/schema -> %s <-\n", w)
		}
	}
	// Ignoro scanner.Err() per non bloccare l'esecuzione in caso di errori di lettura
}

func stampa_parole(d dizionario) {
	fmt.Println("[")
	for parola := range d.parole {
		fmt.Println(parola)
	}
	fmt.Println("]")

}

func stampa_schemi(d dizionario) {
	fmt.Println("[")
	for schema := range d.schemi {
		fmt.Println(schema)
	}
	fmt.Println("]")
}

func elimina(d dizionario, w string) {
	if contieneMaiuscola(w) {
		// È uno schema
		// _, esiste := d.schemi[w]
		// if esiste {
		// 	delete(d.schemi, w)
		// }
		if esisteSchema(d, w) {
			delete(d.schemi, w)
		}
	} else {
		// È una parola
		// _, esiste := d.parole[w]
		// if esiste {
		// 	delete(d.parole, w)
		// }

		if esisteParola(d, w) {
			delete(d.parole, w)			
		}
	}
}

func contieneMaiuscola(s string) bool {
	for _, c := range s {
		if c >= 'A' && c <= 'Z' {
			return true
		}
	}
	return false
}

// func (d dizionario) contiene(w string) bool {
//     _, esiste := d.parole[w]
//     return esiste
// }

func min(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}

func distanzaLevenshtein(s1, s2 string) int {
	m := len(s1)
	n := len(s2)

	if n == 0 {
		return m
	}
	if m == 0 {
		return n
	}

	prev := make([]int, n+1)
	curr := make([]int, n+1)

	for j := 0; j <= n; j++ {
		prev[j] = j
	}

	for i := 1; i <= m; i++ {
		curr[0] = i
		for j := 1; j <= n; j++ {
			costo := 0
			if s1[i-1] != s2[j-1] {
				costo = 1
			}

			curr[j] = min(curr[j-1]+1, prev[j]+1, prev[j-1]+costo)
		}
		// Scambio dei riferimenti tra prev e curr
		prev, curr = curr, prev
	}

	return prev[n]
}

func isMaiuscola(r rune) bool {
	return r >= 'A' && r <= 'Z'
}

func compatibile(parola, schema string) bool {
	if len(parola) != len(schema) {
		return false
	}

	mappa := make(map[rune]rune)

	for i, c := range schema {
		p := rune(parola[i])

		if isMaiuscola(c) {
			val, esiste := mappa[c]
			if esiste {
				if val != p {
					return false
				}
			} else {
				mappa[c] = p
			}
		} else {
			if c != p {
				return false
			}
		}
	}

	return true
}

func ricerca(d dizionario, schema string) {
	fmt.Printf("%s:[\n", schema)
	// fmt.Println(schema,":[")
	for parola := range d.parole {
		if compatibile(parola, schema) {
			fmt.Println(parola)
		}
	}
	fmt.Println("]")
}

func isSimile(x, y string) bool {
	return distanzaLevenshtein(x, y) == 1
}

func catena(d dizionario, x, y string) {

	// Se x e y sono uguali, la catena è triviale
	if x == y {
		fmt.Println("(")
		fmt.Println(x)
		fmt.Println(")")
		return
	}

	// Coda per la BFS
	queue := []string{x}

	// Mappa per tracciare i predecessori e ricostruire il percorso
	predecessore := make(map[string]string)

	// Insieme delle parole visitate
	visitato := make(map[string]bool)
	visitato[x] = true

	for len(queue) > 0 {
		parolaCorrente := queue[0]
		queue = queue[1:]

		// Scorro tutte le parole nel dizionario
		for parolaVicino := range d.parole {
			// Se già visitata, skippo/continuo
			if visitato[parolaVicino] {
				continue
			}
			// Se simile
			if isSimile(parolaCorrente, parolaVicino) {
				// Salvo predecessore
				predecessore[parolaVicino] = parolaCorrente
				// Se arrivo alla destinazione
				if parolaVicino == y {
					// Ricostruisco il percorso
					ricostruisciCatena(predecessore, x, y)
					return
				}
				// Altrimenti aggiungo alla coda e segno come visitata
				queue = append(queue, parolaVicino)
				visitato[parolaVicino] = true
			}
		}
	}

	// Se esco dal ciclo senza aver trovato y
	fmt.Println("non esiste")
}

// Funzione ausiliaria per ricostruire e stampare la catena dal predecessore
func ricostruisciCatena(predecessore map[string]string, inizio, fine string) {
	// Parto dalla fine e risalgo
	catena := []string{fine}
	for parola := fine; parola != inizio; parola = predecessore[parola] {
		catena = append([]string{predecessore[parola]}, catena...)
	}
	// Stampo la catena
	fmt.Println("(")
	for _, parola := range catena {
		fmt.Println(parola)
	}
	fmt.Println(")")
}



func esegui(dizionario dizionario, s string) {
	formatoErrato := "Formato errato per il comando"
	campi := strings.Fields(s)
	if len(campi) == 0 {
		return
	}

	switch campi[0] {
	case "c": // INIZIALIZZA "c", CARICA FILE "c nomeFile", CATENA "c x y"
		if len(campi) == 1 { // CREA
			crea(dizionario)	

		} else if len(campi) == 2 { // CARICA
			carica(dizionario, campi[1])

		} else if len(campi) == 3 { // CATENA
			// "c x y" → catena(x, y)
			x, y := campi[1], campi[2]
			catena(dizionario, x, y)

		} else {
			fmt.Println(formatoErrato, "c")
		}

	case "t": // TERMINA ESECUZIONE
		fmt.Println("Esecuzione terminata")
		os.Exit(0)

	case "p": // STAMPA PAROLE
		if len(campi) != 1 {
			fmt.Println(formatoErrato, "p")
			return
		}
		stampa_parole(dizionario)

	case "s": // STAMPA SCHEMI 
		if len(campi) != 1 {
			fmt.Println(formatoErrato, "s")
			return
		}
		stampa_schemi(dizionario)

	case "i": // INSERISCI PAROLA/SCHEMA
		if len(campi) != 2 {
			fmt.Println(formatoErrato, "i")
			return
		}

		if isValida(campi[1]) {
			inserisci(dizionario, campi[1])
		} else {
			fmt.Println("Parola/schema non valida")
		}

	case "e": // ELIMINA PAROLA/SCHEMA
		if len(campi) != 2 {
			fmt.Println(formatoErrato, "e")
			return
		}
		elimina(dizionario, campi[1])

	case "r": // STAMPA LO SCHEMA E LE PAROLE COMPATIBILI
		if len(campi) == 2 {
			schema := campi[1]
			if esisteSchema(dizionario, schema) {
				ricerca(dizionario, schema)
				
			} else {
				fmt.Println("Schema non esistente nel dizionario")
				
			}
		} else {
			fmt.Println(formatoErrato, "r")

			
		}
	case "d": // STAMPA DISTANZA DI EDITING
		if len(campi) != 3 {
			fmt.Println(formatoErrato, "d")
			return
		}

		x := campi[1]
		y := campi[2]
		distanza := distanzaLevenshtein(x, y)
		fmt.Println(distanza)

	default:
		fmt.Println("Comando non riconosciuto")
	}
}

func main() {

	fmt.Println("\n___   ---   ===   ^^^   ***   |||||   ***   ^^^   ===   ---   ___   ---   ===   ^^^   ***   |||||   ***   ^^^   ===   ---   ___\n",
	"\n	PROGETTO \"PAROLE E CATENE DI PAROLE\", LABORATORIO DI ALGORITMI E STRUTTURE DATI\n",
	"\nGestione di un dizionario di parole e schemi.\n",
	"Comandi:\n",
	"c ------> Crea un nuovo dizionario vuoto (eliminando l'eventuale dizionario già esistente).\n",
	"t ------> Termina esecuzione.\n",
	"c file -> Inserisce nel dizionario le parole e/o gli schemi contenuti nel file \"file\".\n",
	"p ------> Stampa tutte le parole del dizionario.\n",
	"s ------> Stampa tutti gli schemi del dizionario.\n",
	"i w ----> Inserisce nel dizionario la parola / lo schema w.\n",
	"e w ----> Elimina dal dizionario la parola / lo schema w.\n",
	"r S ----> Stampa lo schema S e poi l'insieme di tutte le parole nel dizionario che sono compatibili con lo schema S.\n",
	"d x y --> Stampa la distanza di editing fra le due parole x e y.\n",
	"c x y --> Stampa una catena di lunghezza minima tra x e y di parole nel dizionario.\n",
	"\n‾‾‾   ---   ===   vvv   ***   |||||   ***   vvv   ===   ---   ‾‾‾   ---   ===   vvv   ***   |||||   ***   vvv   ===   ---   ‾‾‾\n",
	"\nInserisci i comandi: ")


	scanner := bufio.NewScanner(os.Stdin)
	var d dizionario = newDizionario()
	for scanner.Scan() {
		linea := scanner.Text()
		if linea == "" {
			continue
		}
		// esegue il comando sulla linea letta
		esegui(d, linea)
	}

	err := scanner.Err();
	if  err != nil {
		fmt.Fprintln(os.Stderr, "Errore di lettura:", err)
	}
}