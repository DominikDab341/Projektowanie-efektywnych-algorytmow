import matplotlib.pyplot as plt


def draw_charts():
    # ==========================================
    # WYKRES 1: Czas wykonania Brute-Force
    # ==========================================

    # Dane (pomijamy N=8, ponieważ czas 0s nie istnieje na skali logarytmicznej)
    N_bf = [9, 10, 11, 12, 13, 14]

    # Czasy przeliczone na sekundy
    # 502.8 µs, 4.5993 ms, 57.8369 ms, 614.8999 ms, 8.2272952 s, 1m54.3125687 s (114.3125687 s)
    time_bf = [502.8e-6, 4.5993e-3, 57.8369e-3, 614.8999e-3, 8.2272952, 114.3125687]

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
    err_nn = [52.92, 57.31, 64.97, 71.81, 67.78]
    err_rnn = [12.22, 18.62, 22.06, 24.66, 23.71]
    err_rand = [82.62, 101.25, 132.81, 150.05, 169.57]

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
    