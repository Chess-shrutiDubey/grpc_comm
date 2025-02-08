import pandas as pd
import matplotlib.pyplot as plt
import seaborn as sns
import json
import numpy as np
from pathlib import Path

def parse_duration(duration_str):
    """Convert duration string to milliseconds"""
    if 'ms' in duration_str:
        return float(duration_str.replace('ms', ''))
    elif 'µs' in duration_str:
        return float(duration_str.replace('µs', '')) / 1000
    return float(duration_str)

def analyze_rtt():
    """Analyze RTT performance"""
    results_dir = Path('../results/rtt')
    results = []
    
    for result_file in results_dir.glob('rtt_*.json'):
        with open(result_file, 'r') as f:
            data = json.load(f)
            # Convert RTTs to milliseconds
            rtts_ms = [parse_duration(rtt) for rtt in data['rtts']]
            results.extend([{
                'message_size': data['message_size'],
                'rtt_ms': rtt,
                'is_first': idx == 0
            } for idx, rtt in enumerate(rtts_ms)])
    
    df = pd.DataFrame(results)
    
    fig, (ax1, ax2, ax3) = plt.subplots(1, 3, figsize=(15, 5))
    
    # RTT by message size
    sns.boxplot(data=df, x='message_size', y='rtt_ms', ax=ax1)
    ax1.set_title('RTT by Message Size')
    ax1.set_xlabel('Message Size (bytes)')
    ax1.set_ylabel('RTT (ms)')
    
    # First vs subsequent RTTs
    sns.boxplot(data=df, x='is_first', y='rtt_ms', ax=ax2)
    ax2.set_title('First vs Subsequent RTTs')
    ax2.set_xticklabels(['Subsequent', 'First'])
    ax2.set_ylabel('RTT (ms)')
    
    # RTT distribution
    sns.histplot(data=df, x='rtt_ms', ax=ax3)
    ax3.set_title('RTT Distribution')
    ax3.set_xlabel('RTT (ms)')
    
    plt.tight_layout()
    plt.savefig('../results/rtt/rtt_analysis.png')
    
    # Calculate statistics by message size
    stats = df.groupby('message_size').agg({
        'rtt_ms': ['count', 'mean', 'std', 'min', '25%', '50%', '75%', 'max']
    })
    
    return stats

def analyze_bandwidth():
    """Analyze bandwidth performance"""
    results_dir = Path('../results/bandwidth')
    results = []
    
    for result_file in results_dir.glob('bandwidth_*.json'):
        with open(result_file, 'r') as f:
            data = json.load(f)
            results.append(data)
    
    df = pd.DataFrame(results)
    
    plt.figure(figsize=(10, 6))
    sns.lineplot(data=df, x='message_size', y='bandwidth_mbps')
    plt.xscale('log')
    plt.title('Bandwidth vs Message Size')
    plt.xlabel('Message Size (bytes)')
    plt.ylabel('Bandwidth (MB/s)')
    plt.grid(True)
    plt.savefig('../results/bandwidth/bandwidth_analysis.png')
    
    return df.groupby('message_size')['bandwidth_mbps'].describe()

def analyze_marshal():
    """Analyze marshalling overhead"""
    results_dir = Path('../results/marshal')
    results = []
    
    for result_file in results_dir.glob('marshal_*.json'):
        with open(result_file, 'r') as f:
            data = json.load(f)
            results.append(data)
    
    df = pd.DataFrame(results)
    
    plt.figure(figsize=(10, 6))
    sns.scatterplot(data=df, x='message_size', y='marshal_time_ns')
    plt.xscale('log')
    plt.yscale('log')
    plt.title('Marshal Time vs Message Size')
    plt.xlabel('Message Size (bytes)')
    plt.ylabel('Marshal Time (ns)')
    plt.grid(True)
    plt.savefig('../results/marshal/marshal_analysis.png')
    
    return df.groupby('message_size').agg({
        'marshal_time_ns': ['mean', 'std', 'min', 'max']
    })

if __name__ == '__main__':
    # Create results directories
    for dir_name in ['rtt', 'bandwidth', 'marshal']:
        Path(f'../results/{dir_name}').mkdir(parents=True, exist_ok=True)
    
    print("\n=== RTT Analysis ===")
    print(analyze_rtt())
    
    print("\n=== Bandwidth Analysis ===")
    print(analyze_bandwidth())
    
    print("\n=== Marshal Analysis ===")
    print(analyze_marshal())