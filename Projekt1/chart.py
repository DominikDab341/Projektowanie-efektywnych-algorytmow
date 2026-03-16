import matplotlib.pyplot as plt


def draw_charts():
    # ==========================================
    # WYKRES 1: Czas wykonania Brute-Force
    # ==========================================

    # Dane (pomijamy N=8, ponieważ czas 0s nie istnieje na skali logarytmicznej)
    N_bf = [9, 10, 11, 12, 13, 14]

    # Czasy przeliczone na sekundy
    # 618.4 µs, 7.0205 ms, 56.3894 ms, 607.946 ms, 7.51 s, 1m 43.6 s (103.628 s)
    time_bf = [618.4e-6, 7.0205e-3, 56.3894e-3, 607.946e-3, 7.512642, 103.6279638]

    plt.figure(figsize=(8, 5))
    plt.plot(N_bf, time_bf, marker='o', linestyle='-', color='#1f77b4', linewidth=2, markersize=8)

    plt.yscale('log') # Skala logarytmiczna dla osi Y
    plt.xlabel('Rozmiar problemu (N)', fontsize=12)
    plt.ylabel('Czas wykonania [s]', fontsize=12)
    plt.title('Czas działania algorytmu dokładnego (Brute-Force)', fontsize=14)

    # Dodanie siatki pomocniczej dla skali logarytmicznej
    plt.grid(True, which="both", ls="--", alpha=0.5)
    plt.xticks(N_bf) # Wymuszenie wyświetlania tylko całkowitych wartości N na osi X

    plt.tight_layout()
    plt.savefig('czas_bruteforce.pdf') # Zapis do wektorowego PDF (idealnego do LaTeX)
    plt.show()

    # ==========================================
    # WYKRES 2: Błędy algorytmów heurystycznych
    # ==========================================

    N_heur = [10, 11, 12, 13, 14]

    # Dane z Twojej drugiej tabeli
    err_nn = [45.68, 52.86, 64.83, 60.34, 68.35]
    err_rnn = [14.49, 17.55, 22.03, 22.54, 24.83]
    err_rand = [89.36, 111.03, 127.12, 154.50, 175.40]

    plt.figure(figsize=(8, 5))

    # Rysowanie trzech linii z różnymi znacznikami
    plt.plot(N_heur, err_nn, marker='s', linestyle='-', color='#ff7f0e', linewidth=2, markersize=8, label='NN')
    plt.plot(N_heur, err_rnn, marker='o', linestyle='-', color='#2ca02c', linewidth=2, markersize=8, label='RNN')
    plt.plot(N_heur, err_rand, marker='^', linestyle='-', color='#d62728', linewidth=2, markersize=8, label='Random')

    plt.xlabel('Rozmiar problemu (N)', fontsize=12)
    plt.ylabel('Średni błąd względny [%]', fontsize=12)
    plt.title('Jakość algorytmów heurystycznych', fontsize=14)

    plt.grid(True, linestyle='--', alpha=0.7)
    plt.xticks(N_heur)
    plt.legend(fontsize=11)

    plt.tight_layout()
    plt.savefig('bledy_heurystyki.pdf')
    plt.show()

if __name__ == "__main__":
    draw_charts()
    