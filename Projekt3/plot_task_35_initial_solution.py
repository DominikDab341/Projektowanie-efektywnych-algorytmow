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
    
    # Regex parsowania - obsługuje różne spacje, łapie parametr jako wariant
    pattern = re.compile(r"^\s*([a-zA-Z0-9_]+)(?:\.atsp)?\s*\((.*)\)\s*\|\s*([\d.]+)\s*\|\s*([\d.]+)\s*\|\s*([\d.]+)%?\s*$")
    
    for line in lines:
        match = pattern.match(line)
        if match:
            instance = match.group(1)
            method = match.group(2).strip()
            avg_cost = float(match.group(3))
            avg_time = float(match.group(4))
            prd = float(match.group(5))
            
            if method in ["Losowe", "Zachłanne"]:
                data.append({
                    "instance": instance,
                    "method": method,
                    "avg_cost": avg_cost,
                    "avg_time_ms": avg_time,
                    "prd": prd
                })
            
    if not data:
        print("Błąd: Nie udało się sparsować żadnych danych. Upewnij się, że format jest poprawny.")
        sys.exit(1)
        
    df = pd.DataFrame(data)
    
    # Oczekiwana kolejność metod
    order = ["Losowe", "Zachłanne"]
    
    os.makedirs('output', exist_ok=True)
    
    # Przefiltruj metody, które faktycznie istnieją w danych (dla bezpieczeństwa)
    cols = [c for c in order if c in df['method'].unique()]
    
    # a) Zgrupowany wykres słupkowy dla PRD
    pivot_prd = df.pivot_table(index='instance', columns='method', values='prd')
    pivot_prd = pivot_prd.reindex(columns=cols) # Zachowaj kolejność
    
    pivot_prd.plot(kind='bar', figsize=(10, 6))
    plt.title('Błąd PRD dla metod rozwiązania początkowego')
    plt.xlabel('Instancja')
    plt.ylabel('Błąd PRD [%]')
    plt.grid(True, axis='y', linestyle='--', alpha=0.7)
    plt.xticks(range(len(pivot_prd.index)), pivot_prd.index, rotation=45, ha='right', fontsize=8)
    plt.legend(title='Metoda')
    plt.tight_layout()
    plt.savefig('output/task_35_initial_solution_1.pdf')
    plt.close()
    
    # b) Zgrupowany wykres czasu (liniowy)
    pivot_time = df.pivot_table(index='instance', columns='method', values='avg_time_ms')
    pivot_time = pivot_time.reindex(columns=cols)
    
    pivot_time.plot(kind='line', marker='o', figsize=(10, 6))
    plt.title('Czas działania dla metod rozwiązania początkowego')
    plt.xlabel('Instancja')
    plt.ylabel('Średni czas [ms]')
    plt.grid(True, axis='y', linestyle='--', alpha=0.7)
    plt.xticks(range(len(pivot_time.index)), pivot_time.index, rotation=45, ha='right', fontsize=8)
    plt.legend(title='Metoda')
    plt.tight_layout()
    plt.savefig('output/task_35_initial_solution_2.pdf')
    plt.close()
    

    
    print(f"Zapisano wykresy rozdzielone: output/task_35_initial_solution_1.pdf, output/task_35_initial_solution_2.pdf")
    
    # Statystyki w konsoli
    print(f"\n--- STATYSTYKI ---")
    avg_prd = df.groupby('method')['prd'].mean()
    avg_time = df.groupby('method')['avg_time_ms'].mean()
    for m in cols:
        if m in avg_prd and m in avg_time:
            print(f"Metoda '{m}': Średni PRD = {avg_prd[m]:.2f}%, Średni czas = {avg_time[m]:.2f} ms")
        
    if "Losowe" in avg_time and "Zachłanne" in avg_time:
        diff_time = (avg_time["Zachłanne"] - avg_time["Losowe"]) / avg_time["Losowe"] * 100
        print(f"\nRóżnica procentowa czasu (Zachłanne względem Losowe): {diff_time:+.2f}%")
        
    if "Losowe" in avg_prd and "Zachłanne" in avg_prd:
        best_method = "Zachłanne" if avg_prd["Zachłanne"] < avg_prd["Losowe"] else "Losowe"
        print(f"Metoda z niższym średnim PRD: {best_method}")

if __name__ == "__main__":
    main()
