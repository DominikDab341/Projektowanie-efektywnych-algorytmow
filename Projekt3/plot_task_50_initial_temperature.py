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
    
    in_section = False
    for line in lines:
        if "TEST 5.0" in line:
            in_section = True
            continue
            
        if not in_section:
            continue
            
        match = pattern.match(line)
        if match:
            instance = match.group(1)
            temp_label = match.group(2).strip()
            avg_cost = float(match.group(3))
            avg_time = float(match.group(4))
            prd = float(match.group(5))
            
            if temp_label in ["100", "1000", "10000", "Auto(Calc)"]:
                data.append({
                    "instance": instance,
                    "temp_label": temp_label,
                    "avg_cost": avg_cost,
                    "avg_time_ms": avg_time,
                    "prd": prd
                })
            
    if not data:
        print("Błąd: Nie udało się sparsować żadnych danych. Upewnij się, że format jest poprawny.")
        sys.exit(1)
        
    df = pd.DataFrame(data)
    
    # Przewidziana kolejność
    order = ["100", "1000", "10000", "Auto(Calc)"]
    
    os.makedirs('output', exist_ok=True)
    
    cols = [c for c in order if c in df['temp_label'].unique()]
    
    # a) Zgrupowany wykres słupkowy PRD
    pivot_prd = df.pivot_table(index='instance', columns='temp_label', values='prd')
    pivot_prd = pivot_prd.reindex(columns=cols)
    
    pivot_prd.plot(kind='bar', figsize=(10, 6))
    plt.title('Błąd PRD dla temperatur początkowych')
    plt.xlabel('Instancja')
    plt.ylabel('Błąd PRD [%]')
    plt.grid(True, axis='y', linestyle='--', alpha=0.7)
    plt.xticks(range(len(pivot_prd.index)), pivot_prd.index, rotation=45, ha='right', fontsize=8)
    plt.legend(title='Temperatura')
    plt.tight_layout()
    plt.savefig('output/task_50_initial_temperature_1.pdf')
    plt.close()
    
    # b) Zgrupowany wykres czasu (liniowy)
    pivot_time = df.pivot_table(index='instance', columns='temp_label', values='avg_time_ms')
    pivot_time = pivot_time.reindex(columns=cols)
    
    pivot_time.plot(kind='line', marker='o', figsize=(10, 6))
    plt.title('Czas działania dla temperatur początkowych')
    plt.xlabel('Instancja')
    plt.ylabel('Średni czas [ms]')
    plt.grid(True, axis='y', linestyle='--', alpha=0.7)
    plt.xticks(range(len(pivot_time.index)), pivot_time.index, rotation=45, ha='right', fontsize=8)
    plt.legend(title='Temperatura')
    plt.tight_layout()
    plt.savefig('output/task_50_initial_temperature_2.pdf')
    plt.close()
    
    avg_prd = df.groupby('temp_label')['prd'].mean().reindex(cols)
    avg_time = df.groupby('temp_label')['avg_time_ms'].mean().reindex(cols)
    

    
    print(f"Zapisano wykresy rozdzielone: output/task_50_initial_temperature_1.pdf, output/task_50_initial_temperature_2.pdf")
    
    # Statystyki
    print(f"\n--- STATYSTYKI ---")
    for m in cols:
        print(f"Temperatura '{m}': Średni PRD = {avg_prd[m]:.2f}%, Średni czas = {avg_time[m]:.2f} ms")
        
    lowest_prd = avg_prd.idxmin()
    fastest = avg_time.idxmin()
    print(f"\nTemperatura dająca najniższy średni PRD: {lowest_prd}")
    print(f"Najszybsza temperatura: {fastest}")
    
    if "100" in avg_time and "10000" in avg_time:
        diff = (avg_time["10000"] - avg_time["100"]) / avg_time["100"] * 100
        print(f"Różnica procentowa czasu (10000 względem 100): {diff:+.2f}%")

if __name__ == "__main__":
    main()
