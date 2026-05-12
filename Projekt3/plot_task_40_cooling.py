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
    
    # Regex uwzględniający nazwy schematów jako parametry w nawiasach
    pattern = re.compile(r"^\s*([a-zA-Z0-9_]+)(?:\.atsp)?\s*\((.*)\)\s*\|\s*([\d.]+)\s*\|\s*([\d.]+)\s*\|\s*([\d.]+)%?\s*$")
    
    for line in lines:
        match = pattern.match(line)
        if match:
            instance = match.group(1)
            cooling = match.group(2).strip()
            avg_cost = float(match.group(3))
            avg_time = float(match.group(4))
            prd = float(match.group(5))
            
            if cooling in ["Geometric", "Linear", "Lundy-Mees"]:
                data.append({
                    "instance": instance,
                    "cooling": cooling,
                    "avg_cost": avg_cost,
                    "avg_time_ms": avg_time,
                    "prd": prd
                })
            
    if not data:
        print("Błąd: Nie udało się sparsować żadnych danych. Upewnij się, że format jest poprawny.")
        sys.exit(1)
        
    df = pd.DataFrame(data)
    
    # Ustalamy kolejność schematów
    order = ["Geometric", "Linear", "Lundy-Mees"]
    
    os.makedirs('output', exist_ok=True)
    
    cols = [c for c in order if c in df['cooling'].unique()]
    
    # a) Zgrupowany wykres słupkowy PRD
    pivot_prd = df.pivot_table(index='instance', columns='cooling', values='prd')
    pivot_prd = pivot_prd.reindex(columns=cols)
    
    pivot_prd.plot(kind='bar', figsize=(10, 6))
    plt.title('Błąd PRD dla schematów chłodzenia')
    plt.xlabel('Instancja')
    plt.ylabel('Błąd PRD [%]')
    plt.grid(True, axis='y', linestyle='--', alpha=0.7)
    plt.xticks(range(len(pivot_prd.index)), pivot_prd.index, rotation=45, ha='right', fontsize=8)
    plt.legend(title='Schemat chłodzenia')
    plt.tight_layout()
    plt.savefig('output/task_40_cooling_1.pdf')
    plt.close()
    
    # b) Zgrupowany wykres czasu (liniowy)
    pivot_time = df.pivot_table(index='instance', columns='cooling', values='avg_time_ms')
    pivot_time = pivot_time.reindex(columns=cols)
    
    pivot_time.plot(kind='line', marker='o', figsize=(10, 6))
    plt.yscale('log')
    plt.title('Czas działania dla schematów chłodzenia (skala log)')
    plt.xlabel('Instancja')
    plt.ylabel('Średni czas [ms]')
    plt.grid(True, axis='y', linestyle='--', alpha=0.7)
    plt.xticks(range(len(pivot_time.index)), pivot_time.index, rotation=45, ha='right', fontsize=8)
    plt.legend(title='Schemat chłodzenia')
    plt.tight_layout()
    plt.savefig('output/task_40_cooling_2.pdf')
    plt.close()
    
    avg_prd = df.groupby('cooling')['prd'].mean().reindex(cols)
    avg_time = df.groupby('cooling')['avg_time_ms'].mean().reindex(cols)
    

    
    print(f"Zapisano wykresy rozdzielone: output/task_40_cooling_1.pdf, output/task_40_cooling_2.pdf")
    
    # Statystyki
    print(f"\n--- STATYSTYKI ---")
    for m in cols:
        print(f"Schemat '{m}': Średni PRD = {avg_prd[m]:.2f}%, Średni czas = {avg_time[m]:.2f} ms")
        
    fastest = avg_time.idxmin()
    lowest_prd = avg_prd.idxmin()
    print(f"\nNajszybszy schemat: {fastest}")
    print(f"Schemat z najniższym średnim PRD: {lowest_prd}")

if __name__ == "__main__":
    main()
