import pandas as pd
import matplotlib.pyplot as plt
import seaborn as sns
import os

def plot_optimization_graphs():
    # Read the CSV file
    df = pd.read_csv('results/optimization_analysis.csv')
    
    # Set style
    sns.set_theme(style="whitegrid")
    os.makedirs('results', exist_ok=True)

    # 1. Bandwidth Comparison
    plt.figure(figsize=(10, 6))
    bandwidth_data = df[df['Test Type'].str.contains('Bandwidth')]
    sizes = [int(size.split()[1].replace('B','')) for size in bandwidth_data['Test Type']]
    
    sns.lineplot(x=sizes, y=bandwidth_data['Unoptimized'], marker='o', label='Unoptimized')
    sns.lineplot(x=sizes, y=bandwidth_data['Optimized'], marker='s', label='Optimized')
    
    plt.xscale('log')
    plt.xlabel('Message Size (bytes)')
    plt.ylabel('Bandwidth (MB/s)')
    plt.title('Bandwidth Comparison: Optimized vs Unoptimized')
    plt.grid(True)
    plt.tight_layout()
    plt.savefig('results/bandwidth_comparison.png', dpi=300, bbox_inches='tight')
    plt.close()

    # 2. Marshal Time Comparison
    plt.figure(figsize=(12, 6))
    marshal_data = df[df['Test Type'].str.contains('Marshal')]
    marshal_times = marshal_data[marshal_data['Metric'].str.contains('Marshal Time')]
    
    # Get unique data types from Test Type column
    data_types = [t.split()[1] for t in marshal_times['Test Type'].unique()]
    x = range(len(data_types))
    width = 0.35
    
    plt.bar([i - width/2 for i in x], marshal_times['Unoptimized'], width, 
            label='Unoptimized', color=sns.color_palette()[0])
    plt.bar([i + width/2 for i in x], marshal_times['Optimized'], width, 
            label='Optimized', color=sns.color_palette()[1])
    
    plt.xlabel('Data Type')
    plt.ylabel('Time (µs)')
    plt.title('Marshal Time Comparison')
    plt.xticks(x, data_types)  # Use actual data types from the data
    plt.legend()
    plt.grid(True, alpha=0.3)
    plt.tight_layout()
    plt.savefig('results/marshal_comparison.png', dpi=300, bbox_inches='tight')
    plt.close()

    # 3. RTT Metrics Comparison
    plt.figure(figsize=(10, 6))
    rtt_data = df[df['Test Type'] == 'RTT']
    x = range(len(rtt_data))
    
    plt.bar([i - width/2 for i in x], rtt_data['Unoptimized'], width, 
            label='Unoptimized', color=sns.color_palette()[0])
    plt.bar([i + width/2 for i in x], rtt_data['Optimized'], width, 
            label='Optimized', color=sns.color_palette()[1])
    
    plt.xlabel('RTT Metric')
    plt.ylabel('Time (µs)')
    plt.title('RTT Metrics Comparison')
    plt.xticks(x, rtt_data['Metric'])
    plt.legend()
    plt.grid(True, alpha=0.3)
    plt.tight_layout()
    plt.savefig('results/rtt_comparison.png', dpi=300, bbox_inches='tight')
    plt.close()

    # 4. Performance Improvement Percentage
    plt.figure(figsize=(12, 6))
    
    # Create a color column for the DataFrame
    df['Color'] = ['Improved' if x >= 0 else 'Degraded' for x in df['Difference (%)']]
    
    # Use proper seaborn syntax with hue parameter
    sns.barplot(
        data=df,
        x=range(len(df)),
        y='Difference (%)',
        hue='Color',
        palette={'Improved': 'green', 'Degraded': 'red'},
        legend=False
    )
    
    plt.axhline(y=0, color='black', linestyle='-', alpha=0.3)
    plt.xlabel('Test Cases')
    plt.ylabel('Performance Difference (%)')
    plt.title('Optimization Impact (%)')
    plt.xticks(
        range(len(df)), 
        df['Test Type'] + ' - ' + df['Metric'],
        rotation=45,
        ha='right'
    )
    plt.tight_layout()
    plt.grid(True, alpha=0.3)
    plt.savefig('results/optimization_impact.png', dpi=300, bbox_inches='tight')
    plt.close()

    print("Graphs have been generated in the results directory")

if __name__ == "__main__":
    plot_optimization_graphs()