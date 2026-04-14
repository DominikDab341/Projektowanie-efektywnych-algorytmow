import matplotlib.pyplot as plt
import numpy as np

def draw_charts():
    # ==========================================
    # Porównanie czasów wykonania algorytmów Branch and Bound
    # ==========================================

    # Dane dla Breadth-First-Search (od N=8)
    N_bfs = [8, 9, 10, 11]
    time_bfs = [1.55, 13.84, 126.65, 1344.36]

    # Obliczenie regresji wykładniczej (liniowa na skali logarytmicznej) dla BFS
    log_time_bfs = np.log10(time_bfs)
    slope, intercept = np.polyfit(N_bfs, log_time_bfs, 1)
    
    # Ekstrapolacja dla N od 11 do 40
    N_bfs_extrapolated = list(range(11, 41))
    time_bfs_extrapolated = 10**(np.array(N_bfs_extrapolated) * slope + intercept)

    # Dane dla Brute-Force (od N=9) przeliczone na ms
    N_bf = [9, 10, 11, 12, 13, 14]
    time_bf_sec = [502.8e-6, 4.5993e-3, 57.8369e-3, 614.8999e-3, 8.2272952, 114.3125687]
    time_bf_ms = [t * 1000 for t in time_bf_sec]

    # Obliczenie regresji wykładniczej dla Brute-Force
    log_time_bf = np.log10(time_bf_ms)
    slope_bf, intercept_bf = np.polyfit(N_bf, log_time_bf, 1)
    
    # Ekstrapolacja dla Brute-Force od N=14 do 40
    N_bf_extrapolated = list(range(14, 41))
    time_bf_extrapolated = 10**(np.array(N_bf_extrapolated) * slope_bf + intercept_bf)

    # Dane dla Best-First-Search (INF) i (NN) od N=8 do 40
    N_best = [8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40]
    time_best_inf = [0.09, 0.09, 0.12, 0.08, 0.41, 0.77, 1.01, 1.82, 2.50, 4.75, 8.95, 16.99, 24.69, 36.14, 60.68, 81.01, 183.23, 217.15, 312.71, 656.85, 708.03, 1647.12, 3466.94, 3913.73, 6705.21, 8863.03, 14379.00, 18416.48, 27598.20, 31270.79, 38708.02, 50970.72, 54902.73]
    time_best_nn = [0.06, 0.10, 0.13, 0.28, 0.30, 0.41, 0.97, 2.07, 2.80, 4.92, 8.89, 16.36, 24.20, 36.16, 59.62, 79.54, 183.00, 213.68, 304.27, 649.27, 699.85, 1593.35, 3427.36, 3866.73, 6624.14, 8482.05, 14033.66, 17930.65, 27157.72, 30642.56, 38181.88, 50129.53, 53923.20]

    plt.figure(figsize=(10, 6))

    # Prawdziwe dane BFS
    plt.plot(N_bfs, time_bfs, marker='o', linestyle='-', color='#d62728', linewidth=2, markersize=6, label='Breadth-First')
    
    # Ekstrapolacja BFS
    plt.plot(N_bfs_extrapolated, time_bfs_extrapolated, marker='', linestyle='--', color='#d62728', linewidth=2, alpha=0.5, label='Breadth-First (Symulacja)')

    # Dane Brute-Force
    plt.plot(N_bf, time_bf_ms, marker='v', linestyle='-', color='#9467bd', linewidth=2, markersize=6, label='Brute-Force')
    plt.plot(N_bf_extrapolated, time_bf_extrapolated, marker='', linestyle='--', color='#9467bd', linewidth=2, alpha=0.5, label='Brute-Force (Symulacja)')

    # Dane Best-First
    plt.plot(N_best, time_best_inf, marker='s', linestyle='-', color='#1f77b4', linewidth=2, markersize=6, label='Best-First (INF)')
    plt.plot(N_best, time_best_nn, marker='^', linestyle='-', color='#2ca02c', linewidth=2, markersize=6, label='Best-First (NN)')

    plt.yscale('log')
    plt.xlabel('Rozmiar problemu (N)', fontsize=12)
    plt.ylabel('Czas wykonania [ms]', fontsize=12)
    plt.title('Porównanie czasów działania algorytmów Branch and Bound', fontsize=14)

    plt.grid(True, which="both", ls="--", alpha=0.5)
    plt.xlim(8, 41)
    plt.ylim(top=10**10)  # Ograniczenie góry wykresu do 10^10
    
    # Dodajmy jeszcze gęstsze podziałki na osi X, żeby N=8 było wyraźnie zaznaczone
    plt.xticks(range(7, 42, 2))

    plt.legend(fontsize=11)

    plt.tight_layout()
    plt.savefig('czas_bb.pdf')
    plt.show()

if __name__ == "__main__":
    draw_charts()