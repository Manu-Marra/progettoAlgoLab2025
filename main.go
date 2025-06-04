package main

import (
    "bufio"
    "fmt"
    "os"
    "strings"
)

// tipo dizionario rappresenta l’intero dizionario
type dizionario struct {
    parole  map[string]struct{}
    schemi  map[string]struct{}
}

// crea() - crea un nuovo dizionario vuoto (azzera il dizionario)
func crea(d *dizionario) {
    d.parole = make(map[string]struct{})
    d.schemi = make(map[string]struct{})
}

// newDizionario() - crea e restituisce un nuovo dizionario vuoto
func newDizionario() dizionario {
    d := dizionario{}
    crea(&d)
    return d
}

func (d dizionario) inserisci(w string) {
    if contieneMaiuscola(w) {
        // è uno schema
        _, exists := d.schemi[w]  // verifica se w è già presente in schemi
        if !exists {
            d.schemi[w] = struct{}{}
        }
    } else {
        // è una parola
        _, exists := d.parole[w]  // verifica se w è già presente in parole
        if !exists {
            d.parole[w] = struct{}{}
        }
    }
}

func (d dizionario) carica(file string) {
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
        d.inserisci(w)
    }
    // Ignoriamo scanner.Err() per non bloccare l'esecuzione in caso di errori di lettura
}

func (d dizionario) stampa_parole() {
    fmt.Println("[")
    for parola := range d.parole {
        fmt.Println(parola)
    }
    fmt.Println("]")

}

func(d dizionario) stampa_schemi() {
    fmt.Println("[")
    for schema := range d.schemi {
        fmt.Println(schema)
    }
    fmt.Println("]")
}

func (d dizionario) elimina(w string) {
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

            curr[j] = min(curr[j-1]+1,prev[j]+1,prev[j-1]+costo)
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


func (d dizionario) ricerca(schema string) {
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
    return distanzaLevenshtein(x,y) == 1
}

func (d dizionario)catena(x, y string) {

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


func esegui(d dizionario, s string) {
    campi := strings.Fields(s)
    if len(campi) == 0 {
        return
    }

    switch campi[0] {
    case "c":
        if len(campi) == 1 {
            // solo "c" → crea dizionario vuoto
            crea(&d)
        } else if len(campi) == 2 {
            d.carica(campi[1])

        } else if len(campi) == 3 {
            // "c x y" → catena(x, y)
            x, y := campi[1], campi[2]
            d.catena(x, y)
        } else {
            fmt.Println("Comando errato")
        }

    case "t":
        fmt.Println("Esecuzione terminata")
        os.Exit(0)
    
    case "p":
        d.stampa_parole()
    case "s":
        d.stampa_schemi()
    case "i":
        d.inserisci(campi[1])
    case "e":
        d.elimina(campi[1])
    case "r":
        if len(campi) == 2 {
            schema := campi[1]
            d.ricerca(schema)
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
    diz := newDizionario()
    scanner := bufio.NewScanner(os.Stdin)

    for scanner.Scan() {
        linea := scanner.Text()
        if linea == "" {
            continue
        }
        // esegue il comando sulla linea letta
        esegui(diz, linea)
    }

    if err := scanner.Err(); err != nil {
        fmt.Fprintln(os.Stderr, "Errore di lettura:", err)
    }
}