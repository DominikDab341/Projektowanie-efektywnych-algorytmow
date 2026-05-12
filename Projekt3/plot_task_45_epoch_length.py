import sys
import os
import re
import pandas as pd
import numpy as np
import matplotlib.pyplot as plt

def main():
    try:
        with open('wyniki.txt', 'r', encoding='utf-8') as f:
            lines = f.readlines()
    except FileNotFoundError:
        print("Błąd: Nie znaleziono pliku 'wyniki.txt'.")
        sys.exit(1)
    data = []
    
    pattern = re.compile(r"^\s*([a-zA-Z0-9_]+)(?:\.atsp)?\s*\((.*)\)\s*\|\s*([\d.]+)\s*\|\s*([\d.]+)\s*\|\s*([\d.]+)%?\s*$")
    
    for line in lines:
        match = pattern.match(line)
        if match:
            instance = match.group(1)
            epoch_label = match.group(2).strip()
            avg_cost = float(match.group(3))
            avg_time = float(match.group(4))
            prd = float(match.group(5))
            
            if "Rozmiar x" in epoch_label:
                data.append({
                    "instance": instance,
                    "epoch_label": epoch_label,
                    "avg_cost": avg_cost,
                    "avg_time_ms": avg_time,
                    "prd": prd
                })
            
    if not data:
        print("Błąd: Nie udało się sparsować żadnych danych. Upewnij się, że format jest poprawny.")
        sys.exit(1)
        
    df = pd.DataFrame(data)
    
    # Funkcja pomocnicza do wyciągania liczby (mnożnika) z tekstu "Rozmiar x 10"
    def get_multiplier(label):
        m = re.search(r'\d+', label)
        return int(m.group(0)) if m else 1
        
    df['epoch_multiplier'] = df['epoch_label'].apply(get_multiplier)
    
    order = ["Rozmiar x 1", "Rozmiar x 10", "Rozmiar x 100"]
    
    os.makedirs('output', exist_ok=True)
    
    cols = [c for c in order if c in df['epoch_label'].unique()]
    

    
    # Średnie dla wykresów liniowych (sortowane według mnożnika)
    avg_prd = df.groupby('epoch_multiplier')['prd'].mean().sort_index()
    avg_time = df.groupby('epoch_multiplier')['avg_time_ms'].mean().sort_index()
    
    # c) Średni PRD vs mnożnik epoki (linia)
    plt.figure(figsize=(8, 5))
    plt.plot(avg_prd.index, avg_prd.values, marker='o', color='tab:green', linewidth=2)
    plt.xscale('log')
    plt.title('Średni błąd PRD w funkcji mnożnika długości epoki')
    plt.xlabel('Mnożnik długości epoki')
    plt.ylabel('Średni błąd PRD [%]')
    plt.grid(True, linestyle='--', alpha=0.7)
    plt.tight_layout()
    plt.savefig('output/task_45_epoch_length_1.pdf')
    plt.close()
    
    # d) Średni Czas vs mnożnik epoki (linia, log-log)
    plt.figure(figsize=(8, 5))
    plt.plot(avg_time.index, avg_time.values, marker='o', color='tab:red', linewidth=2)
    plt.xscale('log')
    plt.yscale('log')
    plt.title('Średni czas w funkcji mnożnika długości epoki (skala log-log)')
    plt.xlabel('Mnożnik długości epoki')
    plt.ylabel('Średni czas [ms]')
    plt.grid(True, linestyle='--', alpha=0.7)
    plt.tight_layout()
    plt.savefig('output/task_45_epoch_length_2.pdf')
    plt.close()
    

    
    print(f"Zapisano wykresy rozdzielone: output/task_45_epoch_length_1.pdf, output/task_45_epoch_length_2.pdf")
    
    # Statystyki w konsoli
    print(f"\n--- STATYSTYKI ---")
    for mult in avg_prd.index:
        print(f"Epoka x{mult}: Średni PRD = {avg_prd[mult]:.2f}%, Średni czas = {avg_time[mult]:.2f} ms")
        
    if 1 in avg_prd and 100 in avg_prd:
        red_prd = avg_prd[1] - avg_prd[100]
        print(f"\nRedukcja błędu PRD między x1 a x100: {red_prd:.2f} pkt. proc.")
        if avg_prd[1] > 0:
            print(f"Względna redukcja błędu PRD: {(red_prd / avg_prd[1] * 100):.2f}%")
            
    if 1 in avg_time and 100 in avg_time:
        inc_time = avg_time[100] / avg_time[1]
        print(f"Przybliżony wzrost czasu między x1 a x100: ~{inc_time:.1f}x")

if __name__ == "__main__":
    main()
