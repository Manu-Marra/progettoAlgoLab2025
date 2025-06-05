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

func inserisci(d dizionario, w string) {

	if contieneMaiuscola(w) {
		// è uno schema
		_, exists := d.schemi[w] // verifica se w è già presente in schemi
		if !exists {
			d.schemi[w] = struct{}{}
		}
	} else {
		// è una parola
		_, exists := d.parole[w] // verifica se w è già presente in parole
		if !exists {
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
	// Ignoriamo scanner.Err() per non bloccare l'esecuzione in caso di errori di lettura
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
		_, found := d.schemi[w]
		if found {
			delete(d.schemi, w)
		}
	} else {
		// È una parola
		_, found := d.parole[w]
		if found {
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
			val, exists := mappa[c]
			if exists {
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
			// Se già visitata, skippo
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

	campi := strings.Fields(s)
	if len(campi) == 0 {
		return
	}

	switch campi[0] {
	case "c":
		if len(campi) == 1 {
			crea(dizionario)
			

		} else if len(campi) == 2 { // CARICA
			carica(dizionario, campi[1])

		} else if len(campi) == 3 { // CATENA
			// "c x y" → catena(x, y)
			x, y := campi[1], campi[2]
			catena(dizionario, x, y)

		} else {
			fmt.Println("Comando errato")
		}

	case "t":
		fmt.Println("Esecuzione terminata")
		os.Exit(0)

	case "p":
		stampa_parole(dizionario)
	case "s":
		stampa_schemi(dizionario)
	case "i":
		if isValida(campi[1]) {
			inserisci(dizionario, campi[1])
		} else {
			fmt.Println("Formato errato")
		}
	case "e":
		elimina(dizionario, campi[1])
	case "r":
		if len(campi) == 2 {
			schema := campi[1]
			ricerca(dizionario, schema)
		}
	case "d":
		if len(campi) < 3 {
			fmt.Println("Manca uno degli argomenti per il comando d x y")
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
	// dizionario := newDizionario()
	scanner := bufio.NewScanner(os.Stdin)
	var d dizionario = newDizionario()
	for scanner.Scan() {
		linea := scanner.Text()
		if linea == "" {
			continue
		}
		// esegue il comando sulla linea letta
		// esegui(dizionario, linea) -------
		esegui(d, linea)
		// fmt.Println("dizionario in main")
		// dizionario.stampa_parole()
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Errore di lettura:", err)
	}
}