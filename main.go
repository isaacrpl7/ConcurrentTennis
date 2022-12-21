package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Player struct {
	estado    string //'Esperando', 'Posse' ou 'Parar'
	pontuacao int
	nome      string
}

func start(player *Player, enemy *Player, command chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case cmd := <-command: //Consumindo o channel e guardando em cmd
			//fmt.Println(cmd)
			switch cmd {
			case "Parar":
				return
			case "Jogar":
				player.estado = "Posse"
			case "Pontuar":
				player.pontuacao = player.pontuacao + 1
				player.estado = "Posse"

				//Determinando que um jogador ganhou do adversário
				if adjust := player.pontuacao - 2; player.pontuacao >= 4 && adjust >= enemy.pontuacao {
					fmt.Println(player.nome, "GANHOU O JOGO!")

					fmt.Print(player.nome)
					fmt.Print(": ")
					fmt.Println(player.pontuacao, "pontos!")

					fmt.Print(enemy.nome)
					fmt.Print(": ")
					fmt.Println(enemy.pontuacao, "pontos!")

					player.estado = "Parar"
					command <- "Parar"
				}
			default:
				player.estado = "Esperando"
			}
		default:
			if player.estado == "Posse" {
				// Jogar
				fmt.Println(player.nome, "jogando")
				time.Sleep(1 * time.Second)

				//Randomicamente determinando se o jogador acertou ou errou a jogada
				s1 := rand.NewSource(time.Now().UnixNano())
				r1 := rand.New(s1)
				n := r1.Intn(2)

				fmt.Println(player.nome, "possui pontuacao de", player.pontuacao) //Se for 0, errou, se for 1, acertou

				if n == 0 {
					fmt.Println(player.nome, "errou!")
					fmt.Println("")
					time.Sleep(1 * time.Second)

					player.estado = "Esperando"
					command <- "Pontuar"
				} else {
					fmt.Println(player.nome, "acertou!")
					fmt.Println("")
					time.Sleep(1 * time.Second)

					player.estado = "Esperando"
					command <- "Jogar"
				}
			}

			if player.estado == "Parar" {
				return
			}
		}
	}
}

func main() {
	var wg sync.WaitGroup //WaitGroup serve para criar um grupo de goRotinas e observar se elas finalizaram ou não a execução
	wg.Add(1)
	wg.Add(1)

	player1 := Player{"Posse", 0, "Jogador 1"}
	player2 := Player{"Esperando", 0, "Jogador 2"}

	command := make(chan string) //Fazendo o channel para passar mensagens entre as gorotinas

	//Iniciando as goRotinas
	go start(&player1, &player2, command, &wg)
	go start(&player2, &player1, command, &wg)

	wg.Wait()
}
