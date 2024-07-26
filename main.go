package main

import (
	"tictactoe/server"
	"tictactoe/stats"
)

func main() {

	stats.LoadStatistics()
	server.StartRouter()

	defer stats.SaveStatistics()
}

/*

Se poate configura dimensiunea ( mereu patrat);
Se poate configura caracterul de baza(default "_").
Metode http pt aceasta configurare
Playerii isi pot inregistra numele si caracterul preferat cu care joaca. Daca acestia isi schimba caracterul preferat in timpul jocului, schimbarea trebuie sa poata fi vizualizata in getState sau alta printare.  Metode http.
3. Se vrea un http call "getStatistics" care enunta fiecare player( per nume!) cate jocuri a castigat.

matrice de pointeri

4. sa se poata rula mai multe jocuri simultan (mai multi playeri)

5. jucatorii sa isi poata schimba simbolul in timpul meciului si simbolul tablei sa poata fi schimbat

6. sa se poata trimite cereri de la administratia jocului pentru urmatoarele meciuri de ce dimensiune sa fie

7. sa poata incepe oricare player jocul


------


*/
