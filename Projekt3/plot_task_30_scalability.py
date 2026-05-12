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
    
    # Gwarantuje, że w nawiasach są tylko cyfry (rozmiar)
    pattern = re.compile(r"^\s*([a-zA-Z0-9_]+)(?:\.atsp)?\s*\(\s*(\d+)\s*\)\s*\|\s*([\d.]+)\s*\|\s*([\d.]+)\s*\|\s*([\d.]+)%?\s*$")
    
    in_section = False
    for line in lines:
        if "TEST 3.0" in line:
            in_section = True
            continue
        elif "TEST 3.5" in line:
            break
            
        if not in_section:
            continue
            
        match = pattern.match(line)
        if match:
            # Wyciągnięcie odpowiednich grup z wyrażenia regularnego
            instance = match.group(1)
            param = match.group(2).strip()
            avg_cost = float(match.group(3))
            avg_time = float(match.group(4))
            prd = float(match.group(5))
            
            data.append({
                "instance": instance,
                "param": param,
                "avg_cost": avg_cost,
                "avg_time_ms": avg_time,
                "prd": prd
            })
            
    # Walidacja, czy jakiekolwiek dane zostały sparsowane
    if not data:
        print("Błąd: Nie udało się sparsować żadnych danych. Upewnij się, że format jest poprawny.")
        sys.exit(1)
        
    df = pd.DataFrame(data)
    
    # Konwersja rozmiaru na liczbę całkowitą
    df['size'] = df['param'].astype(int)
    df = df.sort_values('size')
    
    # Utworzenie katalogu wyjściowego, jeśli nie istnieje
    os.makedirs('output', exist_ok=True)
    
    # Generowanie wykresów do PDF
    # a) Średni czas działania w funkcji rozmiaru
    plt.figure(figsize=(10, 6))
    plt.plot(range(len(df)), df['avg_time_ms'], marker='o', linestyle='-', color='tab:orange')
    plt.yscale('log')
    plt.title('Średni czas działania od rozmiaru instancji (skala log)')
    plt.xticks(range(len(df)), df['instance'], rotation=45, ha='right', fontsize=8)
    plt.xlabel('Instancja')
    plt.ylabel('Średni czas [ms]')
    plt.grid(True)
    plt.tight_layout()
    plt.savefig('output/task_30_scalability_1.pdf')
    plt.close()
    
    # b) Błąd PRD w funkcji rozmiaru instancji
    plt.figure(figsize=(10, 6))
    plt.plot(range(len(df)), df['prd'], marker='o', linestyle='-', color='tab:green')
    plt.title('Błąd PRD od rozmiaru instancji')
    plt.xticks(range(len(df)), df['instance'], rotation=45, ha='right', fontsize=8)
    plt.xlabel('Instancja')
    plt.ylabel('Błąd PRD [%]')
    plt.grid(True)
    plt.tight_layout()
    plt.savefig('output/task_30_scalability_2.pdf')
    plt.close()
    
    print(f"Zapisano wykresy rozdzielone: output/task_30_scalability_1.pdf, output/task_30_scalability_2.pdf")
    
    # Wyświetlenie statystyk w konsoli
    print(f"\n--- STATYSTYKI ---")
    print(f"Liczba instancji: {len(df)}")
    print(f"Średni czas: {df['avg_time_ms'].mean():.2f} ms")
    print(f"Średni PRD: {df['prd'].mean():.2f}%")
    
    max_prd_idx = df['prd'].idxmax()
    print(f"Maksymalny PRD: {df.loc[max_prd_idx, 'prd']:.2f}% (Instancja: {df.loc[max_prd_idx, 'instance']})")

if __name__ == "__main__":
    main()
